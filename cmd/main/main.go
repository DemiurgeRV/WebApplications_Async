package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

const key = "ndscL3Jwp9kMNjknk12"

func timeSleep() {
	time.Sleep(10 * time.Second)
}

func SendImage(id string, url string) {
	timeSleep()
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

	// Декодируем изображение
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Применяем эффект блюра к изображению
	img = imaging.Blur(img, 3)

	// Кодирование обработанного изображения обратно в формат PNG
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, img, imaging.PNG)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	// Создаем FormData и добавляем изображение
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	postKey, err := writer.CreateFormField("key")
	if err != nil {
		fmt.Println("Error creating form field:", err)
		return
	}
	postKey.Write([]byte(key))

	part, err := writer.CreateFormFile("image", ""+id+"_blurred.png")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	part.Write(buf.Bytes()) // Записываем обработанное изображение
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
		getKey := c.PostForm("key")
		if getKey != key {
			fmt.Println("Error")
			return
		}
		id := c.PostForm("id")
		go SendImage(id, "http://localhost:8000/api/orders/"+id+"/image/")
		c.JSON(http.StatusOK, gin.H{"message": "Image retrieval and update initiated"})
	})
	router.Run(":8080")
}
