package main

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/image", saveImage)

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}

func saveImage(c *gin.Context) {
	body := c.Request.Body
	

	data, err := ioutil.ReadAll(body)

	path := "/Users/honglinma/Workspace/FaceRecogination/backend/images"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

    uuid := strings.Replace(uuid.New().String(), "-", "", -1)

	fileName := path + "/" + "image_" + uuid + ".png"
	
	bytes, _ := base64.StdEncoding.DecodeString(string(data))
	err = os.WriteFile(fileName, bytes, os.ModePerm)

	if err != nil {
		panic(err)
	}


	c.JSON(200, gin.H{
		"message": string(data),
	})
}