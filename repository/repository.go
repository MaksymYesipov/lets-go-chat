package repository

import (
	"database/sql"
	"httpServer/repository/domain"
)

type UserRepository interface {
	GetAll() ([]domain.User, error)
	Get(ID string) (*domain.User, error)
	Create(user domain.User) (*sql.Rows, error)
	CloseDB()
}
