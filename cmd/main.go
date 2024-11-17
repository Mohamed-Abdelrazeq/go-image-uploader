package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", health)

	r.POST("/upload", upload)

	r.Run(":8080")
}

func upload(c *gin.Context) {
	// parse image
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	// save image
	err = c.SaveUploadedFile(file, "./images/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save image",
		})
		return
	}

	// return response
	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"url":     "http://localhost:8080/images/" + file.Filename,
	})
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
