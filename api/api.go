package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"protohush"
)

type Api struct {
	Chat protohush.Chat
}

func NewApi(chat protohush.Chat) Api {
	return Api{Chat: chat}
}

type SearchRequest struct {
	Message string `json:"message"`
}

func (api Api) Run() {
	r := gin.Default()

	r.POST("protohush/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		chatResponse, err := api.Chat.Search(req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, chatResponse)
	})

	fmt.Println("Server started at :8080")
	r.Run(":8080")
}
