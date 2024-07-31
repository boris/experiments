from web3 import Web3
import click
import os

INFURA_PROJECT_ID = os.environ.get('INFURA_PROJECT_ID')
ETH_ENDPOINT = os.environ.get('ETH_ENDPOINT', 'sepolia')

@click.command()
@click.option('--wallet_addr', prompt='Enter the wallet address', help='The wallet address to check the balance')
def check_funds(wallet_addr):
    w3 = Web3(Web3.HTTPProvider(f'https://{ETH_ENDPOINT}.infura.io/v3/{INFURA_PROJECT_ID}'))

    balance = w3.eth.get_balance(wallet_addr)
    print(f"Balance: {w3.from_wei(balance, 'ether')} ETH")

    if os.environ.get('ETH_ENDPOINT') == None:
        print(f"See your wallet in Etherscan: https://etherscan.io/address/{wallet_addr}")
    else:
        print(f"See your wallet in Etherscan: https://{ETH_ENDPOINT}.etherscan.io/address/{wallet_addr}")

if __name__ == '__main__':
    check_funds()


#wallet_addr = '0xA33976bD76dE5Dd01Fa7c040927D40034AC28063'
