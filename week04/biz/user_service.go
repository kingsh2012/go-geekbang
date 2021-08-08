package biz

import (
	"github.com/pkg/errors"
	"go-geekbang/week04/data"
)

type IDao interface {
	GetUserData(account string) (user data.UserData, err error)
}

type UserService struct {
	dao IDao
}

func NewUserService(dao IDao) *UserService {
	return &UserService{
		dao: dao,
	}
}

func (s *UserService) GetUserNameByAccount(userAccount string) (userName string, err error) {
	userInfo, err := s.dao.GetUserData(userAccount)
	if err != nil {
		err = errors.Wrap(err, "UserService.GetUserNameByAccount")
		return
	}
	userName = userInfo.UserName
	return
}
