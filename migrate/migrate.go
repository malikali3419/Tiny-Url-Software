package main

import (
	"example.com/m/initializers"
	"example.com/m/models"
)

func init() {
	initializers.LoadEnviromentalVariables()
	initializers.ConnectToDatabse()
}
func main() {
	initializers.DB.AutoMigrate(&models.URlMapping{}, &models.User{}, &models.JWTClaims{})
}
