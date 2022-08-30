package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SDKRequest struct {
	Number  string  `json:"number" binding:"required"`
	Token   string  `json:"token" binding:"required"`
	Amount  int32   `json:"amount" binding:"required"`
	ChildId []int32 `json:"childId" binding:"required"`
	TrID    string  `json:"trid" binding:"required"`
	Save    bool    `json:"save"`
}

func (server *Server) SDKPayment(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var data SDKRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	paymentNumber := db.CreatePaymentNumberParams{
		Number:   data.Number,
		Method:   "ESEWA",
		ParentID: authPayloadKey.ParentId,
		Save:     util.CreateNullBool(true, data.Save),
	}
	err := server.store.CreatePaymentNumber(ctx, paymentNumber)
	if err != nil {
		logrus.Error(err)
		saveErr := util.NewInternalServerError("error when saving the payment details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	count := int32(len(data.ChildId))
	expire := time.Now().AddDate(0, 1, 0)
	for i := 0; i < len(data.ChildId); i++ {
		payment := db.CreatePaymentParams{
			TransactionID:    data.TrID,
			TransactionToken: data.Token,
			Method:           "ESEWA",
			ParentID:         authPayloadKey.ParentId,
			ChildID:          data.ChildId[i],
			Amount:           (data.Amount / count) / 100,
			PayAt:            util.CreateNUllTime(true, time.Now()),
			ExpireAt:         util.CreateNUllTime(true, expire),
		}
		err := server.store.CreatePayment(ctx, payment)
		if err != nil {
			logrus.Error(err)
			saveErr := util.NewInternalServerError("error when saving the payment", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return
		}
	}
	ctx.JSON(http.StatusOK, GenerateResponse("Esewa Payment successful"))

}
