package handler

import (
	"encoding/json"
	"hsLineOpc/api"
	"hsLineOpc/internal/consts"
	"hsLineOpc/internal/opc"
	"log/slog"
	"net/http"
	"time"
)

type MyServer struct {
	opcClient *opc.MyClient
	active    bool
	inDefault bool
}

type Config struct {
	Port string
}

var _ api.ServerInterface = (*MyServer)(nil)

func NewServer() *MyServer {
	return &MyServer{
		opcClient: opc.NewClient(),
		active:    false,
		inDefault: false,
	}
}

func (s *MyServer) makeGripperDefault() {
	for !s.opcClient.GetNodeValue(consts.InputGripperPack) && !s.opcClient.GetNodeValue(consts.InputGripperStart) {
		s.opcClient.WriteNodeValue(consts.OutputGripperLeft, true)
	}
	s.opcClient.WriteNodeValue(consts.OutputGripperLeft, false)

	for !s.opcClient.GetNodeValue(consts.InputGripperConveyor) && !s.opcClient.GetNodeValue(consts.InputGripperStart) {
		s.opcClient.WriteNodeValue(consts.OutputGripperRight, true)
	}
	s.opcClient.WriteNodeValue(consts.OutputGripperRight, false)

	// опустить шайбу
	s.opcClient.WriteNodeValue(consts.OutputGripperOpen, true)
	s.opcClient.WriteNodeValue(consts.OutputGripperUpDown, true)
	time.Sleep(time.Second * 3)

	s.opcClient.WriteNodeValue(consts.OutputGripperOpen, false)
	s.opcClient.WriteNodeValue(consts.OutputGripperUpDown, false)
	time.Sleep(time.Second * 3)
}

// Возврат установки в начальное положение
// (POST /default)
func (s *MyServer) PostDefault(w http.ResponseWriter, r *http.Request) {
	if s.inDefault {
		sendInfoResponse(w, nil, http.StatusOK)
		return
	}

	// TODO: add check default pos
	s.makeGripperDefault()
	s.inDefault = true

	sendInfoResponse(w, nil, http.StatusOK)
}

// Запуск установки
// (POST /start)
func (s *MyServer) PostStart(w http.ResponseWriter, r *http.Request) {
	if !s.inDefault {
		sendErrorResponse(w, "установка не в начальном положении", http.StatusForbidden)
		return
	}

	if s.active {
		sendInfoResponse(w, nil, http.StatusOK)
		return
	}

	s.active = true
	s.inDefault = false

	go func() {
		for s.active {
			// TODO: galaction start
			slog.Info("working...")
			time.Sleep(time.Second)
		}
	}()

	sendInfoResponse(w, nil, http.StatusOK)
}

// Остановка установки
// (POST /stop)
func (s *MyServer) PostStop(w http.ResponseWriter, r *http.Request) {
	if !s.active {
		sendInfoResponse(w, nil, http.StatusOK)
		return
	}

	s.active = false

	sendInfoResponse(w, nil, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	resp, _ := json.Marshal(api.Error{Message: errMsg})
	slog.Error(errMsg)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func sendInfoResponse(w http.ResponseWriter, object any, status int) {
	if object != nil {
		resp, err := json.Marshal(object)
		if err != nil {
			sendErrorResponse(w, "failed to form response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	w.WriteHeader(status)
}
