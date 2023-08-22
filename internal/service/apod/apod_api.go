package apod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/justty/golang-units/internal/model"
)

type ApodAPI struct {

}

func NewApodAPI() ApodAPI {
	return ApodAPI {}
}

func (a ApodAPI) GetApod() (*model.ApodMetadata, error) {
	apiKey := os.Getenv("APOD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("APOD_API_KEY is not set in the .env file")
	}

	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error http get request: %w", err)
	}
	defer resp.Body.Close()

	var metadata model.ApodMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("error decode response: %w", err)
	}

	return &metadata, nil
}