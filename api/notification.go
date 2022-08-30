package api

import (
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//get child  notifications list
type GetChildNotificationListResponse struct {
	ID              int32     `json:"id"`
	UserID          int32     `json:"user_id"`
	SecondaryUserID int32     `json:"secondary_user_id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

func (server *Server) GetChildNotificationList(ctx *gin.Context) {
	var rsp []GetChildNotificationListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get notification
	notification, err := server.store.GetChildNotification(ctx, authPayloadKey.ParentId)
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

		item := GetChildNotificationListResponse{
			ID:              notification[i].ID,
			UserID:          notification[i].UserID,
			SecondaryUserID: notification[i].SecondaryUserID.Int32,
			Title:           notification[i].Title,
			Type:            notification[i].Type,
			Description:     notification[i].Description,
			CreatedAt:       notification[i].CreatedAt,
		}
		rsp = append(rsp, item)
	}

	ctx.JSON(http.StatusOK, rsp)
}
