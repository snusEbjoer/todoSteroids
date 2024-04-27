package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/snusEbjoer/api-gateway/internal/user"
	"github.com/snusEbjoer/todo-utils/rmq"
	"github.com/snusEbjoer/todo-utils/types"
	"github.com/snusEbjoer/todo-utils/utils"
)

func JWTVerify(next func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}
		accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			log.Println(ok)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					return nil, err
				}
			}

			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
		log.Print(token, token.Valid, token.Method)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You're Unauthorized due invalid token"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get claims jwt"))
			return
		}

		ctx := context.Background()
		log.Println(claims)
		ctx = context.WithValue(ctx, types.UsernameKey, claims["username"].(string))
		ctx = context.WithValue(ctx, types.UserIDKey, int(claims["id"].(float64)))
		log.Println(ctx)
		next(w, r.WithContext(ctx))
	})
}

func main() {
	rmq, err := rmq.New("amqp://user:pass@rabbitmq:5672/", "gateway")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	utils.FailOnError(err, "Failed to declare exchage")
	fmt.Print("started")
	userCrtl := user.New(rmq)
	router := http.NewServeMux()

	router.HandleFunc("POST /signUp", userCrtl.SignUp())
	router.HandleFunc("POST /login", userCrtl.Login())
	router.HandleFunc("GET /user", JWTVerify(userCrtl.GetCurrentUser()))
	router.HandleFunc("PATCH /user", JWTVerify(userCrtl.PatchUser()))

	log.Fatal(http.ListenAndServe("0.0.0.0:42069", router))
}
