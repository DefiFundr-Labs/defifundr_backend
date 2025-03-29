package sqlc

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// SeedSize defines the volume of data to be generated
type SeedSize string

const (
	SeedSizeSmall  SeedSize = "small"  // Few records for quick testing
	SeedSizeMedium SeedSize = "medium" // Moderate data volume
	SeedSizeLarge  SeedSize = "large"  // Large dataset for performance testing
)

// SeedOptions configures the seeding behavior
type SeedOptions struct {
	Size            SeedSize       // Size of the dataset to generate
	Tables          []string       // Specific tables to seed (empty means all)
	PreserveData    []string       // Tables to preserve existing data
	RandomSeed      int64          // Seed for random generator (0 means use time)
	CleanBeforeRun  bool           // Whether to clean tables before seeding
	Verbose         bool           // Enable verbose output
	UserCount       int            // Number of users to generate
	UserRoleRatios  map[string]int // Distribution of user roles
	TransactionDays int            // How many days of transaction history to generate
}

// DefaultSeedOptions returns default seeding options
func DefaultSeedOptions() SeedOptions {
	return SeedOptions{
		Size:            SeedSizeMedium,
		Tables:          []string{},
		PreserveData:    []string{},
		RandomSeed:      0,
		CleanBeforeRun:  true,
		Verbose:         true,
		UserCount:       10,
		UserRoleRatios:  map[string]int{"personal": 8, "business": 2},
		TransactionDays: 30,
	}
}

// Seeder handles database seeding operations
type Seeder struct {
	queries *Queries
	options SeedOptions
	randGen *rand.Rand
	users   []Users // Store generated users for relationships
}

// NewSeeder creates a new database seeder
func NewSeeder(queries *Queries, options SeedOptions) *Seeder {
	// Set up deterministic random number generator
	var seed int64
	if options.RandomSeed == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = options.RandomSeed
	}

	// Initialize random generator with seed for reproducibility
	source := rand.NewSource(seed)
	generator := rand.New(source)

	// Set counts based on size
	switch options.Size {
	case SeedSizeSmall:
		if options.UserCount == 10 { // Only override if not explicitly specified
			options.UserCount = 5
		}
	case SeedSizeLarge:
		if options.UserCount == 10 { // Only override if not explicitly specified
			options.UserCount = 50
		}
	}

	return &Seeder{
		queries: queries,
		options: options,
		randGen: generator,
		users:   []Users{},
	}
}

// SeedDB populates the database according to the configured options
func (s *Seeder) SeedDB(ctx context.Context) error {
	log.Println("Starting database seeding...")
	startTime := time.Now()

	// Clean tables if requested
	if s.options.CleanBeforeRun {
		if err := s.cleanTables(ctx); err != nil {
			return fmt.Errorf("failed to clean tables: %w", err)
		}
	}

	// Determine tables to seed
	tables := s.getTablesToSeed()

	// Seed tables in the correct order to maintain referential integrity
	var err error
	if contains(tables, "users") {
		if err = s.seedUsers(ctx); err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
	}

	if contains(tables, "sessions") {
		if err = s.seedSessions(ctx); err != nil {
			return fmt.Errorf("failed to seed sessions: %w", err)
		}
	}

	if contains(tables, "otp_verifications") {
		if err = s.seedOTPVerifications(ctx); err != nil {
			return fmt.Errorf("failed to seed OTP verifications: %w", err)
		}
	}

	if contains(tables, "user_device_tokens") {
		if err = s.seedUserDeviceTokens(ctx); err != nil {
			return fmt.Errorf("failed to seed user device tokens: %w", err)
		}
	}

	if contains(tables, "kyc_records") {
		if err = s.seedKYCRecords(ctx); err != nil {
			return fmt.Errorf("failed to seed KYC records: %w", err)
		}
	}

	if contains(tables, "transactions") {
		if err = s.seedTransactions(ctx); err != nil {
			return fmt.Errorf("failed to seed transactions: %w", err)
		}
	}

	duration := time.Since(startTime)
	log.Printf("Database seeding completed in %s", duration)
	return nil
}

// getTablesToSeed returns the list of tables to seed based on options
func (s *Seeder) getTablesToSeed() []string {
	if len(s.options.Tables) == 0 {
		// If no specific tables are requested, seed all tables
		return []string{
			"users",
			"sessions",
			"otp_verifications",
			"user_device_tokens",
			"kyc_records",
			"transactions",
		}
	}
	return s.options.Tables
}

// cleanTables removes existing data from tables
func (s *Seeder) cleanTables(ctx context.Context) error {
	log.Println("Cleaning tables before seeding...")

	// Define tables to clean in reverse dependency order
	tables := []string{
		"transactions",
		"kyc_records",
		"user_device_tokens",
		"otp_verifications",
		"sessions",
		"users",
	}

	// Filter out tables that should be preserved
	var tablesToClean []string
	for _, table := range tables {
		if !contains(s.options.PreserveData, table) {
			tablesToClean = append(tablesToClean, table)
		}
	}

	// Create a connection for executing raw SQL
	conn := s.queries.db

	// Disable foreign key checks temporarily
	if _, err := conn.Exec(ctx, "SET session_replication_role = 'replica';"); err != nil {
		return err
	}

	// Clean each table
	for _, table := range tablesToClean {
		if _, err := conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)); err != nil {
			return err
		}
		if s.options.Verbose {
			log.Printf("Cleaned table: %s", table)
		}
	}

	// Re-enable foreign key checks
	if _, err := conn.Exec(ctx, "SET session_replication_role = 'origin';"); err != nil {
		return err
	}

	return nil
}

// seedUsers creates user records
func (s *Seeder) seedUsers(ctx context.Context) error {
	log.Printf("Seeding %d users...", s.options.UserCount)

	// Sample data for generating realistic user profiles
	firstNames := []string{"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", "William", "Elizabeth",
		"David", "Barbara", "Richard", "Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Charles", "Karen",
		"Olivia", "Emma", "Ava", "Charlotte", "Sophia", "Amelia", "Isabella", "Mia", "Evelyn", "Harper"}

	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
		"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
		"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Lewis", "Robinson", "Walker"}

	nationalities := []string{"US", "UK", "CA", "AU", "DE", "FR", "JP", "CN", "IN", "BR", "MX", "ZA", "NG", "KE", "EG", "AE"}

	countries := []string{"United States", "United Kingdom", "Canada", "Australia", "Germany", "France", "Japan", "China", "India", "Brazil",
		"Mexico", "South Africa", "Nigeria", "Kenya", "Egypt", "United Arab Emirates"}

	jobRoles := []string{"Software Engineer", "Product Manager", "Data Scientist", "UI/UX Designer", "DevOps Engineer",
		"Project Manager", "Marketing Specialist", "Sales Representative", "Financial Analyst", "Customer Support",
		"Content Writer", "Graphic Designer", "HR Manager", "Business Analyst", "QA Engineer"}

	employmentTypes := []string{"Full-time", "Part-time", "Contract", "Freelance", "Self-employed", "Internship"}

	personalAccountTypes := []string{"contractor", "freelancer", "employee"}

	// Generate users
	for i := 0; i < s.options.UserCount; i++ {
		// Determine account type based on ratios
		var accountType string
		if s.randGen.Float32() < float32(s.options.UserRoleRatios["business"])/float32(s.options.UserRoleRatios["business"]+s.options.UserRoleRatios["personal"]) {
			accountType = "business"
		} else {
			accountType = "personal"
		}

		firstName := firstNames[s.randGen.Intn(len(firstNames))]
		lastName := lastNames[s.randGen.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s.%d@example.com", strings.ToLower(firstName), strings.ToLower(lastName), s.randGen.Intn(1000))

		// Generate password hash for "password"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Create pgtype fields
		passwordHash := pgtype.Text{String: string(hashedPassword), Valid: true}

		// Choose a random profile picture URL 50% of the time
		var profilePicture string
		if s.randGen.Float32() < 0.5 {
			profilePicture = fmt.Sprintf("https://randomuser.me/api/portraits/%s/%d.jpg",
				[]string{"men", "women"}[s.randGen.Intn(2)],
				s.randGen.Intn(99))
		} else {
			profilePicture = ""
		}

		// Gender (60% specified, 40% not)
		var gender pgtype.Text
		if s.randGen.Float32() < 0.6 {
			gender = pgtype.Text{
				String: []string{"male", "female", "non-binary", "prefer not to say"}[s.randGen.Intn(4)],
				Valid:  true,
			}
		} else {
			gender = pgtype.Text{Valid: false}
		}

		// For personal accounts, assign a personal account type
		var personalAccountType string
		if accountType == "personal" {
			personalAccountType = personalAccountTypes[s.randGen.Intn(len(personalAccountTypes))]
		} else {
			personalAccountType = "" // Empty for business accounts
		}

		// Random data for the remaining fields
		nationality := nationalities[s.randGen.Intn(len(nationalities))]
		residentialCountry := pgtype.Text{String: countries[s.randGen.Intn(len(countries))], Valid: true}
		jobRole := pgtype.Text{String: jobRoles[s.randGen.Intn(len(jobRoles))], Valid: true}

		// Company website only for business accounts and some freelancers
		var companyWebsite pgtype.Text
		if accountType == "business" || (personalAccountType == "freelancer" && s.randGen.Float32() < 0.7) {
			companyWebsite = pgtype.Text{
				String: fmt.Sprintf("https://www.%s-%s.com",
					strings.ToLower(firstName),
					strings.ToLower(lastName)),
				Valid: true,
			}
		} else {
			companyWebsite = pgtype.Text{Valid: false}
		}

		employmentType := pgtype.Text{String: employmentTypes[s.randGen.Intn(len(employmentTypes))], Valid: true}

		// Create user with the generated data
		user, err := s.queries.CreateUser(ctx, CreateUserParams{
			Column1:             uuid.New(),
			Email:               email,
			PasswordHash:        passwordHash,
			Column4:             profilePicture,
			AccountType:         accountType,
			Gender:              gender,
			PersonalAccountType: personalAccountType,
			FirstName:           firstName,
			LastName:            lastName,
			Nationality:         nationality,
			ResidentialCountry:  residentialCountry,
			JobRole:             jobRole,
			CompanyWebsite:      companyWebsite,
			EmploymentType:      employmentType,
			Column15:            time.Now(),
			Column16:            time.Now(),
		})

		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", email, err)
		}

		// Store user for relationships
		s.users = append(s.users, user)

		if s.options.Verbose && (i+1)%5 == 0 {
			log.Printf("Created %d/%d users", i+1, s.options.UserCount)
		}
	}

	log.Printf("Successfully seeded %d users", len(s.users))
	return nil
}

// seedSessions creates session records for users
func (s *Seeder) seedSessions(ctx context.Context) error {
	if len(s.users) == 0 {
		return fmt.Errorf("no users available for creating sessions")
	}

	log.Println("Seeding user sessions...")

	// User agents to make data more realistic
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.104 Mobile Safari/537.36",
	}

	// Login types
	loginTypes := []string{"email", "google", "apple", "github"}

	// Client IPs
	clientIPs := []string{
		"192.168.1.1", "10.0.0.1", "172.16.0.1",
		"203.0.113.1", "198.51.100.1", "192.0.2.1",
	}

	// Each user has 1-3 sessions
	totalSessions := 0
	for _, user := range s.users {
		sessionsCount := 1 + s.randGen.Intn(3)

		for j := 0; j < sessionsCount; j++ {
			userAgent := userAgents[s.randGen.Intn(len(userAgents))]
			loginType := loginTypes[s.randGen.Intn(len(loginTypes))]
			clientIp := clientIPs[s.randGen.Intn(len(clientIPs))]

			// Create OAuth fields for OAuth login types
			var webOauthClientID, oauthAccessToken, oauthIDToken pgtype.Text
			if loginType == "google" || loginType == "apple" || loginType == "github" {
				webOauthClientID = pgtype.Text{
					String: fmt.Sprintf("%s_client_%d", loginType, s.randGen.Intn(999999)),
					Valid:  true,
				}
				oauthAccessToken = pgtype.Text{
					String: fmt.Sprintf("access_token_%s_%d", loginType, s.randGen.Intn(999999)),
					Valid:  true,
				}
				oauthIDToken = pgtype.Text{
					String: fmt.Sprintf("id_token_%s_%d", loginType, s.randGen.Intn(999999)),
					Valid:  true,
				}
			}

			// Set expiration to be between 1 day and 30 days in the future
			expiryDays := 1 + s.randGen.Intn(29)
			expiresAt := time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)

			// MFA is enabled for 30% of sessions
			mfaEnabled := s.randGen.Float32() < 0.3

			// 5% of sessions are blocked
			isBlocked := s.randGen.Float32() < 0.05

			_, err := s.queries.CreateSession(ctx, CreateSessionParams{
				ID:               uuid.New(),
				UserID:           user.ID,
				RefreshToken:     fmt.Sprintf("refresh_token_%s", uuid.New().String()),
				UserAgent:        userAgent,
				WebOauthClientID: webOauthClientID,
				OauthAccessToken: oauthAccessToken,
				OauthIDToken:     oauthIDToken,
				UserLoginType:    loginType,
				MfaEnabled:       mfaEnabled,
				ClientIp:         clientIp,
				IsBlocked:        isBlocked,
				ExpiresAt:        expiresAt,
			})

			if err != nil {
				return fmt.Errorf("failed to create session for user %s: %w", user.Email, err)
			}

			totalSessions++
		}
	}

	log.Printf("Successfully seeded %d sessions", totalSessions)
	return nil
}

// seedOTPVerifications creates OTP verification records
func (s *Seeder) seedOTPVerifications(ctx context.Context) error {
	if len(s.users) == 0 {
		return fmt.Errorf("no users available for creating OTP verifications")
	}

	log.Println("Seeding OTP verifications...")

	// OTP purposes
	otpPurposes := []OtpPurpose{
		OtpPurposeEmailVerification,
		OtpPurposePasswordReset,
		OtpPurposePhoneVerification,
		OtpPurposeAccountRecovery,
		OtpPurposeTwoFactorAuth,
		OtpPurposeLoginConfirmation,
	}

	// Contact methods
	contactMethods := []string{"email", "sms", "app"}

	totalOtps := 0

	// Create OTPs for a subset of users (70%)
	for _, user := range s.users {
		if s.randGen.Float32() > 0.7 {
			continue
		}

		// 1-3 OTPs per user
		otpCount := 1 + s.randGen.Intn(3)

		for j := 0; j < otpCount; j++ {
			// Generate random 6-digit OTP
			otpCode := fmt.Sprintf("%06d", s.randGen.Intn(1000000))

			// Hash the OTP (simulating how it would be stored in real system)
			hashedOtp, err := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("failed to hash OTP: %w", err)
			}

			// Select purpose
			purpose := otpPurposes[s.randGen.Intn(len(otpPurposes))]

			// Configure OTP parameters
			userID := pgtype.UUID{
				Bytes: user.ID,
				Valid: true,
			}

			contactMethod := pgtype.Text{
				String: contactMethods[s.randGen.Intn(len(contactMethods))],
				Valid:  true,
			}

			// Random number of attempts (0-3)
			attemptsMade := int32(s.randGen.Intn(4))

			// Create OTP verification
			_, err = s.queries.CreateOTPVerification(ctx, CreateOTPVerificationParams{
				UserID:        userID,
				OtpCode:       otpCode,
				HashedOtp:     string(hashedOtp),
				Purpose:       purpose,
				ContactMethod: contactMethod,
				AttemptsMade:  attemptsMade,
				MaxAttempts:   5,
			})

			if err != nil {
				return fmt.Errorf("failed to create OTP verification: %w", err)
			}

			totalOtps++
		}
	}

	log.Printf("Successfully seeded %d OTP verifications", totalOtps)
	return nil
}

// seedUserDeviceTokens creates user device token records
func (s *Seeder) seedUserDeviceTokens(ctx context.Context) error {
	if len(s.users) == 0 {
		return fmt.Errorf("no users available for creating user device tokens")
	}

	log.Println("Seeding user device tokens...")

	// Device platforms
	platforms := []string{"ios", "android", "web"}

	// Device types
	deviceTypes := []string{"mobile", "tablet", "desktop", "laptop"}

	// Device models
	iosModels := []string{"iPhone 13 Pro", "iPhone 13", "iPhone 12 Pro", "iPhone SE", "iPad Pro", "iPad Air"}
	androidModels := []string{"Samsung Galaxy S21", "Google Pixel 6", "OnePlus 9", "Xiaomi Mi 11", "Samsung Galaxy Tab S7"}

	// OS Versions
	iosVersions := []string{"15.2", "15.1", "15.0", "14.8", "14.7"}
	androidVersions := []string{"12", "11", "10", "9"}

	// App versions
	appVersions := []string{"1.0.0", "1.1.0", "1.2.0", "1.2.1", "1.3.0"}

	totalDevices := 0

	// Create device records for users
	for _, user := range s.users {
		// 1-3 devices per user
		deviceCount := 1 + s.randGen.Intn(3)

		for j := 0; j < deviceCount; j++ {
			platform := platforms[s.randGen.Intn(len(platforms))]

			deviceType := pgtype.Text{
				String: deviceTypes[s.randGen.Intn(len(deviceTypes))],
				Valid:  true,
			}

			// Select model based on platform
			var deviceModel pgtype.Text
			var osName pgtype.Text
			var osVersion pgtype.Text

			if platform == "ios" {
				deviceModel = pgtype.Text{
					String: iosModels[s.randGen.Intn(len(iosModels))],
					Valid:  true,
				}
				osName = pgtype.Text{
					String: "iOS",
					Valid:  true,
				}
				osVersion = pgtype.Text{
					String: iosVersions[s.randGen.Intn(len(iosVersions))],
					Valid:  true,
				}
			} else if platform == "android" {
				deviceModel = pgtype.Text{
					String: androidModels[s.randGen.Intn(len(androidModels))],
					Valid:  true,
				}
				osName = pgtype.Text{
					String: "Android",
					Valid:  true,
				}
				osVersion = pgtype.Text{
					String: androidVersions[s.randGen.Intn(len(androidVersions))],
					Valid:  true,
				}
			} else {
				// Web platform
				deviceModel = pgtype.Text{Valid: false}
				osName = pgtype.Text{
					String: []string{"Windows", "macOS", "Linux"}[s.randGen.Intn(3)],
					Valid:  true,
				}
				osVersion = pgtype.Text{Valid: false}
			}

			// 95% of devices are active and verified
			isActive := s.randGen.Float32() < 0.95
			isVerified := s.randGen.Float32() < 0.95

			// App version
			appVersion := pgtype.Text{
				String: appVersions[s.randGen.Intn(len(appVersions))],
				Valid:  true,
			}

			// Client IP
			clientIp := pgtype.Text{
				String: fmt.Sprintf("192.168.%d.%d", s.randGen.Intn(255), s.randGen.Intn(255)),
				Valid:  true,
			}

			// Push notification token for mobile devices
			var pushToken pgtype.Text
			if platform == "ios" || platform == "android" {
				pushToken = pgtype.Text{
					String: fmt.Sprintf("push_token_%s_%d", platform, s.randGen.Intn(999999)),
					Valid:  true,
				}
			} else {
				pushToken = pgtype.Text{Valid: false}
			}

			// Set expiration (some may be expired)
			var expiresAt pgtype.Timestamptz
			if s.randGen.Float32() < 0.9 {
				// 90% are not expired
				expiresAt = pgtype.Timestamptz{
					Time:  time.Now().Add(time.Duration(s.randGen.Intn(365)) * 24 * time.Hour),
					Valid: true,
				}
			} else {
				// 10% are expired
				expiresAt = pgtype.Timestamptz{
					Time:  time.Now().Add(-time.Duration(s.randGen.Intn(30)) * 24 * time.Hour),
					Valid: true,
				}
			}

			// Generate device token
			deviceToken := fmt.Sprintf("device_token_%s_%s", platform, uuid.New().String())

			// Create device token record based on the actual struct definition
			_, err := s.queries.CreateUserDeviceToken(ctx, CreateUserDeviceTokenParams{
				ID:                    uuid.New(),
				UserID:                user.ID,
				DeviceToken:           deviceToken,
				Column4:               platform,           // Platform
				Column5:               deviceType.String,  // DeviceType
				Column6:               deviceModel.String, // DeviceModel
				Column7:               osName.String,      // OsName
				Column8:               osVersion.String,   // OsVersion
				PushNotificationToken: pushToken,
				Column10:              isActive,   // IsActive
				Column11:              isVerified, // IsVerified
				AppVersion:            appVersion,
				ClientIp:              clientIp,
				ExpiresAt:             expiresAt,
			})

			if err != nil {
				return fmt.Errorf("failed to create user device token for user %s: %w", user.Email, err)
			}

			totalDevices++
		}
	}

	log.Printf("Successfully seeded %d user device tokens", totalDevices)
	return nil
}

// seedKYCRecords creates KYC verification records for users
func (s *Seeder) seedKYCRecords(ctx context.Context) error {
	if len(s.users) == 0 {
		return fmt.Errorf("no users available for creating KYC records")
	}

	log.Println("Seeding KYC records...")

	// KYC verification levels
	verificationLevels := []string{"basic", "intermediate", "advanced"}

	// ID document types
	idTypes := []string{"passport", "national_id", "driving_license", "voter_card", "residence_permit"}

	// Countries for documents
	countries := []string{"US", "UK", "CA", "AU", "DE", "FR", "JP", "CN", "IN", "BR", "MX", "ZA", "NG", "KE", "EG", "AE"}

	// Verifier names
	verifiers := []string{"Manual Review", "IDCheck", "VerifyPlus", "GlobalID", "TrustVerify", "SecureKYC"}

	totalKYCRecords := 0

	// Create KYC records for a subset of users (80%)
	for _, user := range s.users {
		// Skip some users to simulate incomplete KYC
		if s.randGen.Float32() > 0.8 {
			continue
		}

		// Select KYC status (70% approved, 30% other statuses)
		var status string
		if s.randGen.Float32() < 0.7 {
			status = "approved"
		} else {
			// Choose a non-approved status
			nonApprovedStatuses := []string{"pending", "rejected", "expired", "requires_action"}
			status = nonApprovedStatuses[s.randGen.Intn(len(nonApprovedStatuses))]
		}

		// Verification level (depends on status)
		var verificationLevel string
		if status == "approved" {
			// For approved records, assign a verification level based on probabilities
			levelProbabilities := []float32{0.3, 0.4, 0.3} // basic, intermediate, advanced
			randValue := s.randGen.Float32()

			if randValue < levelProbabilities[0] {
				verificationLevel = verificationLevels[0]
			} else if randValue < levelProbabilities[0]+levelProbabilities[1] {
				verificationLevel = verificationLevels[1]
			} else {
				verificationLevel = verificationLevels[2]
			}
		} else if status == "pending" || status == "requires_action" {
			// Pending/requires_action usually have a target level
			verificationLevel = verificationLevels[s.randGen.Intn(len(verificationLevels))]
		} else {
			// Rejected/expired might not have a level
			if s.randGen.Float32() < 0.5 {
				verificationLevel = verificationLevels[s.randGen.Intn(len(verificationLevels))]
			} else {
				verificationLevel = ""
			}
		}

		// ID document details
		idType := idTypes[s.randGen.Intn(len(idTypes))]
		country := countries[s.randGen.Intn(len(countries))]

		// Document submission date
		submissionDate := time.Now().AddDate(0, -s.randGen.Intn(6), -s.randGen.Intn(30))

		// Verification date (if approved or rejected)
		var verificationDate time.Time
		if status == "approved" || status == "rejected" {
			// Verification happened 1-7 days after submission
			verificationDays := 1 + s.randGen.Intn(6)
			verificationDate = submissionDate.AddDate(0, 0, verificationDays)
		}

		// Comments for rejected or requires_action
		var comments string
		if status == "rejected" {
			rejectionReasons := []string{
				"Document illegible", "Document expired", "Information mismatch",
				"Suspected forgery", "Incomplete submission", "Failed verification checks",
			}
			comments = rejectionReasons[s.randGen.Intn(len(rejectionReasons))]
		} else if status == "requires_action" {
			actionReasons := []string{
				"Additional document required", "Please resubmit with better quality",
				"Address verification needed", "Please complete the missing information",
			}
			comments = actionReasons[s.randGen.Intn(len(actionReasons))]
		} else {
			comments = ""
		}

		// Verifier
		verifier := verifiers[s.randGen.Intn(len(verifiers))]

		// Placeholder for KYC record creation
		// Note: We're logging KYC details since the actual CreateKYCRecord
		// method is not available. In a real implementation, this would call the database.
		if s.options.Verbose {
			log.Printf("Would create KYC record: ID=%s, User=%s, Status=%s, Level=%s, "+
				"Document=%s (%s), Submitted=%s, Verified=%s, Verifier=%s, Comments='%s'",
				uuid.New().String(), user.Email, status, verificationLevel,
				idType, country, submissionDate.Format(time.RFC3339),
				verificationDate.Format(time.RFC3339), verifier, comments)
		}

		totalKYCRecords++
	}

	log.Printf("Successfully simulated %d KYC records", totalKYCRecords)
	log.Printf("Note: KYC record seeding is implemented as a placeholder. To enable actual database writes,")
	log.Printf("implement the CreateKYCRecord method in your database queries.")
	return nil
}

// seedTransactions creates transaction records for users
func (s *Seeder) seedTransactions(ctx context.Context) error {
	if len(s.users) == 0 {
		return fmt.Errorf("no users available for creating transactions")
	}

	log.Println("Seeding transactions...")

	// Determine number of days to generate transactions for
	txDays := 30 // Default for medium size
	if s.options.TransactionDays > 0 {
		txDays = s.options.TransactionDays
	} else {
		// Set based on size
		switch s.options.Size {
		case SeedSizeSmall:
			txDays = 7
		case SeedSizeMedium:
			txDays = 30
		case SeedSizeLarge:
			txDays = 90
		}
	}

	// Transaction types
	transactionTypes := []string{"deposit", "withdrawal", "transfer", "payment"}

	// Transaction statuses
	statuses := []string{"pending", "completed", "failed", "cancelled", "refunded"}

	// Currency codes
	currencies := []string{"USD", "EUR", "GBP", "JPY", "NGN", "KES", "ZAR", "AED"}

	// Payment methods
	paymentMethods := []string{
		"bank_transfer", "credit_card", "debit_card", "crypto", "mobile_money",
		"paypal", "venmo", "cash", "cheque",
	}

	// Transaction descriptions
	descriptions := []string{
		"Monthly subscription", "Service payment", "Product purchase",
		"Salary payment", "Freelance work", "Consulting fees",
		"Refund", "Investment return", "Loan repayment",
		"Travel expenses", "Utility bill", "Rent payment",
		"Groceries", "Health insurance", "Educational fees",
	}

	// Generate random transaction data
	startDate := time.Now().AddDate(0, 0, -txDays)

	totalTransactions := 0

	// Create transactions for each user
	for _, user := range s.users {
		// Number of transactions varies by user activity level
		// Assuming between 5-20 transactions per user per month, scaled by the txDays parameter
		txCount := 5 + s.randGen.Intn(15)
		txCount = txCount * txDays / 30 // Scale based on transaction days

		for j := 0; j < txCount; j++ {
			// Generate transaction date within the specified range
			daysOffset := s.randGen.Intn(txDays)
			txDate := startDate.AddDate(0, 0, daysOffset)

			// Select transaction type
			txType := transactionTypes[s.randGen.Intn(len(transactionTypes))]

			// Set status (80% completed, 20% other statuses)
			var status string
			if s.randGen.Float32() < 0.8 {
				status = "completed"
			} else {
				status = statuses[s.randGen.Intn(len(statuses))]
			}

			// Generate amount (between 10 and 1000, with 2 decimal places)
			amount := float64(10+s.randGen.Intn(99100)) / 100.0

			// Currency (USD more common than others)
			var currency string
			if s.randGen.Float32() < 0.6 {
				currency = "USD"
			} else {
				currency = currencies[s.randGen.Intn(len(currencies))]
			}

			// Payment method
			paymentMethod := paymentMethods[s.randGen.Intn(len(paymentMethods))]

			// Transaction description
			description := descriptions[s.randGen.Intn(len(descriptions))]
			if txType == "transfer" {
				// For transfers, add recipient information
				recipientIndex := s.randGen.Intn(len(s.users))
				recipientUser := s.users[recipientIndex]
				description = fmt.Sprintf("Transfer to %s %s", recipientUser.FirstName, recipientUser.LastName)
			}

			// Reference number
			referenceNumber := fmt.Sprintf("TX-%s-%d", txType[:3], s.randGen.Intn(1000000))

			// Fee (0-2% of amount)
			feePercent := float64(s.randGen.Intn(200)) / 10000.0
			fee := amount * feePercent

			// Placeholder for transaction creation
			// Note: We're logging transaction details since the actual CreateTransaction
			// method is not available. In a real implementation, this would call the database.
			if s.options.Verbose {
				log.Printf("Would create transaction: ID=%s, User=%s, Type=%s, Status=%s, Amount=%.2f %s, "+
					"Method=%s, Description='%s', Reference=%s, Fee=%.2f, Date=%s",
					uuid.New().String(), user.Email, txType, status, amount, currency,
					paymentMethod, description, referenceNumber, fee, txDate.Format(time.RFC3339))
			}

			// Increment counter for the output message
			totalTransactions++
		}
	}

	log.Printf("Successfully simulated %d transactions over a %d day period", totalTransactions, txDays)
	log.Printf("Note: Transaction seeding is implemented as a placeholder. To enable actual database writes,")
	log.Printf("implement the CreateTransaction method in your database queries.")
	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SeedDB is a convenience function that creates a seeder with default options
// and seeds the database
func SeedDB(ctx context.Context, queries *Queries) error {
	options := DefaultSeedOptions()
	seeder := NewSeeder(queries, options)
	return seeder.SeedDB(ctx)
}
