package entity

import (
	"net/url"
	"time"
)

type People struct {
	PeopleID   string    `json:"people_id" bson:"_id"`
	TraceID    string    `json:"trace_id" bson:"trace_id"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
	Name       string    `json:"name" bson:"name"`
	PictureURL string    `json:"picture_url" bson:"picture_url"`
	// Bio         string       `json:"bio" bson:"bio"`
	// Occupations []string     `json:"occupations" bson:"occupations"`
	// BirthDate   time.Time    `json:"birth_date" bson:"birth_date"`
	// BirthPlace  string       `json:"birth_place" bson:"birth_place"`
	// DeathDate   sql.NullTime `json:"death_date" bson:"death_date"`
}

func (p *People) GetIMDbUrl() string {
	return "https://www.imdb.com/search/name/?name=" + url.QueryEscape(p.Name)
}

// type Occupation string

// func (o Occupation) String() string {
// 	return string(o)
// }

// const (
// 	OccupationActor    Occupation = "actor"
// 	OccupationDirector Occupation = "director"
// 	OccupationWriter   Occupation = "writer"
// )
