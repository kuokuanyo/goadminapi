package ui

import (
	"goadminapi/modules/service"
	"goadminapi/template/types"
)

// Service struct
type Service struct {
	NavButtons *types.Buttons
}

// GetService get service(struct)
func GetService(srv service.List) *Service {
	if v, ok := srv.Get("ui").(*Service); ok {
		return v
	}
	panic("wrong service")
}

// Name return ui
func (s *Service) Name() string {
	return "ui"
}

// NewService 將參數設置service(struct)
func NewService(btns *types.Buttons) *Service {
	return &Service{
		NavButtons: btns,
	}
}
