FROM docker.elastic.co/beats/metricbeat:6.6.1

USER root

RUN groupadd docker --gid 974 ; usermod --append --group 974 metricbeat

COPY ./metricbeat/metricbeat.yml /usr/share/metricbeat/metricbeat.yml
COPY ./metricbeat/modules.d /usr/share/metricbeat/modules.d

USER metricbeat
