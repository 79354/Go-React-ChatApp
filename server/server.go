package main

import (
	"fmt"
	"log"
	"os"

	"chat-app/config"
	"chat-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/gotdotenv"
)

func main(){
	err := gotdotenv.Load()
	if err != nil{
		log.Fatal("Error loading the environment")
	}

	fmt.Println(
		fmt.Sprintf("%s%s%s%s", "Server will start at http://", os.Getenv("HOST"), ":", os.Getenv("PORT")),
	)

	config.ConnectDatabase()

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(utils.CORSMiddleware())

	routes(router)

	router.Run(":" + os.Getenv("PORT"))
}

func routes(router *gin.Engine) {
	
}