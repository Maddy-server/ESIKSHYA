package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Resp struct {
	Message string `json:"message"`
}

func (server *Server) GameWithRandom(ctx *gin.Context) {

	//variable initialization
	errormsg := "error"
	correctmsg := "CorrectAnswer"
	wrongmsg := "WrongAnswer"
	finalmsg := "gameCompleted"
	questionmsg := "questions"
	connection := "connected"

	//variable decleration

	var data GameWithRandomRequest
	var getScoreListParams db.GetScoreListParams
	var ownQueueInfoUpdate db.GameQueue
	var ownQueueInfo db.GameQueue

	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Error("error get connection ", err)
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()

	//first request in random
	err = ws.ReadJSON(&data)
	if err != nil {
		log.Println("error read json")
		log.Fatal(err)
	}

	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//get own details
	ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "could not get the your information",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		return
	}

	// //sucess resp
	// rsp := Resp{
	// 	Message: "fetch own data",
	// }
	// err = ws.WriteJSON(rsp)
	// if err != nil {
	// 	log.Println("error write json: " + err.Error())
	// }
	// //end sucess resp

	//create fake lobby
	createGameRandomLobbyParams := db.CreateGameRandomLobbyParams{
		UserID:    0,
		OpID:      0,
		Class:     authPayloadKey.ParentId,
		Status:    "none",
		CreatedAt: time.Now(),
	}
	err = server.store.CreateGameRandomLobby(ctx, createGameRandomLobbyParams)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while creating lobby",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())

		}
		return
	}
	// //sucess resp
	// rsp = Resp{
	// 	Message: "fake lobby created",
	// }
	// err = ws.WriteJSON(rsp)
	// if err != nil {
	// 	log.Println("error write json: " + err.Error())
	// }
	// //end sucess resp

	//get lobby info
	fake, err := server.store.GetFakeGameRandomLobbyByClass(ctx, authPayloadKey.ParentId)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while geting lobby",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())

		}
		return
	}
	// //sucess resp
	// rsp = Resp{
	// 	Message: "get fake lobby",
	// }
	// err = ws.WriteJSON(rsp)
	// if err != nil {
	// 	log.Println("error write json: " + err.Error())
	// }
	// //end sucess resp
	//wait in queue
	CreateQueueParams := db.CreateQueueParams{
		UserID:  authPayloadKey.ParentId,
		Status:  "queue",
		LobbyID: fake.ID,
		Subject: data.Subject,
		Grade:   ownInfo.Grade,

		CreatedAt: time.Now(),
	}
	err = server.store.CreateQueue(ctx, CreateQueueParams)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while processing in queue",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())

		}
		err = removefakeLobby(ctx, authPayloadKey.ParentId, server)
		if err != nil {
			return
		}
		err = beOffline(ctx, authPayloadKey.ParentId, server)
		if err != nil {
			return
		}
		return
	}
	// //sucess resp
	// rsp = Resp{
	// 	Message: "queue is created",
	// }
	// err = ws.WriteJSON(rsp)
	// if err != nil {
	// 	log.Println("error write json: " + err.Error())
	// }
	// //end sucess resp
	time.Sleep(10 * time.Second)
	//get won queue info
	ownQueueInfo, err = server.store.GetOwnQueueInfo(ctx, authPayloadKey.ParentId)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while fetching queue status",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		err = removefakeLobby(ctx, authPayloadKey.ParentId, server)
		if err != nil {
			return
		}
		err = beOffline(ctx, authPayloadKey.ParentId, server)
		if err != nil {
			return
		}
		return
	}

	//Check if you are connected to other oponent or not
	if ownQueueInfo.Status == "queue" {
		// //sucess resp
		// rsp = Resp{
		// 	Message: " till you are in queue",
		// }
		// err = ws.WriteJSON(rsp)
		// if err != nil {
		// 	log.Println("error write json: " + err.Error())
		// }
		// //end sucess resp
		var randomPlayer db.GameQueue
		var a db.GameQueue
		//get random player
		getQueueParams := db.GetQueueParams{
			Status:  "queue",
			Subject: data.Subject,
			Grade:   ownInfo.Grade,
			UserID:  authPayloadKey.ParentId,
		}
		for i := 0; i < 30000000000; i++ {
			randomPlayer, err = server.store.GetQueue(ctx, getQueueParams)
			if randomPlayer != a {
				// //sucess resp
				// rsp = Resp{
				// 	Message: "you found one",
				// }
				// err = ws.WriteJSON(rsp)
				// if err != nil {
				// 	log.Println("error write json: " + err.Error())
				// }
				// //end sucess resp
				//update the queue status of other users
				updateQueueParams := db.UpdateQueueParams{
					Status: "connected",
					UserID: randomPlayer.UserID,
				}
				err = server.store.UpdateQueue(ctx, updateQueueParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while updating queue",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					err = removefakeLobby(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					err = beOffline(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					return
				}
				//update the queue status of own
				updateQueueParams = db.UpdateQueueParams{
					Status: "connected",
					UserID: authPayloadKey.ParentId,
				}
				err = server.store.UpdateQueue(ctx, updateQueueParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while updating queue",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					err = removefakeLobby(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					err = beOffline(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					return
				}

				//create a lobby
				createGameRandomLobbyParams := db.CreateGameRandomLobbyParams{
					UserID:    authPayloadKey.ParentId,
					OpID:      randomPlayer.UserID,
					Class:     ownQueueInfo.Grade,
					Status:    "connected",
					CreatedAt: time.Now(),
				}
				err := server.store.CreateGameRandomLobby(ctx, createGameRandomLobbyParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while creating lobby",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					err = removefakeLobby(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					err = beOffline(ctx, authPayloadKey.ParentId, server)
					if err != nil {
						return
					}
					return
				}
				//get lobby
				getGameRandomLobbyParams := db.GetGameRandomLobbyParams{
					UserID: authPayloadKey.ParentId,
					OpID:   randomPlayer.UserID,
				}
				getlobby, err := server.store.GetGameRandomLobby(ctx, getGameRandomLobbyParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while getting lobby",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)

					return
				}
				//updating queue with lobby own id
				updateQueueLobbyParams := db.UpdateQueueLobbyParams{
					LobbyID: getlobby.ID,
					UserID:  authPayloadKey.ParentId,
				}
				err = server.store.UpdateQueueLobby(ctx, updateQueueLobbyParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while updating queue with lobby",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				//updating queue with lobby oponent id
				updateQueueLobbyParams = db.UpdateQueueLobbyParams{
					LobbyID: getlobby.ID,
					UserID:  randomPlayer.UserID,
				}
				err = server.store.UpdateQueueLobby(ctx, updateQueueLobbyParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "Error while updating queue with lobby",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				//get 10 questions
				getQuestionsParams := db.GetQuestionsParams{
					Class:   ownInfo.Grade,
					Subject: data.Subject,
				}
				questions, err := server.store.GetQuestions(ctx, getQuestionsParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "could not get the questions",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				for i := 0; i < len(questions); i++ {
					createRandomLobbyQuestionsParams := db.CreateRandomLobbyQuestionsParams{
						LobbyID:        getlobby.ID,
						Questions:      questions[i].Questions,
						OptionsA:       questions[i].OptionsA,
						OptionsB:       questions[i].OptionsB,
						OptionsC:       questions[i].OptionsC,
						OptionsD:       questions[i].OptionsD,
						CorrectOptions: questions[i].CorrectOptions,
					}
					err := server.store.CreateRandomLobbyQuestions(ctx, createRandomLobbyQuestionsParams)
					if err != nil {
						errors := Errors{
							Type:       errormsg,
							Message:    "could not save questions",
							Error:      err,
							StatusCode: 500,
						}
						err = ws.WriteJSON(errors)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
						_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

						_ = beOffline(ctx, authPayloadKey.ParentId, server)

						_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
						return
					}
				}
				//count score
				getScoreListParams := db.GetScoreListParams{
					Player1ID: authPayloadKey.ParentId,
					Player2ID: randomPlayer.UserID,
				}
				//count score
				scoreList, _ := server.store.GetScoreList(ctx, getScoreListParams)
				indicators := int32(len(scoreList) + 1)
				//create Score ROW in and makes points as 0 for both
				createScoreParams := db.CreateScoreParams{
					Player1ID:    authPayloadKey.ParentId,
					Player2ID:    randomPlayer.UserID,
					Player1Point: 0,
					Player2Point: 0,
					PlayedTime:   time.Now(),
					Indicator:    indicators,
				}
				err = server.store.CreateScore(ctx, createScoreParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "can't create the score board",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				break
			}
		}
		if err != nil {
			//if no user found then
			errors := Errors{
				Type:       errormsg,
				Message:    "No one is in Queue Play with your friends",
				Error:      err,
				StatusCode: 400,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
	}
	for i := 0; i < 30000000000; i++ {
		//refresh your queue data
		//get won queue info
		ownQueueInfoUpdate, err = server.store.GetOwnQueueInfo(ctx, authPayloadKey.ParentId)
		if err != nil {

			errors := Errors{
				Type:       errormsg,
				Message:    "Error while fetching queue status",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return

		}
		if ownQueueInfoUpdate.Status == "connected" {
			break
		}
	}

	// get lobby again
	newLobby, err := server.store.GetGameRandomLobbyById(ctx, ownQueueInfoUpdate.LobbyID)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while fetching lobby status",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

		_ = beOffline(ctx, authPayloadKey.ParentId, server)

		_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
		return
	}
	if newLobby.Status == "connected" {

		if newLobby.UserID == authPayloadKey.ParentId {
			//get op details
			opDetails, err := server.store.GetChildDetail(ctx, newLobby.OpID)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "could not get the your information",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

				_ = beOffline(ctx, authPayloadKey.ParentId, server)

				_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
				return
			}
			//send connection status
			connectionResponse := ConnectionResponse{
				Type:             connection,
				ConnectionStatus: "Connected",
				OPName:           opDetails.FullName,
				OwnName:          ownInfo.FullName,
			}
			err = ws.WriteJSON(connectionResponse)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			//update the queue status of own
			updateQueueParams := db.UpdateQueueParams{
				Status: "playing",
				UserID: authPayloadKey.ParentId,
			}
			err = server.store.UpdateQueue(ctx, updateQueueParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "Error while updating queue",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

				_ = beOffline(ctx, authPayloadKey.ParentId, server)

				_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
				return
			}
			//scoure params
			getScoreListParams = db.GetScoreListParams{
				Player1ID: authPayloadKey.ParentId,
				Player2ID: newLobby.OpID,
			}

		} else if newLobby.OpID == authPayloadKey.ParentId {
			//get op details
			opDetails, err := server.store.GetChildDetail(ctx, newLobby.UserID)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "could not get the your information",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

				_ = beOffline(ctx, authPayloadKey.ParentId, server)

				_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
				return
			}
			//send connection status
			connectionResponse := ConnectionResponse{
				Type:             connection,
				ConnectionStatus: "Connected",
				OPName:           opDetails.FullName,
				OwnName:          ownInfo.FullName,
			}
			err = ws.WriteJSON(connectionResponse)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			//update the queue status of own
			updateQueueParams := db.UpdateQueueParams{
				Status: "playing",
				UserID: authPayloadKey.ParentId,
			}
			err = server.store.UpdateQueue(ctx, updateQueueParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "Error while updating queue",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

				_ = beOffline(ctx, authPayloadKey.ParentId, server)

				_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
				return
			}

			//scoure params
			getScoreListParams = db.GetScoreListParams{
				Player1ID: newLobby.UserID,
				Player2ID: authPayloadKey.ParentId,
			}
		} else {
			errors := Errors{
				Type:       errormsg,
				Message:    "Issue While Connecting",
				Error:      err,
				StatusCode: 400,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}

		}
	} else {
		errors := Errors{
			Type:       errormsg,
			Message:    "You are not connected to any one plz try again",
			Error:      err,
			StatusCode: 400,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

		_ = beOffline(ctx, authPayloadKey.ParentId, server)

		_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
		return
	}
	//check both user status is playing or not
	// //wait for other user
	// rsp = Resp{
	// 	Message: "wait for other user",
	// }
	// err = ws.WriteJSON(rsp)
	// if err != nil {
	// 	log.Println("error write json: " + err.Error())
	// }

	for i := 0; i < 30000000000; i++ {
		getOpStatus, err := server.store.GetOwnQueueInfo(ctx, newLobby.OpID)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "Error while fetching queue info",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		getOwnStatus, err := server.store.GetOwnQueueInfo(ctx, newLobby.UserID)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "Error while fetching queue info",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		if getOpStatus.Status == "playing" && getOwnStatus.Status == "playing" {
			break
		}
	}

	//update status lobby connected to playing
	updateGameRandomLobbyParams := db.UpdateGameRandomLobbyParams{
		Status: "playing",
		ID:     newLobby.ID,
	}
	err = server.store.UpdateGameRandomLobby(ctx, updateGameRandomLobbyParams)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while updating Lobby Status",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

		_ = beOffline(ctx, authPayloadKey.ParentId, server)

		_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
		return
	}
	// get lobby again
	newLobby, err = server.store.GetGameRandomLobbyById(ctx, ownQueueInfoUpdate.LobbyID)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "Error while fetching lobby status",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

		_ = beOffline(ctx, authPayloadKey.ParentId, server)

		_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
		return
	}
	//count score
	scoreList, _ := server.store.GetScoreList(ctx, getScoreListParams)
	indicator := int32(len(scoreList))

	//check if ststus is playing or not
	if newLobby.Status == "playing" {
		//get 10 questions
		questions, err := server.store.GetRandomLobbyQuestions(ctx, newLobby.ID)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the questions",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		//map questions
		for i := 0; i < len(questions); i++ {
			var rsp GameQuestionsResponse
			var ans GameAnswerRequest
			var option string
			// send questions to both users
			if questions[i].CorrectOptions == questions[i].OptionsA {
				option = "a"
			} else if questions[i].CorrectOptions == questions[i].OptionsB {
				option = "b"
			} else if questions[i].CorrectOptions == questions[i].OptionsC {
				option = "c"
			} else if questions[i].CorrectOptions == questions[i].OptionsD {
				option = "d"
			} else {
				option = "none"
			}
			rsp = GameQuestionsResponse{
				Type:      questionmsg,
				ID:        int32(i + 1),
				Questions: questions[i].Questions,
				OptionsA:  questions[i].OptionsA,
				OptionsB:  questions[i].OptionsB,
				OptionsC:  questions[i].OptionsC,
				OptionsD:  questions[i].OptionsD,
				Correct:   option,
			}
			err := ws.WriteJSON(rsp)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}

			if newLobby.UserID == authPayloadKey.ParentId {
				//reade answer from player 1
				err = ws.ReadJSON(&ans)
				if err != nil {
					log.Println("error read json for answer")
					log.Fatal(err)
				}
				//check score information
				getScoreParams := db.GetScoreParams{
					Player1ID: authPayloadKey.ParentId,
					Player2ID: newLobby.OpID,
					Indicator: indicator,
				}
				score, err := server.store.GetScore(ctx, getScoreParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "could not get the score",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				//verify answer
				//do not update point if answer is wrong
				var wrong WrongResponse
				//check correct answer and send the correct options to player
				if questions[i].CorrectOptions == questions[i].OptionsA {
					if ans.CorrectAnswer == "a" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerOneParams := db.UpdateScorePlayerOneParams{
							Player1Point: score.Player1Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerOne(ctx, updateScorePlayerOneParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player1Point + 5,
							Oscore:  score.Player2Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "a",
							Score:         score.Player1Point,
							Oscore:        score.Player2Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsB {
					if ans.CorrectAnswer == "b" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerOneParams := db.UpdateScorePlayerOneParams{
							Player1Point: score.Player1Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerOne(ctx, updateScorePlayerOneParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player1Point + 5,
							Oscore:  score.Player2Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "b",
							Score:         score.Player1Point,
							Oscore:        score.Player2Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsC {
					if ans.CorrectAnswer == "c" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerOneParams := db.UpdateScorePlayerOneParams{
							Player1Point: score.Player1Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerOne(ctx, updateScorePlayerOneParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player1Point + 5,
							Oscore:  score.Player2Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "c",
							Score:         score.Player1Point,
							Oscore:        score.Player2Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsD {
					if ans.CorrectAnswer == "d" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerOneParams := db.UpdateScorePlayerOneParams{
							Player1Point: score.Player1Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerOne(ctx, updateScorePlayerOneParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player1Point + 5,
							Oscore:  score.Player2Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						wrong = WrongResponse{
							//send wrong response
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "d",
							Score:         score.Player1Point,
							Oscore:        score.Player2Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}

					}
				} else {
					wrong = WrongResponse{
						//send wrong response
						Type:          wrongmsg,
						Message:       "worng",
						CorrectAnswer: "NO Options",
						Score:         score.Player1Point,
						Oscore:        score.Player2Point,
					}
					err := ws.WriteJSON(wrong)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
				}
			} else if newLobby.OpID == authPayloadKey.ParentId {
				//reade answer from player 2
				err = ws.ReadJSON(&ans)
				if err != nil {
					log.Println("error read json")
					log.Fatal(err)
				}
				//check the points of player 2
				getScoreParams := db.GetScoreParams{
					Player1ID: newLobby.UserID,
					Player2ID: authPayloadKey.ParentId,
					Indicator: indicator,
				}
				score, err := server.store.GetScore(ctx, getScoreParams)
				if err != nil {
					fmt.Println(err.Error())
					errors := Errors{
						Type:       errormsg,
						Message:    "could not get the score",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				//verify answer
				//do not update point if answer is wrong
				var wrong WrongResponse
				//check correct answer and send the correct options to player
				if questions[i].CorrectOptions == questions[i].OptionsA {
					if ans.CorrectAnswer == "a" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerTwoParams := db.UpdateScorePlayerTwoParams{
							Player2Point: score.Player2Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerTwo(ctx, updateScorePlayerTwoParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player2Point + 5,
							Oscore:  score.Player1Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "a",
							Score:         score.Player2Point,
							Oscore:        score.Player1Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsB {
					if ans.CorrectAnswer == "b" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerTwoParams := db.UpdateScorePlayerTwoParams{
							Player2Point: score.Player2Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerTwo(ctx, updateScorePlayerTwoParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player2Point + 5,
							Oscore:  score.Player1Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "b",
							Score:         score.Player2Point,
							Oscore:        score.Player1Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsC {
					if ans.CorrectAnswer == "c" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerTwoParams := db.UpdateScorePlayerTwoParams{
							Player2Point: score.Player2Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerTwo(ctx, updateScorePlayerTwoParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player2Point + 5,
							Oscore:  score.Player1Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						//send wrong response
						wrong = WrongResponse{
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "c",
							Score:         score.Player2Point,
							Oscore:        score.Player1Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else if questions[i].CorrectOptions == questions[i].OptionsD {
					if ans.CorrectAnswer == "d" {
						// add 5 point in player 1 for correct answer
						updateScorePlayerTwoParams := db.UpdateScorePlayerTwoParams{
							Player2Point: score.Player2Point + 5,
							ID:           score.ID,
						}
						err = server.store.UpdateScorePlayerTwo(ctx, updateScorePlayerTwoParams)
						if err != nil {
							errors := Errors{
								Type:       errormsg,
								Message:    "could not updated the first player score",
								Error:      err,
								StatusCode: 500,
							}
							err = ws.WriteJSON(errors)
							if err != nil {
								log.Println("error write json: " + err.Error())
							}
							_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

							_ = beOffline(ctx, authPayloadKey.ParentId, server)

							_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
							return
						}
						//send correct response
						correct := CorrectResponse{
							Type:    correctmsg,
							Message: "correct",
							Score:   score.Player2Point + 5,
							Oscore:  score.Player1Point,
						}
						err = ws.WriteJSON(correct)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					} else {
						wrong = WrongResponse{
							//send wrong response
							Type:          wrongmsg,
							Message:       "worng",
							CorrectAnswer: "d",
							Score:         score.Player2Point,
							Oscore:        score.Player1Point,
						}
						err := ws.WriteJSON(wrong)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
					}
				} else {
					wrong = WrongResponse{
						//send wrong response
						Type:          wrongmsg,
						Message:       "worng",
						CorrectAnswer: "NO Options",
						Score:         score.Player2Point,
						Oscore:        score.Player1Point,
					}
					err := ws.WriteJSON(wrong)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
				}
			} else {
				errors := Errors{
					Type:       errormsg,
					Message:    "You are not connected to any one plz try again",
					Error:      err,
					StatusCode: 400,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
			}
		}
	} else {
		errors := Errors{
			Type:       errormsg,
			Message:    "You are not connected to any one plz try again",
			Error:      err,
			StatusCode: 400,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

		_ = beOffline(ctx, authPayloadKey.ParentId, server)

		_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
		return
	}
	//After 10 questions completed (Game Completed)
	//send final response
	if newLobby.UserID == authPayloadKey.ParentId {
		//check latest points of player 1
		getScoreParams := db.GetScoreParams{
			Player1ID: authPayloadKey.ParentId,
			Player2ID: newLobby.OpID,
			Indicator: indicator,
		}
		score, err := server.store.GetScore(ctx, getScoreParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the score",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
		}
		//get own details
		ownDetails, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the your information",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
		}
		//get friends details
		opDetails, err := server.store.GetChildDetail(ctx, newLobby.OpID)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the your friends infomation",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
		}
		//check who wins the game
		var msg string
		if score.Player1Point == score.Player2Point {
			msg = "Match Draw"
		} else if score.Player1Point < score.Player2Point {
			msg = fmt.Sprintf("%s scores more then you", opDetails.FullName)
		} else {
			msg = fmt.Sprintf("you scores more then %s", opDetails.FullName)
		}
		//save score to score table
		createScorePointParams := db.CreateScorePointParams{
			UserID:    authPayloadKey.ParentId,
			OwnPoints: score.Player1Point,
			OpID:      newLobby.OpID,
			OpPoints:  score.Player2Point,
			Subject:   data.Subject,
			CreatedAt: time.Now(),
		}
		err = server.store.CreateScorePoint(ctx, createScorePointParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not save score",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
		}
		//send the final response
		fnlp := FinalResponse{
			Type:     finalmsg,
			Message:  msg,
			Score:    score.Player1Point,
			Oscore:   score.Player2Point,
			YourName: ownDetails.FullName,
			OPName:   opDetails.FullName,
		}
		err = ws.WriteJSON(fnlp)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
	} else if newLobby.OpID == authPayloadKey.ParentId {
		//check the final points
		getScoreParams := db.GetScoreParams{
			Player1ID: newLobby.UserID,
			Player2ID: authPayloadKey.ParentId,
			Indicator: indicator,
		}
		score, err := server.store.GetScore(ctx, getScoreParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the score",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		//get own information
		ownDetails, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the your information",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		//get friends information
		opDetails, err := server.store.GetChildDetail(ctx, newLobby.UserID)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not get the your friends information",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
		//check who wins the game
		var msg string
		if score.Player1Point == score.Player2Point {
			msg = "Match Draw"
		} else if score.Player2Point < score.Player1Point {
			msg = fmt.Sprintf("%s scores more then you", opDetails.FullName)
		} else {
			msg = fmt.Sprintf("you scores more then %s", opDetails.FullName)
		}
		//save score to score table
		createScorePointParams := db.CreateScorePointParams{
			UserID:    authPayloadKey.ParentId,
			OwnPoints: score.Player2Point,
			OpID:      newLobby.UserID,
			OpPoints:  score.Player1Point,
			Subject:   data.Subject,
			CreatedAt: time.Now(),
		}
		err = server.store.CreateScorePoint(ctx, createScorePointParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not save score",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
		}
		//send the final response
		fnlp := FinalResponse{
			Type:     finalmsg,
			Message:  msg,
			Score:    score.Player2Point,
			Oscore:   score.Player1Point,
			YourName: ownDetails.FullName,
			OPName:   opDetails.FullName,
		}
		err = ws.WriteJSON(fnlp)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		//update status lobby connected to playing
		updateGameRandomLobbyParams := db.UpdateGameRandomLobbyParams{
			Status: "completed",
			ID:     newLobby.ID,
		}
		err = server.store.UpdateGameRandomLobby(ctx, updateGameRandomLobbyParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "Error while updating Lobby Status",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

			_ = beOffline(ctx, authPayloadKey.ParentId, server)

			_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
			return
		}
	} else {
		errors := Errors{
			Type:       errormsg,
			Message:    "You are not connected",
			Error:      err,
			StatusCode: 400,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
	}
	if newLobby.UserID == authPayloadKey.ParentId {
		for i := 0; i < 3000000000; i++ {
			// get lobby again
			newLobby, _ = server.store.GetGameRandomLobbyById(ctx, ownQueueInfoUpdate.LobbyID)

			if newLobby.Status == "completed" {
				err = server.store.DeleteGameFriendLobby(ctx, newLobby.ID)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "faild to delete the lobby",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

					_ = beOffline(ctx, authPayloadKey.ParentId, server)

					_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)
					return
				}
				break
			}
		}
	}
	_ = removefakeLobby(ctx, authPayloadKey.ParentId, server)

	_ = beOffline(ctx, authPayloadKey.ParentId, server)

	_ = removeRandomLobby(ctx, authPayloadKey.ParentId, server)

}

func (server *Server) RemoveRandomLobby(ctx *gin.Context) {
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//get queue info
	queueInfo, err := server.store.GetOwnQueueInfo(ctx, authPayloadKey.ParentId)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("you are not online"))
		return
	}
	//delete lobby
	err = server.store.DeleteGameRandomLobby(ctx, queueInfo.LobbyID)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("you are not in lobby"))
		return
	}
	//get queue info
	fake, err := server.store.GetGameRandomLobbyById(ctx, authPayloadKey.ParentId)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("you are not online"))
		return
	}
	//delete lobby
	err = server.store.DeleteGameRandomLobby(ctx, fake.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("you are not in lobby"))
		return
	}

	ctx.JSON(http.StatusOK, GenerateResponse("lobby successfully deleted"))
}
func removefakeLobby(ctx *gin.Context, userID int32, server *Server) error {
	//get queue info
	fake, err := server.store.GetGameRandomLobbyById(ctx, userID)
	if err != nil {
		return err
	}
	//delete lobby
	err = server.store.DeleteGameRandomLobby(ctx, fake.ID)
	if err != nil {
		return err
	}

	return nil
}
func removeRandomLobby(ctx *gin.Context, userID int32, server *Server) error {
	//get queue info
	queueInfo, err := server.store.GetOwnQueueInfo(ctx, userID)
	if err != nil {
		return err
	}
	//delete lobby
	err = server.store.DeleteGameRandomLobby(ctx, queueInfo.LobbyID)
	if err != nil {
		return err
	}

	return nil
}
func beOffline(ctx *gin.Context, userID int32, server *Server) error {

	err := server.store.RemoveQueue(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) BeOffline(ctx *gin.Context) {
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.store.RemoveQueue(ctx, authPayloadKey.ParentId)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("you are not online"))
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("you are now offline"))
}
