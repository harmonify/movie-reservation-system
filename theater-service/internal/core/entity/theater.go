package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Theater struct {
	TheaterID   string         `json:"theater_id"`
	TraceID     string         `json:"trace_id"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Website     string         `json:"website"`
	Latitude    float32        `json:"latitude"`
	Longitude   float32        `json:"longitude"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (*Theater) TableName() string {
	return "theater"
}

func NewTheater(saveEntity *SaveTheater) *Theater {
	return &Theater{
		TraceID:     saveEntity.TraceID,
		Name:        saveEntity.Name,
		Address:     saveEntity.Address,
		PhoneNumber: saveEntity.PhoneNumber,
		Email:       saveEntity.Email,
		Website:     saveEntity.Website,
		Latitude:    saveEntity.Latitude,
		Longitude:   saveEntity.Longitude,
	}
}

type TheaterSortBy string

func (s TheaterSortBy) String() string {
	return string(s)
}

func (s TheaterSortBy) IsValid() bool {
	switch s {
	case TheaterSortByNewest, TheaterSortByNearest:
		return true
	}
	return false
}

const (
	TheaterSortByNewest  TheaterSortBy = "NEWEST"
	TheaterSortByNearest TheaterSortBy = "NEAREST"
)

type FindManyTheaters struct {
	Keyword  sql.NullString
	Location *FindManyTheatersLocation
	SortBy   TheaterSortBy
	Page     uint32
	PageSize uint32
}

type FindManyTheatersLocation struct {
	Latitude  float32 // user latitude
	Longitude float32 // user longitude
	Radius    float32 // in meters
}

type FindManyTheatersResult struct {
	Theaters []*Theater
	Metadata *FindManyTheatersMetadata
}

type FindManyTheatersMetadata struct {
	TotalResults int64
}

type FindOneTheater struct {
	TheaterID   sql.NullString
	TraceID     sql.NullString
	PhoneNumber sql.NullString
	Email       sql.NullString
	Website     sql.NullString
}

type SaveTheater struct {
	TraceID     string
	Name        string
	Address     string
	PhoneNumber string
	Email       string
	Website     string
	Latitude    float32
	Longitude   float32
}

type SaveTheaterResult struct {
	TheaterID string
}

type UpdateTheater struct {
	Name        sql.NullString
	Address     sql.NullString
	PhoneNumber sql.NullString
	Email       sql.NullString
	Website     sql.NullString
	Latitude    float32
	Longitude   float32
}
