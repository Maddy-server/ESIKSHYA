package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"Edtech_Golang/util"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//webSocket game with friends
func (server *Server) GameWithFriends(ctx *gin.Context) {
	//variable initialization
	errormsg := "error"
	correctmsg := "CorrectAnswer"
	wrongmsg := "WrongAnswer"
	finalmsg := "gameCompleted"
	questionmsg := "questions"
	connection := "connected"
	//variable decleration
	var data GameWithFriendsRequest
	var getGameFriendLobbyParams db.GetGameFriendLobbyParams
	//var friendsLobby db.GameFriendLobby
	var getScoreListParams db.GetScoreListParams
	//variable initialization
	currentTime := time.Now()
	scoreTime := time.Now()
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Error("error get connection ", err)
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()
	//Read data in ws
	err = ws.ReadJSON(&data)
	if err != nil {
		log.Println("error read json")
		log.Fatal(err)
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//condition for friends who wants to play is "waiting" and friends who gets notification is "connected"
	if data.Status == "waiting" {
		//check a lobby
		getGameFriendLobbyParams = db.GetGameFriendLobbyParams{
			UserID: authPayloadKey.ParentId,
			OpID:   data.OponentId,
		}
		friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
		if err != nil {
			//create a lobby
			createGameFriendLobbyParams := db.CreateGameFriendLobbyParams{
				UserID:    authPayloadKey.ParentId,
				OpID:      data.OponentId,
				Status:    "waiting",
				CreatedAt: currentTime,
			}
			err := server.store.CreateGameFriendLobby(ctx, createGameFriendLobbyParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "can't create game lobby ",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				return
			}
			//get lobby for id
			friendsLobby, err = server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "did not get lobby",
					Error:      err,
					StatusCode: 400,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				return
			}
			ownInfo, err := server.store.GetChildDetail(ctx, authPayloadKey.ParentId)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "could not get own info",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
				return
			}
			//save notifications
			createGameNotificationParams := db.CreateGameNotificationParams{
				UserID:      data.OponentId,
				OponentID:   authPayloadKey.ParentId,
				Title:       "Game Notification",
				Type:        "gamenotification",
				Description: fmt.Sprintf("%s invited to Play a game", ownInfo.FullName),
				CreatedAt:   time.Now(),
				Subject:     data.Subject,
				Status:      "connected",
				Grade:       data.Grade,
			}
			err = server.store.CreateGameNotification(ctx, createGameNotificationParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "Your friend is offline",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
				return
			}
			//send Notification to Player 2
			err = server.store.SendChildNotification(
				ctx,
				data.OponentId,
				"GameNotification",
				fmt.Sprintf("%s invited to Play a game", ownInfo.FullName),
			)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "could not send notification",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
			}
			//get 10 questions
			getQuestionsParams := db.GetQuestionsParams{
				Class:   data.Grade,
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
				RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
				return
			}
			for i := 0; i < len(questions); i++ {
				createFriendsLobbyQuestionsParams := db.CreateFriendsLobbyQuestionsParams{
					LobbyID:        friendsLobby.ID,
					Questions:      questions[i].Questions,
					OptionsA:       questions[i].OptionsA,
					OptionsB:       questions[i].OptionsB,
					OptionsC:       questions[i].OptionsC,
					OptionsD:       questions[i].OptionsD,
					CorrectOptions: questions[i].CorrectOptions,
				}
				err := server.store.CreateFriendsLobbyQuestions(ctx, createFriendsLobbyQuestionsParams)
				if err != nil {
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
						RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
						return
					}
				}
			}

			//wait 15sec for friends

			for i := 0; i < 30000000000; i++ {

				//fetch the lobby information
				friendsLobby, err = server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "could not find game lobby ",
						Error:      err,
						StatusCode: 400,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
					return
				}
				//count score
				getScoreListParams = db.GetScoreListParams{
					Player1ID: authPayloadKey.ParentId,
					Player2ID: data.OponentId,
				}
				//count score
				scoreList, _ := server.store.GetScoreList(ctx, getScoreListParams)
				indicators := int32(len(scoreList) + 1)
				//check the status of lobby as connected
				if friendsLobby.Status == "connected" {
					//create Score ROW in and makes points as 0 for both
					createScoreParams := db.CreateScoreParams{
						Player1ID:    authPayloadKey.ParentId,
						Player2ID:    data.OponentId,
						Player1Point: 0,
						Player2Point: 0,
						PlayedTime:   scoreTime,
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
						RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
						return
					}
					//get op details
					opDetails, err := server.store.GetChildDetail(ctx, data.OponentId)
					if err != nil {
						errors := Errors{
							Type:       errormsg,
							Message:    "could not get the op information",
							Error:      err,
							StatusCode: 500,
						}
						err = ws.WriteJSON(errors)
						if err != nil {
							log.Println("error write json: " + err.Error())
						}
						RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
						return
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
						RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
						return
					}
					connectionResponse := ConnectionResponse{
						Type:             connection,
						ConnectionStatus: "Connected",
						OPName:           opDetails.FullName,
						OwnName:          ownDetails.FullName,
					}
					err = ws.WriteJSON(connectionResponse)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					break
				}
				//make systey sleep for 5sec
				//time.Sleep(30 * time.Second)
			}
		}
		// updateGameFriendLobbyParams := db.UpdateGameFriendLobbyParams{
		// 	Status: "waiting",
		// 	ID:     friendsLobby.ID,
		// }
		// err = server.store.UpdateGameFriendLobby(ctx, updateGameFriendLobbyParams)
		// if err != nil {
		// 	if err != nil {
		// 		errors := Errors{
		// 			Message:    "can't create lobby",
		// 			Error:      err,
		// 			StatusCode: 500,
		// 		}
		// 		err = ws.WriteJSON(errors)
		// 		if err != nil {
		// 			log.Println("error write json: " + err.Error())
		// 		}
		// 		return

		// 	}
		// }

	} else {
		//check is lobby is created or not
		getGameFriendLobbyParams = db.GetGameFriendLobbyParams{
			UserID: data.OponentId,
			OpID:   authPayloadKey.ParentId,
		}
		friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not find game lobby ",
				Error:      err,
				StatusCode: 400,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
			return

		}
		//update game lobby status as connected:- it sows that notification friends is now connected
		updateGameFriendLobbyParams := db.UpdateGameFriendLobbyParams{
			Status: "connected",
			ID:     friendsLobby.ID,
		}
		err = server.store.UpdateGameFriendLobby(ctx, updateGameFriendLobbyParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "can't get the connection with your friends",
				Error:      err,
				StatusCode: 400,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
			return

		}
		//count score
		getScoreListParams = db.GetScoreListParams{
			Player1ID: data.OponentId,
			Player2ID: authPayloadKey.ParentId,
		}
		for i := 0; i < 30000000000; i++ {

			//fetch the lobby information
			friendsLobby, err = server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "could not find game lobby ",
					Error:      err,
					StatusCode: 500,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
				return
			}

			//check the status of lobby as connected
			if friendsLobby.Status == "connected" {
				//get op details
				opDetails, err := server.store.GetChildDetail(ctx, data.OponentId)
				if err != nil {
					errors := Errors{
						Type:       errormsg,
						Message:    "could not get the op information",
						Error:      err,
						StatusCode: 500,
					}
					err = ws.WriteJSON(errors)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
					return
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
					RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
					return
				}
				connectionResponse := ConnectionResponse{
					Type:             connection,
					ConnectionStatus: "Connected",
					OPName:           opDetails.FullName,
					OwnName:          ownDetails.FullName,
				}
				err = ws.WriteJSON(connectionResponse)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
				break
			}
		}

	}

	friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
	if err != nil {
		errors := Errors{
			Type:       errormsg,
			Message:    "could not find game lobby ",
			Error:      err,
			StatusCode: 500,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
		return
	}
	//check status is connected or not
	if friendsLobby.Status == "connected" {
		time.Sleep(5 * time.Second)
		//count score
		scoreList, _ := server.store.GetScoreList(ctx, getScoreListParams)
		indicator := int32(len(scoreList))
		//update lobby as both are playing
		updateGameFriendLobbyParams := db.UpdateGameFriendLobbyParams{
			Status: "playing",
			ID:     friendsLobby.ID,
		}
		err = server.store.UpdateGameFriendLobby(ctx, updateGameFriendLobbyParams)
		if err != nil {
			errors := Errors{
				Type:       errormsg,
				Message:    "could not update the playing status",
				Error:      err,
				StatusCode: 500,
			}
			err = ws.WriteJSON(errors)
			if err != nil {
				log.Println("error write json: " + err.Error())
			}
			RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
			return
		}

		//get 10 questions
		questions, err := server.store.GetFriendsLobbyQuestions(ctx, friendsLobby.ID)
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
			RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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

			if data.Status == "waiting" {
				//reade answer from waiting person
				err = ws.ReadJSON(&ans)
				if err != nil {
					log.Println("error read json for answer")
					log.Fatal(err)
				}

				//check score information
				getScoreParams := db.GetScoreParams{
					Player1ID: authPayloadKey.ParentId,
					Player2ID: data.OponentId,
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
					err = ws.WriteJSON(getScoreParams)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}
					RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
			} else {
				//reade answer from player 2
				err = ws.ReadJSON(&ans)
				if err != nil {
					log.Println("error read json")
					log.Fatal(err)
				}
				//check the points of player 2
				getScoreParams := db.GetScoreParams{
					Player1ID: data.OponentId,
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
					RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
							RemoveGameFriendsLobby(ctx, friendsLobby.ID, server)
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
			}
		}
		//After 10 questions completed (Game Completed)
		if data.Status == "waiting" {
			//check latest points of player 1
			getScoreParams := db.GetScoreParams{
				Player1ID: authPayloadKey.ParentId,
				Player2ID: data.OponentId,
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
			opDetails, err := server.store.GetChildDetail(ctx, data.OponentId)
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
				OpID:      data.OponentId,
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
		} else {
			//check the final points
			getScoreParams := db.GetScoreParams{
				Player1ID: data.OponentId,
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
			}
			//get friends information
			opDetails, err := server.store.GetChildDetail(ctx, data.OponentId)
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
			}
			//save score to score table
			createScorePointParams := db.CreateScorePointParams{
				UserID:    authPayloadKey.ParentId,
				OwnPoints: score.Player1Point,
				OpID:      data.OponentId,
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
			//check who wins the game
			var msg string
			if score.Player1Point == score.Player2Point {
				msg = "Match Draw"
			} else if score.Player2Point < score.Player1Point {
				msg = fmt.Sprintf("%s scores more then you", opDetails.FullName)
			} else {
				msg = fmt.Sprintf("you scores more then %s", opDetails.FullName)
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
		}
		if data.Status == "waiting" {
			//check a lobby
			getGameFriendLobbyParams = db.GetGameFriendLobbyParams{
				UserID: authPayloadKey.ParentId,
				OpID:   data.OponentId,
			}
			friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "faild to get the lobby",
					Error:      err,
					StatusCode: 400,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
			}
			err = server.store.DeleteGameFriendLobby(ctx, friendsLobby.ID)
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
			}
		}
	} else {
		//remove lobby
		if data.Status == "waiting" {
			//check a lobby
			getGameFriendLobbyParams = db.GetGameFriendLobbyParams{
				UserID: authPayloadKey.ParentId,
				OpID:   data.OponentId,
			}
			friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
			if err != nil {
				errors := Errors{
					Type:       errormsg,
					Message:    "faild to get the lobby",
					Error:      err,
					StatusCode: 400,
				}
				err = ws.WriteJSON(errors)
				if err != nil {
					log.Println("error write json: " + err.Error())
				}
			}
			err = server.store.DeleteGameFriendLobby(ctx, friendsLobby.ID)
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
			}
		}
		//friends not connected in the system
		errors := Errors{
			Type:       errormsg,
			Message:    "your friend did not join the game",
			Error:      err,
			StatusCode: 400,
		}
		err = ws.WriteJSON(errors)
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
	}
}

type RemoveLobbyRequest struct {
	OponentId int32 `json:"oponent_id"`
}

func (server *Server) RemoveLobby(ctx *gin.Context) {
	var req RemoveLobbyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get payload data from access_token
	authPayloadKey := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//check a lobby
	getGameFriendLobbyParams := db.GetGameFriendLobbyParams{
		UserID: authPayloadKey.ParentId,
		OpID:   req.OponentId,
	}
	friendsLobby, err := server.store.GetGameFriendLobby(ctx, getGameFriendLobbyParams)
	if err != nil {
		ctx.JSON(http.StatusOK, GenerateResponse("we couldnot found  lobby"))
		return
	}
	err = server.store.DeleteGameFriendLobby(ctx, friendsLobby.ID)
	if err != nil {
		saveErr := util.NewInternalServerError("error while getting child detail", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	ctx.JSON(http.StatusOK, GenerateResponse("lobby successfully deleted"))
}

func RemoveGameFriendsLobby(ctx *gin.Context, lobbyId int32, server *Server) {
	err := server.store.DeleteGameFriendLobby(ctx, lobbyId)
	if err != nil {

		return
	}
}
