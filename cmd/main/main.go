package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

func smoothEffect() int {
	time.Sleep(10 * time.Second)
	return rand.Intn(2)
}

func SendImage(id string, url string) {
	result := smoothEffect()
	fmt.Println(result)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting image:", err)
		return
	}
	defer resp.Body.Close()

	// Прочитаем полученное изображение
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading image data:", err)
		return
	}

	// Создаем FormData и добавляем изображение
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", ""+id+".png")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	part.Write(imgData)
	writer.Close()

	// Отправляем FormData на сервер Django
	resp2, err := http.NewRequest("PUT", "http://localhost:8000/api/orders/"+id+"/image/update/", body)
	if err != nil {
		fmt.Println("Error create PUT request:", err)
		return
	}
	resp2.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp3, err := client.Do(resp2)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer resp3.Body.Close()

	fmt.Println("Image updated in Orders model")
}

func main() {
	router := gin.Default()
	router.POST("/edit_image/", func(c *gin.Context) {
		id := c.PostForm("id")
		go SendImage(id, "http://localhost:8000/api/orders/"+id+"/image/")
		c.JSON(http.StatusOK, gin.H{"message": "Image retrieval and update initiated"})
	})
	router.Run(":8080")
}
