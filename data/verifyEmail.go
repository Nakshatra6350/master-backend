package data

import (
	"math/rand"
	"net/smtp"
	"time"
)

type OTP struct {
	Code      string
	Email     string
	ExpiresAt time.Time
}

var OTPStore = map[string]OTP{}

func SaveOTP(otp OTP) {
	OTPStore[otp.Email] = otp
}

func GetOTP(email string) (OTP, bool) {
	otp, exists := OTPStore[email]
	return otp, exists
}

func DeleteOTP(email string) {
	delete(OTPStore, email)
}

func GenerateOTP(email string) OTP {
	const otpLength = 6
	const otpCharset = "0123456789"
	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = otpCharset[rand.Intn(len(otpCharset))]
	}
	return OTP{
		Code:      string(otp),
		Email:     email,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP valid for 10 minutes
	}
}

func SendOTPEmail(otp OTP) error {
	from := "nakshatragarg678@gmail.com"
	password := "ipas hjzk ezex ggek"
	to := otp.Email
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OTP Verification\n\n" +
		"Your OTP is: " + otp.Code

	auth := smtp.PlainAuth("", from, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	return err
}
