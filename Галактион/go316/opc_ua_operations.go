package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// OPCClient 封装客户端实例
type OPCClient struct {
	client *opcua.Client
}

// NewOPCClient 创建客户端实例（需处理错误）
func NewOPCClient(endpoint string) (*OPCClient, error) {
	opts := []opcua.Option{
		opcua.SecurityMode(ua.MessageSecurityModeNone),
	}
	client, err := opcua.NewClient(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return &OPCClient{client: client}, nil
}

// Connect 建立连接
func (c *OPCClient) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

// Close 关闭连接
func (c *OPCClient) Close(ctx context.Context) {
	c.client.Close(ctx)
}

// // WriteBools 批量写入布尔值
// func (c *OPCClient) WriteBools(ctx context.Context, nodes []*ua.NodeID, values []bool) error {
// 	req := &ua.WriteRequest{
// 		NodesToWrite: make([]*ua.WriteValue, len(nodes)),
// 	}
// 	for i := range nodes {
// 		req.NodesToWrite[i] = &ua.WriteValue{
// 			NodeID:      nodes[i],
// 			AttributeID: ua.AttributeIDValue,
// 			Value: &ua.DataValue{
// 				Value: ua.MustVariant(values[i]),
// 			},
// 		}
// 	}
// 	_, err := c.client.Write(ctx, req)
// 	return err
// }

// func (c *OPCClient) WriteBools(nodes []*ua.NodeID, values []bool) error {
// 	if len(nodes) != len(values) {
// 		return fmt.Errorf("количество узлов и значений должно быть одинаковым")
// 	}

// 	start := time.Now()
// 	writeValues := make([]*ua.WriteValue, len(nodes))
// 	for i, node := range nodes {
// 		writeValues[i] = &ua.WriteValue{
// 			NodeID:      node,
// 			AttributeID: ua.AttributeIDValue,
// 			Value: &ua.DataValue{
// 				EncodingMask: ua.DataValueValue,
// 				Value:        ua.MustVariant(values[i]),
// 			},
// 		}
// 	}

// 	_, err := c.client.Write(context.Background(), &ua.WriteRequest{
// 		NodesToWrite: writeValues,
// 	})

// 	fmt.Printf("Write %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds())
// 	return err
// }

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

	_, err := c.client.Write(context.Background(), &ua.WriteRequest{
		NodesToWrite: writeValues,
	})

	fmt.Printf("Write %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds())
	return err
}

// // ReadBools 批量读取布尔值
// func (c *OPCClient) ReadBools(ctx context.Context, nodes []*ua.NodeID) ([]bool, error) {
// 	req := &ua.ReadRequest{
// 		NodesToRead: make([]*ua.ReadValueID, len(nodes)),
// 	}
// 	for i := range nodes {
// 		req.NodesToRead[i] = &ua.ReadValueID{
// 			NodeID:      nodes[i],
// 			AttributeID: ua.AttributeIDValue,
// 		}
// 	}
// 	resp, err := c.client.Read(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	results := make([]bool, len(resp.Results))
// 	for i, v := range resp.Results {
// 		results[i] = v.Value.Value().(bool)
// 	}
// 	return results, nil
// }

// func (c *OPCClient) ReadBools(nodes []*ua.NodeID) ([]bool, error) {
// 	start := time.Now()
// 	readValues := make([]*ua.ReadValueID, len(nodes))
// 	for i, node := range nodes {
// 		readValues[i] = &ua.ReadValueID{
// 			NodeID:      node,
// 			AttributeID: ua.AttributeIDValue,
// 		}
// 	}

// 	resp, err := c.client.Read(context.Background(), &ua.ReadRequest{
// 		MaxAge:             1000,
// 		NodesToRead:        readValues,
// 		TimestampsToReturn: ua.TimestampsToReturnBoth,
// 	})

// 	fmt.Printf("Read %d values - %d ms\n", len(nodes), time.Since(start).Milliseconds())

// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка чтения значений: %s", err)
// 	}

// 	results := make([]bool, len(nodes))
// 	for i, result := range resp.Results {
// 		if result.Status != ua.StatusOK {
// 			return nil, fmt.Errorf("ошибка статуса для узла %v: %v", nodes[i], result.Status)
// 		}
// 		results[i] = result.Value.Value().(bool)
// 	}

// 	return results, nil
// }

func (c *OPCClient) ReadBools(nodes []*ua.NodeID) ([]bool, error) {
	start := time.Now()
	readValues := make([]*ua.ReadValueID, len(nodes))
	for i, node := range nodes {
		readValues[i] = &ua.ReadValueID{
			NodeID:      node,
			AttributeID: ua.AttributeIDValue,
		}
	}

	resp, err := c.client.Read(context.Background(), &ua.ReadRequest{
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
