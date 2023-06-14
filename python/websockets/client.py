import asyncio
import websockets
from aiohttp import web

async def heartbeat_client():
    uri = "ws://localhost:8765"  # Change the URI if the server is running on a different address

    async with websockets.connect(uri) as websocket:
        while True:
            message = await websocket.recv()
            print(f"{message}")

async def healthcheck(request):
    return web.Response(text="OK", status=200)

# Create server
app = web.Application()
app.add_routes([web.get('/healthz', healthcheck)])

# Start the heartbeat client
#asyncio.get_event_loop().run_until_complete(heartbeat_client())
async def start_server():
    runner = web.AppRunner(app)
    await runner.setup()
    site = web.TCPSite(runner, '0.0.0.0', 8766)
    await site.start()

    # Start the heartbeat client
    await heartbeat_client()

try:
    asyncio.run(start_server())
except KeyboardInterrupt:
    pass
