package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	docs "github.com/rlatmfrl24/applemint-go/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log.Print("starting server...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Println("Error loading .env file")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:  []string{"X-Requested-With", "Content-Type", "Authorization", "origin", "x-csrftoken", "x-access-token"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	// Item API
	itemGroup := router.Group("/item")
	{
		itemGroup.GET("/move/:id", moveItemRequestHandler)
		collection := itemGroup.Group("/:collection")
		{
			collection.GET("/:id", getItemRequestHandler)
			collection.PUT("/:id", updateItemReqeustHandler)
			collection.DELETE("/:id", deleteItemRequestHandler)
		}
		bookmarkGroup := itemGroup.Group("/bookmark")
		{
			bookmarkGroup.PUT("", makeBookmarkRequestHander)
			bookmarkGroup.GET("/list", getBookmarkListRequestHandler)
		}

	}

	// Collection API
	collectionGroup := router.Group("/collection")
	{
		collectionGroup.GET("/list", getCollectionListRequestHandler)
		collectionGroup.GET("/clean", cleanOldItemsRequestHandler)
		targetCollection := collectionGroup.Group("/:target")
		{
			targetCollection.GET("", getItemListRequestHandler)
			targetCollection.GET("/info", getCollectionInfoRequestHandler)
			targetCollection.DELETE("/clear", clearCollectionRequestHandler)
		}
	}

	// Crawl API
	router.GET("/crawl/:target", crawlRequestHandler)

	// Gallery API
	galleryGroup := router.Group("/gallery")
	{
		galleryGroup.GET("", getGalleryItemsRequestHandler)
		galleryGroup.GET("/imgur", analyzeImgurRequestHandler)
	}

	// External App API
	router.GET("/dropbox", dropboxRequestHandler)
	raindropGroup := router.Group("/raindrop")
	{
		raindropGroup.GET("/list", getRaindropCollectionListRequestHandler)
		raindropGroup.PUT("/:collectionId", sendToRaindropRequestHandler)
	}

	// Swagger API
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "AppleMint"
	docs.SwaggerInfo.Version = "1.0"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run()
}
