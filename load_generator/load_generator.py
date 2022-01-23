import requests
import time
import math
import multiprocessing
import os


ENGINE_URL = os.getenv("ENGINE_URL", default = "http://chess.example:7772/")
FEN = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
MAX_REQS = 8
PERIOD = 100

def request(v):
    res = requests.post(ENGINE_URL, json={'fen': FEN, 'budget': 1e7})
    print(str(res.content))


while True:
    n_req = math.ceil( (MAX_REQS * math.sin(2*math.pi*time.time()/PERIOD) + MAX_REQS) / 2 )    
    pool_obj = multiprocessing.Pool()
    answer = pool_obj.map_async(request,range(0,n_req))
    
    print("Richieste: " + str(n_req))

    time.sleep(10)
    