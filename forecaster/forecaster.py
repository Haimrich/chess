from calendar import timegm
from datetime import datetime
import requests
import pandas as pd
import os
from statsmodels.tsa.statespace.sarimax import SARIMAX

import _strptime 
from flask import Flask, request, jsonify

from pymongo import MongoClient

app = Flask(__name__)

app.debug = True

PROMETHEUS_URL = os.getenv("PROMETHEUS_URL")
N_FORECAST = int(os.getenv("N_FORECAST",default="20"))

MONGODB_URI = os.getenv("MONGODB_URI", default="nop")

# Questo script effettua forecasting sull'utilizzo della cpu da parte dell'engine interrogando Prometheus
# ma fornisce anche dati sulle posizioni per cui occorre più tempo di calcolo, prelevando i dati dal database delle partite


@app.route('/')
def health_check():
    return 'This datasource is healthy.'


@app.route('/search', methods=['POST'])
def search():
    metrics = ['engine_prediction']

    if (MONGODB_URI != "nop"):
        metrics.append('engine_difficult_positions')

    return jsonify(metrics)


@app.route('/query', methods=['POST'])
def query():
    grafana_req = request.get_json()
    target = grafana_req['targets'][0]['target']
    
    if target == 'engine_prediction':
        return jsonify(forecast_data())
    elif target == 'engine_difficult_positions':
        return jsonify(position_data())
    else:
        return jsonify()


def forecast_data():
    query = 'sum (rate (container_cpu_usage_seconds_total{container="engine-service"}[80s]))[20m:5s]'
    req = requests.get(PROMETHEUS_URL + "/api/v1/query", params={"query": query})
    json = req.json()

    data = pd.DataFrame(json["data"]["result"][0]["values"], columns = ['time','value'])
    data['time'] = pd.to_numeric(data['time'])*1000
    data['value'] = pd.to_numeric(data['value'])

    my_order = (1, 0, 0)
    my_seasonal_order = (2, 1, 1, 20)

    model = SARIMAX(data.value, order=my_order, seasonal_order=my_seasonal_order, enforce_stationarity=False, enforce_invertibility=False).fit()

    forecast = model.forecast(N_FORECAST)

    data = [
        {
            "target": "engine_prediction",
            "datapoints": [[forecast.iloc[i], int(data['time'].iloc[-1])+(i+1)*5000] for i in range(0, N_FORECAST)]
        }
    ]
    return data

def position_data():
    client = MongoClient(MONGODB_URI)
    data = client['chess']['engine_metrics'].aggregate([
        {
            '$addFields': {
                'position': '$fen', 
                'time': {
                    '$avg': '$measures'
                }
            }
        }, {
            '$project': {
                'position': '$position', 
                'time': '$time'
            }
        }, {
            '$sort': {
                'time': -1
            }
        }, {
            '$limit': 10
        }
    ])
    
    result = [
        {
            "columns":[
                {"text":"Position","type":"string"},
                {"text":"Avg. Search Time","type":"number"}
            ],
            "rows": [[d["position"], d["time"]] for d in data],
            "type":"table"
        }
    ]

    return result

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000)