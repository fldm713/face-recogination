package main

import (
	"backend/customvision"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func main() {

	tagMap := map[string]string{
		"Honglin": "e1ac5ab5-110e-4289-9cce-d3c0928e0438",
		"Others": "a44955a5-8dd3-41fe-9c5a-535f15122c85",
	}

	customvision.Init("3179ff3c-1e72-46d5-9905-c1ad9844ec14", tagMap)
	// projectId := customvision.CreateProject("Face Recognition", "Face Recognition")
	// fmt.Printf("Project id: %v\n", projectId) 

	// honglinTag, _ := customvision.CreateTag(projectId, "Honglin", "Honglin Tag")
	// fmt.Printf("Hemlock tag id: %v\n", honglinTag.ID)
	// othersTag, _ := customvision.CreateTag(projectId, "Others", "others Tag")
	// fmt.Printf("Hemlock tag id: %v\n", othersTag.ID)


	// customvision.UpLoadImages(ctx, trainer, projectId, hemlockTag, cherryTag)
	// customvision.Train(ctx, trainer, projectId)
	// projectId, _ := uuid.FromString("d5421ea1-27fb-4df1-a8f1-41e8d24c95b2")
	// fmt.Printf("project id: %v\n", projectId)
	// customvision.Predict(ctx, projectId)

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	port := "80"
	val, ok := os.LookupEnv("GO_HTTP_PORT")
	if ok {
		port = val
	}
	
	r.POST("/image", saveImage)

	r.POST("/train", train)

	r.POST("/validate", validate)

	r.Run("0.0.0.0:"+port) // listen and serve on 0.0.0.0:8080

	
}

func saveImage(c *gin.Context) {
	body := c.Request.Body
	

	data, err := ioutil.ReadAll(body)

	if err != nil {
		panic(err)
	}

	path := "/images"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	uid, _ := uuid.NewV1()
    uuidString := strings.Replace(uid.String(), "-", "", -1)

	fileName := path + "/" + "image_" + uuidString + ".png"
	
	bytes, _ := base64.StdEncoding.DecodeString(string(data))
	err = os.WriteFile(fileName, bytes, os.ModePerm)

	if err != nil {
		panic(err)
	}

	err = customvision.UpLoadImage(bytes, "Honglin")
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"message": "Upload Success!",
	})
}

func train(c *gin.Context) {
	err := customvision.Train()
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Training Success!",
		})
	} else {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}	
}

func validate(c *gin.Context) {
	body := c.Request.Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}

	path := "/images"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	uid, _ := uuid.NewV1()
    uuidString := strings.Replace(uid.String(), "-", "", -1)

	fileName := path + "/" + "image_" + uuidString + ".png"
	
	bytes, _ := base64.StdEncoding.DecodeString(string(data))
	err = os.WriteFile(fileName, bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}


	result, err := customvision.Predict(fileName)

	fmt.Printf("%v\n", err != nil)
	fmt.Printf("%v\n", result)

	if err == nil {
		c.JSON(200, result)
		// c.JSON(200, gin.H{
		// 	"message": "validation",
		// })
	} else {
		fmt.Println("111")
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}	
}
