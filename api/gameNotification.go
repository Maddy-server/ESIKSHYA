package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//get child  notifications list
type GetGameNotificationListResponse struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	OponentID   int32     `json:"oponent_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Status      string    `json:"status"`
	Grade       int32     `json:"grade"`
	CreatedAt   time.Time `json:"created_at"`
}

func (server *Server) GetGameNotificationList(ctx *gin.Context) {
	var rsp []GetGameNotificationListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get notification
	notification, err := server.store.GetGameNotification(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get notification", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(notification) == 0 {
		getErr := util.NewRestError("Notification List Is Empty", http.StatusOK, "error when trying to get notifications", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(notification); i++ {

		item := GetGameNotificationListResponse{
			ID:          notification[i].ID,
			UserID:      notification[i].UserID,
			OponentID:   notification[i].OponentID,
			Title:       notification[i].Title,
			Type:        notification[i].Type,
			Subject:     notification[i].Subject,
			Status:      notification[i].Status,
			Grade:       notification[i].Grade,
			Description: notification[i].Description,
			CreatedAt:   notification[i].CreatedAt,
		}
		rsp = append(rsp, item)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type DeleteNotificationRequest struct {
	NotificationId int32 `json:"notification_id"`
}

func (server *Server) DeleteGameNotification(ctx *gin.Context) {
	var req DeleteNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	deleteGameNotificationParams := db.DeleteGameNotificationParams{
		ID:     req.NotificationId,
		UserID: authPayloadKey.ParentId,
	}
	//get notification
	err := server.store.DeleteGameNotification(ctx, deleteGameNotificationParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get notification", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	ctx.JSON(http.StatusOK, "Game Notification is deleted")
}
