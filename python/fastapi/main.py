from fastapi import Request, FastAPI
from typing import List, Tuple

app = FastAPI()


def create_tuple(data: dict) -> List[Tuple[str, int]]:
    users = data.get("users", [])
    tuples = [(user["name"], user["amount"]) for user in users]
    return tuples

@app.post("/")
async def get_json(request: Request):
    file = await request.json()

    current_price = file['price']

    # Later: Add validation for each field

    # Return the response
    tuples = create_tuple(file)

    results = []

    for tuple in tuples:
        value = tuple[1] / current_price
        results.append({
            tuple[0] : f'{value:.8f}'
        })

    return results
