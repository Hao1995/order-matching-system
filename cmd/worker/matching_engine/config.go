package main

var cfg Config

type Config struct {
	App   App   `envPrefix:"APP_"`
	Kafka Kafka `envPrefix:"KAFKA_"`

	TickNum int8 `env:"TICK_NUM" envDefault:"5"`
}

type App struct {
	Name string `env:"NAME,required"`
	Port string `env:"PORT,required"`

	OrderTopic    string `env:"ORDER_TOPIC,required"`
	MatchingTopic string `env:"MATCHING_TOPIC,required"`
}

type Kafka struct {
	Brokers []string `env:"BROKERS,required"`
}
