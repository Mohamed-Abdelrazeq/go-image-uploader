package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", health)

	r.POST("/images", uploadImage)
	r.GET("/images/:filename", showImage)

	r.Run(":8080")
}

func uploadImage(c *gin.Context) {
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

func showImage(c *gin.Context) {
	filename := c.Param("filename")
	c.File("./images/" + filename)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
