package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"httpServer/repository/domain"
)

type UserPostgresRepository struct {
	conn *sqlx.DB
}

func (p UserPostgresRepository) GetAll() ([]domain.User, error) {
	var result []domain.User
	rows, err := p.conn.Queryx("SELECT id, username, password FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := domain.User{}
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (p UserPostgresRepository) Get(ID string) (*domain.User, error) {
	var result domain.User
	err := p.conn.Get(&result, "SELECT * FROM users WHERE id = ?", ID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p UserPostgresRepository) Create(user domain.User) (*sql.Rows, error) {
	rows, err := p.conn.Query("INSERT INTO users(id, username, password) VALUES ($1, $2, $3)", user.ID, user.UserName, user.Password)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p UserPostgresRepository) CloseDB() {
	err := p.conn.Close()
	if err != nil {
		panic(err)
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "chat-db"
)

func GetPostgresRepository() (UserRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected!")

	return UserPostgresRepository{conn: db}, nil
}
