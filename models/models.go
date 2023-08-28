package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type URlMapping struct {
	gorm.Model
	LongUrl   string
	ShortCode string
	ShortUrl  string
	UserID    uint
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}
