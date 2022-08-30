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

type GetPaymentResponse struct {
	ID               int32     `json:"id"`
	TransactionID    string    `json:"transaction_id"`
	TransactionToken string    `json:"transaction_token"`
	Method           string    `json:"method"`
	ParentID         int32     `json:"parent_id"`
	ChildID          int32     `json:"child_id"`
	Amount           int32     `json:"amount"`
	PayAt            time.Time `json:"pay_at"`
	ExpireAt         time.Time `json:"expire_at"`
}

func (server *Server) GetPayment(ctx *gin.Context) {
	var rsp []GetPaymentResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	GetPaymentListParams := db.GetPaymentListParams{
		ParentID: authPayloadKey.ParentId,
		ExpireAt: util.CreateNUllTime(true, time.Now().UTC()),
	}
	//get payment
	payment, err := server.store.GetPaymentList(ctx, GetPaymentListParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get payment", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	for i := 0; i < len(payment); i++ {

		item := GetPaymentResponse{
			ID:               payment[i].ID,
			TransactionID:    payment[i].TransactionID,
			TransactionToken: payment[i].TransactionToken,
			Method:           payment[i].Method,
			ParentID:         payment[i].ParentID,
			ChildID:          payment[i].ChildID,
			Amount:           payment[i].Amount,
			PayAt:            payment[i].PayAt.Time,
			ExpireAt:         payment[i].ExpireAt.Time,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}
