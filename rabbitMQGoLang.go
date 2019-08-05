package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	go server()
	go client()
	var hold string
	fmt.Scanln(&hold)
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "failed to register a consumer")
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("rabbitMQ"),
	}
	ch.Publish("", q.Name, false, false, msg)
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://test:test@localhost:5672")
	failOnError(err, "failed to connect to rabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	q, err := ch.QueueDeclare("rabbit", false, false, false, false, nil)
	failOnError(err, "failes to declare a queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}
