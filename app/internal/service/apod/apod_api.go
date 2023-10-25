package apod

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/justty/golang-units/app/internal/model"
)

type ApodAPI struct {
	apiKey string
}

func NewApodAPI(apiKey string) *ApodAPI {
	return &ApodAPI {
		apiKey: apiKey,
	}
}

func (a *ApodAPI) GetApod() (*model.ApodMetadata, error) {
	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s&thumbs=true", a.apiKey)

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