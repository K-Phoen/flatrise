Flatrise
========

Quick and dirty data analysis pipeline, built to help me search for a flat in
cities I know almost nothing about.

The idea is to crawl the main classified ads websites and aggregate their data
somewhere I can visualize it to get an idea of the prices, to look for nice
neighbourhoods, …

In order to do that, each worker/crawler (one per website) will send data to a
RabbitMq instance, which will be relayed into a Kibana instance through
Logstash.

Please keep in mind that this really is a quick and dirty project, more intended
to help me than to be production-ready.
