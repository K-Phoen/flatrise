#!/usr/bin/env python3

import json
import pika
from elasticsearch import Elasticsearch

def receive_offer(channel, method_frame, header_frame, offer_payload):
    es = Elasticsearch(['http://elasticsearch'])

    offer = json.loads(offer_payload)

    res = es.index(index="offers", doc_type='offer', id=offer['identifier'], body=offer)

    channel.basic_ack(delivery_tag=method_frame.delivery_tag)

    print(offer, res)

rabbitmq = pika.BlockingConnection(pika.URLParameters('amqp://admin:admin@rabbitmq/flatrise'))
channel = rabbitmq.channel()
channel.queue_declare(queue='offers')

channel.basic_consume(receive_offer, queue='offers')

print(' [*] Waiting for messages. To exit press CTRL+C')

try:
    channel.start_consuming()
except KeyboardInterrupt:
    channel.stop_consuming()

connection.close()
