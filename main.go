package main

import (
	"context"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	zkC, err := setupZk(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(zkC testcontainers.Container, ctx context.Context) {
		err := zkC.Terminate(ctx)
		if err != nil {
			log.Fatalf("Failed to terminate container: %v", err.Error())
		}
	}(zkC, ctx)

	port, _ := zkC.MappedPort(ctx, "2181/tcp")
	zkTest(port.Int())
}

func setupZk(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "zookeeper:3.6.2",
		ExposedPorts: []string{"2181/tcp"},
		WaitingFor:   wait.ForListeningPort("2181/tcp"),
	}
	zkC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err)
	}

	return zkC, err
}

func zkTest(mappedPort int) {
	zkAddr := []string{fmt.Sprintf("127.0.0.1:%d", mappedPort)}
	conn, ec, err := zk.Connect(zkAddr, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for e := range ec {
		switch e.Type {
		case zk.EventNodeCreated:
			log.Println("zk.EventNodeCreated")
		case zk.EventNodeDeleted:
			log.Println("zk.EventNodeDeleted")
		case zk.EventNodeDataChanged:
			log.Println("zk.EventNodeDataChanged")
		case zk.EventNodeChildrenChanged:
			log.Println("zk.EventNodeChildrenChanged")
		case zk.EventSession:
			log.Println("zk.EventSession")
		case zk.EventNotWatching:
			log.Println("zk.EventNotWatching")
		}
	}
}
