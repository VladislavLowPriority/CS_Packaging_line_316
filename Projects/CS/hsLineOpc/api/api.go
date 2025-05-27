package api

import (
	"hsLineOpc/internal/consts"
	"sync"
	"time"
)

// Этот пакет будет опрашивать сервер на TS,
// а так же заниматься отправкой и приемом данных opc сервера

// opcua.NewClient(connString, opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*5))

type TsClient struct {
	Start       bool
	Stop        bool
	BackToStart bool

	Client *OpcClient
	mu     sync.RWMutex
}

func NewTsClient(tsConnString string) *TsClient {
	return &TsClient{
		Start:       false,
		Stop:        false,
		BackToStart: false,

		Client: NewClient(tsConnString),
		mu:     sync.RWMutex{},
	}
}

func (c *TsClient) SubscribeTs() {
	c.subscribeReadTsTags()
	c.subscribeSendTsTags()
}

func (c *TsClient) subscribeSendTsTags() {
	go func() {
		for {
			c.mu.RLock()
			for key, val := range c.Client.inputTagMap {
				tsKey := key[:3] + "1" + key[4:]
				c.Client.WriteNodeValue(tsKey, val)
			}

			for key, val := range c.Client.outputTagMap {
				tsKey := key[:3] + "1" + key[4:]
				c.Client.WriteNodeValue(tsKey, val)
			}
			c.mu.RUnlock()

			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func (c *TsClient) subscribeReadTsTags() {
	go func() {
		for {
			c.mu.Lock()
			c.Start = c.Client.GetNodeValue(consts.TsStartHs)
			c.Stop = c.Client.GetNodeValue(consts.TsStopHs)
			c.BackToStart = c.Client.GetNodeValue(consts.TsBackToStart)
			c.mu.Unlock()

			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.Sleep(time.Millisecond * 100)
}
