package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AddChildRequest struct {
	Fullname    string `json:"full_name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
	Grade       int32  `json:"grade" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	School      string `json:"school" binding:"required"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state"  binding:"required"`
}

func (server *Server) AddChildDetails(ctx *gin.Context) {
	var req AddChildRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	childDetailParams := db.CreateChildDetailParams{
		ChildrenID:  authPayloadKey.ParentId,
		FullName:    req.Fullname,
		DateOfBirth: req.DateOfBirth,
		Grade:       req.Grade,
		Gender:      req.Gender,
		School:      req.School,
		Country:     util.CreateNullString(true, req.Country),
		State:       util.CreateNullString(true, req.State),
	}

	err := server.store.CreateChildDetail(ctx, childDetailParams)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062:") {
			saveErr := util.NewInternalServerError("child details already exist", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to save child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("child details added successfully"))
}

//request to edit child detail
type EditChildDetailRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
	Grade    int32  `json:"grade" binding:"required"`
	School   string `json:"school" binding:"required"`
	Country  string `json:"country" binding:"required"`
	State    string `json:"state"  binding:"required"`
}

//edit child detail- patch request

func (server *Server) EditChildDetail(ctx *gin.Context) {
	//take request same as child detail table
	var req EditChildDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

		restErr := util.NewBadRequestError("invalid json body") //cannot unmarshal string into sql.NullTime
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get child detail
	current, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//if kunai field update gareko xaina vani empty aauna paryo
	//method patch
	if req.FullName != "" {
		current.FullName = req.FullName
	}

	if req.Gender != "" {
		current.Gender = req.Gender
	}
	if req.School != "" {
		current.School = req.School
	}
	if req.Grade != 0 {
		current.Grade = req.Grade
	}
	if req.Country != "" {
		current.Country = util.CreateNullString(true, req.Country)
	}
	if req.State != "" {
		current.State = util.CreateNullString(true, req.State)
	}

	args := db.EditChildDetailParams{
		ChildrenID: current.ChildrenID,
		FullName:   current.FullName,
		Grade:      current.Grade,
		Gender:     current.Gender,
		School:     current.School,
		Country:    current.Country,
		State:      current.State,
	}

	editErr := server.store.EditChildDetail(ctx, args)
	if editErr != nil {
		saveErr := util.NewInternalServerError("error when trying to update/edit child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("update/edit successfull"))
}

type GetChildDetailsResponse struct {
	Details GetChildDetails `json:"child_detail"`
}
type GetChildDetails struct {
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth" `
	Grade       int32  `json:"grade"`
	Gender      string `json:"gender" `
	School      string `json:"school" `
	Country     string `json:"country" `
	State       string `json:"state" `
}

func getChildResponse(child db.ChildrenDetail) GetChildDetails {
	return GetChildDetails{
		FullName:    child.FullName,
		DateOfBirth: child.DateOfBirth,
		Grade:       child.Grade,
		Gender:      child.Gender,
		School:      child.School,
		Country:     child.Country.String,
		State:       child.State.String,
	}
}

//fetch details of parent
func (server *Server) GetChildDetails(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get children detail
	child, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	rsp := GetChildDetailsResponse{

		Details: getChildResponse(child),
	}
	ctx.JSON(http.StatusOK, rsp)
}
