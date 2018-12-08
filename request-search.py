#!/usr/bin/env python3

import os
import pika
import sys

REQUESTS_QUEUES = {
    'leboncoin': 'leboncoin_searchs',
    'boligportal': 'boligportal_searchs',
}

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: request-search.py [worker]", file=sys.stderr)
        sys.exit(1)

    engine = sys.argv[1]

    if not engine in REQUESTS_QUEUES:
        print("Unknown engine '%s'. Available engines are: %s" % (engine, ', '.join(REQUESTS_QUEUES.keys())), file=sys.stderr)
        sys.exit(1)

    search_queue = REQUESTS_QUEUES[engine]

    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    rabbitmq = pika.BlockingConnection(pika.URLParameters(rabbitMqUrl))
    channel = rabbitmq.channel()
    channel.queue_declare(queue=search_queue, durable=True)

    print('Requesting search for %s… ' % engine, end='')

    channel.basic_publish(exchange='', routing_key=search_queue, body='irrelevant', properties=pika.BasicProperties(
        delivery_mode = 2, # make message persistent
    ))

    print('Done')

    rabbitmq.close()
