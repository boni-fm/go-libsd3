package kafkautil

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	logging "github.com/boni-fm/go-libsd3/pkg/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

// ProducerConfig holds the configuration for the Kafka producer
type ProducerConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
	Topic   string   `yaml:"topic" json:"topic"`
}

// Producer wraps the franz-go client
type Producer struct {
	client *kgo.Client
	config *ProducerConfig
	logger *log.Logger
}

// Message represents a message to be sent to Kafka
type Message struct {
	Key     []byte            `json:"key,omitempty"`
	Value   []byte            `json:"value"`
	Headers map[string]string `json:"headers,omitempty"`
}

// ProducerStats holds statistics about the producer
type ProducerStats struct {
	Messages        int64         `json:"messages"`
	Bytes           int64         `json:"bytes"`
	Errors          int64         `json:"errors"`
	ProduceTime     time.Duration `json:"produce_time"`
	RecordsBuffered int           `json:"records_buffered"`
}

// NewProducer creates a new Kafka producer with the given configuration
func NewProducer(config *ProducerConfig, logger *logging.Logger) (*Producer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("brokers cannot be empty")
	}

	if config.Topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}

	// Create franz-go client options
	opts := []kgo.Opt{
		kgo.SeedBrokers(config.Brokers...),
		kgo.DefaultProduceTopic(config.Topic),
		kgo.ProducerBatchMaxBytes(1024 * 1024),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RequiredAcks(kgo.LeaderAck()),
		kgo.ClientID("franz-kafka-producer"),
		kgo.ConnIdleTimeout(5 * time.Minute),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	}

	// Create the client
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create franz-go client: %w", err)
	}

	return &Producer{
		client: client,
		config: config,
		logger: log.New(log.Writer(), "[Franz Producer] ", log.LstdFlags),
	}, nil
}

// SendMessage sends a single message to Kafka
func (p *Producer) SendMessage(ctx context.Context, message *Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Create franz-go record
	record := &kgo.Record{
		Topic:     p.config.Topic,
		Key:       message.Key,
		Value:     message.Value,
		Timestamp: time.Now(),
	}

	// Convert headers
	if message.Headers != nil {
		record.Headers = make([]kgo.RecordHeader, 0, len(message.Headers))
		for key, value := range message.Headers {
			record.Headers = append(record.Headers, kgo.RecordHeader{
				Key:   key,
				Value: []byte(value),
			})
		}
	}

	// Send the record and wait for result
	results := p.client.ProduceSync(ctx, record)

	// Check for errors
	if err := results.FirstErr(); err != nil {
		p.logger.Printf("Failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// SendMessages sends multiple messages to Kafka in a batch
func (p *Producer) SendMessages(ctx context.Context, messages []*Message) error {
	if len(messages) == 0 {
		return nil
	}

	// Create records slice
	records := make([]*kgo.Record, len(messages))
	for i, message := range messages {
		if message == nil {
			return fmt.Errorf("message at index %d cannot be nil", i)
		}

		records[i] = &kgo.Record{
			Topic:     p.config.Topic,
			Key:       message.Key,
			Value:     message.Value,
			Timestamp: time.Now(),
		}

		// Convert headers
		if message.Headers != nil {
			records[i].Headers = make([]kgo.RecordHeader, 0, len(message.Headers))
			for key, value := range message.Headers {
				records[i].Headers = append(records[i].Headers, kgo.RecordHeader{
					Key:   key,
					Value: []byte(value),
				})
			}
		}
	}

	// Send all records and wait for results
	results := p.client.ProduceSync(ctx, records...)

	// Check for errors
	if err := results.FirstErr(); err != nil {
		p.logger.Printf("Failed to send %d messages: %v", len(messages), err)
		return fmt.Errorf("failed to send messages: %w", err)
	}

	return nil
}

// SendJSON sends a JSON-encoded message to Kafka
func (p *Producer) SendJSON(ctx context.Context, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	message := &Message{
		Key:   []byte(key),
		Value: jsonData,
		Headers: map[string]string{
			"content-type": "application/json",
			"timestamp":    time.Now().Format(time.RFC3339),
			"producer":     "franz-go",
		},
	}

	return p.SendMessage(ctx, message)
}

// SendJSONWithHeaders sends a JSON-encoded message with custom headers
func (p *Producer) SendJSONWithHeaders(ctx context.Context, key string, data interface{}, headers map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Merge headers with defaults
	allHeaders := map[string]string{
		"content-type": "application/json",
		"timestamp":    time.Now().Format(time.RFC3339),
		"producer":     "franz-go",
	}
	for k, v := range headers {
		allHeaders[k] = v
	}

	message := &Message{
		Key:     []byte(key),
		Value:   jsonData,
		Headers: allHeaders,
	}

	return p.SendMessage(ctx, message)
}

// SendString sends a string message to Kafka
func (p *Producer) SendString(ctx context.Context, key, value string) error {
	message := &Message{
		Key:   []byte(key),
		Value: []byte(value),
		Headers: map[string]string{
			"content-type": "text/plain",
			"timestamp":    time.Now().Format(time.RFC3339),
			"producer":     "franz-go",
		},
	}

	return p.SendMessage(ctx, message)
}

// SendBytes sends raw bytes to Kafka
func (p *Producer) SendBytes(ctx context.Context, key, value []byte) error {
	message := &Message{
		Key:   key,
		Value: value,
		Headers: map[string]string{
			"content-type": "application/octet-stream",
			"timestamp":    time.Now().Format(time.RFC3339),
			"producer":     "franz-go",
		},
	}

	return p.SendMessage(ctx, message)
}

// SendAsync sends a message asynchronously without waiting for confirmation
func (p *Producer) SendAsync(ctx context.Context, message *Message, callback func(error)) {
	if message == nil {
		if callback != nil {
			callback(fmt.Errorf("message cannot be nil"))
		}
		return
	}

	// Create franz-go record
	record := &kgo.Record{
		Topic:     p.config.Topic,
		Key:       message.Key,
		Value:     message.Value,
		Timestamp: time.Now(),
	}

	// Convert headers
	if message.Headers != nil {
		record.Headers = make([]kgo.RecordHeader, 0, len(message.Headers))
		for key, value := range message.Headers {
			record.Headers = append(record.Headers, kgo.RecordHeader{
				Key:   key,
				Value: []byte(value),
			})
		}
	}

	// Send asynchronously
	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			p.logger.Printf("Async send failed: %v", err)
		}
		if callback != nil {
			callback(err)
		}
	})
}

// SendToPartition sends a message to a specific partition
func (p *Producer) SendToPartition(ctx context.Context, message *Message, partition int32) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Create franz-go record with specific partition
	record := &kgo.Record{
		Topic:     p.config.Topic,
		Partition: partition,
		Key:       message.Key,
		Value:     message.Value,
		Timestamp: time.Now(),
	}

	// Convert headers
	if message.Headers != nil {
		record.Headers = make([]kgo.RecordHeader, 0, len(message.Headers))
		for key, value := range message.Headers {
			record.Headers = append(record.Headers, kgo.RecordHeader{
				Key:   key,
				Value: []byte(value),
			})
		}
	}

	// Send the record and wait for result
	results := p.client.ProduceSync(ctx, record)

	// Check for errors
	if err := results.FirstErr(); err != nil {
		p.logger.Printf("Failed to send message to partition %d: %v", partition, err)
		return fmt.Errorf("failed to send message to partition %d: %w", partition, err)
	}

	return nil
}

// Flush waits for all buffered records to be sent
func (p *Producer) Flush(ctx context.Context) error {
	p.client.Flush(ctx)
	return nil
}

// Close closes the producer and releases resources
func (p *Producer) Close() error {
	p.logger.Println("Closing Franz-go producer...")
	p.client.Close()
	return nil
}

// GetStats returns producer statistics
func (p *Producer) GetStats() ProducerStats {
	// Franz-go doesn't expose detailed stats like other libraries
	// We'll return a basic structure - you could extend this by tracking stats yourself
	return ProducerStats{
		Messages:        0, // Would need to track manually
		Bytes:           0, // Would need to track manually
		Errors:          0, // Would need to track manually
		ProduceTime:     0, // Would need to track manually
		RecordsBuffered: 0, // Could get from client metrics if needed
	}
}

// GetConfig returns the current configuration
func (p *Producer) GetConfig() *ProducerConfig {
	return p.config
}

// SetLogger sets a custom logger for the producer
func (p *Producer) SetLogger(logger *log.Logger) {
	p.logger = logger
}

// GetTopic returns the topic this producer is writing to
func (p *Producer) GetTopic() string {
	return p.config.Topic
}

// GetBrokers returns the list of brokers
func (p *Producer) GetBrokers() []string {
	return p.config.Brokers
}
