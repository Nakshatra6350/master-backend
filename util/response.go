package util

import (
	"encoding/json"
	"fmt"
	"os"
	"user-service/models"

	"github.com/gin-gonic/gin"

	logs "github.com/sirupsen/logrus"
)

var L *logs.Logger

func init() {
	L = &logs.Logger{}
	L.SetFormatter(&logs.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})
	L.SetOutput(os.Stdout)
	L.SetLevel(logs.TraceLevel)
}

func SuccessResponse(c *gin.Context, responseCode int, body []byte, messages string) {
	L.Infoln(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c, c.Request.UserAgent(), messages)
	respbody := map[string]interface{}{}
	err := json.Unmarshal(body, &respbody)
	if err != nil {
		fmt.Printf("Error Occured in Unmarshalling in Success Response")
	}
	c.JSON(responseCode, respbody)
	c.Abort()
}

func ErrorResponse(c *gin.Context, errorCode int, message string, status string) {
	L.Errorln(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.UserAgent(), status)
	L.Error(message)
	response := models.NewHTTPError(nil, errorCode, message)
	respbody := map[string]interface{}{}
	err := json.Unmarshal(response, &respbody)
	if err != nil {
		fmt.Printf("Error Occured in Unmarshalling in Success Response")
	}
	c.JSON(errorCode, respbody)
	c.Abort()

}

// func SuccessResponseDecentroCallback(rw http.ResponseWriter, r *http.Request, responseCode int, messages string) {
// 	L.Infoln(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), messages)
// 	rw.Header().Add("Content-Type", "application/json")
// 	rw.WriteHeader(responseCode)
// 	rw.Write([]byte(`{"response_code": "CB_S00000"}`))
// }

// func ErrorResponseDecentroCallback(rw http.ResponseWriter, r *http.Request, errorCode int, message string, status string) {
// 	L.Errorln(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status)
// 	L.Error(message)
// 	rw.Header().Add("Content-Type", "application/json")
// 	rw.WriteHeader(errorCode)
// 	if errorCode == 401 || errorCode == 403 {
// 		rw.Write([]byte(`{"response_code": "CB_E00013"}`))
// 	} else if errorCode == 500 {
// 		rw.Write([]byte(`{"response_code": "CB_E00000"}`))
// 	} else {
// 		rw.Write([]byte(`{"response_code": "CB_E00009"}`))
// 	}
// }

// func TraceLog(c *gin.Context, message string) {
// 	L.Traceln(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, message)
// }

func DebugLog(c *gin.Context, message string) {
	L.Debugln(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, message)
}

func LocalDebugLog(err string) {
	L.Debugln(err)
}
