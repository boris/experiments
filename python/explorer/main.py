import requests
import json
import os

api_key = 'e1558291f84749d7a861822fe8a6284d' or os.environ.get("INFURA_API_KEY")

url = f"https://mainnet.infura.io/v3/{api_key}"


headers = {
    "Content-Type": "application/json"
}

def get_block_number() -> int:
    payload = {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["latest", False],
        "id": 1
    }
    response = requests.post(url, data=json.dumps(payload), headers=headers)
    return int(response.json()['result']['number'], 0)

def get_block_by_number(block_number: int) -> dict:
    payload = {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [hex(block_number), True],
        "id": 1
    }
    response = requests.post(url, data=json.dumps(payload), headers=headers)
    return response.json()


#print(int(response.json()['result'], 0))
#print(response.json()['result'])
print(get_block_number())
print(get_block_by_number(get_block_number()))
