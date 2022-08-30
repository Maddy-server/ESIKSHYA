package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SendFriendRequestRequest struct {
	ReceiverId int32 `json:"receiver_id" binding:"required"`
}

//send friend request
func (server *Server) SendFriendRequest(ctx *gin.Context) {
	var req SendFriendRequestRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ReceiverId == authPayloadKey.ParentId {
		saveErr := util.NewInternalServerError("you cannot send request to your self", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	details, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error while getting child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//get children detail
	_, err = server.store.GetChildForVerify(ctx, req.ReceiverId)
	if err != nil {
		saveErr := util.NewInternalServerError("user doesnot exists", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//check frends/pending or not
	args := db.CheckFriendsListParams{
		ID:         authPayloadKey.ParentId,
		ReceiverID: authPayloadKey.ParentId,
		SenderID:   authPayloadKey.ParentId,
	}
	//get friends list
	friends, err := server.store.CheckFriendsList(ctx, args)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get friends list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(friends) != 0 {
		for i := 0; i < len(friends); i++ {
			if friends[i].ID == req.ReceiverId {
				saveErr := util.NewInternalServerError("you are already friends or request is in pending", errors.New("database error"))
				ctx.JSON(saveErr.Status(), saveErr)
				return
			}

		}
	}

	//send request
	SendFriendRequestParams := db.SendFriendRequestParams{
		SenderID:   authPayloadKey.ParentId,
		ReceiverID: req.ReceiverId,
		Status:     "pending",
		FriendsAt:  util.CreateNUllTime(true, time.Now()),
	}

	err = server.store.SendFriendRequest(ctx, SendFriendRequestParams)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062:") {
			saveErr := util.NewInternalServerError("Already Send Friend Request", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to send friend request", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	err = server.store.SendChildNotification(
		ctx,
		req.ReceiverId,
		"FrendRequest",
		fmt.Sprintf("%s send you a friend request ", details.FullName),
	)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, GenerateResponse("error while sending notification "))
		return
	}

	createChildNotificationParams := db.CreateChildNotificationParams{
		UserID:          req.ReceiverId,
		Title:           "SendFriendRequest",
		Type:            "friendrequest",
		Description:     fmt.Sprintf("%s send you a friend request ", details.FullName),
		SecondaryUserID: util.CreateNullInt32(true, authPayloadKey.ParentId),
		CreatedAt:       time.Now(),
	}
	err = server.store.CreateChildNotification(ctx, createChildNotificationParams)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, GenerateResponse("error while trying to save notification "))
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("friend request send successfully"))
}

//request to update friend request
type AcceptFriendRequestRequest struct {
	SenderId       int32 `json:"sender_id" binding:"required"`
	NotificationId int32 `json:"notification_id" `
}

//accept frends request

func (server *Server) AcceptFriendRequest(ctx *gin.Context) {

	var req AcceptFriendRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body") //cannot unmarshal string into sql.NullTime
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	details, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error while getting child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	args := db.GetFriendParams{
		SenderID:   req.SenderId,
		ReceiverID: authPayloadKey.ParentId,
		Status:     "pending",
	}
	pending, err := server.store.GetFriend(ctx, args)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get pending request list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	arg := db.AcceptFriendRequestParams{
		ID:     pending.ID,
		Status: "friends",
	}
	editErr := server.store.AcceptFriendRequest(ctx, arg)
	if editErr != nil {
		saveErr := util.NewInternalServerError("error when trying to accept friends", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	err = server.store.SendChildNotification(
		ctx,
		req.SenderId,
		"FrendRequest",
		fmt.Sprintf("%s accepted you a friend request ", details.FullName),
	)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, GenerateResponse("error while sending notification "))
		return
	}
	err = server.store.DeleteChildNotification(ctx, req.NotificationId)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, GenerateResponse("error while deleting notification "))
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("friends request accepted successfully"))
}

//get friends list
type GetFriendsListResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (server *Server) GetFriendsList(ctx *gin.Context) {
	var rsp []GetFriendsListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.GetFriendsListParams{
		ID:         authPayloadKey.ParentId,
		Status:     "friends",
		ReceiverID: authPayloadKey.ParentId,
		SenderID:   authPayloadKey.ParentId,
	}
	//get friends list
	friends, err := server.store.GetFriendsList(ctx, args)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get friends list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(friends) == 0 {
		getErr := util.NewRestError("No friends", http.StatusOK, "error when trying to get friends list", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(friends); i++ {

		item := GetFriendsListResponse{
			ID:   friends[i].ID,
			Name: friends[i].Username,
		}
		rsp = append(rsp, item)
	}

	ctx.JSON(http.StatusOK, rsp)
}

//reject friend request
type RejectFriendRequestRequest struct {
	SenderId       int32 `json:"sender_id" binding:"required"`
	NotificationId int32 `json:"notification_id" `
}

//reject frends request
func (server *Server) RejectFriendRequest(ctx *gin.Context) {

	var req RejectFriendRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body") //cannot unmarshal string into sql.NullTime
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.GetFriendParams{
		SenderID:   req.SenderId,
		ReceiverID: authPayloadKey.ParentId,
		Status:     "pending",
	}
	pending, err := server.store.GetFriend(ctx, args)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get pending request list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	err = server.store.RejectFriendRequest(ctx, pending.ID)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to reject request", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	err = server.store.DeleteChildNotification(ctx, req.NotificationId)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusOK, GenerateResponse("error while deleting notification "))
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("friends request rejected successfully"))
}
