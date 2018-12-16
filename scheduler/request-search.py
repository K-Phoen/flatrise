#!/usr/bin/env python3

import logging
import os
import scheduler
import sys

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: request-search.py [worker]", file=sys.stderr)
        sys.exit(1)

    logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)

    engine = sys.argv[1]
    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    scheduler = scheduler.Scheduler(rabbitMqUrl)
    scheduler.schedule_search(engine)
