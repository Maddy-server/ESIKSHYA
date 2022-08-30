package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/token"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (server *Server) SinglePlayerGame(ctx *gin.Context) {

	//variable initialization
	errormsg := "error"
	correctmsg := "CorrectAnswer"
	wrongmsg := "WrongAnswer"
	finalmsg := "gameCompleted"
	questionmsg := "questions"

	//variable decleration

	var data GameWithRandomRequest
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
	//count score
	getScoreListParams := db.GetScoreListParams{
		Player1ID: authPayloadKey.ParentId,
		Player2ID: 0,
	}
	//count score
	scoreList, _ := server.store.GetScoreList(ctx, getScoreListParams)
	indicators := int32(len(scoreList) + 1)
	//create Score ROW in and makes points as 0 for both
	createScoreParams := db.CreateScoreParams{
		Player1ID:    authPayloadKey.ParentId,
		Player2ID:    0,
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
		return
	}
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
		//reade answer
		err = ws.ReadJSON(&ans)
		if err != nil {
			log.Println("error read json for answer")
			log.Fatal(err)
		}
		//check score information
		getScoreParams := db.GetScoreParams{
			Player1ID: authPayloadKey.ParentId,
			Player2ID: 0,
			Indicator: indicators,
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
	}

	//check latest points of player
	getScoreParams := db.GetScoreParams{
		Player1ID: authPayloadKey.ParentId,
		Player2ID: 0,
		Indicator: indicators,
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

	msg := fmt.Sprintf("you score %d out of 50", score.Player1Point)

	//save score to score table
	createScorePointParams := db.CreateScorePointParams{
		UserID:    authPayloadKey.ParentId,
		OwnPoints: score.Player1Point,
		OpID:      0,
		OpPoints:  0,
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
		YourName: ownInfo.FullName,
	}
	err = ws.WriteJSON(fnlp)
	if err != nil {
		log.Println("error write json: " + err.Error())
	}
}
