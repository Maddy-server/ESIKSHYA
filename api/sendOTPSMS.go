package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// type Response struct{
// 	Count int32 `json:"count"`
// 	ResponseCode int32 `json:"response_code"`
// 	Response string `json:"response"`
// 	Message_id int64 `json:"message_id"`
// 	Credit_Consumed int32 `json:"credit_consumed"`
// 	Credit_Available int32 `json:"credit_available"`
// }

func (server *Server) SendOTPSMS(to string, otp string) error {
	// var smsResponse Response
	const myUrl = "http://api.sparrowsms.com/v2/sms/"
	textSms := fmt.Sprintf("%s is your edtech verification code \n Ref: ", otp)
	data := url.Values{}

	data.Add("token", server.config.SMS.Token)
	data.Add("from", server.config.SMS.From)
	data.Add("to", to)
	data.Add("text", textSms)

	response, err := http.PostForm(myUrl, data)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	return nil
}
