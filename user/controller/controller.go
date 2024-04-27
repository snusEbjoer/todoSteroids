package controller

import (
	"encoding/json"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/snusEbjoer/todo-utils/myerrors"
	"github.com/snusEbjoer/todo-utils/rmq"
	"github.com/snusEbjoer/todo-utils/types"
	"github.com/snusEbjoer/user/auth"
	"github.com/snusEbjoer/user/repo"
)

type Controller struct {
	Repo *repo.Repository
}

func (c *Controller) Login(msg amqp.Delivery) rmq.Message {
	var body types.UserCreds
	headers := amqp.Table{}
	err := json.Unmarshal(msg.Body, &body)

	if err != nil {
		headers["status"] = http.StatusBadRequest
		return rmq.Message{Body: []byte("Bad Request"), Headers: headers}
	}

	user, err := c.Repo.GetUserByUsername(body.Username)
	if err != nil {
		headers["status"] = http.StatusUnauthorized
		return rmq.Message{Body: []byte("err: User does not exists"), Headers: headers}
	}

	ok := auth.CheckPasswordHash(body.Password, user.HashedPassword)
	if !ok {
		headers["status"] = http.StatusUnauthorized
		return rmq.Message{Body: []byte("Wrong Password"), Headers: headers}
	}

	accessToken, err := auth.GenerateAccessToken(user.Username, user.Id)
	if err != nil {
		headers["status"] = http.StatusInternalServerError
		return rmq.Message{Body: []byte("Internal Server Error"), Headers: headers}
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		headers["status"] = http.StatusInternalServerError
		return rmq.Message{Body: []byte("Internal Server Error"), Headers: headers}
	}

	data, err := json.Marshal(types.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		headers["status"] = http.StatusInternalServerError
		return rmq.Message{Body: []byte("Internal Server Error"), Headers: headers}
	}
	headers["status"] = http.StatusOK
	return rmq.Message{Body: data, Headers: headers}
}

func (c *Controller) SignUp(msg amqp.Delivery) rmq.Message {
	var body types.UserCreds

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusBadRequest,
		}
		return rmq.Message{Body: []byte("Bad Request"), Headers: headers}
	}

	user, err := c.Repo.CreateUser(body)
	if err == myerrors.ErrUsernameNotUnique {
		headers := amqp.Table{
			"status": http.StatusBadRequest,
		}
		return rmq.Message{Body: []byte(err.Error()), Headers: headers}
	}
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusInternalServerError,
		}
		return rmq.Message{Body: []byte("Internal Server Error(pizdec)"), Headers: headers}
	}
	data, err := json.Marshal(user)
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusInternalServerError,
		}
		return rmq.Message{Body: []byte("Internal Server Error(pizdec)"), Headers: headers}
	}
	headers := amqp.Table{
		"status": http.StatusOK,
	}
	return rmq.Message{Body: data, Headers: headers}
}

func (c *Controller) GetUser(msg amqp.Delivery) rmq.Message {
	user, err := c.Repo.GetUserDtoByUsername(string(msg.Body))
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusBadRequest,
		}
		return rmq.Message{Body: []byte("User not found"), Headers: headers}
	}
	data, err := json.Marshal(user)
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusInternalServerError,
		}
		return rmq.Message{Body: []byte("Internal Server Error(pizdec)"), Headers: headers}
	}
	headers := amqp.Table{
		"status": http.StatusOK,
	}
	return rmq.Message{Body: data, Headers: headers}
}

func (c *Controller) PatchUser(msg amqp.Delivery) rmq.Message {
	var user types.UpdateUser
	err := json.Unmarshal(msg.Body, &user)

	if err != nil {
		headers := amqp.Table{
			"status": http.StatusBadRequest,
		}
		return rmq.Message{Body: []byte("Bad Request"), Headers: headers}
	}

	if user.Username != nil {
		if ok := c.Repo.CheckUnique(*user.Username); !ok {
			headers := amqp.Table{
				"status": http.StatusBadRequest,
			}
			return rmq.Message{Body: []byte("Username alredy taken"), Headers: headers}
		}
	}

	data, err := c.Repo.UpdateUser(user, *user.Id)
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusInternalServerError,
		}
		return rmq.Message{Body: []byte("Internal Server Error(pizdec)"), Headers: headers}
	}

	body, err := json.Marshal(data)
	if err != nil {
		headers := amqp.Table{
			"status": http.StatusInternalServerError,
		}
		return rmq.Message{Body: []byte("Internal Server Error(pizdec)"), Headers: headers}
	}
	headers := amqp.Table{
		"status": http.StatusOK,
	}
	return rmq.Message{Body: body, Headers: headers}
}
