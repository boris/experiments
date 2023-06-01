import asyncio
import random
import websockets

# Create handler for each connection
async def handler(websocket, path):
    # rand_number simualtes the BTC price
    rand_number = random.randint(10000, 100000)
    print(f"BTC price: {rand_number}")

    await websocket.send(str(rand_number))

# Create server
start_server = websockets.serve(handler, "localhost", 8765)

# Start server
asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
