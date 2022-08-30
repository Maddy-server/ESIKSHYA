package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type timeTableReq struct {
	Class      int32  `json:"grade" binding:"required"`
	Section    string `json:"section" binding:"required"`
	Descrition string `json:"description" binding:"required"`
	Day        string `json:"day" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

func (server *Server) AddTimeTable(ctx *gin.Context) {
	var req timeTableReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := util.NewBadRequestError("invalid json body")
		ctx.JSON(err.Status(), err)
		return
	}
	// startTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("2022-02-02T%sZ", req.StartTime))
	// endTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("2022-02-02T%sZ", req.EndTime))
	// get child id from payload
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// prepare params
	arg := db.AddTimeTableParams{
		ChildrenID:  authPayloadKey.ParentId,
		Class:       req.Class,
		Section:     req.Section,
		Description: req.Descrition,
		Day:         req.Day,
		StartTime:   util.CreateNullString(true, req.StartTime),
		EndTime:     util.CreateNullString(true, req.EndTime),
	}

	//add timetable
	err := server.store.AddTimeTable(ctx, arg)
	if err != nil {
		saveErr := util.NewInternalServerError("error while saving time table", errors.New("internal server error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}

	//send response
	ctx.JSON(http.StatusCreated, GenerateResponse("time table added successfully"))
}

type TimeTableResponse struct {
	ID          int32  `json:"id"`
	ChildrenID  int32  `json:"children_id"`
	Class       int32  `json:"grade"`
	Section     string `json:"section"`
	Description string `json:"description"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

func (server *Server) GetTimeTable(ctx *gin.Context) {
	var rsp []TimeTableResponse
	// get payload
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get time table
	timeTable, err := server.store.GetTimeTable(ctx, authPayloadKey.ParentId) //should create new payoad for child
	if err != nil {
		fmt.Println(err)
		getErr := util.NewInternalServerError("error while getting timetable", errors.New("internal server error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	for i := 0; i < len(timeTable); i++ {
		item := TimeTableResponse{
			ID:          timeTable[i].ID,
			ChildrenID:  timeTable[i].ChildrenID,
			Class:       timeTable[i].Class,
			Section:     timeTable[i].Section,
			Description: timeTable[i].Description,
			Day:         timeTable[i].Day,
			StartTime:   timeTable[i].StartTime.String,
			EndTime:     timeTable[i].EndTime.String,
		}
		rsp = append(rsp, item)
	}

	//return response
	ctx.JSON(http.StatusOK, rsp)

}

type GetTimeTableByParentReq struct {
	Child_Id int32 `json:"child_id" binding:"required"`
}

func (server *Server) GetTimeTableByParent(ctx *gin.Context) {
	var rsp []TimeTableResponse

	var req GetTimeTableByParentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err := util.NewBadRequestError("invalid json body")
		ctx.JSON(err.Status(), err)
		return
	}
	// get payload
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	parent, err := server.store.GetParentId(ctx, req.Child_Id)
	if err != nil {
		fmt.Println(err)
		getErr := util.NewInternalServerError("error while getting parents id", errors.New("internal server error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	if parent != authPayloadKey.ParentId {
		fmt.Println(err)
		getErr := util.NewInternalServerError("You can only view time table of you  own child", errors.New("internal server error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}

	//get time table
	timeTable, err := server.store.GetTimeTable(ctx, req.Child_Id) //should create new payoad for child
	if err != nil {
		fmt.Println(err)
		getErr := util.NewInternalServerError("error while getting timetable", errors.New("internal server error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	for i := 0; i < len(timeTable); i++ {
		item := TimeTableResponse{
			ID:          timeTable[i].ID,
			ChildrenID:  timeTable[i].ChildrenID,
			Class:       timeTable[i].Class,
			Section:     timeTable[i].Section,
			Description: timeTable[i].Description,
			Day:         timeTable[i].Day,
			StartTime:   timeTable[i].StartTime.String,
			EndTime:     timeTable[i].EndTime.String,
		}
		rsp = append(rsp, item)
	}

	//return response
	ctx.JSON(http.StatusOK, rsp)

}
