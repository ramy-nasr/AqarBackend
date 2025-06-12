package consumer

import (
	"context"
	"encoding/json"
	"transaction-backend/application"
	"transaction-backend/domain"

	"github.com/streadway/amqp"
)

func StartRabbitConsumer(amqpURL, queue, exchange, routingKey string, service *application.TransactionService) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	_ = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var txn domain.Transaction
			if err := json.Unmarshal(d.Body, &txn); err == nil {
				_ = service.HandleNewTransaction(context.Background(), txn)
			}
		}
	}()

	return nil
}
