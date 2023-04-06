package kafka

import (
	"crypto/tls"
	"math/big"
	"time"

	"github.com/Shopify/sarama"
	"github.com/btc-scan/config"
	"github.com/sirupsen/logrus"
)

type PushKafkaService struct {
	Producer   sarama.SyncProducer
	TopicTx    string
	TopicMatch string
}

func NewSyncProducer(config config.Kafka) (pro sarama.SyncProducer, e error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Timeout = 5 * time.Second
	cfg.Producer.MaxMessageBytes = config.ProducerMax
	//none/gzip/snappy/lz4/ZStandard
	cfg.Producer.Compression = 2
	if config.Password != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = config.User
		cfg.Net.SASL.Password = config.Password
		cfg.Net.SASL.Handshake = true
		cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	}
	if config.Tls {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{}
	}
	brokers := config.Brokers
	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		logrus.Infof("NewSyncProducer error=%v", err)
	}
	return p, err
}

func NewPushKafkaService(config *config.Config, p sarama.SyncProducer) (*PushKafkaService, error) {
	b := &PushKafkaService{}

	b.TopicTx = config.Kafka.TopicTx
	b.TopicMatch = config.Kafka.TopicMatch

	b.Producer = p

	return b, nil
}
func bigIntToStr(v *big.Int) string {
	if v == nil {
		return ""
	}
	return v.String()
}

func (b *PushKafkaService) Pushkafka(value []byte, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	length := msg.Value.Length()
	logrus.Infof("before compension push block data size :%d", length)

	pid, offset, err := b.Producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf("block-send err:%v", err)
	}
	logrus.Infof("send success pid:%v offset:%v\n", pid, offset)
	return err
}
