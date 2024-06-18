package handlers

var INTERNAL_SERVER_ERROR string
var FORBIDDEN string
var BAD_REQUEST string
var UNAUTHORIZED string
var PAYMENT_REQUIRED string
var OK string
var NOT_FOUND string

func init() {
	// Used in Virtual Accounts
	INTERNAL_SERVER_ERROR = "500 Internal Server Error"
	FORBIDDEN = "403 Forbidden"
	BAD_REQUEST = "400 Bad Request"
	UNAUTHORIZED = "401 Unauthorised"
	OK = "200 OK"
	PAYMENT_REQUIRED = "402 Payment Required"
	NOT_FOUND = "404 Not Found"
}
