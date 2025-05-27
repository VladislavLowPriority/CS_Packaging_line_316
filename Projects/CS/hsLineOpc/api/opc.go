package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type OpcClient struct {
	*opcua.Client
	inputTagMap  map[string]bool
	outputTagMap map[string]bool
}

func NewClient(connString string) *OpcClient {
	client := &OpcClient{
		opcua.NewClient(connString, opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*5)),
		make(map[string]bool),
		make(map[string]bool),
	}

	return client
}

func (c *OpcClient) GetNodeValue(nodeId string) bool {
	id, err := ua.ParseNodeID(nodeId)
	if err != nil {
		log.Fatal(err)
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	var resp *ua.ReadResponse
	for {
		resp, err = c.Read(req)
		if err == nil {
			break
		}

		// Following switch contains known errors that can be retried by the user.
		// Best practice is to do it on read operations.
		log.Println(err.Error())
		switch {
		case err == io.EOF && c.State() != opcua.Closed:
			// has to be retried unless user closed the connection
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionIDInvalid):
			// Session is not activated has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionNotActivated):
			// Session is invalid has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSecureChannelIDInvalid):
			// secure channel will be recreated internally.
			time.After(1 * time.Second)
			continue

		default:
			log.Fatalf("Read failed: %s", err)
		}
	}

	if resp != nil && resp.Results[0].Status != ua.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}

	c.inputTagMap[nodeId] = resp.Results[0].Value.Bool()

	return resp.Results[0].Value.Bool()
}

func (c *OpcClient) WriteNodeValue(nodeId string, value bool) error {
	id, err := ua.ParseNodeID(nodeId)
	if err != nil {
		log.Fatal(err)
	}

	v, err := ua.NewVariant(value)
	if err != nil {
		log.Fatalf("invalid value: %v", err)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        v,
				},
			},
		},
	}

	resp, err := c.Write(req)
	if err != nil || resp.Results[0] != ua.StatusOK {
		slog.Error(fmt.Sprintf("Write failed: %s", err))
	}

	c.outputTagMap[nodeId] = value

	return nil
}

// отправляет false на все теги для остановки установки
func (c *OpcClient) SendAllFalses() {
	for key := range c.outputTagMap {
		c.WriteNodeValue(key, false)
	}
}

func (c *OpcClient) WriteBools(nodes []*ua.NodeID, values []bool) error {
	if len(nodes) != len(values) {
		return fmt.Errorf("количество узлов и значений должно быть одинаковым")
	}

	start := time.Now()
	for i, node := range nodes {
		c.WriteNodeValue(node.String(), values[i])
	}
	slog.Info(fmt.Sprintf("Write %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds()))

	return nil
}

func (c *OpcClient) ReadBools(nodes []*ua.NodeID) ([]bool, error) {
	start := time.Now()
	results := make([]bool, 0, len(nodes))
	for _, node := range nodes {
		results = append(results, c.GetNodeValue(node.String()))
	}

	slog.Info(fmt.Sprintf("Read %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds()))

	return results, nil
}
