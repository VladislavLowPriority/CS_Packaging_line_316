package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/razzle131/hs316go/internal/opc"
	"github.com/razzle131/hs316go/internal/tags"
)

func main() {
	ctx := context.TODO()

	client := opc.NewClient()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// uncomment on usage
	//moveGripperToStart(c, false)

	// for i := 0; i < 10; i++ {
	// 	testGripperSkip(client)
	// }
}

func testGripperSkip(c *opc.MyClient) {
	if !c.GetNodeValue(tags.InputGripperStart) {
		log.Fatal("make gripper move to start")
	}

	c.WriteNodeValue(tags.OutputGripperRight, true)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+1000))
	c.WriteNodeValue(tags.OutputGripperRight, false)

	moveGripperToStart(c, true)
}

// true -> left, false -> right
func moveGripperToStart(c *opc.MyClient, direction bool) error {
	if c.GetNodeValue(tags.InputGripperStart) {
		return errors.New("already at start")
	}

	tag := ""
	if direction {
		tag = tags.OutputGripperLeft
	} else {
		tag = tags.OutputGripperRight
	}

	c.WriteNodeValue(tag, true)
	for !c.GetNodeValue(tags.InputGripperStart) {
		time.Sleep(time.Millisecond)
	}
	c.WriteNodeValue(tag, false)

	return nil
}
