/*
Package router : authorize and authenticate HTTP Request using HTTP Header.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package router

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

/*
Handler : a struct to handle HTTP Request and check its Header.
	Handler encloses github.com/gin-gonic/gin.Engine.
	Handler authorizes and authenticates all HTTP Requests using its HTTP Header.
*/
type Handler struct {
	Engine *gin.Engine
}

/*
NewHandler : a factory method to create Handler.
*/
func NewHandler() *Handler {
	engine := gin.Default()
	c := &checker{
		payloads: make(map[string]string),
	}

	engine.POST("/distinct/", func(context *gin.Context) {
		distinctMessage(context, c)
	})

	router := &Handler{
		Engine: engine,
	}
	return router
}

/*
Run : start listening HTTP Request using enclosed gin.Engine.
*/
func (router *Handler) Run(port string) {
	router.Engine.Run(port)
}

type bodyType struct {
	Payload string `json:"payload" binding:"required"`
}

func distinctMessage(context *gin.Context, checker *checker) {
	var body bodyType

	if err := context.ShouldBindWith(&body, binding.JSON); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"result": "failure",
			"error":  err.Error(),
		})
		return
	}
	if checker.IsDuplicate(body.Payload) {
		log.Printf(formatter("distinctMessage", fmt.Sprintf("duplicate payload = %s", body.Payload)))
		context.JSON(http.StatusConflict, gin.H{
			"result":  "duplicate",
			"payload": body.Payload,
		})
	} else {
		log.Printf(formatter("distinctMessage", fmt.Sprintf("new payload = %s", body.Payload)))
		context.JSON(http.StatusOK, gin.H{
			"result":  "success",
			"payload": body.Payload,
		})
	}
}

type checker struct {
	payloads map[string]string
}

func (c *checker) IsDuplicate(payload string) bool {
	_, exists := c.payloads[payload]
	c.payloads[payload] = payload
	return exists
}

func formatter(group string, msg string) string {
	return fmt.Sprintf("[APP] %s | [%s] %s", time.Now().Format("2006/01/02 - 15:04:05"), group, msg)
}
