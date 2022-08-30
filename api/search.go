package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//search friends list
type SearchResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (server *Server) Search(ctx *gin.Context) {
	var rsp []SearchResponse
	searchQuery := ctx.Request.URL.Query().Get("q")
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//get children detail
	child, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	searchParam := db.SearchListParams{
		Grade:    child.Grade,
		Username: "%" + strings.ToLower(searchQuery) + "%",
	}
	search, err := server.store.SearchList(ctx, searchParam)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to search", errors.New("database error"))
		ctx.JSON(saveErr.Status(), err)
		return
	}

	if len(search) == 0 {
		getErr := util.NewRestError("Search is empty", http.StatusOK, "error when trying to search", nil)
		ctx.JSON(getErr.Status(), getErr)
		return

	}
	for i := 0; i < len(search); i++ {

		item := SearchResponse{
			ID:   search[i].ID,
			Name: search[i].Username,
		}
		rsp = append(rsp, item)
	}

	ctx.JSON(http.StatusOK, rsp)
}
