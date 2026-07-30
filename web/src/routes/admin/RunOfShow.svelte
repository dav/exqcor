<script>
  import { onMount, onDestroy } from 'svelte'
  import AdminGate from '../../lib/AdminGate.svelte'
  import { api } from '../../lib/api.js'
  import { subscribe } from '../../lib/sse.js'
  import { clockOffset, secondsLeft, formatClock } from '../../lib/timer.js'

  let show = $state(null)
  let scripts = $state([])
  let turns = $state({}) // sectionId -> turn context
  let audience = $state([])
  let error = $state('')
  let offset = $state(0)
  let now = $state(Date.now())
  let pollHandle, tickHandle, unsubscribe

  const waiting = $derived(audience.filter((m) => m.status === 'waiting'))
  const called = $derived(audience.filter((m) => m.status === 'called'))

  async function poll() {
    try {
      ;[show, audience] = await Promise.all([api.get('/api/show'), api.get('/api/audience')])
      if (show.server_now) offset = clockOffset(show.server_now)
      if (show.open) {
        const entries = await Promise.all(
          show.sections.map(async (s) => [s.id, await api.get(`/api/sections/${s.id}/turn`)])
        )
        turns = Object.fromEntries(entries)
      }
      error = ''
    } catch (err) {
      error = err.message
    }
  }

  onMount(async () => {
    scripts = await api.get('/api/scripts').catch(() => [])
    await poll()
    pollHandle = setInterval(poll, 5000)
    tickHandle = setInterval(() => (now = Date.now()), 500)
    unsubscribe = subscribe({
      queue_changed: poll,
      turn_started: poll,
      turn_ended: poll,
      show_state: poll,
    })
  })
  onDestroy(() => {
    clearInterval(pollHandle)
    clearInterval(tickHandle)
    unsubscribe?.()
  })

  async function callNext(sectionId) {
    error = ''
    try {
      await api.post('/api/queue/call-next', { section_id: sectionId })
    } catch (err) {
      error = err.message
    }
  }

  async function setStatus(memberId, action) {
    error = ''
    try {
      await api.post(`/api/queue/${memberId}/${action}`)
    } catch (err) {
      error = err.message
    }
  }

  async function openShow(scriptId) {
    error = ''
    try {
      await api.post('/api/show/open', { script_id: scriptId })
      await poll()
    } catch (err) {
      error = err.message
    }
  }

  async function endTurn(turnId) {
    if (!confirm('Cut this writer off now?')) return
    try {
      await api.post(`/api/turns/${turnId}/done`)
      await poll()
    } catch (err) {
      error = err.message
    }
  }

  async function markComplete(section) {
    try {
      await api.put(`/api/sections/${section.id}`, { status: 'complete' })
      await poll()
    } catch (err) {
      error = err.message
    }
  }
</script>

<AdminGate>
  <div class="wrap">
    <p class="cue"><a href="#/admin">Production desk</a> / Run of show</p>
    <h1>Run of show</h1>
    {#if error}<p class="error">{error}</p>{/if}

    {#if !show?.open}
      <p class="muted">No show is open. Pick tonight's script:</p>
      {#each scripts as s (s.id)}
        <div class="panel spread">
          <span>{s.title} <span class="muted">{s.theme}</span></span>
          <button class="primary" onclick={() => openShow(s.id)}>Open as tonight's show</button>
        </div>
      {/each}
      {#if scripts.length === 0}
        <p class="muted"><a href="#/admin/scripts">Create a script first.</a></p>
      {/if}
    {:else}
      <div class="panel spread">
        <div>
          <p class="cue">Tonight</p>
          <h2 style="margin:0.25rem 0 0">{show.script.title}</h2>
        </div>
        <a class="btn" href="#/script/{show.script.id}">View script</a>
      </div>

      {#each show.sections as sec (sec.id)}
        {@const t = turns[sec.id]?.turn}
        <div class="panel">
          <div class="spread">
            <h3 style="margin:0">{sec.name}</h3>
            <span class="cue">{sec.status} · {sec.turn_count} turns</span>
          </div>
          {#if t}
            <p style="margin:0.75rem 0 0">
              Writer at work —
              <strong class="mono">{formatClock(secondsLeft(t.ends_at, offset, now))}</strong>
              left.
              {#if turns[sec.id].my_lines?.length}
                <span class="muted">{turns[sec.id].my_lines.length} lines so far.</span>
              {/if}
            </p>
            <p class="row" style="margin-top:0.5rem">
              <button class="danger" onclick={() => endTurn(t.id)}>End turn now</button>
            </p>
          {:else}
            <p class="muted" style="margin:0.75rem 0 0">No writer at the moment.</p>
            <p class="row" style="margin-top:0.5rem">
              <button class="primary" onclick={() => callNext(sec.id)} disabled={waiting.length === 0}>
                Call next writer here
              </button>
              {#if sec.status !== 'complete'}
                <button onclick={() => markComplete(sec)}>Mark section complete</button>
              {/if}
            </p>
          {/if}
        </div>
      {/each}

      <div class="panel">
        <div class="spread">
          <h3 style="margin:0">Audience queue</h3>
          <span class="cue">{waiting.length} waiting · {called.length} called</span>
        </div>
        {#if audience.length === 0}
          <p class="muted" style="margin:0.75rem 0 0">
            No one has scanned in yet. Numbers appear here as people join.
          </p>
        {:else}
          <table style="margin-top:0.75rem">
            <thead><tr><th>#</th><th>Name</th><th>Status</th><th></th></tr></thead>
            <tbody>
              {#each audience as m (m.id)}
                <tr>
                  <td><strong>{m.number}</strong></td>
                  <td>{m.name || '—'}</td>
                  <td class="muted">{m.status}</td>
                  <td class="row" style="justify-content:flex-end">
                    {#if m.status === 'called'}
                      <button onclick={() => setStatus(m.id, 'skip')}>No-show</button>
                    {:else if m.status === 'skipped' || m.status === 'done'}
                      <button onclick={() => setStatus(m.id, 'requeue')}>Back in line</button>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
    {/if}
  </div>
</AdminGate>

<style>
  .mono {
    font-family: var(--font-script);
    color: var(--tungsten);
  }
</style>
