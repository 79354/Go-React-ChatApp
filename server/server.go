package main

import (
	"fmt"
	"log"
	"os"

	"chat-app/config"
	"chat-app/handlers"
	"chat-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
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
	lobby := handlers.NewLobby()
	go lobby.Run()

	router.GET("/", handlers.RenderHome())

	router.GET("/isUsernameAvailable/:username", handlers.IsUsernameAvailable())

	router.POST("/login", handlers.Login())
	router.POST("/registration", handlers.Registration())

	router.GET("/UserSessionCheck/:userID", handlers.UserSessionCheck())
	router.GET("/getConversation/:toUserID/:fromUserID", handlers.GetMessagesHandler())

}