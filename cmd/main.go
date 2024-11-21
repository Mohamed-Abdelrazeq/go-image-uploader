package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	s3Session *s3.S3
	bucket    string
)

func initAWS() {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		),
	)
	s3Session = s3.New(sess)
	bucket = os.Getenv("AWS_BUCKET_NAME")
}

func main() {
	r := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initAWS()

	r.GET("/health", health)

	r.POST("/images", uploadImage)
	r.GET("/images/:filename", showImage)

	r.Run(":8080")
}
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open image",
		})
		return
	}
	defer src.Close()

	_, err = s3Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Filename),
		Body:   src,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image to S3",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"url":     "https://" + bucket + ".s3.amazonaws.com/" + file.Filename,
	})
}

func showImage(c *gin.Context) {
	filename := c.Param("filename")
	c.Redirect(http.StatusMovedPermanently, "https://"+bucket+".s3.amazonaws.com/"+filename)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
