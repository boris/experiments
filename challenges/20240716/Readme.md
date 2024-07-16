# Objective:
Create a FastAPI application that exposes an endpoint to fetch the latest
Ethereum block number and the balance of a given Ethereum address using Infura.

## Requirements:
1. An [Infura](https://infura.io/) account.
1. Set up a FastAPI application.
1. Create an endpoint /balance that accepts an Ethereum address as a query
   parameter.
1. Use the Infura Ethereum API to fetch the latest block number and the balance
   of the provided address.
1. Return the block number and balance in the response.

## Solution
```
# using python 3.10.12
pip install -r requirements.txt
make run
```

To test:
```
 make check
curl -s "localhost:8000/balance?address=0xA33976bD76dE5Dd01Fa7c040927D40034AC28063" | jq
{
  "latest_block": 6319119,
  "balance": 0.22
}

```
