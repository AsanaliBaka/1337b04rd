package main

import (
	"fmt"
	"log"
	"math/rand"

	"1337b04rd/internal/infrastructure/api"
	"1337b04rd/pkg/logger"
)

func main() {
	logger, err := logger.NewCustomLogger()
	if err != nil {
		log.Fatalln("Failed to initialize logger:", err)
	}
	logger.Info("Logger initialized")

	rickmortyApi := api.NewRickAndMortyAPI(logger)
	fmt.Println(rickmortyApi)
	image, name, err := rickmortyApi.GetRandomCharacter(rand.Intn(800))
	if err != nil {
		logger.Error("Failed to get random character:", err)
		return
	}
	fmt.Println(image)
	fmt.Println(name)
}
