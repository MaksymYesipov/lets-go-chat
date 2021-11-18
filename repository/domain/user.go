package domain

type User struct {
	ID       string `db:"id"`
	UserName string `db:"username"`
	Password string `db:"password"`
}
