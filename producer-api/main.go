package main

import (
	"context"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	kafkaUrl := "localhost:9092"
	topic := "test-topic"
	writer := getWriter(kafkaUrl, topic)

	defer writer.Close()

	app := fiber.New()

	app.Get("/:value", func(c *fiber.Ctx) error {
		return c.SendString(producer(writer, c.Params("value")))
	})
	app.Listen(":3000")
}

func getWriter(url, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{url},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func producer(writer *kafka.Writer, message string) string {
	msg := kafka.Message{
		Value: []byte(message),
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return "ok"
}
