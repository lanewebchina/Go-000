package main

import (
	"fmt"
	service "week02/service"

	"github.com/pkg/errors"
)

func main() {
	User, err := new(service.UserService).GetByUserID(100)
	if err != nil {
		//使用errors.Cause来获取根因
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		//使用%+v打印堆栈详细记录
		fmt.Printf("stack trace:\n%+v\n", err)
	}
	fmt.Printf("User=%v", User)
}
