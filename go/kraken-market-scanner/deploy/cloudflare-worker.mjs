export default {
  async scheduled(event, env, ctx) {
    ctx.waitUntil(triggerScan(env))
  },

  async fetch(request, env) {
    if (new URL(request.url).pathname !== '/trigger') {
      return new Response('not found', { status: 404 })
    }

    if (!env.WORKER_TRIGGER_TOKEN) {
      return Response.json({ error: 'worker trigger token is not configured' }, { status: 500 })
    }

    const authHeader = request.headers.get('authorization') || ''
    if (authHeader !== `Bearer ${env.WORKER_TRIGGER_TOKEN}`) {
      return Response.json({ error: 'unauthorized' }, { status: 401 })
    }

    try {
      const payload = await triggerScan(env)
      return Response.json(payload)
    } catch (error) {
      return Response.json({ error: String(error) }, { status: 502 })
    }
  },
}

async function triggerScan(env) {
  const response = await fetch(`${env.SCANNER_URL}?send=1`, {
    headers: {
      Authorization: `Bearer ${env.SCANNER_API_TOKEN}`,
    },
  })

  const text = await response.text()
  if (!response.ok) {
    throw new Error(`scanner request failed: ${response.status} ${text}`)
  }

  return {
    ok: true,
    status: response.status,
    body: safeJSON(text),
  }
}

function safeJSON(text) {
  try {
    return JSON.parse(text)
  } catch {
    return { raw: text }
  }
}
