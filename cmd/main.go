package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang-kafka/internal/handler"
	"golang-kafka/internal/service"
	"net/http"
	"time"
)

// @title           Kafka producer
// @version         1.0
// @description     fio sender

// @host      localhost:8080
// @BasePath  /
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Info("Starting kafka producer")

	// Read config
	if err := initConfig(); err != nil {
		logrus.Fatalf("Failed to read config: %s", err.Error())
	}

	// Configure logger
	if err := configureLogger(viper.GetString("logger.log-level")); err != nil {
		logrus.Fatalf("Failed to configure logger: %s", err.Error())
	}

	// kafka service
	produceService, err := service.NewKafkaProduceService(
		viper.GetString("kafka.topic"),
		viper.GetString("kafka.config-path"))
	if err != nil {
		logrus.Fatalf("Failed create kafka produce service: %s", err.Error())
	}
	defer produceService.Shutdown()

	// handler + service
	ser := service.NewService(produceService)
	han := handler.NewHandler(ser)

	// Routes
	routes := han.InitRoutes()

	// Server
	server := createServer(viper.GetString("server.port"), routes)
	logrus.Infof("Server running on http://localhost%s", server.Addr)
	logrus.Infof("Swagger: http://localhost%s/swagger/index.html", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start server: %s", err.Error())
	}
}

// Configuration
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func configureLogger(logLevel string) error {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func createServer(port string, routes *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           routes,
		ReadHeaderTimeout: 2 << 20,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
}
