package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/justty/golang-units/internal/model"
)

type ApodMetadataStore struct {
	db *sqlx.DB
}

func NewApodMetadataStore(dbConnection *sqlx.DB) ApodMetadataStore {
	return ApodMetadataStore{
		db: dbConnection,
	}
}

func (a ApodMetadataStore) GetAll() ([]model.ApodMetadata, error) {
	var arrApodMetadata []model.ApodMetadata
	err := a.db.Select(&arrApodMetadata, "SELECT * FROM apod_metadata")
	if err != nil {
		return nil, err
	}

	return arrApodMetadata, nil
}

func (a ApodMetadataStore) GetByID(id int) (*model.ApodMetadata, error) {
	var apodMetadata model.ApodMetadata
	err := a.db.Get(&apodMetadata, "SELECT * FROM apod_metadata WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &apodMetadata, nil
}

func (a ApodMetadataStore) Create(apodMetadata *model.ApodMetadata) error {
	sqlStatement := `
	INSERT INTO apod_metadata
	(
		explanation,
		media_type,
		service_version,
		title,
		hdurl,
		url,
		image_path,
		date
	)
	VALUES 
	(
		:explanation,
		:media_type,
		:service_version,
		:title,
		:hdurl,
		:url,
		:image_path,
		:date
	)
	`

	if _, err := a.db.NamedExec(sqlStatement, apodMetadata); err != nil {
		return err
	}

	return nil
}

func (a ApodMetadataStore) Update(apodMetadata *model.ApodMetadata) error {
	sqlStatement := `
	UPDATE apod_metadata
	SET 
		explanation=:explanation,
		media_type=:media_type,
		service_version=:service_version,
		title=:title,
		hdurl=:hdurl,
		url=:url,
		image_path=:image_path,
		date=:date
	WHERE id=:id
	`

	if _, err := a.db.NamedExec(sqlStatement, apodMetadata); err != nil {
		return err
	}

	return nil
}
