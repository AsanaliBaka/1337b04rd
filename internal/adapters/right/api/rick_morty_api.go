package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"1337b04rd/internal/domain"
)

const (
	baseURL      = "https://rickandmortyapi.com/api"
	maxCharacter = 826 // Максимальное количество персонажей в API
)

type charactersResponse struct {
	Info struct {
		Count int `json:"count"`
	} `json:"info"`
	Results []CharacterResponse `json:"results"`
}

type CharacterResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type RickAndMortyAPI struct {
	client          *http.Client
	totalCharacters int
}

func NewRickAndMortyAPI() (*RickAndMortyAPI, error) {
	api := &RickAndMortyAPI{
		client: &http.Client{Timeout: 10 * time.Second},
	}

	// При инициализации узнаём общее количество персонажей
	if err := api.fetchTotalCharacters(); err != nil {
		return nil, fmt.Errorf("failed to init RickAndMortyAPI: %w", err)
	}

	return api, nil
}

func (c *RickAndMortyAPI) fetchTotalCharacters() error {
	url := fmt.Sprintf("%s/character", baseURL)
	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("fetch characters count error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RickMorty API status: %d", resp.StatusCode)
	}

	var data charactersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("decode characters count error: %w", err)
	}

	c.totalCharacters = data.Info.Count
	if c.totalCharacters == 0 {
		return errors.New("no characters available")
	}

	return nil
}

func (c *RickAndMortyAPI) GetRandomUser() (*domain.User, error) {
	if c.totalCharacters == 0 {
		return nil, errors.New("no characters available")
	}

	// Генерируем случайный ID в пределах доступных персонажей
	randomID := rand.Intn(c.totalCharacters) + 1
	return c.GetUserByID(randomID)
}

func (c *RickAndMortyAPI) GetUserByID(id int) (*domain.User, error) {
	url := fmt.Sprintf("%s/character/%d", baseURL, id)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("RickMorty GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RickMorty API status: %d", resp.StatusCode)
	}

	var data CharacterResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if data.Image == "" || data.Name == "" {
		return nil, errors.New("character missing name or image")
	}

	return &domain.User{
		Name:      data.Name,
		AvatarURL: data.Image,
	}, nil
}
