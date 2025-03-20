package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/gopcua/opcua"
 "github.com/gopcua/opcua/ua"
)

// Контроллер технологического процесса
type ProcController struct {
 // Состояние системы
 counter      int  // Счетчик шагов процесса
 finished     bool // Флаг завершения всего процесса
 holeDetected bool // Результат проверки отверстия
 m5Status     bool // Текущее состояние цилиндра M5

 // Конфигурация OPC UA клиента
 client   *opcua.Client
 endpoint string

 // Идентификаторы узлов OPC UA
 nodes struct {
  carouselRotation *ua.NodeID // Вращение карусели
  m5Sensor         *ua.NodeID // Датчик M5
  colorSensor      *ua.NodeID // Цветовой датчик
  silverSensor     *ua.NodeID // Датчик серебра
  outputDrill      *ua.NodeID // Управление сверлом
  outputCarousel   *ua.NodeID // Управление каруселью
  outputDrillDown  *ua.NodeID // Опускание сверла
  outputDrillUp    *ua.NodeID // Подъем сверла
  outputM4Clamp    *ua.NodeID // Зажим M4
  outputM5Cylinder *ua.NodeID // Цилиндр M5
  redTag           *ua.NodeID // Метка красного цвета
  silverTag        *ua.NodeID // Метка серебряного цвета
  blackTag         *ua.NodeID // Метка черного цвета
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
   redTag           *ua.NodeID
   silverTag        *ua.NodeID
   blackTag         *ua.NodeID
  }{
   // Инициализация идентификаторов узлов OPC UA
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
   redTag:           ua.MustParseNodeID("ns=4;i=24"),
   silverTag:        ua.MustParseNodeID("ns=4;i=26"),
   blackTag:         ua.MustParseNodeID("ns=4;i=25"),
  },
 }
}

// Подключение к OPC UA серверу
func (pc *ProcController) Connect() error {
 ctx := context.Background()
 client := opcua.NewClient(pc.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
 if err := client.Connect(ctx); err != nil {
  return fmt.Errorf("ошибка подключения: %w", err)
 }
 pc.client = client
 return nil
}

// Основной цикл управления процессом
func (pc *ProcController) Run() error {
 defer pc.client.Close()

 // Инициализация вращения карусели
 if err := pc.activateCarousel(); err != nil {
  return err
 }

 // Главный цикл обработки
 for !pc.finished {
  // Логика перехода между шагами процесса
  if pc.counter < 4 {
   pc.counter++
  }

  switch {
  case pc.counter == 4:
   // Этап определения цвета и проверки отверстия
   if err := pc.handleColorDetection(); err != nil {
    return err
   }
   pc.counter++

  case pc.counter == 5:
   // Этап сверления и завершения процесса
   if err := pc.handleDrilling(); err != nil {
    return err
   }
   pc.resetState()
  }

  time.Sleep(800 * time.Millisecond)
 }
 return nil
}

// Активация вращения карусели
func (pc *ProcController) activateCarousel() error {
 return pc.writeBool(pc.nodes.outputCarousel, true)
}

// Обработка этапа определения цвета
func (pc *ProcController) handleColorDetection() error {
 // Остановка карусели для проведения измерений
 if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
  return err
 }
 time.Sleep(1 * time.Second) // Ожидание стабилизации

 // Определение цвета детали
 colorVal, err := pc.readColor()
 if err != nil {
  return err
 }

 // Активация соответствующей метки
 if err := pc.
markColor(colorVal); err != nil {
  return err
 }

 // Проверка наличия отверстия
 return pc.detectHole()
}

// Чтение показаний цветовых датчиков
func (pc *ProcController) readColor() (string, error) {
 resp, err := pc.client.Read(&ua.ReadRequest{
  NodesToRead: []*ua.ReadValueID{
   {NodeID: pc.nodes.colorSensor},
   {NodeID: pc.nodes.silverSensor},
  },
 })
 if err != nil || len(resp.Results) < 2 {
  return "", fmt.Errorf("ошибка чтения датчиков: %w", err)
 }

 colorActive := resp.Results[0].Value.Value().(bool)
 silverActive := resp.Results[1].Value.Value().(bool)

 // Логика определения цвета по комбинации датчиков
 switch {
 case colorActive && silverActive:
  return "silver", nil
 case colorActive:
  return "red", nil
 default:
  return "black", nil
 }
}

// Активация метки в зависимости от цвета
func (pc *ProcController) markColor(color string) error {
 var targetNode *ua.NodeID
 switch color {
 case "silver":
  targetNode = pc.nodes.silverTag
 case "red":
  targetNode = pc.nodes.redTag
 default:
  targetNode = pc.nodes.blackTag
 }
 return pc.writeBool(targetNode, true)
}

// Процедура проверки наличия отверстия
func (pc *ProcController) detectHole() error {
 // Активация цилиндра для проверки
 if err := pc.writeBool(pc.nodes.outputM5Cylinder, true); err != nil {
  return err
 }
 time.Sleep(500 * time.Millisecond) // Время на срабатывание

 // Чтение состояния датчика
 val, err := pc.readBool(pc.nodes.m5Sensor)
 if err != nil {
  return err
 }
 pc.holeDetected = val

 // Деактивация цилиндра
 return pc.writeBool(pc.nodes.outputM5Cylinder, false)
}

// Обработка этапа сверления
func (pc *ProcController) handleDrilling() error {
 if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
  return err
 }
 time.Sleep(1 * time.Second) // Ожидание остановки

 if pc.holeDetected {
  if err := pc.performDrilling(); err != nil {
   return err
  }
 }

 // Сброс состояния карусели
 return pc.resetCarousel()
}

// Последовательность операций сверления
func (pc *ProcController) performDrilling() error {
 // Фиксация детали
 if err := pc.writeBool(pc.nodes.outputM4Clamp, true); err != nil {
  return err
 }
 time.Sleep(100 * time.Millisecond) // Время на зажим

 // Последовательность операций сверления
 if err := pc.drillDown(); err != nil {
  return err
 }
 if err := pc.activateDrill(); err != nil {
  return err
 }
 return pc.drillUp()
}

// Управление положением сверла
func (pc *ProcController) drillDown() error {
 if err := pc.writeBool(pc.nodes.outputDrillDown, true); err != nil {
  return err
 }
 time.Sleep(600 * time.Millisecond) // Время опускания
 return pc.writeBool(pc.nodes.outputDrillDown, false)
}

// Активация вращения сверла
func (pc *ProcController) activateDrill() error {
 if err := pc.writeBool(pc.nodes.outputDrill, true); err != nil {
  return err
 }
 time.Sleep(1 * time.Second) // Время сверления
 return pc.writeBool(pc.nodes.outputDrill, false)
}

// Подъем сверла
func (pc *ProcController) drillUp() error {
 if err := pc.writeBool(pc.nodes.outputDrillUp, true); err != nil {
  return err
 }
 time.Sleep(600 * time.Millisecond) // Время подъема
 return pc.writeBool(pc.nodes.outputDrillUp, false)
}

// Сброс карусели в исходное состояние
func (pc *ProcController) resetCarousel() error {
 if err := pc.writeBool(pc.nodes.outputCarousel, true); err != nil {
  return err
 }
 time.Sleep(200 * time.Millisecond) // Кратковременная активация
 return pc.writeBool(pc.nodes.outputCarousel, false)
}

// Сброс состояния контроллера
func (pc *ProcController) resetState() {
 pc.counter = -1
 pc.finished = true
}

// Вспомогательные методы работы с OPC UA

// Запись булевого значения в узел
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

// Чтение булевого значения из узла
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

 fmt.Println("Технологический процесс завершен")
}