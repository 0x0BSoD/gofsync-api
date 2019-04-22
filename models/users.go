package models

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Password   string `json:"password"`
	Username   string `json:"username"`
	RememberMe bool   `json:"remember_me,omitempty"`
}
type Claims struct {
	Username   string `json:"username"`
	RememberMe bool   `json:"remember_me,omitempty"`
	jwt.StandardClaims
}
type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
