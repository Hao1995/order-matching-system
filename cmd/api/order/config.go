package main

var cfg Config

type Config struct {
	App   App   `envPrefix:"APP_"`
	Kafka Kafka `envPrefix:"KAFKA_"`
}

type App struct {
	Name string `env:"NAME,required"`
	Port string `env:"PORT,required"`

	OrderTopic string `env:"ORDER_TOPIC,required"`
	Symbol     string `env:"SYMBOL,required"`
}

type Kafka struct {
	Brokers []string `env:"BROKERS,required"`
}
