FROM docker.elastic.co/beats/metricbeat:6.5.2

USER root

RUN groupadd docker --gid 999 ; usermod --append --group 999 metricbeat

COPY ./metricbeat/modules.d /usr/share/metricbeat/modules.d

USER metricbeat
