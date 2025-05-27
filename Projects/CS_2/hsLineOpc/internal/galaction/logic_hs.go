package galaction

import (
	"context"
	"fmt"
	"hsLineOpc/api"

	"time"

	"github.com/gopcua/opcua/ua"
)

type HS struct {
	client *api.OpcClient

	// 输入节点
	gripperStartSensor    *ua.NodeID
	gripperPackSensor     *ua.NodeID
	gripperConveyorSensor *ua.NodeID

	// 输出节点
	gripperToggleUpDown *ua.NodeID
	gripperOpen         *ua.NodeID
	gripperMoveLeft     *ua.NodeID
	gripperMoveRight    *ua.NodeID
	dropPuck            *ua.NodeID
	greenTag            *ua.NodeID
	yellowTag           *ua.NodeID
	pushBox             *ua.NodeID
	fixBoxUpperSide     *ua.NodeID
}

func NewHS(client *api.OpcClient) *HS {
	return &HS{
		client:                client,
		gripperStartSensor:    MustParseNodeID("ns=4;i=31"),
		gripperPackSensor:     MustParseNodeID("ns=4;i=30"),
		gripperConveyorSensor: MustParseNodeID("ns=4;i=32"),
		gripperToggleUpDown:   MustParseNodeID("ns=4;i=39"),
		gripperOpen:           MustParseNodeID("ns=4;i=40"),
		gripperMoveLeft:       MustParseNodeID("ns=4;i=38"),
		gripperMoveRight:      MustParseNodeID("ns=4;i=37"),
		dropPuck:              MustParseNodeID("ns=4;i=41"),
		greenTag:              MustParseNodeID("ns=4;i=34"),
		yellowTag:             MustParseNodeID("ns=4;i=35"),
		pushBox:               MustParseNodeID("ns=4;i=43"),
		fixBoxUpperSide:       MustParseNodeID("ns=4;i=44"),
	}
}

func (h *HS) GrDown(ctx context.Context, duration time.Duration) error {
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperToggleUpDown}, []bool{true}); err != nil {
		return fmt.Errorf("夹爪下降失败: %w", err)
	}
	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *HS) GrUp(ctx context.Context) error {
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperToggleUpDown}, []bool{false}); err != nil {
		return fmt.Errorf("夹爪上升失败: %w", err)
	}
	time.Sleep(1500 * time.Millisecond)
	return nil
}

func (h *HS) GrMovePuckToCarousel(ctx context.Context) error {
	// 初始化输出状态
	if err := h.client.WriteBools([]*ua.NodeID{
		h.yellowTag,
		h.greenTag,
		h.dropPuck,
	}, []bool{false, true, true}); err != nil {
		return err
	}

	time.Sleep(700 * time.Millisecond)
	if err := h.client.WriteBools([]*ua.NodeID{h.dropPuck}, []bool{false}); err != nil {
		return err
	}

	// 抓取流程
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{true}); err != nil {
		return err
	}
	if err := h.GrDown(ctx, 1500*time.Millisecond); err != nil {
		return err
	}
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{false}); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := h.GrUp(ctx); err != nil {
		return err
	}

	// 向左移动
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveLeft}, []bool{true}); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveLeft}, []bool{false}); err != nil {
		return err
	}

	// 放置物体
	if err := h.GrDown(ctx, 2*time.Second); err != nil {
		return err
	}
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{true}); err != nil {
		return err
	}
	return h.GrUp(ctx)
}

func (h *HS) GrMovePuckToPack(ctx context.Context) error {
	// 抓取物体
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{true}); err != nil {
		return err
	}
	if err := h.GrDown(ctx, 3*time.Second); err != nil {
		return err
	}
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{false}); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := h.GrUp(ctx); err != nil {
		return err
	}

	// 向右移动直到传感器触发
	start := time.Now()
	for {
		vals, err := h.client.ReadBools([]*ua.NodeID{h.gripperPackSensor})
		if err != nil {
			return fmt.Errorf("读取包装传感器失败: %w", err)
		}
		if vals[0] {
			break
		}

		if time.Since(start) > 10*time.Second {
			return fmt.Errorf("移动至包装位置超时")
		}

		if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveRight}, []bool{true}); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 推送盒子
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveRight, h.pushBox}, []bool{false, true}); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return h.client.WriteBools([]*ua.NodeID{h.pushBox}, []bool{false})
}

func (h *HS) GrMovePuckToConveyor(ctx context.Context) error {
	// 抓取物体
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{true}); err != nil {
		return err
	}
	if err := h.GrDown(ctx, 2500*time.Millisecond); err != nil {
		return err
	}
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{false}); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	if err := h.GrUp(ctx); err != nil {
		return err
	}

	// 向右移动
	startTime := time.Now()
	for {
		sensorVals, err := h.client.ReadBools([]*ua.NodeID{h.gripperConveyorSensor})
		if err != nil {
			return fmt.Errorf("读取传送带传感器失败: %w", err)
		}

		if sensorVals[0] {
			break
		}

		if time.Since(startTime) > 10*time.Second {
			return fmt.Errorf("移动至传送带超时")
		}

		if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveRight}, []bool{true}); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 放置物体
	if err := h.GrDown(ctx, 2*time.Second); err != nil {
		return err
	}
	if err := h.client.WriteBools([]*ua.NodeID{h.gripperOpen}, []bool{true}); err != nil {
		return err
	}
	time.Sleep(300 * time.Millisecond)
	return h.GrUp(ctx)
}

func (h *HS) GrMoveToStart(ctx context.Context) error {

	// 向左移动直到起始传感器触发
	start := time.Now()

	h.client.WriteBools([]*ua.NodeID{h.gripperMoveRight}, []bool{false})
	for {
		vals, err := h.client.ReadBools([]*ua.NodeID{h.gripperStartSensor})
		if err != nil {
			return fmt.Errorf("读取起始传感器失败: %w", err)
		}
		if vals[0] {
			break
		}

		if time.Since(start) > 15*time.Second {
			return fmt.Errorf("返回起始位置超时")
		}

		if err := h.client.WriteBools([]*ua.NodeID{h.gripperMoveLeft}, []bool{true}); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 停止移动并更新状态灯
	if err := h.client.WriteBools([]*ua.NodeID{
		h.gripperMoveRight,
		h.gripperMoveLeft,
		h.greenTag,
		h.yellowTag,
	}, []bool{false, false, false, true}); err != nil {
		return err
	}
	return nil
}
