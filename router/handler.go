/*
Package router : routing http request and check message duplication using Checker.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/checker"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/conf"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/utils"
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
func NewHandler(config *conf.Config) (*Handler, error) {
	engine := gin.Default()
	c, err := checker.NewChecker(config)
	if err != nil {
		return nil, err
	}

	engine.POST("/distinct/", func(context *gin.Context) {
		distinctMessage(context, c)
	})

	router := &Handler{
		Engine: engine,
	}
	return router, nil
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

func distinctMessage(context *gin.Context, checker *checker.Checker) {
	logger := utils.NewLogger("distinctMessage")
	var body bodyType

	cType := context.GetHeader("Content-Type")

	if cType != "application/json" {
		logger.Errorf("header failed: Content-Type=%s", cType)
		context.JSON(http.StatusBadRequest, gin.H{
			"result": "failure",
			"error":  "Content-Type not allowd: " + cType,
		})
		return
	}

	if err := context.ShouldBindWith(&body, binding.JSON); err != nil {
		logger.Errorf("validate failed: %s", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"result": "failure",
			"error":  err.Error(),
		})
		return
	}
	isDup, err := checker.IsDuplicate(body.Payload)
	if isDup || err != nil {
		logger.Infof("duplicate payload = %s", body.Payload)
		context.JSON(http.StatusConflict, gin.H{
			"result":  "duplicate",
			"payload": body.Payload,
		})
	} else {
		logger.Infof("new payload = %s", body.Payload)
		context.JSON(http.StatusOK, gin.H{
			"result":  "success",
			"payload": body.Payload,
		})
	}
}
