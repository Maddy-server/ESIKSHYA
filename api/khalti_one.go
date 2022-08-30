package api

import (
	"Edtech_Golang/util"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type KhaltiStepOneRequest struct {
	PublicKey       string `json:"public_key" binding:"required"`
	Mobile          string `json:"mobile" binding:"required"`
	TransactionPIN  string `json:"transaction_pin" binding:"required"`
	Amount          int64  `json:"amount" binding:"required"`
	ProductIdentity string `json:"product_identity," binding:"required"`
	ProductName     string `json:"product_name" binding:"required"`
}

func (server *Server) KhaltiStepOne(ctx *gin.Context) {

	var req KhaltiStepOneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	marshalRequest, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		saveErr := util.NewInternalServerError("cannot marshal json", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	resp, err := http.Post("https://khalti.com/api/v2/payment/initiate/", "application/json", bytes.NewBuffer(marshalRequest))
	if err != nil {
		logrus.Error(err)
	}

	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		ctx.JSON(http.StatusOK, res)
		return
	}
	logrus.Error(err)
	saveErr := util.NewInternalServerError("can't complete transaction", errors.New("database error"))
	ctx.JSON(saveErr.Status(), saveErr)

}
