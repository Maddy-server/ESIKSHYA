package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")

	if err != nil {

		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//routs with auth
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	adminRoutes := router.Group("/admin").Use(CORSMiddleware())
	paymentRoutes := router.Group("/payment").Use(authMiddleware(server.tokenMaker))
	videoRoutes := router.Group("/video").Use(authMiddleware(server.tokenMaker))
	friendRoutes := router.Group("/friend").Use(authMiddleware(server.tokenMaker))
	gameRoutes := router.Group("/game").Use(authMiddleware(server.tokenMaker))
	bookRoutes := router.Group("/book").Use(authMiddleware(server.tokenMaker))
	scoreRoutes := router.Group("/score").Use(authMiddleware(server.tokenMaker))

	router.GET("/hello", server.Hello)
	gameRoutes.GET("/withfriends", server.GameWithFriends)
	//lobby
	gameRoutes.POST("/remove/lobby", server.RemoveLobby)

	gameRoutes.GET("/withrandom", server.GameWithRandom)
	//removelobby
	gameRoutes.GET("/withrandom/remove/lobby", server.RemoveRandomLobby)
	gameRoutes.GET("/withrandom/offline", server.BeOffline)
	//get game notification
	gameRoutes.GET("/notification", server.GetGameNotificationList)
	gameRoutes.POST("/notification/remove", server.DeleteGameNotification)

	//single player game
	gameRoutes.GET("/singleplayer", server.SinglePlayerGame)
	//score api
	scoreRoutes.GET("/list/details", server.GetScoreDetailsList)
	scoreRoutes.GET("/list/country", server.GetRankListOfCountry)
	scoreRoutes.GET("/list/state", server.GetRankListOfState)
	scoreRoutes.GET("/leaderboard/country", server.GetRankOfCountry)
	scoreRoutes.GET("/leaderboard/state", server.GetRankOfState)
	scoreRoutes.POST("/stats", server.GetStatsOfScore)
	//Parent SignUp and Login
	router.POST("/parent-sign-up", server.ParentSignUp)
	router.POST("/parent-verify", server.ParentVerify)
	router.POST("/parent-set-password", server.ParentSetPassword)
	router.POST("/parent-login", server.ParentLogin)
	router.POST("/parent-resend-code", server.ParentResendCode)
	authRoutes.GET("/parent-logout", server.ParentLogout)
	//Children SignUp and Login
	router.POST("/child-sign-up-check-parent", server.ChildSignUpCheckParent)
	router.POST("/child-check-username", server.CheckUsernameAvailability)
	router.POST("/child-register-with-password", server.ChildRegisterWithPassword)
	router.POST("/child-login", server.ChildLogin)
	authRoutes.GET("/child-logout", server.ChildLogout)
	//admin routes
	adminRoutes.POST("/videos/add", server.AddVideo)
	adminRoutes.POST("/book/add", server.AddBook)
	adminRoutes.POST("/book/update", server.UpdateBook)
	adminRoutes.POST("/book/image", server.UploadBookImage)

	//book apis
	bookRoutes.GET("/home", server.FetchBookListHome)
	bookRoutes.POST("/section", server.FetchBookListByTypes)
	bookRoutes.GET("/popular", server.FetchPopularBook)
	bookRoutes.GET("/new", server.FetchNewAddedBook)
	bookRoutes.GET("/history", server.FetchHistoryBook)
	bookRoutes.GET("/saved", server.FetchSavedBook)
	bookRoutes.POST("/details", server.FetchBookDetailsById)
	bookRoutes.POST("/content", server.FetchBookContent)
	bookRoutes.POST("/save", server.SaveBook)
	bookRoutes.POST("/saved/remove", server.RemoveSavedBook)

	//payment apis
	//Khalti
	paymentRoutes.POST("/khalti/one", server.KhaltiStepOne)
	paymentRoutes.POST("/khalti/two", server.KhaltiStepTwo)
	paymentRoutes.POST("/khalti/three", server.KhaltiStepThree)
	//SDK Esewa
	paymentRoutes.POST("/sdk", server.SDKPayment)
	//Get Payment Details
	paymentRoutes.GET("/details", server.GetPayment)
	//parent apis
	authRoutes.POST("/add-parent-detail", server.AddParentDetails)
	authRoutes.PATCH("/edit-parent-detail", server.EditParentsDetail)
	authRoutes.GET("/children-details", server.GetChildrenDetails)
	authRoutes.GET("/parent-detail", server.GetParentDetails)
	authRoutes.POST("/get-time-table-parent", server.GetTimeTableByParent)
	//parents token apis
	authRoutes.POST("/parents/token", server.AddParentsToken)
	authRoutes.GET("/parents/token/remove", server.RemoveParentsToken)
	//child apis
	authRoutes.GET("/child-check-detail", server.CheckChildDetail)
	authRoutes.GET("/child-detail", server.GetChildDetails)
	authRoutes.POST("/add-child-detail", server.AddChildDetails)
	authRoutes.PATCH("/edit-child-detail", server.EditChildDetail)
	//child token apis
	authRoutes.POST("/child/token", server.AddChildToken)
	authRoutes.GET("/child/token/remove", server.RemoveChildToken)
	//child friends api for game
	friendRoutes.POST("/send", server.SendFriendRequest)
	friendRoutes.POST("/accepts", server.AcceptFriendRequest)
	friendRoutes.GET("/list", server.GetFriendsList)
	friendRoutes.POST("/reject", server.RejectFriendRequest)
	//get notification by child
	authRoutes.GET("/child/notification", server.GetChildNotificationList)
	//search friends
	authRoutes.GET("/search", server.Search)
	//video apis
	videoRoutes.GET("/class", server.GetClassVideo)
	videoRoutes.POST("/subject", server.GetSubjectVideo)
	videoRoutes.POST("/single", server.GetVideo)

	authRoutes.GET("/get-time-table", server.GetTimeTable)
	authRoutes.POST("/add-time-table", server.AddTimeTable)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
