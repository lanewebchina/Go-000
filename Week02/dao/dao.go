package dao

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

var ErrRecordNotFound = errors.New("record not found")

type User struct {
	UserID   int
	UserName string
}

func (u *User) GetByUserID(userID int) (user User, err error) {
	err = DB.Table(UserTable).Where("id=?", userID).First(user).Error

	if errors.Is(err, sql.ErrNoRows) {
		err = ErrRecordNotFound
	}
	if err != nil {
		//其实只是底层的方法需要Wrap
		err = errors.Wrap(err, fmt.Sprintf("find user by userid: %v ", userID))
	}
	return
}
