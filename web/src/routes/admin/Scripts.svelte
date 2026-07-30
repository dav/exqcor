<script>
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import AdminGate from '../../lib/AdminGate.svelte'
  import { api } from '../../lib/api.js'

  let scripts = $state([])
  let title = $state('')
  let theme = $state('')
  let error = $state('')

  async function load() {
    scripts = await api.get('/api/scripts')
  }
  onMount(load)

  async function create(e) {
    e.preventDefault()
    error = ''
    try {
      const s = await api.post('/api/scripts', { title, theme })
      push(`/admin/scripts/${s.id}`)
    } catch (err) {
      error = err.message
    }
  }

  async function duplicate(id) {
    try {
      const s = await api.post(`/api/scripts/${id}/duplicate`, {})
      push(`/admin/scripts/${s.id}`)
    } catch (err) {
      error = err.message
    }
  }

  async function remove(s) {
    if (!confirm(`Delete "${s.title}" and everything written in it? This cannot be undone.`)) return
    try {
      await api.del(`/api/scripts/${s.id}`)
      await load()
    } catch (err) {
      error = err.message
    }
  }
</script>

<AdminGate>
  <div class="wrap">
    <p class="cue"><a href="#/admin">Production desk</a> / Scripts</p>
    <h1>Scripts</h1>
    {#if error}<p class="error">{error}</p>{/if}

    {#if scripts.length === 0}
      <p class="muted">No scripts yet. A script is one production: its theme, characters, and acts.</p>
    {:else}
      <table>
        <thead>
          <tr><th>Title</th><th>Theme</th><th></th></tr>
        </thead>
        <tbody>
          {#each scripts as s (s.id)}
            <tr>
              <td><a href="#/admin/scripts/{s.id}">{s.title}</a></td>
              <td class="muted">{s.theme}</td>
              <td class="row" style="justify-content:flex-end">
                <button onclick={() => duplicate(s.id)}>Duplicate</button>
                <button class="danger" onclick={() => remove(s)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}

    <div class="panel" style="margin-top:1.5rem">
      <p class="cue">New script</p>
      <form onsubmit={create}>
        <label for="title">Title</label>
        <input id="title" type="text" bind:value={title} required placeholder="Corpse! The Soap Opera" />
        <label for="theme">Theme</label>
        <input id="theme" type="text" bind:value={theme} placeholder="Soap Opera, Film Noir, SciFi…" />
        <p><button class="primary" type="submit">Create script</button></p>
      </form>
    </div>
  </div>
</AdminGate>
