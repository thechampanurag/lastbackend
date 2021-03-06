package user

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SignInCmd() {

	var (
		login    string
		password string
		ctx      = context.Get()
	)

	fmt.Print("Login: ")
	fmt.Scan(&login)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
	password = string(pass)

	fmt.Print("\r\n")

	err = SignIn(login, password)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Login successful")
}

type userLoginS struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SignIn(login, password string) error {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
	)

	res := struct {
		Token string `json:"token"`
	}{}

	_, _, err = ctx.HTTP.
		POST("/session").
		AddHeader("Content-Type", "application/json").
		BodyJSON(userLoginS{login, password}).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return e.LoginErrorMessage
	}

	if er.Code == 404 {
		return e.LoginErrorMessage
	}

	if er.Code == 500 {
		return e.UnknownMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	err = ctx.Storage.Set("session", res)
	if err != nil {
		return err
	}

	return nil
}
