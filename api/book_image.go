package api

import (
	"Edtech_Golang/util"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (server *Server) UploadBookImage(ctx *gin.Context) {

	var r struct {
		Data   string `json:"data"`
		BookId int32  `json:"book_id"`
	}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	decode, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid base64 data")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	filename := fmt.Sprintf("books/%d", r.BookId)

	err = server.store.UploadToS3(filename, decode)
	if err != nil {
		logrus.Error(err)
		getErr := util.NewInternalServerError("could not upload book image", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("Image added successfully"))

}
