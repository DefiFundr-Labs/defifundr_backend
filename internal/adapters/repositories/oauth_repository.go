// repositories/oauth_repository.go
package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/core/domain"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type OAuthRepository struct {
	store       db.Queries
	logger      logging.Logger
	jwksCache   map[string]*keyfunc.JWKS
	cacheExpiry map[string]time.Time
	cacheMutex  sync.RWMutex
}

func NewOAuthRepository(store db.Queries, logger logging.Logger) *OAuthRepository {
	return &OAuthRepository{
		store:       store,
		logger:      logger,
		jwksCache:   make(map[string]*keyfunc.JWKS),
		cacheExpiry: make(map[string]time.Time),
	}
}

// getJWKS retrieves and caches a JWKS for token validation
func (r *OAuthRepository) getJWKS(jwksURL string) (*keyfunc.JWKS, error) {
	// Check cache first with read lock
	r.cacheMutex.RLock()
	jwks, found := r.jwksCache[jwksURL]
	expiry, _ := r.cacheExpiry[jwksURL]
	r.cacheMutex.RUnlock()

	// Return cached JWKS if still valid
	if found && time.Now().Before(expiry) {
		return jwks, nil
	}

	// Acquire write lock to update cache
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	// Double-check after acquiring lock
	jwks, found = r.jwksCache[jwksURL]
	expiry, _ = r.cacheExpiry[jwksURL]
	if found && time.Now().Before(expiry) {
		return jwks, nil
	}

	// Fetch new JWKS
	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			r.logger.Error("Error refreshing JWKS", err, map[string]interface{}{
				"jwks_url": jwksURL,
			})
		},
	}

	newJWKS, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	// Update cache
	r.jwksCache[jwksURL] = newJWKS
	r.cacheExpiry[jwksURL] = time.Now().Add(time.Hour)

	return newJWKS, nil
}

// ValidateWebAuthToken validates a Web3Auth token
func (r *OAuthRepository) ValidateWebAuthToken(ctx context.Context, tokenString string) (*domain.Web3AuthClaims, error) {
	jwksURL := "https://api-auth.web3auth.io/jwks"
	jwks, err := r.getJWKS(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	// Parse the token with claims
	claims := &domain.Web3AuthClaims{}
	parser := jwtv4.NewParser(jwtv4.WithValidMethods([]string{"ES256"}))
	token, err := parser.ParseWithClaims(tokenString, claims, jwks.Keyfunc)

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, errors.New("token has expired, please re-authenticate")
		}
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify required claims
	if claims.Verifier == "" || claims.VerifierID == "" {
		return nil, errors.New("missing required Web3Auth claims")
	}

	// Verify the issuer
	if claims.Issuer != "https://api-auth.web3auth.io" {
		return nil, fmt.Errorf("invalid issuer: %v", claims.Issuer)
	}

	// Log successful validation
	r.logger.Info("Successfully validated Web3Auth token", map[string]interface{}{
		"email":        claims.Email,
		"verifier":     claims.Verifier,
		"wallet_count": len(claims.Wallets),
	})

	return claims, nil
}

// GetUserInfoFromProviderToken extracts user information from an OAuth provider token
func (r *OAuthRepository) GetUserInfoFromProviderToken(ctx context.Context, provider string, token string) (*domain.User, error) {
	// For Web3Auth, validate token first
	if provider == string(domain.Web3AuthProvider) {
		claims, err := r.ValidateWebAuthToken(ctx, token)
		if err != nil {
			return nil, err
		}

		// Extract name parts
		firstName, lastName := extractNameFromClaims(claims)

		// Create user object with information from claims
		profileImage := claims.ProfileImage
		user := &domain.User{
			Email:          claims.Email,
			FirstName:      firstName,
			LastName:       lastName,
			ProfilePicture: &profileImage,
			AuthProvider:   string(mapVerifierToProvider(claims.Verifier)),
			ProviderID:     claims.VerifierID,
		}

		return user, nil
	}

	// For other providers, implement specific token validation
	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

// Helper function to extract name from claims
func extractNameFromClaims(claims *domain.Web3AuthClaims) (string, string) {
	if claims.Name == "" {
		return "User", ""
	}

	nameParts := strings.Split(claims.Name, " ")
	firstName := nameParts[0]

	var lastName string
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	return firstName, lastName
}

// Helper function to map Web3Auth verifier to provider
func mapVerifierToProvider(verifier string) domain.AuthProvider {
	lowerVerifier := strings.ToLower(verifier)

	if strings.Contains(lowerVerifier, "google") {
		return domain.GoogleProvider
	} else if strings.Contains(lowerVerifier, "facebook") {
		return domain.FacebookProvider
	} else if strings.Contains(lowerVerifier, "apple") {
		return domain.AppleProvider
	} else if strings.Contains(lowerVerifier, "twitter") {
		return domain.TwitterProvider
	} else if strings.Contains(lowerVerifier, "discord") {
		return domain.DiscordProvider
	}

	return domain.Web3AuthProvider
}
