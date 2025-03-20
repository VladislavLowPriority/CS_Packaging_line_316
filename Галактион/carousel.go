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
	// 状态变量
	counter      int  //流程步骤计数器
	finished     bool //总流程完成标志
	holeDetected bool //孔位检测结果
	m5Status     bool //暂存 M5 气缸的实时状态

	// OPC UA 客户端配置
	client   *opcua.Client
	endpoint string

	// 节点ID集合
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
	ctx := context.Background()                                                            //创建一个空的上下文对象，用于后续的客户端连接操作。
	client := opcua.NewClient(pc.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone)) //pc.endpoint是OPC UA服务器的地址	opcua.SecurityMode(ua.MessageSecurityModeNone)指定了连接的安全模式为无加密
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	pc.client = client
	return nil
}

func (pc *ProcController) Run() error {
	defer pc.client.Close() //使用defer关键字确保在Run方法执行完毕后（无论是正常结束还是因错误退出），调用pc.client.Close方法关闭OPC UA客户端连接。

	// 初始化转台
	if err := pc.activateCarousel(); err != nil {
		return err
	}

	for !pc.finished {
		if pc.counter < 4 {
			pc.counter++
		}

		switch {
		case pc.counter == 4:
			if err := pc.handleColorDetection(); err != nil { //处理颜色检测
				return err
			}
			pc.counter++

		case pc.counter == 5:
			if err := pc.handleDrilling(); err != nil { //处理钻孔操作
				return err
			}
			pc.resetState() //重置状态
		}

		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

func (pc *ProcController) activateCarousel() error {
	return pc.writeBool(pc.nodes.outputCarousel, true)
}

func (pc *ProcController) handleColorDetection() error {
	// 停止转台
	if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	// 读取颜色传感器
	colorVal, err := pc.readColor()
	if err != nil {
		return err
	}

	// 执行颜色标记
	if err := pc.markColor(colorVal); err != nil {
		return err
	}

	// M5气缸检测
	return pc.detectHole()
}

func (pc *ProcController) readColor() (string, error) {

	// 读取双传感器状态
	resp, err := pc.client.Read(&ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: pc.nodes.colorSensor},
			{NodeID: pc.nodes.silverSensor},
		},
	})
	if err != nil || len(resp.Results) < 2 {
		return "", fmt.Errorf("传感器读取失败: %w", err)
	}

	colorActive := resp.Results[0].Value.Value().(bool)
	silverActive := resp.Results[1].Value.Value().(bool)

	switch {
	case colorActive && silverActive:
		return "silver", nil
	case colorActive:
		return "red", nil
	default:
		return "black", nil
	}
}

func (pc *ProcController) markColor(color string) error {
	var targetNode *ua.NodeID
	switch color {
	case "silver":
		targetNode = pc.nodes.silver_tag
	case "red":
		targetNode = pc.nodes.red_tag
	default:
		targetNode = pc.nodes.black_tag
	}
	return pc.writeBool(targetNode, true)
}

func (pc *ProcController) detectHole() error { //检测孔
	// 下压气缸
	if err := pc.writeBool(pc.nodes.outputM5Cylinder, true); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	// 读取传感器
	val, err := pc.readBool(pc.nodes.m5Sensor)
	if err != nil {
		return err
	}
	pc.holeDetected = val

	// 升起气缸
	return pc.writeBool(pc.nodes.outputM5Cylinder, false)
}

func (pc *ProcController) handleDrilling() error { //处理钻孔操作
	// 停止转台
	if err := pc.writeBool(pc.nodes.outputCarousel, false); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if pc.holeDetected {
		if err := pc.performDrilling(); err != nil {
			return err
		}
	}

	// 复位转台
	return pc.resetCarousel()
}

func (pc *ProcController) performDrilling() error {
	// 夹紧夹具
	if err := pc.writeBool(pc.nodes.outputM4Clamp, true); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	// 下钻操作
	if err := pc.drillDown(); err != nil {
		return err
	}

	// 钻孔操作
	if err := pc.activateDrill(); err != nil {
		return err
	}

	// 抬钻操作
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

// 通用读写方法
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
	//controller := NewProcController("opc.tcp://10.160.160.61:4840)

	if err := controller.Connect(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	if err := controller.Run(); err != nil {
		log.Fatalf("运行错误: %v", err)
	}

	fmt.Println("工艺流程执行完毕")
}
