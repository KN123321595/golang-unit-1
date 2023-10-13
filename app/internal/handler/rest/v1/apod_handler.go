package rest_v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/justty/golang-units/app/internal/model"
)

type ApodMetadataStore interface {
	GetAll() ([]model.ApodMetadata, error)
	GetByDate(date string) (*model.ApodMetadata, error)
}

type ApodHandler struct {
	apodMetadataStore ApodMetadataStore
}

func NewApodHandler(apodMetadataStore ApodMetadataStore) *ApodHandler {
	return &ApodHandler{
		apodMetadataStore: apodMetadataStore,
	}
}

func (a *ApodHandler) GetApodMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	date := r.URL.Query().Get("date")
	if date == "" {
		arrApodMetadata, err := a.apodMetadataStore.GetAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("error get all apod metadata: %s", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&arrApodMetadata)
	} else {
		apodMetadata, err := a.apodMetadataStore.GetByDate(date)
		if err != nil {
			http.Error(w, fmt.Sprintf("error get apod metadata by date: %s", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&apodMetadata)
	}
}
