package apod

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/justty/golang-units/internal/model"
)

const WorkerServiceName = "apod_worker"

type ApodMetadataStore interface {
	Create(apodMetadata *model.ApodMetadata) error
}

type ApodAPIInterface interface {
	GetApod() (*model.ApodMetadata, error)
}

type ApodWorker struct {
	apodMetadataStore ApodMetadataStore
	apodAPIInterface  ApodAPIInterface
}

func NewApodWorker(apodMetadataStore ApodMetadataStore, apodAPIInterface ApodAPIInterface) ApodWorker {
	return ApodWorker{
		apodMetadataStore: apodMetadataStore,
		apodAPIInterface:  apodAPIInterface,
	}
}

func (a ApodWorker) Process() error {
	fmt.Printf("Start service %s\n", WorkerServiceName)

	fmt.Println("Get metadata by APOD API")
	apodMetadata, err := a.apodAPIInterface.GetApod()
	if err != nil {
		return fmt.Errorf("error get apod: %w", err)
	}

	fmt.Println("Save image from APOD metadata")
	imagePath, err := a.SaveImage(apodMetadata.Hdurl)
	if err != nil {
		return fmt.Errorf("error save image path: %w", err)
	}
	apodMetadata.ImagePath = imagePath

	if err := a.apodMetadataStore.Create(apodMetadata); err != nil {
		return fmt.Errorf("error create apodMetadata in db: %w", err)
	}

	fmt.Println("End service")
	return nil
}

func (a ApodWorker) SaveImage(url string) (string, error) {
	imageResponse, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching image: %w", err)
	}
	defer imageResponse.Body.Close()

	imageData, err := io.ReadAll(imageResponse.Body)
	if err != nil {
		return "", fmt.Errorf("error reading image data: %w", err)
	}

	imagePath := fmt.Sprintf("images/%s.jpg", time.Now().Format("20060102_150405"))
	if err = os.WriteFile(imagePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("error saving image: %w", err)
	}

	return imagePath, nil
}
