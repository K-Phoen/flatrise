#!/usr/bin/env python3

from bs4 import BeautifulSoup
import http.client
import json

class LocationNotFound(Exception):
    pass

class Crawler:
    API_HOST = 'www.blocket.se'
    API_LIST_URL = '/karta/items?ca=11&ca=11&st=s&cg=3020&sort=&ps=&pe=&ss=&se=&ros=&roe=&mre=&q=&is=1&f=b&w=3&ac=0MNXXY7CTORXWG23IN5WG2000&zl=12&ne=59.39389826993069%2C18.441925048828125&sw=59.2802650449542%2C17.865142822265625'

    def offers(self):
        for offer_data in self._query_offers():
            try:
                yield self._query_offer_details(offer_data)
            except LocationNotFound:
                pass

    def _query_offer_details(self, offer_data):
        offer_id = 'https://www.blocket.se/stockholm/seo-friendly-slug_%s.htm' % offer_data['identifier']
        source = self._get_url(offer_id)

        offer_data['identifier'] = offer_id

        soup = BeautifulSoup(source, 'html.parser')
        map_links = soup.find_all(id='hitta-map-broker')

        if len(map_links) == 0:
            raise LocationNotFound('No map link found')

        map_link = map_links[0]['src']
        lat, lng = map_link.split('/')[7].split('?')[0].split(':')

        offer_data['location']['lat'] = lat
        offer_data['location']['lon'] = lng

        return offer_data

    def _query_offers(self):
        data = self._get_json(self.API_LIST_URL)

        for offer_data in data['list_items']:
            yield {
                'identifier': offer_data['id'],
                'title': offer_data['address'],
                'description': '',
                'price': offer_data['monthly_rent'],
                'location': {
                    'lat': 0,
                    'lon': 0,
                },
                'rooms': int(float(offer_data['rooms'].replace(',', '.'))),
                'area': int(offer_data['sqm']),
            }

    def _get_url(self, url):
        conn = http.client.HTTPSConnection(self.API_HOST)
        conn.request('GET', url)
        response = conn.getresponse()

        if response.status != 200:
            raise Exception('Error while calling the API')

        return response.read().decode('utf-8')

    def _get_json(self, url):
        payload = self._get_url(url)

        return json.loads(payload)
