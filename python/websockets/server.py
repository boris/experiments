import asyncio
import datetime as dt
import random
import websockets

# Create a heartbeat task
async def heartbeat(websocket):
    while True:
        await asyncio.sleep(3)
        await websocket.send(f"{dt.datetime.now()} - Heartbeat")

# Create handler for each connection
async def handler(websocket, path):
    if path == "/healthz":
        await websocket.send("HTTP/1.1 200 OK\r\n\r\n")
    else:
        await websocket.send(f"{dt.datetime.now()} - WS connection established")

        # Start the heartbeat task
        asyncio.create_task(heartbeat(websocket))

        async for message in websocket:
            # Handle incoming messages
            if message == "close":
                await websocket.close()
            else:
                await websocket.send(f"Received: {message}")


# Create server
start_server = websockets.serve(handler, "0.0.0.0", 8765)

# Start server
asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
