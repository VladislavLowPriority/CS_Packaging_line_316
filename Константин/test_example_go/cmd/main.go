package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/razzle131/hs316go/internal/opc"
	"github.com/razzle131/hs316go/internal/tags"
)

const connString = "opc.tcp://10.160.160.61:4840"

func main() {
	ctx := context.TODO()

	client := opcua.NewClient(connString, opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*5))

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Print("test read gripper start val: ")
	opc.GetNodeValue(tags.InputGripperStart, client)
	log.Println()

	log.Print("test read box down val: ")
	opc.GetNodeValue(tags.InputBoxIsDown, client)
	log.Println()

	testWrite(true, false, client)
	testWrite("true", "false", client)
	testWrite(1, 0, client)
}

func testWrite[T any](startVal, endVal T, client *opcua.Client) {
	log.Printf("start example of type %s\n", reflect.TypeOf(startVal))
	defer log.Printf("end example of type %s\n", reflect.TypeOf(startVal))

	opc.WriteNodeValue(tags.OutputConveyorRight, startVal, client)
	time.Sleep(time.Second * 5)
	opc.WriteNodeValue(tags.OutputConveyorRight, endVal, client)
}
