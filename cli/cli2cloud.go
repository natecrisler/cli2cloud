package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/leonwind/cli2cloud/service/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func sendPipedMessages(c proto.Cli2CloudClient, ctx context.Context) error {
	stream, err := c.Publish(ctx)
	if err != nil {
		return err
	}

	client := proto.Client{
		Encrypted: false,
		Salt:      nil,
		Iv:        nil,
		Timestamp: nil,
	}

	clientId, err := c.RegisterClient(ctx, &client)
	fmt.Printf("Your client ID: %s\n", clientId.Id)
	fmt.Printf("Share and monitor it live from cli2cloud.com/%s\n\n\n", clientId.Id)
	// Wait 3 seconds for user to copy the client ID
	//time.Sleep(3 * time.Second)

	// TODO: Scan Stderr as well
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		row := scanner.Text()

		// Print original input to client as well
		fmt.Println(row)

		content := proto.PublishRequest{
			Payload:  &proto.Payload{Body: row},
			ClientId: clientId,
		}

		if err := stream.Send(&content); err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Unable to connect to grpc", err)
	}

	client := proto.NewCli2CloudClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := sendPipedMessages(client, ctx); err != nil {
		log.Fatal("Error while sending to server", err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatal("Unable to close connection", err)
	}
}
