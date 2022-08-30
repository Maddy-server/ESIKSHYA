package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type childLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ChildDetails struct {
	ChildID     int32     `json:"child_id"`
	FullName    string    `json:"full_name"`
	UserName    string    `json:"username"`
	DateOfBirth string    `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Grade       int32     `json:"grade"`
	School      string    `json:"school"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
}
type childLoginResponse struct {
	AccessToken string       `json:"access_token"`
	Child       ChildDetails `json:"child_details"`
}

func newChildResponse(child db.GetChildRow) ChildDetails {
	return ChildDetails{
		ChildID:     child.ID,
		UserName:    child.Username,
		FullName:    child.FullName.String,
		DateOfBirth: child.DateOfBirth.String,
		Gender:      child.Gender.String,
		Grade:       child.Grade.Int32,
		School:      child.School.String,
		Country:     child.Country.String,
		State:       child.State.String,
		CreatedAt:   child.CreatedAt.Time,
	}
}
func (server *Server) ChildLogin(ctx *gin.Context) {
	var req childLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//getting child from db
	child, err := server.store.GetChild(ctx, req.Username)
	if err != nil {
		// check if there is no row error
		if err == sql.ErrNoRows {
			getErr := util.NewUnauthorizedError("child with the username doesnot exists")
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to get child", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//checking if child is verified
	if !child.Isverified.Bool {
		err := util.NewUnauthorizedError("child not verified")
		ctx.JSON(err.Status(), err)
		return
	}

	// checking password
	err = util.CheckPassword(req.Password, child.Password)
	if err != nil {
		err := util.NewUnauthorizedError("incorrect password")
		ctx.JSON(err.Status(), err)
		return
	}
	//generating accesstoken
	accessToken, err := server.tokenMaker.CreateToken(child.Username, child.ID, time.Hour*90000)
	if err != nil {
		err := util.NewInternalServerError("error when creating access token", errors.New("token error"))
		ctx.JSON(err.Status(), err)
		return
	}

	//creating response
	rsp := childLoginResponse{
		AccessToken: accessToken,
		Child:       newChildResponse(child),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) ChildLogout(ctx *gin.Context) {
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	err := server.store.RemoveChildToken(ctx, authPayloadKey.ParentId)
	if err != nil {
		err := util.NewInternalServerError("error when trying to logout", errors.New("logout error"))
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("sucessfully logout"))
}
