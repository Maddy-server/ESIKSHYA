package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetScoreDetailsListResponse struct {
	YourPoints int32  `json:"your_points"`
	Message    string `json:"message"`
	OpPoints   int32  `json:"op_points"`
}

//Score Details of child
func (server *Server) GetScoreDetailsList(ctx *gin.Context) {
	var rsp []GetScoreDetailsListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	score, err := server.store.ScoreDetailsList(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get score", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	for i := 0; i < len(score); i++ {
		var msg string
		opChild, err := server.store.GetChildDetail(ctx, score[i].OpID)
		if err != nil {
			saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
			ctx.JSON(saveErr.Status(), saveErr)
			return
		}
		if score[i].OwnPoints > score[i].OpPoints {
			msg = fmt.Sprintf("You scores more than %s in %s", opChild.FullName, score[i].Subject)
		} else if score[i].OwnPoints < score[i].OpPoints {
			msg = fmt.Sprintf("%s scores more than you in %s", opChild.FullName, score[i].Subject)
		} else {
			msg = fmt.Sprintf("Match Draw with %s in %s", opChild.FullName, score[i].Subject)
		}
		item := GetScoreDetailsListResponse{
			YourPoints: score[i].OwnPoints,
			OpPoints:   score[i].OpPoints,
			Message:    msg,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

type GetRankListResponse struct {
	Rank   int32  `json:"rank"`
	UserID int32  `json:"user_id"`
	Name   string `json:"name"`
	Grade  int32  `json:"grade"`
	Score  int32  `json:"score"`
	You    bool   `json:"you"`
}

//Rank Of overal Nepal
func (server *Server) GetRankListOfCountry(ctx *gin.Context) {
	var rsp []GetRankListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	countryRank, err := server.store.ScoreDetailsListByCountry(ctx, ownInfo.Country)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get rank list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), err)
		return
	}
	for i := 0; i < len(countryRank); i++ {
		you := false
		if countryRank[i].UserID == authPayloadKey.ParentId {
			you = true
		}
		item := GetRankListResponse{
			Rank:   int32(i + 1),
			UserID: countryRank[i].UserID,
			Name:   countryRank[i].FullName.String,
			Grade:  countryRank[i].Grade.Int32,
			Score:  countryRank[i].Max,
			You:    you,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

//Rank Of own state
func (server *Server) GetRankListOfState(ctx *gin.Context) {
	var rsp []GetRankListResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	stateRank, err := server.store.ScoreDetailsListByState(ctx, ownInfo.State)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get rank list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	for i := 0; i < len(stateRank); i++ {
		you := false
		if stateRank[i].UserID == authPayloadKey.ParentId {
			you = true
		}
		item := GetRankListResponse{
			Rank:   int32(i),
			UserID: stateRank[i].UserID,
			Name:   stateRank[i].FullName.String,
			Grade:  stateRank[i].Grade.Int32,
			Score:  stateRank[i].Max,
			You:    you,
		}
		rsp = append(rsp, item)
	}
	ctx.JSON(http.StatusOK, rsp)
}

type GetCountryRankResponse struct {
	CountryRank int32  `json:"country_rank"`
	Country     string `json:"country"`
	Points      int32  `json:"points"`
}

//Rank Of overal Nepal
func (server *Server) GetRankOfCountry(ctx *gin.Context) {
	var rsp GetCountryRankResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	countryRank, err := server.store.ScoreDetailsListByCountry(ctx, ownInfo.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			rsp = GetCountryRankResponse{
				CountryRank: 0,
				Country:     ownInfo.Country.String,
				Points:      0,
			}
			ctx.JSON(http.StatusOK, rsp)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to get rank list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(countryRank) == 0 {
		rsp = GetCountryRankResponse{
			CountryRank: 0,
			Country:     ownInfo.Country.String,
			Points:      0,
		}
		ctx.JSON(http.StatusOK, rsp)
		return

	}
	for i := 0; i < len(countryRank); i++ {

		if countryRank[i].UserID == authPayloadKey.ParentId {
			rsp = GetCountryRankResponse{
				CountryRank: int32(i),
				Country:     countryRank[i].Country.String,
				Points:      countryRank[i].Max,
			}
			ctx.JSON(http.StatusOK, rsp)
			return
		}

	}
	rsp = GetCountryRankResponse{
		CountryRank: 0,
		Country:     ownInfo.Country.String,
		Points:      0,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetStateRankResponse struct {
	StateRank int32  `json:"state_rank"`
	State     string `json:"state"`
	Points    int32  `json:"points"`
}

//Rank Of own state
func (server *Server) GetRankOfState(ctx *gin.Context) {
	var rsp GetStateRankResponse
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get children details", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	stateRank, err := server.store.ScoreDetailsListByState(ctx, ownInfo.State)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get rank list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	if len(stateRank) == 0 {
		rsp = GetStateRankResponse{
			StateRank: 0,
			State:     ownInfo.Country.String,
			Points:    0,
		}
		ctx.JSON(http.StatusOK, rsp)
		return

	}
	for i := 0; i < len(stateRank); i++ {

		if stateRank[i].UserID == authPayloadKey.ParentId {
			rsp = GetStateRankResponse{
				StateRank: int32(i),
				State:     stateRank[i].State.String,
				Points:    stateRank[i].Max,
			}
			ctx.JSON(http.StatusOK, rsp)
			return
		}

	}
	rsp = GetStateRankResponse{
		StateRank: 0,
		State:     ownInfo.Country.String,
		Points:    0,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type GetStatsOfScoreResponse struct {
	GamePlayed int32 `json:"game_played"`
	QuizWon    int32 `json:"quiz_won"`
}
type GetStatsOfScoreRequest struct {
	Year  int32 `json:"year"`
	Month int32 `json:"month"`
}

//Rank Of own state
func (server *Server) GetStatsOfScore(ctx *gin.Context) {
	var rsp GetStatsOfScoreResponse
	var req GetStatsOfScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	layout := "2009-11-10 23:00:00 +0000 UTC m=+0.000000001"
	str := fmt.Sprintf("%d-%d-01 00:00:00", req.Year, req.Month)
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}

	scoreDetailsSatsParams := db.ScoreDetailsSatsParams{
		UserID:    authPayloadKey.ParentId,
		CreatedAt: t,
	}
	stats, err := server.store.ScoreDetailsSats(ctx, scoreDetailsSatsParams)
	if err != nil {
		saveErr := util.NewInternalServerError("error when trying to get rank list", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	var count int32
	var gamePlayed int32
	count = 0
	gamePlayed = 0
	for i := 0; i < len(stats); i++ {
		gamePlayed = gamePlayed + 1
		if stats[i].OwnPoints > stats[i].OpPoints {
			count = count + 1
		} else {
			count = count + 0
		}

	}
	rsp = GetStatsOfScoreResponse{
		GamePlayed: gamePlayed,
		QuizWon:    count,
	}
	ctx.JSON(http.StatusOK, rsp)
}
