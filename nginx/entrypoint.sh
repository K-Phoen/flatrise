#!/bin/sh

USER=${USER:-'test'}
PASSWD=${PASSWD:-'changeme'}

htpasswd -bc /etc/nginx/conf.d/nginx.htpasswd $USER $PASSWD

nginx -g "daemon off;"
