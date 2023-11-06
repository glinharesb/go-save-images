package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

type ImageInfo struct {
	FolderPath string `json:"folderPath"`
	ImageURL   string `json:"imageURL"`
	ImageName  string `json:"imageName,omitempty"`
}

const maxGoroutines = 25

var semaphore = make(chan struct{}, maxGoroutines)

func main() {
	err := os.MkdirAll("dist", 0755)
	if err != nil {
		log.Fatalf("Failed to create 'dist' directory: %v", err)
	}

	router := gin.Default()
	router.POST("/save-images", saveImagesHandler)
	router.Run(":8080")
}

func saveImagesHandler(c *gin.Context) {
	var images []ImageInfo
	if err := c.BindJSON(&images); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	for _, img := range images {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(img ImageInfo) {
			defer wg.Done()
			defer func() { <-semaphore }()

			fullFolderPath := filepath.Join("dist", img.FolderPath)
			if err := os.MkdirAll(fullFolderPath, os.ModePerm); err != nil {
				log.Println(err)
				return
			}

			resp, err := http.Get(img.ImageURL)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			fileName := img.ImageName
			if fileName == "" {
				fileName = filepath.Base(img.ImageURL)
			}
			filePath := filepath.Join(fullFolderPath, fileName)
			file, err := os.Create(filePath)
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()

			if _, err = io.Copy(file, resp.Body); err != nil {
				log.Println(err)
				return
			}
		}(img)
	}

	wg.Wait()
	c.Status(http.StatusOK)
}
