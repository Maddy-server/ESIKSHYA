package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddBookRequest struct {
	BookName    string `json:"book_name" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Section     string `json:"section" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type AddBookResponse struct {
	Id          int32  `json:"id" `
	BookName    string `json:"book_name"`
	Author      string `json:"author"`
	Section     string `json:"section"`
	Description string `json:"description"`
}

//Add book and return details for image upload
func (server *Server) AddBook(ctx *gin.Context) {
	var req AddBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	randomUnique := util.RandomString(10)
	addBookParams := db.AddBookParams{
		BookName:     req.BookName,
		Writer:       req.Author,
		Section:      req.Section,
		CreatedAt:    time.Now(),
		Randomunique: randomUnique,
	}
	err := server.store.AddBook(ctx, addBookParams)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to Add Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}

	book, err := server.store.FetchBookAfterCreated(ctx, randomUnique)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	decode, err := base64.StdEncoding.DecodeString(req.Description)
	if err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid base64 data")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	filename := fmt.Sprintf("pdf/%d", book.ID)

	err = server.store.UploadToS3(filename, decode)
	if err != nil {
		logrus.Error(err)
		getErr := util.NewInternalServerError("could not upload book pdf", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	pdf := server.store.GetBookPdfUrl(book.ID)
	rsp := AddBookResponse{
		Id:          book.ID,
		BookName:    book.BookName,
		Author:      book.Writer,
		Section:     book.Section,
		Description: pdf,
	}
	ctx.JSON(http.StatusCreated, rsp)
}

//Update book by image id and content
type UpdateBookRequest struct {
	Id      int32  `json:"book_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (server *Server) UpdateBook(ctx *gin.Context) {
	var req UpdateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	updateBookParams := db.UpdateBookParams{
		ID:      req.Id,
		Content: util.CreateNullString(true, req.Content),
	}
	err := server.store.UpdateBook(ctx, updateBookParams)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to update book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("Updated successfully"))
}

//fetch book list by its section
type FetchBookListByTypesRequestJson struct {
	Section string `json:"section" binding:"required"`
}
type FetchBookListByTypesRequestUri struct {
	Page int32 `uri:"page" `
}

type AllBookListByTypesResponse struct {
	PopularFetchBookListByTypesResponse []FetchBookListByTypesResponse `json:"popular_book"`
	NormalFetchBookListByTypesResponse  []FetchBookListByTypesResponse `json:"book_by_section"`
}
type FetchBookListByTypesResponse struct {
	Id          int32  `json:"id" `
	BookName    string `json:"book_name"`
	ImageUrl    string `json:"image_url"`
	SavedStatus bool   `json:"saved_status"`
	Count       int32  `json:"count"`
}

func (server *Server) FetchBookListByTypes(ctx *gin.Context) {
	var rsp AllBookListByTypesResponse
	var rspbook []FetchBookListByTypesResponse
	var rsppopular []FetchBookListByTypesResponse

	reqUri, err := strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	if err != nil || reqUri < 0 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	var reqJson FetchBookListByTypesRequestJson
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// popular books

	popular, err := server.store.FetchPopularBookListBySection(ctx, reqJson.Section)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get popular Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(popular) == 0 {
		getErr := util.NewRestError("No popular Books Available", http.StatusOK, "No Data", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(popular); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: popular[i].ID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		img := server.store.GetBookImageUrl(popular[i].ID)
		items := FetchBookListByTypesResponse{
			Id:          popular[i].ID,
			BookName:    popular[i].BookName,
			Count:       popular[i].Count,
			SavedStatus: savedStatus,
			ImageUrl:    img,
		}
		rsppopular = append(rsppopular, items)
	}
	//normal book list
	fetchBookListBySectionParams := db.FetchBookListBySectionParams{
		Section: reqJson.Section,
		ID:      int32((reqUri * 20) + 1),
	}
	book, err := server.store.FetchBookListBySection(ctx, fetchBookListBySectionParams)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(book) != 0 {
		for i := 0; i < len(book); i++ {
			savedStatus := true
			fetchSavedBookParams := db.FetchSavedBookParams{
				UserID: authPayloadKey.ParentId,
				BookID: book[i].ID,
			}
			_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
			if err != nil {
				savedStatus = false
			}
			img := server.store.GetBookImageUrl(book[i].ID)
			item := FetchBookListByTypesResponse{
				Id:          book[i].ID,
				BookName:    book[i].BookName,
				Count:       book[i].Count,
				SavedStatus: savedStatus,
				ImageUrl:    img,
			}
			rspbook = append(rspbook, item)
		}
	} else {
		rspbook = []FetchBookListByTypesResponse{}

	}
	rsp = AllBookListByTypesResponse{
		NormalFetchBookListByTypesResponse:  rspbook,
		PopularFetchBookListByTypesResponse: rsppopular,
	}
	ctx.JSON(http.StatusOK, rsp)
}

//fetch book home page
func (server *Server) FetchBookListHome(ctx *gin.Context) {
	var rsp []FetchBookListByTypesResponse
	reqUri, err := strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	if err != nil || reqUri < 0 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	ID := (reqUri * 20) + 1
	book, err := server.store.FetchBookListHome(ctx, int32(ID))
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(book) == 0 {
		getErr := util.NewNotFoundError("no book is available")
		ctx.JSON(http.StatusOK, getErr)
		return

	}
	for i := 0; i < len(book); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: book[i].ID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		img := server.store.GetBookImageUrl(book[i].ID)
		item := FetchBookListByTypesResponse{
			Id:          book[i].ID,
			BookName:    book[i].BookName,
			Count:       book[i].Count,
			SavedStatus: savedStatus,
			ImageUrl:    img,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) FetchPopularBook(ctx *gin.Context) {
	var rsp []FetchBookListByTypesResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	book, err := server.store.FetchPopularBook(ctx)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(book) == 0 {
		getErr := util.NewNotFoundError("no popular book is available")
		ctx.JSON(http.StatusOK, getErr)
		return

	}
	for i := 0; i < len(book); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: book[i].ID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		img := server.store.GetBookImageUrl(book[i].ID)
		item := FetchBookListByTypesResponse{
			Id:          book[i].ID,
			BookName:    book[i].BookName,
			Count:       book[i].Count,
			SavedStatus: savedStatus,
			ImageUrl:    img,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) FetchNewAddedBook(ctx *gin.Context) {
	var rsp []FetchBookListByTypesResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	book, err := server.store.FetchNewBook(ctx)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(book) == 0 {
		getErr := util.NewNotFoundError("no new book is added")
		ctx.JSON(http.StatusOK, getErr)
		return

	}
	for i := 0; i < len(book); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: book[i].ID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		img := server.store.GetBookImageUrl(book[i].ID)
		item := FetchBookListByTypesResponse{
			Id:          book[i].ID,
			BookName:    book[i].BookName,
			Count:       book[i].Count,
			SavedStatus: savedStatus,
			ImageUrl:    img,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

//fetch bookdetails by id
type FetchBookContentRequest struct {
	Id int32 `json:"book_id"  binding:"required"`
}
type FetchBookDetailsByIdResponse struct {
	Id          int32  `json:"id" `
	BookName    string `json:"book_name"`
	ImageUrl    string `json:"image_url"`
	SavedStatus bool   `json:"saved_status"`
	Count       int32  `json:"count"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (server *Server) FetchBookDetailsById(ctx *gin.Context) {
	var req FetchBookContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	book, err := server.store.FetchBookDetailsById(ctx, req.Id)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	savedStatus := true
	fetchSavedBookParams := db.FetchSavedBookParams{
		UserID: authPayloadKey.ParentId,
		BookID: book.ID,
	}
	_, err = server.store.FetchSavedBook(ctx, fetchSavedBookParams)
	if err != nil {
		savedStatus = false
	}
	img := server.store.GetBookImageUrl(book.ID)
	pdf := server.store.GetBookPdfUrl(book.ID)
	resp := FetchBookDetailsByIdResponse{
		Id:          book.ID,
		BookName:    book.BookName,
		ImageUrl:    img,
		Count:       book.Count,
		Author:      book.Writer,
		Description: pdf,
		SavedStatus: savedStatus,
	}

	ctx.JSON(http.StatusOK, resp)
}

//fetch book content

type FetchBookContentResponse struct {
	Content string `json:"content" `
}

func (server *Server) FetchBookContent(ctx *gin.Context) {
	var req FetchBookContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	createBookHistoryParams := db.CreateBookHistoryParams{
		BookID:    req.Id,
		UserID:    authPayloadKey.ParentId,
		CreatedAt: time.Now(),
	}
	pdf := server.store.GetBookPdfUrl(req.Id)
	resp := FetchBookContentResponse{
		Content: pdf,
	}
	err := server.store.CreateBookHistory(ctx, createBookHistoryParams)
	if err != nil {
		ctx.JSON(http.StatusOK, resp)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) FetchHistoryBook(ctx *gin.Context) {
	var rsp []FetchBookListByTypesResponse
	reqUri, err := strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	if err != nil || reqUri < 0 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	fetchBookHistoryListParams := db.FetchBookHistoryListParams{
		UserID: authPayloadKey.ParentId,
		ID:     int32((reqUri * 15) + 1),
	}
	history, err := server.store.FetchBookHistoryList(ctx, fetchBookHistoryListParams)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get history Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(history) == 0 {
		getErr := util.NewNotFoundError("you have no history")
		ctx.JSON(http.StatusOK, getErr)
		return

	}
	for i := 0; i < len(history); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: history[i].BookID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		book, _ := server.store.FetchBookById(ctx, history[i].BookID)
		img := server.store.GetBookImageUrl(book.ID)
		items := FetchBookListByTypesResponse{
			Id:          book.ID,
			BookName:    book.BookName,
			SavedStatus: savedStatus,
			ImageUrl:    img,
			Count:       book.Count,
		}
		rsp = append(rsp, items)
	}
	ctx.JSON(http.StatusOK, rsp)
}
func (server *Server) FetchSavedBook(ctx *gin.Context) {
	var rsp []FetchBookListByTypesResponse
	reqUri, err := strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	if err != nil || reqUri < 0 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	FetchSavedBookListParams := db.FetchSavedBookListParams{
		UserID: authPayloadKey.ParentId,
		ID:     int32((reqUri * 15) + 1),
	}
	savedBook, err := server.store.FetchSavedBookList(ctx, FetchSavedBookListParams)
	if err != nil {
		getErr := util.NewInternalServerError("error when trying to get saved Book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if len(savedBook) == 0 {
		getErr := util.NewNotFoundError("you have not saved any book")
		ctx.JSON(http.StatusOK, getErr)
		return

	}
	for i := 0; i < len(savedBook); i++ {
		savedStatus := true
		fetchSavedBookParams := db.FetchSavedBookParams{
			UserID: authPayloadKey.ParentId,
			BookID: savedBook[i].BookID,
		}
		_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
		if err != nil {
			savedStatus = false
		}
		book, _ := server.store.FetchBookById(ctx, savedBook[i].BookID)
		img := server.store.GetBookImageUrl(book.ID)
		items := FetchBookListByTypesResponse{
			Id:          book.ID,
			BookName:    book.BookName,
			SavedStatus: savedStatus,
			ImageUrl:    img,
			Count:       book.Count,
		}
		rsp = append(rsp, items)
	}
	ctx.JSON(http.StatusOK, rsp)
}

type SaveBookRequest struct {
	Id int32 `json:"book_id"  binding:"required"`
}

func (server *Server) SaveBook(ctx *gin.Context) {
	var req SaveBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	fetchSavedBookParams := db.FetchSavedBookParams{
		UserID: authPayloadKey.ParentId,
		BookID: req.Id,
	}
	_, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			createSaveBookParams := db.CreateSaveBookParams{
				BookID:    req.Id,
				UserID:    authPayloadKey.ParentId,
				CreatedAt: time.Now(),
			}
			err := server.store.CreateSaveBook(ctx, createSaveBookParams)
			if err != nil {
				getErr := util.NewInternalServerError("error while saving book", errors.New("database error"))
				ctx.JSON(getErr.Status(), getErr)
				return
			}
			ctx.JSON(http.StatusOK, "saved")
		}
		getErr := util.NewInternalServerError("error while getting saved book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	getErr := util.NewInternalServerError("Already Save", errors.New("database error"))
	ctx.JSON(getErr.Status(), getErr)
}

//remove saved book
func (server *Server) RemoveSavedBook(ctx *gin.Context) {
	var req SaveBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	fetchSavedBookParams := db.FetchSavedBookParams{
		UserID: authPayloadKey.ParentId,
		BookID: req.Id,
	}
	book, err := server.store.FetchSavedBook(ctx, fetchSavedBookParams)
	if err != nil {
		getErr := util.NewInternalServerError("error while getting saved book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	err = server.store.RemovedSavedBook(ctx, book.ID)
	if err != nil {
		getErr := util.NewInternalServerError("error while removeing saved book", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("removed successfully"))
}
