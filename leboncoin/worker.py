#!/usr/bin/env python3

import crawler
import json
import logging
import os
import pika
import sys

OFFERS_QUEUE = 'offers'
LEBONCOIN_SEARCHS_QUEUE = 'leboncoin_searchs'

def search_leboncoin(rabbit_channel):
    leboncoin = crawler.Crawler()

    for offer in leboncoin.offers():
        logging.info('Found offer "%s" -- %s', offer['title'], offer['identifier'])
        rabbit_channel.basic_publish(exchange='', routing_key=OFFERS_QUEUE, body=json.dumps(offer), properties=pika.BasicProperties(
            delivery_mode = 2, # make message persistent
        ))

def search_request_received(channel, method, properties, body):
    logging.info('Received search request!')

    search_leboncoin(channel)

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)
    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    rabbitmq = pika.BlockingConnection(pika.URLParameters(rabbitMqUrl))
    channel = rabbitmq.channel()
    channel.queue_declare(queue=OFFERS_QUEUE, durable=True)
    channel.queue_declare(queue=LEBONCOIN_SEARCHS_QUEUE, durable=True)

    logging.info('Waiting for search requests…')
    channel.basic_consume(search_request_received, queue=LEBONCOIN_SEARCHS_QUEUE)

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        rabbitmq.close()
