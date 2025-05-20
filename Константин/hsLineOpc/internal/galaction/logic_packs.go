// logic_packs.go
package galaction

import (
	"context"
	"hsLineOpc/api"
	"time"

	"github.com/gopcua/opcua/ua"
)

type PackS struct {
	client *api.OpcClient

	FixBoxTongue    *ua.NodeID
	FixBoxUpperSide *ua.NodeID
	PackBox         *ua.NodeID
}

func NewPackS(client *api.OpcClient) *PackS {
	return &PackS{
		client:          client,
		FixBoxTongue:    MustParseNodeID("ns=4;i=45"),
		FixBoxUpperSide: MustParseNodeID("ns=4;i=44"),
		PackBox:         MustParseNodeID("ns=4;i=46"),
	}
}

func (p *PackS) Start(ctx context.Context) error {
	steps := []struct {
		nodes []*ua.NodeID
		vals  []bool
		delay time.Duration
	}{
		{[]*ua.NodeID{p.FixBoxTongue}, []bool{true}, 500 * time.Millisecond},
		{[]*ua.NodeID{p.FixBoxUpperSide}, []bool{false}, 1 * time.Second},
		{[]*ua.NodeID{p.PackBox}, []bool{true}, 2 * time.Second},
		{[]*ua.NodeID{p.PackBox}, []bool{false}, 500 * time.Millisecond},
		{[]*ua.NodeID{p.FixBoxTongue}, []bool{false}, 0},
	}

	for _, step := range steps {
		if err := p.client.WriteBools(step.nodes, step.vals); err != nil {
			return err
		}
		if step.delay > 0 {
			select {
			case <-time.After(step.delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}
