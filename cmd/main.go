package main

import (
	"fmt"
	"task-engine/config"
	"task-engine/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
	}

	fmt.Println("Database URL:", cfg.GetDatabaseURL())
	fmt.Println("Redis URL:", cfg.GetRedisURL())

	r := gin.New()
	r.Use(logger.GinLogger())
}
