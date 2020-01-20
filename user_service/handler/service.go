package handler

import (
	"context"
	"github.com/micro/go-micro/util/log"
	//"time"
	us "Micro_API_Framework/user_service/model/user"
	s "Micro_API_Framework/user_service/proto/user"
)

type Service struct{}

// Init 初始化handler
func Init() {

	var err error
	if err != nil {
		log.Fatal("[Init] 初始化Handler错误")
		return
	}
}

// QueryUserByName 通过参数中的名字返回用户
func (e *Service) QueryUserByName(ctx context.Context, req *s.Request, rsp *s.Response) error {

	user, err := us.QueryUserByName(req.UserName)

	if err != nil {
		rsp.Error = &s.Error{
			Code:   500,
			Detail: err.Error(),
		}

		return nil
	}

	//测试超时容错hystrix
	//time.Sleep(3 * time.Second)

	rsp.User = user
	return nil
}
