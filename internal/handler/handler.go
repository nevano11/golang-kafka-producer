package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "golang-kafka/docs"
	"golang-kafka/internal/entity"
	"golang-kafka/internal/service"
	"io"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	engine := gin.New()

	engine.GET("/", h.welcome)
	engine.POST("/send-fio", h.produceFio)
	engine.POST("/send-fio/list", h.produceFioList)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

// @Summary      	produce
// @Description  	method to send fio
// @Accept       	json
// @Produce      	json
// @Param 			full_name body entity.Fio true "The input todo struct"
// @Success         200 {string} string "ok"
// @Router       	/send-fio [post]
func (h *Handler) produceFio(context *gin.Context) {
	logrus.Info("Handle produceFio")

	var inputFio entity.Fio

	if err := context.BindJSON(&inputFio); err != nil {
		logrus.Errorf("Failed to bind Json: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}
	err := h.service.Produce(inputFio)
	if err != nil {
		logrus.Errorf("Failed to produce fio: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to send")
	}
	context.JSON(http.StatusOK, inputFio)
}

// @Summary      	produce
// @Description  	method to send fio
// @Accept       	json
// @Produce      	json
// @Param 			full_name body []entity.Fio true "The input todo struct"
// @Success         200 {string} string "ok"
// @Router       	/send-fio/list [post]
func (h *Handler) produceFioList(context *gin.Context) {
	logrus.Info("Handle produceFioList")

	var inputFio = make([]entity.Fio, 10)

	if err := context.BindJSON(&inputFio); err != nil {
		logrus.Errorf("Failed to bind Json: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}
	for _, v := range inputFio {
		err := h.service.Produce(v)
		if err != nil {
			logrus.Errorf("Failed to produce fio: %s", err.Error())
			context.JSON(http.StatusBadRequest, "Failed to send")
			return
		}
	}
	context.JSON(http.StatusOK, inputFio)
}

func (h *Handler) welcome(context *gin.Context) {
	logrus.Info("Handle welcome")
	_, _ = io.WriteString(context.Writer, "Welcome")
}
