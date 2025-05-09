package api

import (
	"log"
	"math/rand/v2"
	"testing"
)

func TestRickyAndMorty_GetAllCharacters(t *testing.T) {
	r, err := NewRickAndMortyAPI()
	if err != nil {
		t.Fatalf("failed to create Rick and Morty API client: %v", err)
	}

	data, err := r.GetRandomAvatar()
	if err != nil {
		t.Errorf("failed to get random avatar: %v", err)
	}
	log.Println(data.ImageURL, data.Username)
}

func TestRickyAndMorty_GetCharacterByID(t *testing.T) {
	r, err := NewRickAndMortyAPI()
	if r == nil {
		t.Fatalf("failed to create Rick and Morty API client: %v", err)
	}
	data, err := r.GetRandomAvatarByID(rand.IntN(820))
	if err != nil {
		t.Errorf("failed to get random avatar by ID: %v", err)
	}
	log.Println(data.ImageURL, data.Username)
}

func TestFetchTotalCharachter(t *testing.T) {
	r, err := NewRickAndMortyAPI()
	if err != nil {
		t.Fatalf("failed to create Rick and Morty API client: %v", err)
	}

	data := r.fetchTotalCharacters()

	log.Println(data)
}
