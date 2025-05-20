package galaction

import (
	"context"
	"fmt"
	"hsLineOpc/api"
	"strconv"
	"time"

	"github.com/gopcua/opcua/ua"
)

// OPCClient 封装客户端实例
type OPCClient api.OpcClient

// NewOPCClient 创建客户端实例（需处理错误）
func NewOPCClient(endpoint string) (*OPCClient, error) {
	client := OPCClient(*api.NewClient(endpoint))

	return &client, nil
}

// Connect 建立连接
func (c *OPCClient) Connect(ctx context.Context) error {
	return c.Client.Connect(ctx)
}

// Close 关闭连接
func (c *OPCClient) Close(ctx context.Context) {
	c.Client.Close()
}

func (c *OPCClient) WriteBools(nodes []*ua.NodeID, values []bool) error {
	if len(nodes) != len(values) {
		return fmt.Errorf("количество узлов и значений должно быть одинаковым")
	}

	start := time.Now()
	writeValues := make([]*ua.WriteValue, len(nodes))
	for i, node := range nodes {
		writeValues[i] = &ua.WriteValue{
			NodeID:      node,
			AttributeID: ua.AttributeIDValue,
			Value: &ua.DataValue{
				EncodingMask: ua.DataValueValue,
				Value:        ua.MustVariant(values[i]),
			},
		}
	}

	_, err := c.Client.Write(&ua.WriteRequest{
		NodesToWrite: writeValues,
	})

	fmt.Printf("Write %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds())
	return err
}

func (c *OPCClient) ReadBools(nodes []*ua.NodeID) ([]bool, error) {
	start := time.Now()
	readValues := make([]*ua.ReadValueID, len(nodes))
	for i, node := range nodes {
		readValues[i] = &ua.ReadValueID{
			NodeID:      node,
			AttributeID: ua.AttributeIDValue,
		}
	}

	resp, err := c.Client.Read(&ua.ReadRequest{
		MaxAge:             1000,
		NodesToRead:        readValues,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	})

	fmt.Printf("Read %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds())

	if err != nil {
		return nil, fmt.Errorf("ошибка чтения значений: %s", err)
	}

	results := make([]bool, len(nodes))
	for i, result := range resp.Results {
		if result.Status != ua.StatusOK {
			return nil, fmt.Errorf("ошибка статуса для узла %v: %v", nodes[i], result.Status)
		}

		// 检查值的类型并转换
		switch v := result.Value.Value().(type) {
		case bool:
			results[i] = v
		case string:
			boolVal, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("узел %v содержит несоответствующее значение строки: %v", nodes[i], v)
			}
			results[i] = boolVal
		default:
			return nil, fmt.Errorf("узел %v содержит несоединимый тип значения: %T", nodes[i], v)
		}
	}

	return results, nil
}

func MustParseNodeID(s string) *ua.NodeID {
	return ua.MustParseNodeID(s)
}
