package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"path/filepath"

	"github.com/SalvadorScafati/SovelFilter/filters"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/filter", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("img")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
			return
		}
		defer file.Close()

		fileExtension := filepath.Ext(header.Filename)[1:]
		if fileExtension != "png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format. Only PNG images are supported."})
			return
		}

		inputImg, _, err := image.Decode(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
			return
		}

		resultImg := filters.ApplySobelFilter(inputImg)

		var buf bytes.Buffer
		err = png.Encode(&buf, resultImg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
			return
		}

		c.Data(http.StatusOK, "image/png", buf.Bytes())
	})

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	r.Run(fmt.Sprintf(":%d", port))
}
