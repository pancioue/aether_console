package health

import "context"

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Check(ctx context.Context) (string, error) {
	// 之後你要加 DB ping / dependencies check 都在這裡做
	return "ok", nil
}