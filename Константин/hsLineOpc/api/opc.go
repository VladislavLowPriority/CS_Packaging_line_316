package api

import (
	"errors"
	"io"
	"log"
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

	client.inputTagMap["ns=4;i=9"] = client.GetNodeValue("ns=4;i=9")
	client.inputTagMap["ns=4;i=10"] = client.GetNodeValue("ns=4;i=10")
	client.inputTagMap["ns=4;i=24"] = client.GetNodeValue("ns=4;i=24")
	client.inputTagMap["ns=4;i=25"] = client.GetNodeValue("ns=4;i=25")
	client.inputTagMap["ns=4;i=26"] = client.GetNodeValue("ns=4;i=26")

	client.inputTagMap["ns=4;i=3"] = client.GetNodeValue("ns=4;i=3")
	client.inputTagMap["ns=4;i=4"] = client.GetNodeValue("ns=4;i=4")
	client.inputTagMap["ns=4;i=6"] = client.GetNodeValue("ns=4;i=6")
	client.inputTagMap["ns=4;i=7"] = client.GetNodeValue("ns=4;i=7")

	client.inputTagMap["ns=4;i=31"] = client.GetNodeValue("ns=4;i=31")
	client.inputTagMap["ns=4;i=30"] = client.GetNodeValue("ns=4;i=30")
	client.inputTagMap["ns=4;i=32"] = client.GetNodeValue("ns=4;i=32")

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
		log.Fatalf("Write failed: %s", err)
	}

	return nil
}

// отправляет false на все теги для остановки установки
func (c *OpcClient) SendAllFalses() {
	for key := range c.outputTagMap {
		c.WriteNodeValue(key, false)
	}
}
