package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type ProcController struct {
	// Состояние переменных
	counter      int  // Счётчик шагов процесса
	finished     bool // Флаг завершения процесса
	holeDetected bool // Результат обнаружения отверстия
	m5Status     bool // Текущее состояние цилиндра M5

	// Конфигурация OPC UA клиента
	client   *opcua.Client
	endpoint string

	// Коллекция NodeID
	nodes struct {
		carouselRotation *ua.NodeID
		m5Sensor         *ua.NodeID
		colorSensor      *ua.NodeID
		silverSensor     *ua.NodeID
		outputDrill      *ua.NodeID
		outputCarousel   *ua.NodeID
		outputDrillDown  *ua.NodeID
		outputDrillUp    *ua.NodeID
		outputM4Clamp    *ua.NodeID
		outputM5Cylinder *ua.NodeID
		red_tag          *ua.NodeID
		silver_tag       *ua.NodeID
		black_tag        *ua.NodeID
	}
}

func NewProcController(endpoint string) *ProcController {
	return &ProcController{
		endpoint: endpoint,
		nodes: struct {
			carouselRotation *ua.NodeID
			m5Sensor         *ua.NodeID
			colorSensor      *ua.NodeID
			silverSensor     *ua.NodeID
			outputDrill      *ua.NodeID
			outputCarousel   *ua.NodeID
			outputDrillDown  *ua.NodeID
			outputDrillUp    *ua.NodeID
			outputM4Clamp    *ua.NodeID
			outputM5Cylinder *ua.NodeID
			red_tag          *ua.NodeID
			silver_tag       *ua.NodeID
			black_tag        *ua.NodeID
		}{
			carouselRotation: ua.MustParseNodeID("ns=4;i=3"),
			m5Sensor:         ua.MustParseNodeID("ns=4;i=4"),
			colorSensor:      ua.MustParseNodeID("ns=4;i=6"),
			silverSensor:     ua.MustParseNodeID("ns=4;i=7"),
			outputDrill:      ua.MustParseNodeID("ns=4;i=12"),
			outputCarousel:   ua.MustParseNodeID("ns=4;i=13"),
			outputDrillDown:  ua.MustParseNodeID("ns=4;i=14"),
			outputDrillUp:    ua.MustParseNodeID("ns=4;i=15"),
			outputM4Clamp:    ua.MustParseNodeID("ns=4;i=16"),
			outputM5Cylinder: ua.MustParseNodeID("ns=4;i=17"),
			red_tag:          ua.MustParseNodeID("ns=4;i=24"),
			silver_tag:       ua.MustParseNodeID("ns=4;i=26"),
			black_tag:        ua.MustParseNodeID("ns=4;i=25"),
		},
	}
}

func (pc *ProcController) Connect() error {
	ctx := context.Background()                                                            // Создать пустой контекст
	client := opcua.NewClient(pc.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone)) // Без шифрования
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("Ошибка подключения: %w", err)
	}
	pc.client = client
	return nil
}

func (pc *ProcController) Run() error {
	defer pc.client.Close()

	// Инициализация поворотного стола
	if err := pc.activateCarousel(); err != nil {
		return err
	}

	for !pc.finished {
		if pc.counter < 4 {
			pc.counter++
		}

		switch {
		case pc.counter == 4:
			if err := pc.handleColorDetection(); err != nil { // Обнаружение цвета
				return err
			}
			pc.counter++

		case pc.counter == 5:
			if err := pc.handleDrilling(); err != nil { // Обработка сверления
				return err
			}
			pc.resetState() // Сброс состояния
		}

		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

func (pc *ProcController) activateCarousel() error {
	return pc.writeBool(pc.nodes.outputCarousel, true)
}

func (pc *ProcController) handleColorDetection() error {
	// Остановка стола
	if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	// Чтение датчика цвета
	colorVal, err := pc.readColor()
	if err != nil {
		return err
	}

	// Маркировка цвета
	if err := pc.markColor(colorVal); err != nil {
		return err
	}

	// Проверка цилиндра M5
	return pc.detectHole()
}

func (pc *ProcController) readColor() (string, error) {
	resp, err := pc.client.Read(&ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: pc.nodes.colorSensor},
			{NodeID: pc.nodes.silverSensor},
		},
	})
	if err != nil || len(resp.Results) < 2 {
		return "", fmt.Errorf("Ошибка датчика: %w", err)
	}

	colorActive := resp.Results[0].Value.Value().(bool)
	silverActive := resp.Results[1].Value.Value().(bool)

	switch {
	case colorActive && silverActive:
		return "серебристый", nil
	case colorActive:
		return "красный", nil
	default:
		return "чёрный", nil
	}
}

func (pc *ProcController) markColor(color string) error {
	var targetNode *ua.NodeID
	switch color {
	case "серебристый":
		targetNode = pc.nodes.silver_tag
	case "красный":
		targetNode = pc.nodes.red_tag
	default:
		targetNode = pc.nodes.black_tag
	}
	return pc.writeBool(targetNode, true)
}

func (pc *ProcController) detectHole() error {
	// Активация цилиндра
	if err := pc.writeBool(pc.nodes.outputM5Cylinder, true); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	// Чтение датчика
	val, err := pc.readBool(pc.nodes.m5Sensor)
	if err != nil {
		return err
	}
	pc.holeDetected = val

	// Деактивация цилиндра
	return pc.writeBool(pc.nodes.outputM5Cylinder, false)
}

func (pc *ProcController) handleDrilling() error {
	// Остановка стола
	if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if pc.holeDetected {
		if err := pc.performDrilling(); err != nil {
			return err
		}
	}

	// Сброс стола
	return pc.resetCarousel()
}

func (pc *ProcController) performDrilling() error {
	// Фиксация зажимом
	if err := pc.writeBool(pc.nodes.outputM4Clamp, true); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	// Опускание сверла
	if err := pc.drillDown(); err != nil {
		return err
	}

	// Активация сверления
	if err := pc.activateDrill(); err != nil {
		return err
	}

	// Подъём сверла
	return pc.drillUp()
}

func (pc *ProcController) drillDown() error {
	if err := pc.writeBool(pc.nodes.outputDrillDown, true); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return pc.writeBool(pc.nodes.outputDrillDown, false)
}

func (pc *ProcController) activateDrill() error {
	if err := pc.writeBool(pc.nodes.outputDrill, true); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return pc.writeBool(pc.nodes.outputDrill, false)
}

func (pc *ProcController) drillUp() error {
	if err := pc.writeBool(pc.nodes.outputDrillUp, true); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return pc.writeBool(pc.nodes.outputDrillUp, false)
}

func (pc *ProcController) resetCarousel() error {
	if err := pc.writeBool(pc.nodes.outputCarousel, true); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	return pc.writeBool(pc.nodes.outputCarousel, false)
}

func (pc *ProcController) resetState() {
	pc.counter = -1
	pc.finished = true
}

// Утилиты записи/чтения
func (pc *ProcController) writeBool(node *ua.NodeID, value bool) error {
	_, err := pc.client.Write(&ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      node,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					Value: ua.MustVariant(value),
				},
			},
		},
	})
	return err
}

func (pc *ProcController) readBool(node *ua.NodeID) (bool, error) {
	resp, err := pc.client.Read(&ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: node},
		},
	})
	if err != nil || len(resp.Results) == 0 {
		return false, err
	}
	return resp.Results[0].Value.Value().(bool), nil
}

func main() {
	controller := NewProcController("opc.tcp://localhost:4840")

	if err := controller.Connect(); err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	if err := controller.Run(); err != nil {
		log.Fatalf("Ошибка выполнения: %v", err)
	}

	fmt.Println("Процесс завершен успешно.")
}
