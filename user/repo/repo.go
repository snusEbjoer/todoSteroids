package repo

import (
	"database/sql"
	"log"

	"github.com/snusEbjoer/todo-utils/dto"
	"github.com/snusEbjoer/todo-utils/myerrors"
	"github.com/snusEbjoer/todo-utils/types"
	"github.com/snusEbjoer/user/auth"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db}
}
func (r *Repository) CheckUnique(username string) bool {
	if err := r.db.QueryRow(`SELECT username FROM "user" WHERE username=($1);`, username).Scan(&username); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err) // fix later
		}
		return true
	}
	return false
}

func (r *Repository) CreateUser(body types.UserCreds) (dto.UserDto, error) {
	var user dto.UserDto
	ok := r.CheckUnique(body.Username)
	if !ok {
		return dto.UserDto{}, myerrors.ErrUsernameNotUnique
	}
	hPass, err := auth.HashPassword(body.Password)
	if err != nil {
		return dto.UserDto{}, err
	}
	err = r.db.QueryRow(`
	INSERT INTO "user" (username, hashed_password)
	VALUES ($1,$2) RETURNING id, username, description, created_at;
	`, body.Username, hPass).Scan(&user.Id, &user.Username, &user.Description, &user.Created_at)

	if err != nil {
		return dto.UserDto{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByUsername(username string) (types.User, error) {
	var user types.User
	err := r.db.QueryRow(`SELECT * FROM "user" WHERE username=($1);`, username).
		Scan(&user.Id, &user.Username, &user.HashedPassword, &user.Description, &user.Created_at)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

func (r *Repository) GetUserDtoByUsername(username string) (dto.UserDto, error) {
	var user dto.UserDto
	err := r.db.QueryRow(`SELECT id, username, description, created_at FROM "user" WHERE username=($1);`, username).
		Scan(&user.Id, &user.Username, &user.Description, &user.Created_at)
	if err != nil {
		return dto.UserDto{}, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(body types.UpdateUser, id int) (dto.UserDto, error) {
	var user dto.UserDto
	err := r.db.QueryRow(`
	UPDATE "user" SET username=COALESCE($1, username), description=COALESCE($2, description)
	WHERE id=$3 RETURNING id, username, description, created_at;
	`, body.Username, body.Description, id).Scan(&user.Id, &user.Username, &user.Description, &user.Created_at)
	if err != nil {
		return dto.UserDto{}, err
	}
	return user, nil
}
