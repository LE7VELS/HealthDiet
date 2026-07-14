package service

import "context"

// AppService 提供当前应用最小的业务输出。
// 持久化业务 Service 后续通过构造函数接收 Repository 接口。
type AppService struct {
	message string
}

func NewAppService(message string) *AppService {
	return &AppService{message: message}
}

func (s *AppService) Message(context.Context) string {
	return s.message
}
