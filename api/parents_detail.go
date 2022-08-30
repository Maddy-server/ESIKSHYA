package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddParentRequest struct {
	Full_name string `json:"full_name" binding:"required"`
}

func (server *Server) AddParentDetails(ctx *gin.Context) {
	var req AddParentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var randomKey string
	for i := 0; i < 10; i++ {
		randomKey = util.RandomString(10)
		_, err := server.store.CompairKey(ctx, randomKey)
		if err != nil {
			break
		}

	}
	parentDetailParams := db.CreateParentDetailParams{
		ParentID:  authPayloadKey.ParentId,
		FullName:  req.Full_name,
		Address:   " ",
		RandomKey: randomKey,
	}

	err := server.store.CreateParentDetail(ctx, parentDetailParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to save parent detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("parents details added successfully"))
}

//request to edit parent detail
type EditParentsDetailRequest struct {
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

//edit parents detail- patch request

func (server *Server) EditParentsDetail(ctx *gin.Context) {
	//take request same as parent detail table
	var req EditParentsDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

		restErr := util.NewBadRequestError("invalid json body") //cannot unmarshal string into sql.NullTime
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get parent detail
	current, err := server.store.GetParentDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get parent detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//if kunai field update gareko xaina vani empty aauna paryo
	//method patch
	if req.FullName != "" {
		current.FullName = req.FullName
	}

	if req.Address != "" {
		current.Address = req.Address
	}

	args := db.EditParentDetailParams{
		ParentID: current.ParentID,
		FullName: current.FullName,
		Address:  current.Address,
	}

	editErr := server.store.EditParentDetail(ctx, args)
	if editErr != nil {
		saveErr := util.NewInternalServerError("error when trying to update/edit parent detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("update/edit successfull"))
}

//get Children detail by parent- get request
type GetChildrenDetailsResponse struct {
	ChildID     int32  `json:"child_id"`
	FullName    string `json:"full_name"`
	UserName    string `json:"username"`
	DateOfBirth string `json:"date_of_birth"`
	Grade       int32  `json:"grade"`
	Gender      string `json:"gender"`
	School      string `json:"school"`
	Country     string `json:"country"`
	State       string `json:"state"`
}

func (server *Server) GetChildrenDetails(ctx *gin.Context) {
	var rsp []GetChildrenDetailsResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get parent detail
	child, err := server.store.GetChildrenDetails(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(child) == 0 {
		getErr := util.NewRestError("No clid exist", http.StatusOK, "error when trying to get children details", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(child); i++ {

		item := GetChildrenDetailsResponse{
			ChildID:     child[i].ID,
			UserName:    child[i].Username,
			FullName:    child[i].FullName.String,
			DateOfBirth: child[i].DateOfBirth.String,
			Grade:       child[i].Grade.Int32,
			Gender:      child[i].Gender.String,
			School:      child[i].School.String,
			Country:     child[i].Country.String,
			State:       child[i].State.String,
		}
		rsp = append(rsp, item)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetParentDetailsResponse struct {
	Details GetParentDetails `json:"parent_detail"`
}
type GetParentDetails struct {
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	ParentsCode string `json:"parents_code"`
}

func getParentResponse(parent db.ParentsDetail) GetParentDetails {
	return GetParentDetails{
		FullName:    parent.FullName,
		Address:     parent.Address,
		ParentsCode: parent.RandomKey,
	}
}

//fetch details of parent
func (server *Server) GetParentDetails(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get parent detail
	parent, err := server.store.GetParentDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	rsp := GetParentDetailsResponse{

		Details: getParentResponse(parent),
	}
	ctx.JSON(http.StatusOK, rsp)
}
