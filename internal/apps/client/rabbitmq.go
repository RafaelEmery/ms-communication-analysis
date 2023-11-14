package client

import (
	"encoding/json"
	"log"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/streadway/amqp"
)

type Message struct {
	Content  string            `json:"content"`
	SentAt   time.Time         `json:"sent_at"`
	Metadata map[string]string `json:"metadata"`
}

// TODO: some dependencies on rabbitmq approach
// - Get consume message error
// - Async paradigm
func HandleRabbitMQ(ch *amqp.Channel, q amqp.Queue, data InteractionData) {
	for i := 0; i < data.RequestQuantity; i++ {
		jsonBody, err := getMessageBody(data.Resource)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		if err := ch.Publish(
			"",     // Exchange
			q.Name, // Routing key (nome da fila)
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonBody,
				Headers: amqp.Table{
					"resource": data.Resource,
				},
			}); err != nil {
			log.Default().Println(err)
			continue
		}
	}
}

func getMessageBody(resource string) ([]byte, error) {
	strContent := ""
	if resource == createResource {
		var product domain.Product

		content, err := json.Marshal(product.Fake())
		if err != nil {
			return []byte{}, err
		}

		strContent = string(content)
	}

	m := Message{
		Content: strContent,
		SentAt:  time.Now(),
		Metadata: map[string]string{
			"resource": resource,
		},
	}

	body, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
