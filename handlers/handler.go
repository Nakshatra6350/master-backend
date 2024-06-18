package handlers

import "log"

type Utility struct {
	l *log.Logger
}

// NewDocuments creates a product handler withthe given logger
func NewUtility(l *log.Logger) *Utility {
	return &Utility{l}
}

type User struct {
	l *log.Logger
}

// NewDocuments creates a product handler withthe given logger
func NewUser(l *log.Logger) *User {
	return &User{l}
}
