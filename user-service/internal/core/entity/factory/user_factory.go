package entityfactory

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type UserRaw struct {
	Password string
}

type UserFactory interface {
	GenerateUser() (user *entity.User, raw *UserRaw, err error)
	GenerateUserV2() (user *entity.User, raw *UserRaw, err error)
}

func NewUserFactory(util *util.Util) UserFactory {
	return &userFactoryImpl{
		util: util,
	}
}

type userFactoryImpl struct {
	util *util.Util
}

func (f *userFactoryImpl) GenerateUser() (*entity.User, *UserRaw, error) {
	var user entity.User
	if err := faker.FakeData(&user); err != nil {
		fmt.Printf("error faking data: %v\n", err)
		return nil, nil, err
	}

	user.Password = "$argon2id$v=19$m=65536,t=1,p=8$idhUhR61RiIephSttaskBA$qVXDMG91UIJr5qduxs5CDO1FC4A8Y8F0QwJhuWOE+tw" // user1234
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.DeletedAt = sql.NullTime{Time: time.Time{}, Valid: false}

	return &user, &UserRaw{Password: "user1234"}, nil
}

func (f *userFactoryImpl) GenerateUserV2() (*entity.User, *UserRaw, error) {
	var user entity.User
	var err error
	if err = faker.FakeData(&user); err != nil {
		fmt.Printf("error faking data: %v\n", err)
		return nil, nil, err
	}

	password := faker.Password()
	user.Password, err = f.util.EncryptionUtil.Argon2Hasher.Hash(password)
	if err != nil {
		fmt.Printf("error hashing password: %v\n", err)
		return nil, nil, err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.DeletedAt = sql.NullTime{Time: time.Time{}, Valid: false}

	return &user, &UserRaw{Password: password}, nil
}
