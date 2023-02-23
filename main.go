package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/yujinlim/wechat-article-scraper/controllers"
	"github.com/yujinlim/wechat-article-scraper/models"
)

func main() {
	// Connect to the database
	models.ConnectDB()
	defer models.CloseDB()

	// Set up the web server
	router := gin.Default()

	// Define the routes
	router.GET("/public-accounts", controllers.GetPublicAccounts)
	router.POST("/public-accounts", controllers.CreatePublicAccount)
	router.GET("/articles", controllers.GetArticles)
	router.POST("/articles", controllers.CreateArticles)

	// Set up the cron job to scrape the articles
	c := cron.New()
	c.AddFunc("@every 1m", controllers.ScrapeArticles)
	c.Start()

	// Start the web server
	log.Fatal(router.Run(":8080"))
}
