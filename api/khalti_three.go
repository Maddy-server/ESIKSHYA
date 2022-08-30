package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"bytes"
	"encoding/json"
	"errors"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type KhaltiStepThreeRequest struct {
	Number  string  `json:"number" binding:"required"`
	Token   string  `json:"token" binding:"required"`
	Amount  int32   `json:"amount" binding:"required"`
	ChildId []int32 `json:"childId" binding:"required"`
	TrID    string  `json:"trid" binding:"required"`
	Save    bool    `json:"save"`
}
type KhaltiRequest struct {
	Token  string `json:"token" binding:"required"`
	Amount int32  `json:"amount" binding:"required"`
}

func (server *Server) KhaltiStepThree(ctx *gin.Context) {

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var data KhaltiStepThreeRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	khaltiData := KhaltiRequest{
		Token:  data.Token,
		Amount: data.Amount,
	}

	marshalRequest, err := json.Marshal(khaltiData)
	if err != nil {
		logrus.Error(err)
		saveErr := util.NewInternalServerError("cannot marshal json", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	req, err := http.NewRequest("POST", "https://khalti.com/api/v2/payment/verify/", bytes.NewBuffer(marshalRequest))
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Set("Authorization", "Key test_secret_key_13dfb348bb1048b998d13f9c8b6a6fe6")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		paymentNumber := db.CreatePaymentNumberParams{
			Number:   data.Number,
			Method:   "KHALTI",
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
				Method:           "KHALTI",
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
		ctx.JSON(http.StatusOK, res)
		return
	}
	logrus.Error(err)
	saveErr := util.NewInternalServerError("can't complete transaction", errors.New("database error"))
	ctx.JSON(saveErr.Status(), resp.StatusCode)
}
