import requests
import json
import os

BASE_URL = "https://api.warframestat.us/pc/"
CULTURE = "?language=en"

# This file makes requests to the warframe api and stores the data in json files
def make_cache_dir():
    try:
        os.makedirs("data", exist_ok=True)
        print("Successfully created cache directory")
    except Exception as e:
        print(f"Failed: {e}")
    
def get_archon_hunt():
    url = BASE_URL + "archonHunt" + CULTURE
    response = requests.get(url)
    data = response.json()

    filename = "data/archon_hunt.json"

    dump_data_to_json(filename, data)

def get_day_nigt_cetus():
    url = BASE_URL + "cetusCycle" + CULTURE
    response = requests.get(url)
    data = response.json()

    filename = "data/cetus_cycle.json"

    dump_data_to_json(filename, data)

def get_alerts():
    url = BASE_URL + "alerts" + CULTURE
    response = requests.get(url)
    data = response.json()

    filename = "data/alerts.json"

    dump_data_to_json(filename, data)

def get_news():
    url = BASE_URL + CULTURE
    response = requests.get(url)
    data = response.json()

    filename = "data/news.json"

    dump_data_to_json(filename, data)

def dump_data_to_json(filename, data):
    with open(filename, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=4)

    print(f"Updated {filename}")

def get_all_data():
    make_cache_dir()
    get_archon_hunt()
    get_day_nigt_cetus()
    get_alerts()
    get_news()

# get_all_data()