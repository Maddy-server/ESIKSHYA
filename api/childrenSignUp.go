package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func GenerateResponse(msg string) Response {
	return Response{
		Message: msg,
	}
}

//childCheckParentReq is request to check if parent exists
type childCheckParentReq struct {
	RandomKey string `json:"random_key" binding:"required"`
}

// checkUsernameAvailabilityReq is request to check if username is availale or not
type checkUsernameAvailabilityReq struct {
	Username string `json:"username" binding:"required"`
}

//ChildRegisterWithPasswordReq is request to register child
type ChildRegisterWithPasswordReq struct {
	RandomKey string `json:"random_key" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

//ChildSignUpCheckParent checks if parent exists to register new child account
func (server *Server) ChildSignUpCheckParent(ctx *gin.Context) {
	var req childCheckParentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//fetch parent here
	parent, err := server.store.GetParentByRandomKey(ctx, req.RandomKey)
	if err != nil {
		if err == sql.ErrNoRows {
			getErr := util.NewRestError("parent with provided phone doesnot exists", http.StatusOK, "error when trying to get parent", nil)
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	//check if parent is verified
	if !parent.Isverified.Bool {
		err := util.NewUnauthorizedError("parent not verified")
		ctx.JSON(err.Status(), err)
		return
	}
	//send response
	ctx.JSON(http.StatusOK, GenerateResponse("Parents is verified"))

}

// CheckUsernameAvailability checks if username is availale or not
func (server *Server) CheckUsernameAvailability(ctx *gin.Context) {
	var req checkUsernameAvailabilityReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//check if username exists
	count, err := server.store.CheckUsernameAvailability(ctx, req.Username)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to check username", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}

	if count != 0 {
		restErr := util.NewBadRequestError("this username already exists please choose another username")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	ctx.JSON(http.StatusOK, GenerateResponse("username available"))
}

// ChildRegisterWithPassword registers new child
func (server *Server) ChildRegisterWithPassword(ctx *gin.Context) {

	var req ChildRegisterWithPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get parent
	parent, err := server.store.GetParentByRandomKey(ctx, req.RandomKey)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}

	//hash the provided password
	hashedPassword, hashErr := util.HashPassword(req.Password)
	if hashErr != nil {
		hashErr := util.NewInternalServerError("error when trying to hash password", errors.New("internal error"))
		ctx.JSON(hashErr.Status(), hashErr)
		return
	}
	//get data from request
	argCreateChild := db.CreateChildParams{
		ParentID: parent.ID.Int32,
		Username: req.Username,
		Password: hashedPassword,
	}

	//save username and password to db
	createErr := server.store.CreateChild(ctx, argCreateChild)
	if createErr != nil {
		saveErr := util.NewInternalServerError("error when trying to create child", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//set verification for child
	setErr := server.store.SetVerificationChild(ctx, parent.ID.Int32)
	if setErr != nil {
		setErr := util.NewInternalServerError("error when trying to set verification of child", errors.New("database error"))
		ctx.JSON(setErr.Status(), setErr)
		return
	}
	//send back response
	ctx.JSON(http.StatusCreated, GenerateResponse("child account created successfully"))
}

//ChildCheckChildDetail checks if child has detail or not
func (server *Server) CheckChildDetail(ctx *gin.Context) {

	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	count, err := server.store.CheckChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to check child detail", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if count == 0 {
		err := util.NewNotFoundError("child detail is not filled")
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("detail already filled"))
}
