package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AddVideoRequest struct {
	Class   int32  `json:"grade" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Topic   string `json:"topic" binding:"required"`
	Url     string `json:"url" binding:"required"`
	ImgUrl  string `json:"img_url" binding:"required"`
	VideoID string `json:"video_id" binding:"required"`
}

func (server *Server) AddVideo(ctx *gin.Context) {
	var req AddVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	createdAt := time.Now().UTC()
	CreateVideoParams := db.CreateVideoParams{
		Class:     req.Class,
		Subject:   req.Subject,
		Topic:     util.CreateNullString(true, req.Topic),
		Url:       req.Url,
		ImgUrl:    util.CreateNullString(true, req.ImgUrl),
		VideoID:   util.CreateNullString(true, req.VideoID),
		CreatedAt: util.CreateNUllTime(true, createdAt),
	}

	err := server.store.CreateVideo(ctx, CreateVideoParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to craete video", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("video added successfully"))
}

type GetVideoRequest struct {
	Id int32 `json:"video_id" binding:"required"`
}
type VideoResponse struct {
	Video Video `json:"video"`
}
type Video struct {
	Id        int32     `json:"id" `
	ImgUrl    string    `json:"img_url"`
	VideoId   string    `json:"video_id"`
	Class     int32     `json:"grade" `
	Subject   string    `json:"subject" `
	Topic     string    `json:"topic" `
	Url       string    `json:"url" `
	CreatedAt time.Time `json:"created_at" `
}

func getvideoResponse(video db.Video) Video {
	return Video{
		Id:        video.ID,
		Class:     video.Class,
		Subject:   video.Subject,
		ImgUrl:    video.ImgUrl.String,
		VideoId:   video.VideoID.String,
		Topic:     video.Topic.String,
		Url:       video.Url,
		CreatedAt: video.CreatedAt.Time,
	}
}

//get video by id
func (server *Server) GetVideo(ctx *gin.Context) {
	var req GetVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	// authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// getPaymentParams := db.GetPaymentParams{
	// 	ChildID:  authPayloadKey.ParentId,
	// 	ExpireAt: util.CreateNUllTime(true, time.Now()),
	// }
	// _, err := server.store.GetPayment(ctx, getPaymentParams)
	// if err != nil {
	// 	// check if there is no row error
	// 	if err == sql.ErrNoRows {
	// 		saveErr := util.NewUnauthorizedError("No Have Not Upgraded")
	// 		ctx.JSON(saveErr.Status(), saveErr)
	// 		return
	// 	}
	// 	saveErr := util.NewInternalServerError("error when trying to get payment details", errors.New("database error"))
	// 	ctx.JSON(saveErr.Status(), saveErr)
	// 	return
	// }
	//get video
	video, err := server.store.GetVideo(ctx, req.Id)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get video", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	rsp := VideoResponse{

		Video: getvideoResponse(video),
	}
	ctx.JSON(http.StatusOK, rsp)
}

//fetch details of parent
func (server *Server) GetClassVideo(ctx *gin.Context) {

	var rsp []Video

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	student, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	getPaymentParams := db.GetPaymentParams{
		ChildID:  authPayloadKey.ParentId,
		ExpireAt: util.CreateNUllTime(true, time.Now()),
	}
	_, err = server.store.GetPayment(ctx, getPaymentParams)
	if err != nil {
		// check if there is no row error
		if err == sql.ErrNoRows {
			//get video
			video, err := server.store.GetClassVideoFree(ctx, student.Grade)
			if err != nil {

				saveErr := util.NewInternalServerError("error when trying to get video", errors.New("database error"))
				ctx.JSON(saveErr.Status(), saveErr)
				return
			}
			if len(video) == 0 {
				getErr := util.NewRestError("List Is Empty", http.StatusOK, "error when trying to get video", nil)
				ctx.JSON(getErr.Status(), getErr)
				return

			}
			for i := 0; i < len(video); i++ {

				item := Video{
					Id:        video[i].ID,
					Class:     video[i].Class,
					Subject:   video[i].Subject,
					ImgUrl:    video[i].ImgUrl.String,
					VideoId:   video[i].VideoID.String,
					Topic:     video[i].Topic.String,
					Url:       video[i].Url,
					CreatedAt: video[i].CreatedAt.Time,
				}
				rsp = append(rsp, item)
			}
			ctx.JSON(http.StatusOK, rsp)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to get payment details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//get video
	video, err := server.store.GetClassVideo(ctx, student.Grade)
	if err != nil {

		saveErr := util.NewInternalServerError("error when trying to get video", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(video) == 0 {
		getErr := util.NewRestError("List Is Empty", http.StatusOK, "error when trying to get video", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(video); i++ {

		item := Video{
			Id:        video[i].ID,
			Class:     video[i].Class,
			Subject:   video[i].Subject,
			ImgUrl:    video[i].ImgUrl.String,
			VideoId:   video[i].VideoID.String,
			Topic:     video[i].Topic.String,
			Url:       video[i].Url,
			CreatedAt: video[i].CreatedAt.Time,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

type GetSubjectVideoRequest struct {
	Subject string `json:"subject" binding:"required"`
}

//fetch details of parent
func (server *Server) GetSubjectVideo(ctx *gin.Context) {
	var req GetSubjectVideoRequest
	var rsp []Video
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	student, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	getPaymentParams := db.GetPaymentParams{
		ChildID:  authPayloadKey.ParentId,
		ExpireAt: util.CreateNUllTime(true, time.Now()),
	}
	_, err = server.store.GetPayment(ctx, getPaymentParams)
	if err != nil {
		// check if there is no row error
		if err == sql.ErrNoRows {
			GetSubjectVideoFreeParams := db.GetSubjectVideoFreeParams{
				Class:   student.Grade,
				Subject: req.Subject,
			}
			//get video
			video, err := server.store.GetSubjectVideoFree(ctx, GetSubjectVideoFreeParams)
			if err != nil {

				saveErr := util.NewInternalServerError("error when trying to get video", errors.New("database error"))
				ctx.JSON(saveErr.Status(), saveErr)
				return
			}
			if len(video) == 0 {
				getErr := util.NewRestError("List Is Empty", http.StatusOK, "error when trying to get video", nil)
				ctx.JSON(getErr.Status(), getErr)
				return

			}
			for i := 0; i < len(video); i++ {

				item := Video{
					Id:        video[i].ID,
					Class:     video[i].Class,
					Subject:   video[i].Subject,
					Topic:     video[i].Topic.String,
					ImgUrl:    video[i].ImgUrl.String,
					VideoId:   video[i].VideoID.String,
					Url:       video[i].Url,
					CreatedAt: video[i].CreatedAt.Time,
				}
				rsp = append(rsp, item)
			}
			ctx.JSON(http.StatusOK, rsp)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to get payment details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	GetSubjectVideoParams := db.GetSubjectVideoParams{
		Class:   student.Grade,
		Subject: req.Subject,
	}
	//get video
	video, err := server.store.GetSubjectVideo(ctx, GetSubjectVideoParams)
	if err != nil {

		saveErr := util.NewInternalServerError("error when trying to get video", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(video) == 0 {
		getErr := util.NewRestError("List Is Empty", http.StatusOK, "error when trying to get video", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(video); i++ {

		item := Video{
			Id:        video[i].ID,
			Class:     video[i].Class,
			Subject:   video[i].Subject,
			Topic:     video[i].Topic.String,
			ImgUrl:    video[i].ImgUrl.String,
			VideoId:   video[i].VideoID.String,
			Url:       video[i].Url,
			CreatedAt: video[i].CreatedAt.Time,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)

}
