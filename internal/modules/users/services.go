package users

import (
	"context"
	"errors"
	"fmt"
	"os"

	"attendance/pkg"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB   *sqlx.DB
	sqlx *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		sqlx: db,
	}
}

func (u *Service) RegisterUser(ctx context.Context, arg CreateUserParams) (UserQuery, error) {
	hash, _ := pkg.HashPassword(arg.Password)

	user := UserQuery{
		Username:  arg.Username,
		FirstName: &arg.FirstName,
		LastName:  &arg.LastName,
		Age:       arg.Age,
		Password:  hash,
	}
	_, err := u.sqlx.NamedExec("INSERT INTO tbl_users (username, first_name, last_name, age, password) VALUES (:username, :first_name, :last_name, :age, :password)", user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return UserQuery{}, fmt.Errorf("username %s already exists", arg.Username)
			}
		}
		return UserQuery{}, fmt.Errorf("failed to create user")
	}
	return user, nil

}

func (u *Service) Login(ctx context.Context, arg LoginParams) (LoginPayload, error) {
	var dbUser struct {
		ID       int    `db:"id"`
		Password string `db:"password"`
	}
	err := u.sqlx.Get(&dbUser, "SELECT id, password FROM tbl_users WHERE username = $1", arg.Username)
	if err != nil {
		return LoginPayload{}, fmt.Errorf("invalid username or password")
	}

	if !pkg.VerifyPassword(arg.Password, dbUser.Password) {
		return LoginPayload{}, fmt.Errorf("invalid username or password")
	}

	// generate token
	secret := os.Getenv("JWT_SECRET_KEY")
	token, err := pkg.GenerateToken(pkg.UserClaim{
		ID:       dbUser.ID,
		Username: arg.Username,
	}, secret)

	if err != nil {
		return LoginPayload{}, err
	}

	loginPayload := LoginPayload{
		Username: arg.Username,
		ID:       dbUser.ID,
		Token:    token,
	}
	return loginPayload, nil
}

func (u *Service) GetAllUsers(ctx context.Context) ([]UserQuery, error) {
	users := []UserQuery{}
	err := u.sqlx.Select(&users, "SELECT username, first_name, last_name, age FROM tbl_users")
	if err != nil {
		return nil, err
	}
	return users, nil
}
