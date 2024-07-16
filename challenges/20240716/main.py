from fastapi import FastAPI, HTTPException
import requests

app = FastAPI()

INFURA_PROJECT_ID = "TBD"

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/balance")
async def get_balance(address: str):
    infura_url = f"https://sepolia.infura.io/v3/{INFURA_PROJECT_ID}"

    latest_block_response = requests.post(infura_url, json={
        "jsonrpc": "2.0",
        "id": 1,
        "method": "eth_blockNumber",
        "params": []
    })

    latest_block = latest_block_response.json().get("result")

    balance_response = requests.post(infura_url, json={
        "jsonrpc": "2.0",
        "id": 1,
        "method": "eth_getBalance",
        "params": [address, latest_block]
    })

    balance = balance_response.json().get("result")
    # convert balance from Wei to ETH
    balance_int = int(balance, 16) / 10**18

    # Convert balance to int as the result is in hex
    return {"latest_block": int(latest_block, 16), "balance": balance_int}
