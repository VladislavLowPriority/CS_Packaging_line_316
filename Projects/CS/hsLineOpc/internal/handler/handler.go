package handler

import (
	"context"
	"errors"
	"hsLineOpc/api"
	"hsLineOpc/internal/consts"
	"hsLineOpc/internal/galaction"
	"time"
)

// import (
// 	"encoding/json"
// 	"hsLineOpc/api"
// 	"hsLineOpc/internal/consts"
// 	"hsLineOpc/internal/opc"
// 	"log/slog"
// 	"net/http"
// 	"time"
// )

type ControlSystem struct {
	IsActive  bool
	IsDefault bool

	hsClient *api.OpcClient
}

func NewControlSystem(hsClient *api.OpcClient) *ControlSystem {
	return &ControlSystem{
		IsActive:  false,
		IsDefault: false,

		hsClient: hsClient,
	}
}

func (s *ControlSystem) makeGripperDefault() {
	for !s.hsClient.GetNodeValue(consts.InputGripperPack) && !s.hsClient.GetNodeValue(consts.InputGripperStart) {
		s.hsClient.WriteNodeValue(consts.OutputGripperLeft, true)
	}
	s.hsClient.WriteNodeValue(consts.OutputGripperLeft, false)

	for !s.hsClient.GetNodeValue(consts.InputGripperConveyor) && !s.hsClient.GetNodeValue(consts.InputGripperStart) {
		s.hsClient.WriteNodeValue(consts.OutputGripperRight, true)
	}
	s.hsClient.WriteNodeValue(consts.OutputGripperRight, false)

	// опустить шайбу
	s.hsClient.WriteNodeValue(consts.OutputGripperOpen, true)
	s.hsClient.WriteNodeValue(consts.OutputGripperUpDown, true)
	time.Sleep(time.Second * 3)

	s.hsClient.WriteNodeValue(consts.OutputGripperOpen, false)
	s.hsClient.WriteNodeValue(consts.OutputGripperUpDown, false)
	time.Sleep(time.Second * 3)
}

// Возврат установки в начальное положение
func (s *ControlSystem) Default() {
	if s.IsDefault || s.IsActive {
		return
	}

	// TODO: add check default pos
	s.makeGripperDefault()
	s.IsDefault = true
}

// Запуск установки
func (s *ControlSystem) Start(ctx context.Context) {
	if !s.IsDefault {
		return
	}

	if s.IsActive {
		return
	}

	s.IsActive = true
	s.IsDefault = false

	go func(ctx context.Context) {
		for s.IsActive {
			galaction.EntryStartHs(ctx, s.hsClient)
		}
	}(ctx)
}

// Остановка установки
func (s *ControlSystem) Stop() error {
	if !s.IsActive {
		return errors.New("system alredy stopped")
	}

	s.hsClient.SendAllFalses()
	s.IsActive = false

	return nil
}
