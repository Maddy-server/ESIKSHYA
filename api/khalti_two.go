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

type KhaltiStepTwoRequest struct {
	PublicKey        string `json:"public_key" binding:"required"`
	Token            string `json:"token" binding:"required"`
	ConfirmationCode string `json:"confirmation_code" binding:"required"`
	TransactionPIN   string `json:"transaction_pin" binding:"required"`
}

func (server *Server) KhaltiStepTwo(ctx *gin.Context) {
	var data KhaltiStepTwoRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	marshalRequest, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		saveErr := util.NewInternalServerError("cannot marshal json", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	resp, err := http.Post("https://khalti.com/api/v2/payment/confirm/", "application/json", bytes.NewBuffer(marshalRequest))
	if err != nil {
		logrus.Error(err)
	}

	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logrus.Error(err)
	// }

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		ctx.JSON(http.StatusOK, res)
		return
	}

	logrus.Error(err)
	saveErr := util.NewInternalServerError("can't complete transaction", errors.New("database error"))
	ctx.JSON(saveErr.Status(), saveErr)

}
