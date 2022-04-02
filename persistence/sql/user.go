package sql

import (
	"errors"
	"time"

	"github.com/thebluefowl/suckerfish/db"
	"github.com/thebluefowl/suckerfish/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	client *db.PGClient
}

func NewUserRepository(client *db.PGClient) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}

type User struct {
	ID        string `gorm:"size:32"`
	Name      string `gorm:"size:128"`
	AvatarURL string `gorm:"size:256"`
	Email     string `gorm:"size:256"`
	Provider  string `gorm:"size:16"`
	Token     string `gorm:"size:512"`
	Company   string `gorm:"size:256"`
	Location  string `gorm:"size:256"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsStaff   bool
}

func (u *User) Domain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Email:     u.Email,
		Provider:  u.Provider,
		Token:     u.Token,
		Company:   u.Company,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsStaff:   u.IsStaff,
	}
}

func getDBUser(u *domain.User) *User {
	return &User{
		ID:        u.ID,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Email:     u.Email,
		Provider:  u.Provider,
		Token:     u.Token,
		Company:   u.Company,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsStaff:   u.IsStaff,
	}
}

func (repository *userRepository) GetByEmail(email string) (*domain.User, error) {
	dbUser := &User{}
	if err := repository.client.DB.Where(&User{Email: email}).First(dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dbUser.Domain(), nil
}

func (repository *userRepository) Create(user *domain.User) error {
	dbUser := getDBUser(user)
	if err := repository.client.DB.Create(dbUser).Error; err != nil {
		return err
	}
	return nil
}
