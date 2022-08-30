package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddChildTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func (server *Server) AddChildToken(ctx *gin.Context) {
	var req AddChildTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	createChildTokenParams := db.CreateChildTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  req.Token,
	}
	_, err := server.store.GetChildToken(ctx, authPayloadKey.ParentId)
	if err != nil {
		err := server.store.CreateChildToken(ctx, createChildTokenParams)
		if err != nil {

			saveErr := util.NewInternalServerError("error when trying to create token", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return

		}
		ctx.JSON(http.StatusCreated, GenerateResponse("token Created successfully"))
		return
	}
	updateChildTokenParams := db.UpdateChildTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  req.Token,
	}
	err = server.store.UpdateChildToken(ctx, updateChildTokenParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to update token", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("token updated successfully"))

}

type AddParentsTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func (server *Server) AddParentsToken(ctx *gin.Context) {
	var req AddParentsTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	createParentsTokenParams := db.CreateParentsTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  req.Token,
	}
	_, err := server.store.GetParentsToken(ctx, authPayloadKey.ParentId)
	if err != nil {
		err := server.store.CreateParentsToken(ctx, createParentsTokenParams)
		if err != nil {

			saveErr := util.NewInternalServerError("error when trying to create token", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return

		}
		ctx.JSON(http.StatusCreated, GenerateResponse("token Created successfully"))
		return
	}

	updateParentsTokenParams := db.UpdateParentsTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  req.Token,
	}
	err = server.store.UpdateParentsToken(ctx, updateParentsTokenParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to update token", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("token updated successfully"))
}

//remove child token
func (server *Server) RemoveChildToken(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	updateChildTokenParams := db.UpdateChildTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  "",
	}
	err := server.store.UpdateChildToken(ctx, updateChildTokenParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to update token", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("token removed successfully"))
}

//remove Parents token
func (server *Server) RemoveParentsToken(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	updateParentsTokenParams := db.UpdateParentsTokenParams{
		UserID: authPayloadKey.ParentId,
		Token:  "",
	}
	err := server.store.UpdateParentsToken(ctx, updateParentsTokenParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to update token", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("token removed successfully"))
}
