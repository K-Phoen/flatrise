FROM python:3.7.1-slim

WORKDIR /usr/src/app

RUN pip install 'pika==0.12.0'

COPY ./docker /usr/src/app/docker/
COPY ./leboncoin /usr/src/app/leboncoin/
