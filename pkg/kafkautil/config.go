package kafkautil

import (
	"fmt"
	"os"

	logging "github.com/boni-fm/go-libsd3/pkg/log"
	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the root structure of the YAML configuration
type YAMLConfig struct {
	Kafka ProducerConfig `yaml:"kafka"`
}

// LoadConfigFromYAML loads producer configuration from a YAML file
func LoadConfigFromYAML(filePath string) (*ProducerConfig, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	return ParseYAMLConfig(data)
}

// ParseYAMLConfig parses YAML data and returns a ProducerConfig
func ParseYAMLConfig(data []byte) (*ProducerConfig, error) {
	var yamlConfig YAMLConfig

	// Unmarshal the YAML data
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&yamlConfig.Kafka); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &yamlConfig.Kafka, nil
}

// validateConfig validates the producer configuration
func validateConfig(config *ProducerConfig) error {
	if len(config.Brokers) == 0 {
		return fmt.Errorf("brokers list cannot be empty")
	}

	if config.Topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}

	// Validate broker format (basic check)
	for i, broker := range config.Brokers {
		if broker == "" {
			return fmt.Errorf("broker at index %d cannot be empty", i)
		}
	}

	return nil
}

// NewProducerFromYAML creates a new producer from YAML configuration file
func NewProducerFromYAML(filePath string, logger *logging.Logger) (*Producer, error) {
	config, err := LoadConfigFromYAML(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from YAML: %w", err)
	}

	return NewProducer(config, logger)
}

// NewProducerFromYAMLBytes creates a new producer from YAML bytes
func NewProducerFromYAMLBytes(yamlData []byte, logger *logging.Logger) (*Producer, error) {
	config, err := ParseYAMLConfig(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return NewProducer(config, logger)
}

// SaveConfigToYAML saves the current configuration to a YAML file
func (p *Producer) SaveConfigToYAML(filePath string) error {
	yamlConfig := YAMLConfig{
		Kafka: *p.config,
	}

	data, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filePath, err)
	}

	return nil
}

// PrintConfig prints the current configuration in a readable format
func (p *Producer) PrintConfig() {
	fmt.Println("=== Franz-go Kafka Producer Configuration ===")
	fmt.Printf("Topic: %s\n", p.config.Topic)
	fmt.Printf("Brokers:\n")
	for i, broker := range p.config.Brokers {
		fmt.Printf("  [%d] %s\n", i+1, broker)
	}
	fmt.Println("==============================================")
}
