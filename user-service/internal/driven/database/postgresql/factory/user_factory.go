package factory

import (
	"time"

	"github.com/google/uuid"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"gorm.io/gorm"
)

type UserFactory interface {
	CreateTestUser(p CreateTestUserParam) *model.User
}

func NewUserFactory() UserFactory {
	return &userFactoryImpl{}
}

type userFactoryImpl struct {
}

type CreateTestUserParam struct {
	HashPassword bool
}

func (f *userFactoryImpl) CreateTestUser(p CreateTestUserParam) *model.User {
	pass := "user1234"
	if p.HashPassword {
		pass = "$argon2id$v=19$m=65536,t=1,p=8$idhUhR61RiIephSttaskBA$qVXDMG91UIJr5qduxs5CDO1FC4A8Y8F0QwJhuWOE+tw" // user1234
	}

	return &model.User{
		UUID:                  uuid.MustParse("868b606b-26d5-4c8d-ba45-9587919e059f"),
		TraceID:               uuid.NewString(),
		Username:              "user1234",
		Password:              pass,
		Email:                 "user1234@example.com",
		PhoneNumber:           "+6281234567890",
		FirstName:             "Example",
		LastName:              "User",
		IsEmailVerified:       true,
		IsPhoneNumberVerified: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
}
