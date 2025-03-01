package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	fail(err, "Error Connecting")
	defer conn.Close()

	ch, err := conn.Channel()
	fail(err, "Error opening channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"YOO",
		false,
		false,
		false,
		false,
		nil,
	)
	fail(err, "Error declaring queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("[x] To exit press CTRL+C")
		fmt.Print("MESSAGE: ")
		body, _ := reader.ReadString('\n')
		body = strings.TrimSpace(body)
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		fail(err, "Error publishing message")
		log.Printf("[x] sent %s\n", body)

	}

}

func fail(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
