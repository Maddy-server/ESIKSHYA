package api

import (
	"fmt"
	"log"
	"net/smtp"
)

// SendOTP sends otp by email
func (server *Server) SendOTP(email, otpCode string) error {
	from := server.config.Mail.Username
	pass := server.config.Mail.Password

	to := email
	msg := "From: Edtech < namaste@mahajodi.space >\r\n" +
		"To:" + email + "\r\n" +
		"Subject: One Time Password \r\n\r\n" +
		"Welcome To the Edtech\r\n" +
		"Dear Sir/Madam," + "\r\n" +
		"Your OTP code to verify Email Id is : " + otpCode + "\r\n\r\n" +
		"Please DO NOT SHARE your otp with anyone" + "\n\n" +

		"Warm regards,\r\n" +
		"Edtech"
	err := smtp.SendMail(server.config.Mail.Host+":"+server.config.Mail.Port,
		smtp.PlainAuth("", from, pass, server.config.Mail.Host),
		from, []string{to}, []byte(msg))

	// handling the errors
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)

	}
	return nil
}
