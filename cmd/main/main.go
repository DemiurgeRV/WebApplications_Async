package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendStatus(pk string, url string) {
	data := []byte(`{"name": "John", "age": 30}`)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status sent successfully for pk:", pk)
}

func main() {
	router := gin.Default()
	router.POST("/edit_image/", func(c *gin.Context) {
		id := c.PostForm("id")
		go SendStatus(id, "http://localhost:8000/api/edit_image/"+id+"/")
		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":8080")
}
