package main

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	zkTest()
}

func zkTest() {
	zkAddr := []string{"127.0.0.1:2181"}
	conn, _, err := zk.Connect(zkAddr, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		println("Try to watch /serverset")
		// XXX: What happens if the node does not exist, or node is deleted?
		_, _, events, err := conn.ChildrenW("/serverset")
		if err != nil {
			log.Fatal(err)
		}

		/*
			type Event struct {
				Type   EventType
				State  State
				Path   string // For non-session events, the path of the watched node.
				Err    error
				Server string // For connection events
			}
		*/

		for e := range events {
			println("===============")
			println("EventNode: " + e.Path)
			println("State: " + e.State.String())

			if e.Err != nil {
				log.Println("Err: " + e.Err.Error())
			}

			switch e.Type {

			// EventNode
			case zk.EventNodeCreated:
				log.Println("zk.EventNodeCreated")
			case zk.EventNodeDeleted:
				log.Println("zk.EventNodeDeleted")
			case zk.EventNodeDataChanged:
				log.Println("zk.EventNodeDataChanged")
			case zk.EventNodeChildrenChanged:
				log.Println("zk.EventNodeChildrenChanged")

			// etc event
			case zk.EventSession:
				log.Println("zk.EventSession")
			case zk.EventNotWatching:
				log.Println("zk.EventNotWatching")
			}

			children, _, err := conn.Children(e.Path)
			if err != nil {
				log.Println(err)
			}

			for _, child := range children {
				data, _, err := conn.Get(e.Path + "/" + child)
				if err != nil {
					log.Println(err)
				}

				log.Println("Child: " + child + ", Data: " + string(data))
			}
		}
	}
}
