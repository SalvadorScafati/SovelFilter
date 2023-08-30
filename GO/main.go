package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/SalvadorScafati/SovelFilter/filters"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/filter", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("img")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
			return
		}
		defer file.Close()

		fileExtension := getFileExtension(file)
		if fileExtension != "jpg" && fileExtension != "jpeg" && fileExtension != "png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format"})
			return
		}

		inputImg, _, err := image.Decode(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
			return
		}

		resultImg := filters.ApplySobelFilter(inputImg)

		var buf bytes.Buffer
		switch fileExtension {
		case "jpg", "jpeg":
			err = jpeg.Encode(&buf, resultImg, nil)
		case "png":
			err = png.Encode(&buf, resultImg)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
			return
		}

		c.Data(http.StatusOK, "image/"+fileExtension, buf.Bytes())
	})

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	r.Run(fmt.Sprintf(":%d", port))
}

func getFileExtension(file multipart.File) string {
	// Seek back to the beginning of the file to read the header
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return ""
	}

	// Read the first 512 bytes to detect the file type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return ""
	}

	// Determine the image format from the buffer
	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}
