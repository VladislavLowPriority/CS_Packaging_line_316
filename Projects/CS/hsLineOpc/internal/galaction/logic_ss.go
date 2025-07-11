// logic_ss.go
package galaction

import (
	"context"
	"hsLineOpc/api"
	"time"

	"github.com/gopcua/opcua/ua"
)

type SS struct {
	client *api.OpcClient

	// Определение узлов OPC UA
	BoxOnConveyorTag  *ua.NodeID
	BoxIsDownTag      *ua.NodeID
	RedTag            *ua.NodeID
	SilverTag         *ua.NodeID
	BlackTag          *ua.NodeID
	MoveConveyorRight *ua.NodeID
	MoveConveyorLeft  *ua.NodeID
	PushSilver        *ua.NodeID
	PushRed           *ua.NodeID
}

func NewSS(client *api.OpcClient) *SS {
	return &SS{
		client:            client,
		BoxOnConveyorTag:  MustParseNodeID("ns=4;i=9"),
		BoxIsDownTag:      MustParseNodeID("ns=4;i=10"),
		RedTag:            MustParseNodeID("ns=4;i=24"),
		SilverTag:         MustParseNodeID("ns=4;i=26"),
		BlackTag:          MustParseNodeID("ns=4;i=25"),
		MoveConveyorRight: MustParseNodeID("ns=4;i=19"),
		MoveConveyorLeft:  MustParseNodeID("ns=4;i=20"),
		PushSilver:        MustParseNodeID("ns=4;i=21"),
		PushRed:           MustParseNodeID("ns=4;i=22"),
	}
}

func (s *SS) Start(ctx context.Context) error {
	// Чтение состояния датчиков
	results, err := s.client.ReadBools([]*ua.NodeID{
		s.BoxOnConveyorTag,
		s.BoxIsDownTag,
		s.RedTag,
		s.SilverTag,
		s.BlackTag,
	})
	if err != nil {
		return err
	}

	boxOnConveyor := results[0]
	boxIsDown := results[1]
	red := results[2]
	silver := results[3]

	if boxOnConveyor {
		// Определение цвета для направления
		var pushNode *ua.NodeID
		switch {
		case red:
			pushNode = s.PushRed
		case silver:
			pushNode = s.PushSilver
		}
		if pushNode != nil {
			if err := s.client.WriteBools([]*ua.NodeID{pushNode}, []bool{true}); err != nil {
				return err
			}
		}

		// Запуск конвейера
		if err := s.client.WriteBools([]*ua.NodeID{s.MoveConveyorRight}, []bool{true}); err != nil {
			return err
		}

		// Ожидание падения коробки
		for !boxIsDown {
			results, err = s.client.ReadBools([]*ua.NodeID{s.BoxIsDownTag})
			if err != nil {
				return err
			}
			boxIsDown = results[0]
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
			}
		}
	}

	// Остановка всех операций
	return s.client.WriteBools([]*ua.NodeID{
		s.MoveConveyorRight,
		s.PushRed,
		s.PushSilver,
	}, []bool{false, false, false})
}
