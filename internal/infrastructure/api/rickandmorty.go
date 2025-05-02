package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"1337b04rd/pkg/logger"
)

type CharacterResponse struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type RickAndMortyAPI struct {
	client *http.Client
}

func NewRickAndMortyAPI(l *logger.CustomLogger) *RickAndMortyAPI {
	return &RickAndMortyAPI{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

var baseURL = "https://rickandmortyapi.com/api"

func (c *RickAndMortyAPI) GetRandomCharacter(id int) (avatarURL string, name string, err error) {
	url := fmt.Sprintf("%s/character/%d", baseURL, id)
	resp, err := c.client.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("RickMorty GET error: %w", err)
	}
	defer resp.Body.Close()

	// на случай если рик и морти не робит или инет дерьмо
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("RickMorty API status: %d", resp.StatusCode)
	}

	var data CharacterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", fmt.Errorf("decode error: %w", err)
	}

	if data.Image == "" || data.Name == "" {
		return "", "", errors.New("character missing name or image")
	}

	return data.Image, data.Name, nil
}
