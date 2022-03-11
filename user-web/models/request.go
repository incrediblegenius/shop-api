package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	Id          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
