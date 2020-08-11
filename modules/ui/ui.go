package ui

import (
	"goadminapi/modules/service"
	"goadminapi/template/types"
)

type Service struct {
	NavButtons *types.Buttons
}

func GetService(srv service.List) *Service {
	if v, ok := srv.Get("ui").(*Service); ok {
		return v
	}
	panic("wrong service")
}

func (s *Service) Name() string {
	return "ui"
}

func NewService(btns *types.Buttons) *Service {
	return &Service{
		NavButtons: btns,
	}
}
