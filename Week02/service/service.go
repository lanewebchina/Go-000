package service

import "week02/dao"

type UserService struct {
}

func (s *UserService) GetByUserID(userID int) (dao.User, error) {
	return new(dao.User).GetByUserID(userID)
}
