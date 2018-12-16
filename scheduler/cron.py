#!/usr/bin/env python3

import logging
import os
import scheduler
import sys

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)
    rabbitMqUrl = os.environ.get('RABBITMQ_URL')

    if rabbitMqUrl is None:
        print("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.", file=sys.stderr)
        sys.exit(1)

    scheduler = scheduler.Scheduler(rabbitMqUrl)
    scheduler.run()
