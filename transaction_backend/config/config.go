package config

var (
	RabbitMQURL = "amqps://fhfjqkfv:DkLtNmo4wftpArdGLesI0CrYTlMhPXqk@toucan.lmq.cloudamqp.com/fhfjqkfv"
	Exchange    = "amq.direct"
	RoutingKey  = "new"
	QueueName   = "AddTransactions"
	PostgresDSN = "postgres://user:pass@db:5432/transactions?sslmode=disable"

)
