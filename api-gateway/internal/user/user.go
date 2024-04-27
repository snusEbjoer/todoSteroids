package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/snusEbjoer/todo-utils/rmq"
	"github.com/snusEbjoer/todo-utils/types"
)

type Controller struct {
	Rmq *rmq.Rmq
}
type Handler = func(w http.ResponseWriter, r *http.Request)

func New(rmq *rmq.Rmq) *Controller {
	return &Controller{rmq}
}

func (c *Controller) Login() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var body types.UserCreds
		json.NewDecoder(r.Body).Decode(&body)
		bodyJson, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		ctx := context.Background()
		data, err := c.Rmq.Send(ctx, "user.login", bodyJson)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		val := data.Headers["status"].(int32)
		w.WriteHeader(int(val))
		w.Write(data.Body)
	}
}

func (c *Controller) SignUp() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var body types.UserCreds
		json.NewDecoder(r.Body).Decode(&body)
		bodyJson, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		ctx := context.Background()
		data, err := c.Rmq.Send(ctx, "user.signUp", bodyJson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.WriteHeader(int(data.Headers["status"].(int32))) // im not a psycho
		w.Write(data.Body)
	}
}

func (c *Controller) GetCurrentUser() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Context())
		username := r.Context().Value(types.UsernameKey).(string)
		ctx := context.Background()
		data, err := c.Rmq.Send(ctx, "user.current", []byte(username))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.WriteHeader(int(data.Headers["status"].(int32)))
		w.Write(data.Body)
	}
}

func (c *Controller) PatchUser() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.UpdateUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if user.Description == nil && user.Username == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		id := r.Context().Value(types.UserIDKey).(int)
		user.Id = &id
		ctx := context.Background() // add timeout
		bodyJson, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		data, err := c.Rmq.Send(ctx, "user.patch", bodyJson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.WriteHeader(int(data.Headers["status"].(int32)))
		w.Write(data.Body)
	}
}
