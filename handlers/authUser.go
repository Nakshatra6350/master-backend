package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"user-service/auth"
	"user-service/data"
	"user-service/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *User) SignupHandler(c *gin.Context) {
	u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), "Inside Create user credentials")
	var userDetails data.UserCredentials
	if err := c.ShouldBindJSON(&userDetails); err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, err.Error())
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), BAD_REQUEST)
		return
	}
	requestBodyErr := userDetails.Validate()
	if requestBodyErr != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "error in validating request body")
		util.ErrorResponse(c, http.StatusBadRequest, "error in validating request body", BAD_REQUEST)
		return
	}
	var mailid string
	emailParts := strings.Split(userDetails.Email, "@")
	if len(emailParts) == 1 {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "Please include @ in the email id")
		util.ErrorResponse(c, http.StatusBadRequest, "Please include @ in the email id", BAD_REQUEST)
		return
	}
	mailid = strings.ToLower(emailParts[0]) + "@" + strings.ToLower(emailParts[1])
	userDetails.Email = mailid
	filter := bson.M{"$or": []bson.M{
		{"email": userDetails.Email},
		{"username": userDetails.Username},
	}}
	var existingUser data.UserCredentials
	err := data.GetMongoDB().FindOne(c, filter).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "Email or username already exists")
		util.ErrorResponse(c, http.StatusBadRequest, "Email or username already exists", BAD_REQUEST)
		return
	}
	password := userDetails.Password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "Error in hashing password")
		util.ErrorResponse(c, http.StatusBadRequest, "Error in hashing password", BAD_REQUEST)
		return
	}
	userDetails.Password = hashedPassword
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusInternalServerError, "Error in loading IST location")
		util.ErrorResponse(c, http.StatusInternalServerError, "Error in loading IST location", INTERNAL_SERVER_ERROR)
		return
	}
	userDetails.CreatedAt = time.Now().In(istLocation)

	_, err = data.GetMongoDB().InsertOne(c, userDetails)
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), "Error in creating credentials")
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), BAD_REQUEST)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User credentials created successfully",
	})
}

func (u *User) LoginHandler(c *gin.Context) {
	u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), "Inside Login user")
	var userDetails data.LoginUser
	if err := c.ShouldBindJSON(&userDetails); err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, err.Error())
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), BAD_REQUEST)
		return
	}

	requestBodyErr := userDetails.Validate()
	if requestBodyErr != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "error in validating request body")
		util.ErrorResponse(c, http.StatusBadRequest, "error in validating request body", BAD_REQUEST)
		return
	}

	filter := bson.M{"email": userDetails.Email}
	var user data.UserCredentials
	err := data.GetMongoDB().FindOne(c, filter).Decode(&user)
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "Invalid Credentials, Email does not exists")
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid Credentials, Email does not exists", BAD_REQUEST)
		return
	}
	if !auth.VerifyPassword(userDetails.Password, user.Password) {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusBadRequest, "Invalid Credentials, Password not verify")
		util.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", UNAUTHORIZED)
		return
	}

	tokenResponse, err := auth.GetJWT(user.Username, userDetails.Password)
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusInternalServerError, "Could not generate token")
		util.ErrorResponse(c, http.StatusInternalServerError, "Could not generate token", INTERNAL_SERVER_ERROR)
		return
	}
	jwt, exists := tokenResponse["access_token"]
	if !exists {
		fmt.Printf("Key not found: %s\n", jwt)
		util.ErrorResponse(c, http.StatusInternalServerError, "Could not get token", INTERNAL_SERVER_ERROR)
	}

	var clientDetails = make(map[string]interface{})
	clientDetails["username"] = user.Username
	clientDetails["email"] = user.Email
	clientDetails["jwt"] = jwt

	clientDetailsByte, err := json.Marshal(clientDetails)
	if err != nil {
		u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), http.StatusInternalServerError, "Error in marshalling login response body")
		util.ErrorResponse(c, http.StatusInternalServerError, "Error in marshalling login response body", INTERNAL_SERVER_ERROR)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("jwt", jwt, 4*60*60, "/", "localhost", false, false)
	util.SuccessResponse(c, http.StatusOK, clientDetailsByte, OK)
}

func (u *User) LogoutHandler(c *gin.Context) {
	u.l.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), "logoutHandler")
	c.SetCookie("jwt", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully logged out and Cookie cleared successfully"})
}
