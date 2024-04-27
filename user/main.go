package main

import (
	_ "github.com/lib/pq"
	"github.com/snusEbjoer/todo-utils/rmq"
	"github.com/snusEbjoer/todo-utils/utils"
	"github.com/snusEbjoer/user/controller"
	"github.com/snusEbjoer/user/db"
	"github.com/snusEbjoer/user/repo"
)

func main() {
	database, err := db.Connect()
	utils.FailOnError(err, "failed to connect to db")

	err = db.Prepare(database)
	utils.FailOnError(err, "failed to prepare database")

	rabbit, err := rmq.New("amqp://user:pass@rabbitmq:5672/", "user")
	utils.FailOnError(err, "failed to connect to rabbitmq")

	repository := repo.New(database)
	ctrl := controller.Controller{Repo: repository}

	rabbit.HandleMessage("user.signUp", ctrl.SignUp)
	rabbit.HandleMessage("user.login", ctrl.Login)
	rabbit.HandleMessage("user.current", ctrl.GetUser)
	rabbit.HandleMessage("user.patch", ctrl.PatchUser)

	rabbit.Listen()
}
