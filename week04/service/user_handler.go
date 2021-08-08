package service

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type IUserService interface {
	GetUserNameByAccount(userAccount string) (string, error)
}

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) ServeHTTP(rsp http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userAccount := r.FormValue("userAccount")
	userName, err := h.userService.GetUserNameByAccount(userAccount)
	rspStr := ""
	if err == nil {
		rspStr = fmt.Sprintf("Get your name is %s", userName)
	} else {
		err = errors.Wrap(err, "UserHandler")
		rspStr = err.Error()
	}
	rsp.Write([]byte(rspStr))
}
