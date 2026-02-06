const defaultPodSpec = `apiVersion: v1
kind: Pod
metadata:
  name: challenge
spec:
  containers:
    - name: app
      image: nginx:stable
      ports:
        - containerPort: 80
          protocol: TCP
      resources:
        requests:
          cpu: "100m"
          memory: "128Mi"
        limits:
          cpu: "100m"
          memory: "128Mi"`

function byId(id) {
    return document.getElementById(id)
}

function apiBase() {
    const raw = byId('apiBase').value.trim()
    if (raw.length > 0) {
        return raw.replace(/\/+$/, '')
    }
    return window.location.origin
}

function readInt(id) {
    return Number.parseInt(byId(id).value, 10)
}

function syncStackID(stackID) {
    if (!stackID) {
        return
    }
    byId('lastStackId').value = stackID
    byId('stackIdInput').value = stackID
}

function showResponse(title, payload) {
    const el = byId('response')
    const ts = new Date().toISOString()
    const body = typeof payload === 'string' ? payload : JSON.stringify(payload, null, 2)
    el.textContent = `[${ts}] ${title}\n\n${body}`
}

async function request(method, path, body) {
    const url = `${apiBase()}${path}`
    const options = { method, headers: {} }
    if (body !== undefined) {
        options.headers['Content-Type'] = 'application/json'
        options.body = JSON.stringify(body)
    }

    const res = await fetch(url, options)
    const raw = await res.text()
    let parsed = raw
    try {
        parsed = raw.length ? JSON.parse(raw) : {}
    } catch (e) {
        
    }

    const responsePayload = {
        method,
        url,
        status: res.status,
        ok: res.ok,
        body: parsed,
    }

    if (!res.ok) {
        throw responsePayload
    }

    return responsePayload
}

async function execute(title, fn) {
    try {
        const result = await fn()
        showResponse(title, result)
    } catch (err) {
        if (err instanceof Error) {
            showResponse(`${title} (ERROR)`, { error: err.message })
            return
        }
        showResponse(`${title} (ERROR)`, err)
    }
}

function getStackID() {
    const stackID = byId('stackIdInput').value.trim()
    if (!stackID) {
        throw new Error('stack_id is required')
    }
    return stackID
}

function wireEvents() {
    byId('btnHealth').addEventListener('click', () => execute('GET /healthz', () => request('GET', '/healthz')))

    byId('btnListAllStacks').addEventListener('click', () => execute('GET /stacks', () => request('GET', '/stacks')))

    byId('btnStats').addEventListener('click', () => execute('GET /stats', () => request('GET', '/stats')))

    byId('btnCreateStack').addEventListener('click', () =>
        execute('POST /stacks', async () => {
            const payload = {
                target_port: readInt('createTargetPort'),
                pod_spec: byId('createPodSpec').value,
            }
            const result = await request('POST', '/stacks', payload)
            const stackID = result && result.body ? result.body.stack_id : ''
            syncStackID(stackID)
            return result
        }),
    )

    byId('btnGetStack').addEventListener('click', () =>
        execute('GET /stacks/{stack_id}', () => request('GET', `/stacks/${getStackID()}`)),
    )

    byId('btnGetStackStatus').addEventListener('click', () =>
        execute('GET /stacks/{stack_id}/status', () => request('GET', `/stacks/${getStackID()}/status`)),
    )

    byId('btnDeleteStack').addEventListener('click', () =>
        execute('DELETE /stacks/{stack_id}', () => request('DELETE', `/stacks/${getStackID()}`)),
    )
}

function boot() {
    byId('apiBase').value = window.location.origin
    byId('createPodSpec').value = defaultPodSpec
    wireEvents()
}

window.addEventListener('DOMContentLoaded', boot)
