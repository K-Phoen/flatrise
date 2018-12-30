#!/usr/bin/env python3

import http.client
import json

class Crawler:
    API_HOST = 'api.leboncoin.fr'
    API_URL = '/finder/search'
    API_KEY = 'ba0c2dad52b3ec'

    def offers(self):
        data = self._query()

        for offer_data in data['ads']:
            yield {
                'identifier': offer_data['url'],
                'title': offer_data['subject'],
                'description': offer_data['body'],
                'price': offer_data['price'][0],
                'currency': 'EUR',
                'price_eur': offer_data['price'][0],
                'location': {
                    'lat': offer_data['location']['lat'],
                    'lon': offer_data['location']['lng'],
                },
                'rooms': int(self._extract_attr(offer_data, 'rooms', 0)),
                'area': int(self._extract_attr(offer_data, 'square', 0)),
            }

    def _extract_attr(self, offer_data, attr_name, default = None):
        attribute_value = [attribute_data['value'] for attribute_data in offer_data['attributes'] if attribute_data['key'] == attr_name]

        try:
            return attribute_value[0]
        except IndexError:
            return default

    def _query(self):
        request_payload = '{"limit":35,"limit_alu":3,"filters":{"category":{"id":"10"},"enums":{"real_estate_type":["2"],"ad_type":["offer"]},"location":{"city_zipcodes":[{"city":"Lyon","label":"Lyon (toute la ville)","departement_id":69, "locationType": "city"}],"regions":["22"]}}}'
        headers = {"Content-type": "application/json", "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0", "api_key": self.API_KEY}

        conn = http.client.HTTPSConnection(self.API_HOST)
        conn.request('POST', self.API_URL, request_payload, headers)
        response = conn.getresponse()

        if response.status != 200:
            raise Exception('Error while calling the API')

        return json.loads(response.read().decode('utf-8'))
