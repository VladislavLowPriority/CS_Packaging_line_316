// logic_procs.go
package galaction

import (
	"context"
	"hsLineOpc/api"
	"log"
	"time"

	"github.com/gopcua/opcua/ua"
)

type ProcS struct {
	client *api.OpcClient

	// 节点定义
	RedTag         *ua.NodeID
	SilverTag      *ua.NodeID
	BlackTag       *ua.NodeID
	CarouselRotate *ua.NodeID
	M5Tag          *ua.NodeID
	RedAndSilvery  *ua.NodeID
	Silvery        *ua.NodeID
	Drilling       *ua.NodeID
	DrillDown      *ua.NodeID
	DrillUp        *ua.NodeID
	M4Toggle       *ua.NodeID
	M5Toggle       *ua.NodeID

	counter int
}

func NewProcS(client *api.OpcClient) *ProcS {
	return &ProcS{
		client:         client,
		RedTag:         MustParseNodeID("ns=4;i=24"),
		SilverTag:      MustParseNodeID("ns=4;i=26"),
		BlackTag:       MustParseNodeID("ns=4;i=25"),
		CarouselRotate: MustParseNodeID("ns=4;i=13"),
		M5Tag:          MustParseNodeID("ns=4;i=4"),
		RedAndSilvery:  MustParseNodeID("ns=4;i=6"),
		Silvery:        MustParseNodeID("ns=4;i=7"),
		Drilling:       MustParseNodeID("ns=4;i=12"),
		DrillDown:      MustParseNodeID("ns=4;i=14"),
		DrillUp:        MustParseNodeID("ns=4;i=15"),
		M4Toggle:       MustParseNodeID("ns=4;i=16"),
		M5Toggle:       MustParseNodeID("ns=4;i=17"),
		counter:        -1,
	}
}

func (p *ProcS) Start(ctx context.Context) error {
	if err := p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{true}); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			p.counter++
			log.Printf("当前计数: %d", p.counter)

			switch {
			case p.counter == 4:
				time.Sleep(100 * time.Millisecond)
				if err := p.handleColorDetection(ctx); err != nil {
					return err
				}
			case p.counter == 5:
				if err := p.handleDrilling(ctx); err != nil {
					return err
				}
				return nil // 完成流程
			case p.counter > 5:
				p.counter = -1
			}

			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func (p *ProcS) handleColorDetection(ctx context.Context) error {
	// 停止旋转
	if err := p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{false}); err != nil {

		return err
	}
	time.Sleep(1000 * time.Millisecond)
	// 颜色检测
	results, err := p.client.ReadBools([]*ua.NodeID{p.RedAndSilvery, p.Silvery})
	if err != nil {
		return err
	}

	var targetNode *ua.NodeID
	switch {
	case results[0] && results[1]:
		targetNode = p.SilverTag
		log.Println("检测到银色")
	case results[0]:
		targetNode = p.RedTag
		log.Println("检测到红色")
	default:
		targetNode = p.BlackTag
		log.Println("检测到黑色")
	}

	// 设置颜色标签
	if err := p.client.WriteBools([]*ua.NodeID{targetNode}, []bool{true}); err != nil {
		return err
	}

	// M5操作
	if err := p.client.WriteBools([]*ua.NodeID{p.M5Toggle}, []bool{true}); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := p.client.WriteBools([]*ua.NodeID{p.M5Toggle}, []bool{false}); err != nil {
		return err
	}

	// 恢复旋转
	return p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{true})
}

func (p *ProcS) handleDrilling(ctx context.Context) error {
	// 停止旋转
	if err := p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{false}); err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond)
	// 钻孔操作
	steps := []struct {
		nodes []*ua.NodeID
		vals  []bool
		delay time.Duration
	}{
		{[]*ua.NodeID{p.M4Toggle, p.DrillDown}, []bool{true, true}, 600 * time.Millisecond},
		{[]*ua.NodeID{p.DrillDown}, []bool{false}, 0},
		{[]*ua.NodeID{p.Drilling}, []bool{true}, 1 * time.Second},
		{[]*ua.NodeID{p.Drilling}, []bool{false}, 0},
		{[]*ua.NodeID{p.DrillUp}, []bool{true}, 600 * time.Millisecond},
		{[]*ua.NodeID{p.DrillUp, p.M4Toggle}, []bool{false, false}, 0},
	}

	for _, step := range steps {
		if err := p.client.WriteBools(step.nodes, step.vals); err != nil {
			return err
		}
		time.Sleep(step.delay)

	}

	p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{true})
	time.Sleep(200 * time.Millisecond)
	p.client.WriteBools([]*ua.NodeID{p.CarouselRotate}, []bool{false})
	p.counter = -1
	time.Sleep(800 * time.Millisecond)
	return nil
}
