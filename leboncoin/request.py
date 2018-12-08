#!/usr/bin/env python3

import os
import pika
import sys

LEBONCOIN_SEARCHS_QUEUE = 'leboncoin_searchs'

if __name__ == '__main__':
    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    rabbitmq = pika.BlockingConnection(pika.URLParameters(rabbitMqUrl))
    channel = rabbitmq.channel()
    channel.queue_declare(queue=LEBONCOIN_SEARCHS_QUEUE, durable=True)

    print('Requesting Leboncoin search… ', end='')

    channel.basic_publish(exchange='', routing_key=LEBONCOIN_SEARCHS_QUEUE, body='irrelevant', properties=pika.BasicProperties(
        delivery_mode = 2, # make message persistent
    ))

    print('Done')

    rabbitmq.close()
