package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/demola234/defiraise/interfaces"
	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/token"
	"github.com/gin-gonic/gin"
)

// @Summary Create user Wallet Address
// @Description Create a new user wallet address
// @Accept  json
// @Produce  json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body interfaces.CreateWalletAddressRequest true "Wallet address Creation Data"
// @Success 200 {object} interfaces.DocSuccessResponse{data=interfaces.CreateWalletAddressResponse} "success"
// @Router /wallet-address/create [post]
func (server *Server) createWalletAddress(ctx *gin.Context) {
	var req interfaces.CreateWalletAddressRequest

	// Validate request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid request"), http.StatusBadRequest))
		return
	}

	// Extract authenticated user from context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// Create the wallet entry
	wallet, err := server.store.CreateUserWallet(ctx, db.CreateUserWalletParams{
		UserID:        user.Username,
		WalletAddress: req.WalletAddress,
		Chain:         req.Chain,
		Status:        db.UserWalletAddressesStatusesActive, // Default status: active
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// Response
	response := interfaces.CreateWalletAddressResponse{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		WalletAddress: wallet.WalletAddress,
		Chain:         wallet.Chain,
		Status:        string(wallet.Status),
		CreatedAt:     wallet.CreatedAt,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, response))
}
