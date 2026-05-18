package example

// ======================================================
// Contoh pemakaian package kafkautil (pkg/kafkautil)
// ======================================================
//
// Package kafkautil nyediain Kafka producer berbasis franz-go.
//
// Fitur yang di-demo:
//   - NewProducer — buat producer dari config struct
//   - NewProducerFromYAML — buat producer dari file YAML
//   - NewProducerFromYAMLBytes — buat producer dari bytes YAML
//   - LoadConfigFromYAML — load konfigurasi dari file
//   - ParseYAMLConfig — parse YAML bytes
//   - SendMessage — kirim satu pesan (sync)
//   - SendMessages — kirim banyak pesan sekaligus (batch, sync)
//   - SendJSON — kirim struct/map sebagai JSON
//   - SendJSONWithHeaders — kirim JSON dengan header kustom
//   - SendString — kirim string biasa
//   - SendBytes — kirim raw bytes
//   - SendAsync — kirim tanpa nunggu konfirmasi (async + callback)
//   - SendToPartition — kirim ke partisi tertentu
//   - Flush — tunggu sampai semua pesan terkirim
//   - Close — tutup koneksi producer
//   - GetStats — statistik producer
//   - GetConfig — ambil konfigurasi aktif
//   - GetTopic / GetBrokers — helper getter
//   - SaveConfigToYAML — simpan konfigurasi ke file YAML
//   - PrintConfig — cetak konfigurasi ke stdout
//
// Note: Fitur ini butuh Kafka broker yang beneran untuk jalan.

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/kafkautil"
	"github.com/boni-fm/go-libsd3/pkg/log"
)

// ContohNewProducer mendemonstrasikan pembuatan producer dari config struct.
func ContohNewProducer() {
	logger := log.NewLoggerWithFilename("kafka-demo")

	// Konfigurasi producer
	cfg := &kafkautil.ProducerConfig{
		Brokers: []string{
			"kafka-1:9092",
			"kafka-2:9092",
			"kafka-3:9092",
		},
		Topic: "my-events-topic",
	}

	producer, err := kafkautil.NewProducer(cfg, logger)
	if err != nil {
		fmt.Printf("gagal buat producer: %v\n", err)
		return
	}
	defer producer.Close()

	fmt.Printf("producer siap, topic: %s\n", producer.GetTopic())
	fmt.Printf("brokers: %v\n", producer.GetBrokers())
}

// ContohNewProducerDariYAML mendemonstrasikan pembuatan producer dari file YAML.
//
// Format YAML yang didukung:
//
//	kafka:
//	  brokers:
//	    - "localhost:9092"
//	  topic: "my-topic"
func ContohNewProducerDariYAML() {
	logger := log.NewLoggerWithFilename("kafka-demo")

	// Buat file YAML sementara buat demo
	yamlContent := `
kafka:
  brokers:
    - "localhost:9092"
    - "localhost:9093"
  topic: "order-events"
`
	tmpFile, err := os.CreateTemp("", "kafka_*.yaml")
	if err != nil {
		fmt.Printf("gagal buat temp: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(yamlContent)
	tmpFile.Close()

	// ─── LoadConfigFromYAML ───
	cfg, err := kafkautil.LoadConfigFromYAML(tmpFile.Name())
	if err != nil {
		fmt.Printf("LoadConfigFromYAML gagal: %v\n", err)
		return
	}
	fmt.Printf("config loaded: topic=%s brokers=%v\n", cfg.Topic, cfg.Brokers)

	// ─── ParseYAMLConfig dari bytes ───
	yamlBytes := []byte(yamlContent)
	cfg2, err := kafkautil.ParseYAMLConfig(yamlBytes)
	if err != nil {
		fmt.Printf("ParseYAMLConfig gagal: %v\n", err)
		return
	}
	fmt.Printf("config dari bytes: topic=%s\n", cfg2.Topic)

	// ─── NewProducerFromYAML ───
	// producer, err := kafkautil.NewProducerFromYAML(tmpFile.Name(), logger)
	// jika Kafka beneran ada, uncomment baris di atas

	// ─── NewProducerFromYAMLBytes ───
	// producer, err := kafkautil.NewProducerFromYAMLBytes(yamlBytes, logger)

	_ = logger // buat demo tanpa Kafka beneran
}

// ContohSendMessage mendemonstrasikan pengiriman pesan ke Kafka.
func ContohSendMessage(producer *kafkautil.Producer) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ─── SendMessage — kirim satu pesan ───
	msg := &kafkautil.Message{
		Key:   []byte("user-123"),
		Value: []byte(`{"event":"login","user_id":123}`),
		Headers: map[string]string{
			"source":      "auth-service",
			"content-type": "application/json",
		},
	}
	if err := producer.SendMessage(ctx, msg); err != nil {
		fmt.Printf("SendMessage gagal: %v\n", err)
		return
	}
	fmt.Println("pesan berhasil dikirim")

	// ─── SendString — kirim string biasa ───
	if err := producer.SendString(ctx, "key-001", "hello dari Go!"); err != nil {
		fmt.Printf("SendString gagal: %v\n", err)
	} else {
		fmt.Println("string berhasil dikirim")
	}

	// ─── SendBytes — kirim raw bytes ───
	rawData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" dalam hex
	if err := producer.SendBytes(ctx, []byte("raw-key"), rawData); err != nil {
		fmt.Printf("SendBytes gagal: %v\n", err)
	} else {
		fmt.Println("raw bytes berhasil dikirim")
	}
}

// ContohSendJSON mendemonstrasikan pengiriman data sebagai JSON.
func ContohSendJSON(producer *kafkautil.Producer) {
	ctx := context.Background()

	// Struct yang akan di-marshal ke JSON
	type OrderEvent struct {
		OrderID    string    `json:"order_id"`
		CustomerID int       `json:"customer_id"`
		Total      float64   `json:"total"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
	}

	event := OrderEvent{
		OrderID:    "ORD-20240115-001",
		CustomerID: 42,
		Total:      150000.50,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
	}

	// ─── SendJSON ─── otomatis marshal ke JSON + set header content-type
	if err := producer.SendJSON(ctx, event.OrderID, event); err != nil {
		fmt.Printf("SendJSON gagal: %v\n", err)
		return
	}
	fmt.Println("order event berhasil dikirim sebagai JSON")

	// ─── SendJSONWithHeaders ─── JSON dengan header kustom tambahan
	customHeaders := map[string]string{
		"trace-id":   "trace-xyz-789",
		"service":    "order-service",
		"version":    "v2",
		"x-priority": "high",
	}
	if err := producer.SendJSONWithHeaders(ctx, event.OrderID, event, customHeaders); err != nil {
		fmt.Printf("SendJSONWithHeaders gagal: %v\n", err)
	} else {
		fmt.Println("order event dengan header kustom berhasil dikirim")
	}
}

// ContohSendBatch mendemonstrasikan pengiriman banyak pesan sekaligus.
func ContohSendBatch(producer *kafkautil.Producer) {
	ctx := context.Background()

	// ─── SendMessages — batch kirim, satu round-trip ───
	messages := []*kafkautil.Message{
		{
			Key:   []byte("batch-1"),
			Value: []byte(`{"id":1,"nama":"Budi"}`),
		},
		{
			Key:   []byte("batch-2"),
			Value: []byte(`{"id":2,"nama":"Siti"}`),
			Headers: map[string]string{"batch": "true"},
		},
		{
			Key:   []byte("batch-3"),
			Value: []byte(`{"id":3,"nama":"Ahmad"}`),
		},
	}

	if err := producer.SendMessages(ctx, messages); err != nil {
		fmt.Printf("SendMessages gagal: %v\n", err)
		return
	}
	fmt.Printf("batch %d pesan berhasil dikirim\n", len(messages))
}

// ContohSendAsync mendemonstrasikan pengiriman asynchronous.
// Fire-and-forget — ga nunggu konfirmasi, hasilnya dikasih ke callback.
func ContohSendAsync(producer *kafkautil.Producer) {
	ctx := context.Background()

	msg := &kafkautil.Message{
		Key:   []byte("async-key"),
		Value: []byte(`{"event":"user_viewed_product","product_id":99}`),
	}

	// ─── SendAsync — non-blocking, callback dipanggil setelah ack/error ───
	producer.SendAsync(ctx, msg, func(err error) {
		if err != nil {
			fmt.Printf("async send gagal: %v\n", err)
		} else {
			fmt.Println("async send sukses (callback)")
		}
	})

	// SendAsync tanpa callback (benar-benar fire-and-forget)
	producer.SendAsync(ctx, msg, nil)

	fmt.Println("async send dijadwalkan, ga nunggu konfirmasi")
}

// ContohSendToPartition mendemonstrasikan pengiriman ke partisi tertentu.
func ContohSendToPartition(producer *kafkautil.Producer) {
	ctx := context.Background()

	msg := &kafkautil.Message{
		Key:   []byte("partition-key"),
		Value: []byte(`{"region":"jawa","data":"penting"}`),
	}

	// Kirim ke partisi 0 (misalnya buat data Jawa)
	if err := producer.SendToPartition(ctx, msg, 0); err != nil {
		fmt.Printf("SendToPartition gagal: %v\n", err)
		return
	}
	fmt.Println("pesan berhasil dikirim ke partisi 0")

	// Kirim ke partisi lain
	msg2 := &kafkautil.Message{
		Key:   []byte("sumatra-key"),
		Value: []byte(`{"region":"sumatra","data":"lain"}`),
	}
	if err := producer.SendToPartition(ctx, msg2, 1); err != nil {
		fmt.Printf("SendToPartition partisi 1 gagal: %v\n", err)
	} else {
		fmt.Println("pesan berhasil dikirim ke partisi 1")
	}
}

// ContohFlushDanClose mendemonstrasikan Flush, Close, GetStats, GetConfig.
func ContohFlushDanClose(producer *kafkautil.Producer) {
	ctx := context.Background()

	// ─── Flush — tunggu semua pesan yang di-buffer sampai terkirim ───
	if err := producer.Flush(ctx); err != nil {
		fmt.Printf("Flush gagal: %v\n", err)
	} else {
		fmt.Println("semua pesan di-buffer sudah terkirim")
	}

	// ─── GetStats — ambil statistik producer ───
	stats := producer.GetStats()
	fmt.Printf("stats: messages=%d bytes=%d errors=%d\n",
		stats.Messages, stats.Bytes, stats.Errors)

	// ─── GetConfig — ambil konfigurasi aktif ───
	cfg := producer.GetConfig()
	fmt.Printf("config aktif: topic=%s brokers=%v\n", cfg.Topic, cfg.Brokers)

	// ─── GetTopic / GetBrokers ───
	fmt.Println("topic:", producer.GetTopic())
	fmt.Println("brokers:", producer.GetBrokers())

	// ─── PrintConfig — cetak konfigurasi ke stdout ───
	producer.PrintConfig()

	// ─── SaveConfigToYAML — simpan konfigurasi ke file ───
	if err := producer.SaveConfigToYAML("/tmp/kafka_config_backup.yaml"); err != nil {
		fmt.Printf("SaveConfigToYAML gagal: %v\n", err)
	} else {
		fmt.Println("konfigurasi berhasil disimpan ke YAML")
		defer os.Remove("/tmp/kafka_config_backup.yaml")
	}

	// ─── Close — tutup koneksi, pastikan selalu dipanggil ───
	// producer.Close()
	// (ga dipanggil di sini karena udah ada defer di fungsi pemanggil)
}

// ContohProducerLengkap mendemonstrasikan penggunaan producer dari awal sampai akhir.
func ContohProducerLengkap() {
	logger := log.NewLoggerWithFilename("kafka-full-demo")

	cfg := &kafkautil.ProducerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "demo-topic",
	}

	producer, err := kafkautil.NewProducer(cfg, logger)
	if err != nil {
		fmt.Printf("gagal buat producer: %v\n", err)
		return
	}
	defer producer.Close() // selalu tutup di akhir

	ctx := context.Background()

	// Kirim berbagai jenis pesan
	producer.SendString(ctx, "greet", "Halo Kafka!")
	producer.SendJSON(ctx, "user-1", map[string]interface{}{
		"user_id": 1,
		"action":  "login",
	})
	producer.SendBytes(ctx, []byte("bin-key"), []byte{0x01, 0x02, 0x03})

	// Flush sebelum shutdown biar ga ada yang nyasar
	producer.Flush(ctx)
	fmt.Println("semua pesan selesai dikirim")
}
