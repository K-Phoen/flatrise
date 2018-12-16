#!/usr/bin/env python3

import crawler
import json
import logging
import os
import pika
import sys

OFFERS_QUEUE = 'offers'
BLOCKET_SEARCHS_QUEUE = 'blocket_searchs'

def search_blocket(rabbit_channel):
    blocket = crawler.Crawler()

    for offer in blocket.offers():
        logging.info('Found offer "%s" -- %s', offer['title'], offer['identifier'])
        rabbit_channel.basic_publish(exchange='', routing_key=OFFERS_QUEUE, body=json.dumps(offer), properties=pika.BasicProperties(
            delivery_mode = 2, # make message persistent
        ))

def search_request_received(channel, method, properties, body):
    print('Received search request!')

    search_blocket(channel)

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)
    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    rabbitmq = pika.BlockingConnection(pika.URLParameters(rabbitMqUrl))
    channel = rabbitmq.channel()
    channel.queue_declare(queue=OFFERS_QUEUE, durable=True)
    channel.queue_declare(queue=BLOCKET_SEARCHS_QUEUE, durable=True)

    logging.info('Waiting for search requests…')
    channel.basic_consume(search_request_received, queue=BLOCKET_SEARCHS_QUEUE)

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        rabbitmq.close()
