import asyncio
import websockets

async def hello():
    uri = "ws://localhost:8765"
    async with websockets.connect(uri) as websocket:
        latest_price = await websocket.recv()
        print(f"< {latest_price}")

asyncio.get_event_loop().run_until_complete(hello())
