#!/usr/bin/env python3

from functools import wraps
import logging
import pika
import schedule
import time

def _needs_rabbitmq(method):
    @wraps(method)
    def _impl(self, *method_args, **method_kwargs):
        self.rabbitmq = pika.BlockingConnection(pika.URLParameters(self.rabbitMqUrl))
        self.channel = self.rabbitmq.channel()

        self._declare_queues()

        try:
            result = method(self, *method_args, **method_kwargs)
        finally:
            self.channel.close()
            self.rabbitmq.close()

        return result

    return _impl

class Scheduler:
    REQUESTS_QUEUES = {
        'leboncoin': 'leboncoin_searchs',
        'boligportal': 'boligportal_searchs',
        'blocket': 'blocket_searchs',
    }

    def __init__(self, rabbitMqUrl):
        self.rabbitMqUrl = rabbitMqUrl

    def run(self):
        schedule.every(1).hour.do(self.schedule_all)

        try:
            while True:
                schedule.run_pending()
                time.sleep(1)
        except KeyboardInterrupt:
            pass

    @_needs_rabbitmq
    def schedule_all(self):
        for engine in Scheduler.REQUESTS_QUEUES:
            self.schedule_search(engine)

    @_needs_rabbitmq
    def schedule_search(self, engine):
        logging.info('Requesting search for %s… ' % engine)

        search_queue = Scheduler.REQUESTS_QUEUES[engine]

        self.channel.basic_publish(exchange='', routing_key=search_queue, body='irrelevant', properties=pika.BasicProperties(
            delivery_mode = 2, # make message persistent
        ))

        logging.info('Search for %s requested.' % engine)

    def _declare_queues(self):
        for queue in Scheduler.REQUESTS_QUEUES.values():
            self.channel.queue_declare(queue=queue, durable=True)
