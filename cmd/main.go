package main

import (
	"fmt"
	"task-engine/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
	}

	fmt.Println("Database URL:", cfg.GetDatabaseURL())
	fmt.Println("Redis URL:", cfg.GetRedisURL())
}
