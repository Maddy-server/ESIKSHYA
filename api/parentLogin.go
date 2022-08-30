package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// parentLoginRequest is request for login
type parentLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// parentLoginResponse is response for login
type parentLoginResponse struct {
	AccessToken string        `json:"access_token"`
	Parent      ParentDetails `json:"parent_details"`
}

type ParentDetails struct {
	Phone     string    `json:"phone"`
	ParentID  int32     `json:"parent_id"`
	FullName  string    `json:"full_name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at" `
}

func newParentResponse(parent db.GetParentForLoginRow) ParentDetails {
	return ParentDetails{
		ParentID:  parent.ID,
		Phone:     parent.Phone,
		FullName:  parent.FullName.String,
		Address:   parent.Address.String,
		CreatedAt: parent.CreatedAt.Time,
	}
}

// ParentLogin logs a parent in
func (server *Server) ParentLogin(ctx *gin.Context) {
	var req parentLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	// getting parent data fron db
	parent, err := server.store.GetParentForLogin(ctx, req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "converting NULL to int32 is unsupported") {
			getParent, getErr := server.store.GetParentInfo(ctx, req.Phone)
			if getErr != nil {
				// check if there is no row error
				if err == sql.ErrNoRows {
					getErr := util.NewRestError("parent with provided mobile doesnot exists", http.StatusBadRequest, "error when trying to get parent", nil)
					ctx.JSON(getErr.Status(), getErr)
					return
				}
				getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
				fmt.Println(err)
				ctx.JSON(getErr.Status(), getErr)
				return
			}
			// checking password
			err = util.CheckPassword(req.Password, getParent.Password.String)
			if err != nil {
				err := util.NewUnauthorizedError("incorrect password")
				ctx.JSON(err.Status(), err.Message())
				return
			}
			// // generating access token
			// accessToken, err := server.tokenMaker.CreateToken("", getParent.ID, time.Hour)
			// if err != nil {
			// 	err := util.NewInternalServerError("error when creating access token", errors.New("token error"))
			// 	ctx.JSON(err.Status(), err)
			// 	return
			// }
			// type parentUndetailedResponse struct {
			// 	AccessToken string `json:"access_token"`
			// 	Phone       string `json:"phone"`
			// 	ParentId    int32  `json:"parent_id"`
			// 	FirstName   string `json:"first_name"`
			// }
			// // creating response
			// rsp := parentUndetailedResponse{
			// 	AccessToken: accessToken,
			// 	Phone:       getParent.Phone,
			// 	ParentId:    getParent.ID,
			// 	FirstName:   "",
			// }
			// ctx.JSON(http.StatusOK, rsp)
			return
		}
		log.Println(err)
		// check if there is no row error
		if err == sql.ErrNoRows {
			getErr := util.NewRestError("parent with provided mobile doesnot exists", http.StatusBadRequest, "error when trying to get parent", nil)
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		fmt.Println(err)
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	// checking password
	err = util.CheckPassword(req.Password, parent.Password.String)
	if err != nil {
		err := util.NewUnauthorizedError("incorrect password")
		ctx.JSON(err.Status(), err)
		return
	}
	// generating access token
	accessToken, err := server.tokenMaker.CreateToken(parent.FullName.String, parent.ID, time.Hour*90000)
	if err != nil {
		err := util.NewInternalServerError("error when creating access token", errors.New("token error"))
		ctx.JSON(err.Status(), err)
		return
	}
	// creating response
	rsp := parentLoginResponse{
		AccessToken: accessToken,
		Parent:      newParentResponse(parent),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) ParentLogout(ctx *gin.Context) {
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	err := server.store.RemoveParentsToken(ctx, authPayloadKey.ParentId)
	if err != nil {
		err := util.NewInternalServerError("error when trying to logout", errors.New("logout error"))
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, GenerateResponse("sucessfully logout"))
}
