package kafka

import (
	"testing"

	"github.com/TerrexTech/go-eventstore-models/model"
	"github.com/TerrexTech/go-kafkautils/consumer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe("Producer test - Kafka", func() {

	var (
		kaAdapter     *KafAdapter
		address       string
		responseTopic string
		proxyConsumer *Consumer
	)

	input := make(chan *model.KafkaResponse)

	address = "kafka:9092"
	responseTopic = "KafkaProducerTest"

	BeforeEach(func() {
		kaAdapter = &KafAdapter{
			Address:         []string{address},
			ProducerResChan: input,
			ResponseTopic:   responseTopic,
		}

		config := consumer.Config{
			ConsumerGroup: "test",
			KafkaBrokers:  []string{"localhost:9092"},
			Topics:        []string{"test"},
		}

		proxyconsumer, err := consumer.New(&config)
		if err != nil {
			panic(err)
		}

		// Read Messages

	})

	// AfterEach(func() {
	// 	defer close(input)
	// })

	It("Should receive response after successfully consuming", func() {

		// _, err := CreateProducer(*kaAdapter)
		// log.Println("Testing - does this work")
		// go func() {
		// 	input <- &model.KafkaResponse{
		// 		AggregateID: 91,
		// 		Input:       "BOOOO",
		// 	}
		// }()

		// go func() {
		// 	for msg := proxyConsumer.Messages() {
		// 	  log.Println(msg)
		// 	}
		//   }()

		// Expect(err).ToNot(HaveOccurred())
	})

	// It("Should receive an error", func() {
	// 	kaAdapter = &KafAdapter{
	// 		Address:         []string{""},
	// 		ProducerResChan: input,
	// 		ResponseTopic:   responseTopic,
	// 	}
	// })

})
