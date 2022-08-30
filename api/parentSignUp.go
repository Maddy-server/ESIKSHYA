package api

import (
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/util"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var OTP string

type parentSignUpRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type parentVerifyRequest struct {
	Phone string `json:"phone" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

type setPasswordRequst struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AleadyRegisteredParentResponse struct {
	Message    string `json:"message"`
	IsVerified bool   `json:"is_verified"`
}

//ParentSignUp with phone number
func (server *Server) ParentSignUp(ctx *gin.Context) {
	var req parentSignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	// save phone number into db
	err := server.store.CreateParent(ctx, req.Phone)
	if err != nil {
		// log.Println(err)
		if strings.Contains(err.Error(), "Error 1062:") {
			//generate new otp
			OTP = util.GenerateOTP(6)
			//converting otp string to type sql.nullstring
			otpNS := util.CreateNullString(true, OTP)
			//save otp arguments
			saveOtpArg := db.SaveOTPParams{
				Otp:   otpNS,
				Phone: req.Phone,
			}
			//update otp into db
			saveErr := server.store.SaveOTP(ctx, saveOtpArg)
			if saveErr != nil {
				saveErr := util.NewInternalServerError("error when trying to save otp", errors.New("database error"))
				ctx.JSON(saveErr.Status(), saveErr)
				return
			}

			//get that parent
			par, err := server.store.GetParent(ctx, req.Phone)
			if err != nil {
				fmt.Println(err)
				getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
				ctx.JSON(getErr.Status(), getErr)
				return
			}
			if !par.Isverified.Bool {
				//parent is not verified go to verification
				//send otp
				sendSMSErr := server.SendOTPSMS(req.Phone, OTP)
				if sendSMSErr != nil {
					sendSMSErr := util.NewInternalServerError("error when trying to send otp", errors.New("send otp error"))
					ctx.JSON(sendSMSErr.Status(), sendSMSErr)
					return
				}
				rsp := AleadyRegisteredParentResponse{
					Message:    "otp is sent to provided mobile number",
					IsVerified: false,
				}
				//send response
				ctx.JSON(http.StatusOK, rsp)
				return
			}
			//parent is verified
			//send response
			rsp := AleadyRegisteredParentResponse{
				Message:    "parent already registered and verified",
				IsVerified: true,
			}
			//send response
			ctx.JSON(http.StatusOK, rsp)
			return
		}
		saveErr := util.NewInternalServerError("error when trying to save parent", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//generate otp
	OTP = util.GenerateOTP(6)
	//converting otp string to type sql.nullstring
	otpNS := util.CreateNullString(true, OTP)
	//save otp arguments
	saveOtpArg := db.SaveOTPParams{
		Otp:   otpNS,
		Phone: req.Phone,
	}
	//save otp into db
	saveErr := server.store.SaveOTP(ctx, saveOtpArg)
	if saveErr != nil {
		saveErr := util.NewInternalServerError("error when trying to save otp", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//send back otp through sms
	sendSMSErr := server.SendOTPSMS(req.Phone, OTP)
	if sendSMSErr != nil {
		sendSMSErr := util.NewInternalServerError("error when trying to send otp", errors.New("send otp error"))
		ctx.JSON(sendSMSErr.Status(), sendSMSErr)
		return
	}
	//sending back respone
	ctx.JSON(http.StatusCreated, GenerateResponse("Otp is sent to provided mobile number"))
}

//ParentVerify verify number with otp
func (server *Server) ParentVerify(ctx *gin.Context) {

	//post phone and otp
	var req parentVerifyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//get phone and otp from db
	i, err := server.store.GetParent(ctx, req.Phone)
	if err != nil {
		// fmt.Println(err)
		// if parent doesnot exists
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			getErr := util.NewRestError("parent with provided mobile doesnot exists", http.StatusOK, "error when trying to get parent", nil)
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	//check if parent is already verified
	if i.Isverified.Bool {
		ctx.JSON(http.StatusOK, GenerateResponse("mobile number already verified"))
		return
	}
	//make sure otp match
	p := strings.Compare(i.Phone, req.Phone)
	o := strings.Compare(i.Otp.String, req.Otp)
	//set isVerified column to 1
	if p == 0 && o == 0 {
		err := server.store.SetVerification(ctx, i.Phone)
		if err != nil {
			setErr := util.NewInternalServerError("error when trying to save otp", errors.New("database error"))
			ctx.JSON(setErr.Status(), setErr)
			return
		}
		//delete otp after verification
		removeErr := server.store.RemoveOTP(ctx, i.Phone)
		if removeErr != nil {
			removeErr := util.NewInternalServerError("error when trying to remove otp", errors.New("database error"))
			ctx.JSON(removeErr.Status(), removeErr)
			return
		}
		ctx.JSON(http.StatusOK, GenerateResponse("mobile number verified successfully"))
		return

	} else {
		ctx.JSON(http.StatusForbidden, GenerateResponse("cannot verify mobile number"))
		return
	}
	//return accesstoken
}

//ParentSetPassword create password
func (server *Server) ParentSetPassword(ctx *gin.Context) {
	//post password associated with provided phone
	var req setPasswordRequst
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		reqErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(reqErr.Status(), reqErr)
		return
	}
	//get parent from db
	i, err := server.store.GetParent(ctx, req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			getErr := util.NewRestError("parent with provided mobile doesnot exists", http.StatusOK, "error when trying to get parent", nil)
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	//before saving check if parent is verified or not
	if !i.Isverified.Bool {
		verifyErr := util.NewUnauthorizedError("provided parent not verified")
		ctx.JSON(verifyErr.Status(), verifyErr)
		return
	}
	//hash the provided password
	hashedPassword, hashErr := util.HashPassword(req.Password)
	if hashErr != nil {
		hashErr := util.NewInternalServerError("error when trying to hash password", errors.New("internal error"))
		ctx.JSON(hashErr.Status(), hashErr)
		return
	}
	//converting hashed password to sql.nullstring type
	hpw := util.CreateNullString(true, hashedPassword)
	//preparing arguments to save password
	savePasswordArguments := db.SavepasswordParams{
		Password: hpw,
		Phone:    i.Phone,
	}
	//store hashed password
	savePwErr := server.store.Savepassword(ctx, savePasswordArguments)
	if savePwErr != nil {
		savePwErr := util.NewInternalServerError("error when trying to save password", errors.New("internal error"))
		ctx.JSON(savePwErr.Status(), savePwErr)
		return
	}
	ctx.JSON(http.StatusCreated, GenerateResponse("password set successfully"))

}

//ParentSignUp with phone number
func (server *Server) ParentResendCode(ctx *gin.Context) {
	var req parentSignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		restErr := util.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	//fetch parent here
	parent, err := server.store.GetParent(ctx, req.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			getErr := util.NewRestError("parent with provided phone doesnot exists", http.StatusOK, "error when trying to get parent", nil)
			ctx.JSON(getErr.Status(), getErr)
			return
		}
		getErr := util.NewInternalServerError("error when trying to get parent", errors.New("database error"))
		ctx.JSON(getErr.Status(), getErr)
		return
	}
	//check if parent is verified
	if !parent.Isverified.Bool {
		err := util.NewUnauthorizedError("parent not verified")
		ctx.JSON(err.Status(), err)
		return
	}
	//generate otp
	OTP = util.GenerateOTP(6)
	//converting otp string to type sql.nullstring
	otpNS := util.CreateNullString(true, OTP)
	//save otp arguments
	saveOtpArg := db.SaveOTPParams{
		Otp:   otpNS,
		Phone: req.Phone,
	}
	//save otp into db
	saveErr := server.store.SaveOTP(ctx, saveOtpArg)
	if saveErr != nil {
		saveErr := util.NewInternalServerError("error when trying to save otp", errors.New("database error"))
		ctx.JSON(saveErr.Status(), saveErr)
		return
	}
	//send back otp through sms
	sendSMSErr := server.SendOTPSMS(req.Phone, OTP)
	if sendSMSErr != nil {
		sendSMSErr := util.NewInternalServerError("error when trying to send otp", errors.New("send otp error"))
		ctx.JSON(sendSMSErr.Status(), sendSMSErr)
		return
	}
	//sending back respone
	ctx.JSON(http.StatusCreated, GenerateResponse("Otp is sent to provided mobile number"))
}
