FROM nginx:1.15.8-alpine

RUN apk add --no-cache apache2-utils

COPY nginx.conf /etc/nginx/nginx.conf

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

CMD ["/usr/local/bin/entrypoint.sh"]
