from web3 import Web3
import os

def get_latest_block():
    infura_url = f"https://mainnet.infura.io/v3/{os.getenv('INFURA_TOKEN')}"
    web3 = Web3(Web3.HTTPProvider(infura_url))

    if web3.is_connected():
        latest_block_n = web3.eth.block_number
        print(latest_block_n)

    else:
        print("Failed to connect")


if __name__ == "__main__":
    get_latest_block()
