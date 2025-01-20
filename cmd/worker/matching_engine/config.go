package main

var cfg Config

type Config struct {
	App   App   `envPrefix:"APP_"`
	Kafka Kafka `envPrefix:"KAFKA_"`

	TickNum int `env:"TICK_NUM",envDefault:"5"`
}

type App struct {
	Name string `env:"NAME,required"`
	Port string `env:"PORT,required"`
}

type Kafka struct {
	Brokers []string `env:"BROKERS,required"`
	Topic   string   `env:"TOPIC,required"`
}
