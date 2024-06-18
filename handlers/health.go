package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"user-service/data"
	"user-service/util"

	"github.com/gin-gonic/gin"
)

func (h *Utility) HealthCheckUp(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"status":  "up",
			"message": "Server running perfectly",
		})
}
func (h *Utility) SendOTP(c *gin.Context) {
	requestBodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Error reading request body", BAD_REQUEST)
		return
	}
	var requestBody map[string]interface{}
	err = json.Unmarshal(requestBodyByte, &requestBody)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Error unmarshalling request body for sending otp on email", BAD_REQUEST)
		return
	}
	email := requestBody["email"].(string)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email is required"})
		return
	}

	otp := data.GenerateOTP(email)
	err = data.SendOTPEmail(otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to send OTP"})
		return
	}

	data.SaveOTP(otp)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP sent successfully"})
}

func (h *Utility) VerifyOTP(c *gin.Context) {
	requestBodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Error reading request body", BAD_REQUEST)
		return
	}
	var requestBody map[string]interface{}
	err = json.Unmarshal(requestBodyByte, &requestBody)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Error unmarshalling request body for sending otp on email", BAD_REQUEST)
		return
	}
	email := requestBody["email"].(string)
	if email == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "EmailId is required", BAD_REQUEST)
		return
	}
	code := requestBody["code"].(string)
	if code == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "Otp is required", BAD_REQUEST)
		return
	}

	otp, exists := data.GetOTP(email)
	if !exists || otp.Code != code || otp.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid or expired OTP"})
		return
	}

	data.DeleteOTP(email)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP verified successfully"})
}
