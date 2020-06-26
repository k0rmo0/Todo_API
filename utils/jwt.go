package utils

import "github.com/dgrijalva/jwt-go"

//Claims ...
type Claims struct {
	Username    string `db:"username" json:"username"`
	Tokenstring string `db:"token_string" json:"token_string"`
	jwt.StandardClaims
}
