package scheduler

import (
	cron "github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type goCronScheduler struct {
	rabbitMqUrl string
	queues      map[string]string
}

func NewGoCronScheduler(rabbitMqUrl string, queues map[string]string) goCronScheduler {
	return goCronScheduler{
		rabbitMqUrl: rabbitMqUrl,
		queues:      queues,
	}
}

func (s goCronScheduler) Run() {
	gocron := cron.NewScheduler()
	gocron.Every(1).Hour().Do(s.SearchAll)

	<-gocron.Start()
}

func (s goCronScheduler) SearchAll() error {
	for engine := range s.queues {
		err := s.Search(engine)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s goCronScheduler) Search(engine string) error {
	queue, exists := s.queues[engine]
	if !exists {
		return errors.New("Unknown engine")
	}

	return s.withRabbitMq(func(channel *amqp.Channel) error {
		err := s.declareQueues(channel)
		if err != nil {
			return err
		}

		err = channel.Publish(
			"",    // exchange
			queue, // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte("irrelevant"),
			},
		)
		if err != nil {
			return errors.Wrap(err, "Could not schedule search")
		}

		return nil
	})
}

func (s goCronScheduler) withRabbitMq(do func(*amqp.Channel) error) error {
	conn, err := amqp.Dial(s.rabbitMqUrl)
	if err != nil {
		return errors.Wrap(err, "Could not connect to RabbitMq")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		errors.Wrap(err, "Could not open a RabbitMq channel")
	}
	defer channel.Close()

	return do(channel)
}

func (s goCronScheduler) declareQueues(channel *amqp.Channel) error {
	for _, queue := range s.queues {
		_, err := channel.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return errors.Wrap(err, "Could not declare queue")
		}
	}

	return nil
}
