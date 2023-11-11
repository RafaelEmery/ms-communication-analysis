package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/streadway/amqp"
)

const (
	createResource        = "create"
	reportResource        = "report"
	getByDiscountResource = "get_by_discount"
)

type Message struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type Consumer struct {
	q   amqp.Queue
	c   Creator
	rg  ReportGenerator
	pdg ProductByDiscountGetter
}

func NewConsumer(q amqp.Queue, c Creator, rg ReportGenerator, pdg ProductByDiscountGetter) Consumer {
	return Consumer{q: q, c: c, rg: rg, pdg: pdg}
}

func (c Consumer) Start(ctx context.Context, ch *amqp.Channel) {
	for {
		msgs, err := ch.Consume(c.q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("could not consume %s", err.Error())
		}

		if len(msgs) == 0 {
			log.Println("no messages")
			continue
		}

		for msg := range msgs {
			var m Message
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				log.Default().Printf("error on reading message %s", err.Error())
				continue
			}

			resource, ok := m.Metadata["resource"]
			if !ok {
				log.Default().Println("skipping message")
				continue
			}

			if err := c.handleUseCases(ctx, resource, m.Content); err != nil {
				log.Default().Printf("error on processing message %s", err.Error())
				continue
			}
		}
	}
}

func (c Consumer) handleUseCases(ctx context.Context, resource, content string) error {
	switch resource {
	case createResource:
		var p domain.Product
		if err := json.Unmarshal([]byte(content), &p); err != nil {
			return err
		}

		o, err := c.c.Create(ctx, p)
		if err != nil {
			return err
		}

		log.Default().Printf("successfully created with id %s", o.ID)
	case reportResource:
		o, err := c.rg.GenerateReport(ctx)
		if err != nil {
			return err
		}

		log.Default().Printf("successfully generated with name %s", o)
	case getByDiscountResource:
		o, err := c.pdg.GetByDiscount(ctx)
		if err != nil {
			return err
		}

		log.Default().Printf("successfully returned with count %d", len(o))
	default:
		return fmt.Errorf("invalid resource %s", resource)
	}

	return nil
}
