package main

import (
	v1 "fitbyte/api/v1"
	"fitbyte/pkg/database"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-2"),
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})

	if err != nil {
		log.Fatal("Error creating AWS session: ", err)
	}

	s3Client := s3.New(sess)

	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Println("Error listing S3 buckets:", err)
		return
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf(" - %s\n", *bucket.Name)
	}

	router := gin.Default()

	v1Group := router.Group("/v1")
	{
		v1.RegisterUserRouter(v1Group)
		//v1.RegisterActivityRoutes(v1Group, handlerInstance)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server started on http://localhost:%s", port)
	err = router.Run(":" + port)
	if err != nil {
		return
	}
}
