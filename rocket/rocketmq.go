package rocket

import (
	"github.com/artisanhe/tools/conf"
	rocketmq "github.com/apache/rocketmq-client-go/core"
)

type Producer struct {
	GroupID        string
	NameServer     string
	AccessKey      string
	SecretKey      string
	SendMsgTimeout int
	CompressLevel  int
	MaxMessageSize int
	ProducerModel  int
	producer       rocketmq.Producer
}

func (Producer) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{}
}

func (Producer) MarshalDefaults(v interface{}) {
	if producer, ok := v.(*Producer); ok {
		if producer.NameServer == "" {
			producer.NameServer = "rocketmq"
		}
		if producer.GroupID == "" {
			producer.GroupID = "g7pay"
		}
		if producer.ProducerModel == 0 {
			producer.ProducerModel = 1
		}

	}
}

func (p *Producer) Init() {

	pc := &rocketmq.ProducerConfig{
		ClientConfig: rocketmq.ClientConfig{
			GroupID:    p.GroupID,
			NameServer: p.NameServer,
		},
		ProducerModel: rocketmq.ProducerModel(p.ProducerModel),
	}
	if p.AccessKey != "" && p.SecretKey != "" {
		pc.ClientConfig.Credentials = &rocketmq.SessionCredentials{
			AccessKey: p.AccessKey,
			SecretKey: p.SecretKey,
		}
	}

	if p.SendMsgTimeout != 0 {
		pc.SendMsgTimeout = p.SendMsgTimeout
	}

	if p.CompressLevel != 0 {
		pc.CompressLevel = p.CompressLevel
	}

	if p.MaxMessageSize != 0 {
		pc.MaxMessageSize = p.MaxMessageSize
	}

	producer, err := rocketmq.NewProducer(pc)
	if err != nil {
		panic(err)

	}
	p.producer = producer

}

func (p *Producer) Get() rocketmq.Producer {
	return p.producer
}

type PushConsumer struct {
	GroupID       string
	NameServer    string
	AccessKey     string
	SecretKey     string
	MessageModel  int
	ConsumerModel int
	pushConsumer  rocketmq.PushConsumer
}

func (PushConsumer) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{}
}

func (PushConsumer) MarshalDefaults(v interface{}) {
	if pushConsumer, ok := v.(*PushConsumer); ok {
		if pushConsumer.NameServer == "" {
			pushConsumer.NameServer = "rocketmq"
		}
		if pushConsumer.GroupID == "" {
			pushConsumer.GroupID = "g7pay"
		}
		if pushConsumer.ConsumerModel == 0 {
			pushConsumer.ConsumerModel = 1
		}
		if pushConsumer.MessageModel == 0 {
			pushConsumer.MessageModel = 1
		}
	}
}

func (p *PushConsumer) Init() {

	config := &rocketmq.PushConsumerConfig{
		ClientConfig: rocketmq.ClientConfig{
			GroupID:    p.GroupID,
			NameServer: p.NameServer,
		},
		Model:         rocketmq.MessageModel(p.MessageModel),
		ConsumerModel: rocketmq.ConsumerModel(p.ConsumerModel),
	}
	if p.AccessKey != "" && p.SecretKey != "" {
		config.ClientConfig.Credentials = &rocketmq.SessionCredentials{
			AccessKey: p.AccessKey,
			SecretKey: p.SecretKey,
		}
	}
	consumer, err := rocketmq.NewPushConsumer(config)
	if err != nil {
		panic(err)

	}
	p.pushConsumer = consumer
}

func (p *PushConsumer) Get() rocketmq.PushConsumer {
	return p.pushConsumer
}

type PullConsumer struct {
	GroupID    string
	NameServer string
	AccessKey  string
	SecretKey  string

	pullConsumer rocketmq.PullConsumer
}

func (PullConsumer) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{}
}

func (PullConsumer) MarshalDefaults(v interface{}) {
	if pullConsumer, ok := v.(*PullConsumer); ok {
		if pullConsumer.NameServer == "" {
			pullConsumer.NameServer = "rocketmq"
		}
		if pullConsumer.GroupID == "" {
			pullConsumer.GroupID = "g7pay"
		}
	}
}

func (p *PullConsumer) Init() {
	config := &rocketmq.PullConsumerConfig{
		ClientConfig: rocketmq.ClientConfig{
			GroupID:    p.GroupID,
			NameServer: p.NameServer,
		},
	}
	if p.AccessKey != "" && p.SecretKey != "" {
		config.ClientConfig.Credentials = &rocketmq.SessionCredentials{
			AccessKey: p.AccessKey,
			SecretKey: p.SecretKey,
		}
	}

	pullConsumer, err := rocketmq.NewPullConsumer(config)
	if err != nil {
		panic(err)

	}
	p.pullConsumer = pullConsumer
}

func (p *PullConsumer) Get() rocketmq.PullConsumer {
	return p.pullConsumer
}
