package api

import (
	"fmt"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/docs"
	"github.com/demola234/defiraise/token"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewTokenMaker("beb4118e1bdc8020df695ceec7e464a5")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %s", err.Error())
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	server.setUpRouter()
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "DefiFundr API"
	docs.SwaggerInfo.Description = "Decentralized Crowdfunding Platform for DeFi Projects"
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "defifundr-hyper.koyeb.app"
	docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetTrustedProxies([]string{"localhost"})
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/user/verify", server.verifyUser)
	router.POST("/user/verify/resend", server.resendVerificationCode)
	router.POST("/user/password/reset", server.resetPassword)
	router.POST("/user/password/reset/verify", server.verifyPasswordResetCode)
	router.POST("/user/password", server.createPassword)
	router.POST("/user/checkUsername", server.checkUsernameExists)
	router.POST("/token/renewAccess", server.renewAccessToken)
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRoutes.GET("/user", server.getUser)
	authRoutes.POST("/user/update", server.updateUser)
	authRoutes.POST("/userAddress", server.getUserByAddress)
	authRoutes.GET("/user/avatar", server.setProfileAvatar)
	authRoutes.POST("/user/avatar/set", server.selectAvatar)
	authRoutes.POST("/user/biometrics", server.setBiometrics)
	authRoutes.POST("/user/logout", server.logoutUser)
	authRoutes.POST("/user/password/change", server.changePassword)
	authRoutes.POST("/user/privatekey", server.getPrivateKey)
	authRoutes.GET("/campaigns/latestCampaigns", server.getLatestActiveCampaigns)
	authRoutes.GET("/campaigns", server.getCampaigns)
	authRoutes.POST("/campaigns", server.createCampaign)
	authRoutes.GET("/campaigns/:id", server.getCampaign)
	authRoutes.GET("/campaigns/categories/:id", server.getCampaignsByCategory)
	authRoutes.GET("/campaigns/owner", server.getCampaignsByOwner)
	authRoutes.GET("/campaignsTypes", server.getCampaignTypes)
	authRoutes.GET("/campaigns/donation/:id", server.getCampaignDonors)
	authRoutes.POST("/campaigns/donate", server.donateToCampaign)
	authRoutes.POST("/campaigns/withdraw", server.withdrawFromCampaign)
	authRoutes.GET("/campaigns/myDonations", server.getMyDonations)
	authRoutes.GET("/currentPrice", server.currentEthPrice)
	authRoutes.GET("/campaigns/categories", server.getCategories)
	authRoutes.GET("/campaigns/search", server.searchCampaignByName)
	authRoutes.POST("/wallet-address/create", server.createWalletAddress)

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
