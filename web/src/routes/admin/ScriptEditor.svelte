<script>
  import { onMount } from 'svelte'
  import AdminGate from '../../lib/AdminGate.svelte'
  import SectionCard from './SectionCard.svelte'
  import { api } from '../../lib/api.js'

  let { params = {} } = $props()

  let script = $state(null)
  let sections = $state([])
  let characters = $state([])
  let actors = $state([])
  let error = $state('')
  let saved = $state(false)

  let newSection = $state('')
  let newCharacter = $state('')
  let newActor = $state('')

  async function load() {
    const id = params.id
    ;[script, sections, characters, actors] = await Promise.all([
      api.get(`/api/scripts/${id}`),
      api.get(`/api/scripts/${id}/sections`),
      api.get(`/api/scripts/${id}/characters`),
      api.get('/api/actors'),
    ])
  }
  onMount(load)

  async function saveScript(e) {
    e.preventDefault()
    error = ''
    saved = false
    try {
      script = await api.put(`/api/scripts/${script.id}`, {
        title: script.title,
        description: script.description,
        theme: script.theme,
        writing_seconds: Number(script.writing_seconds),
        station_mode: script.station_mode,
      })
      saved = true
      setTimeout(() => (saved = false), 2000)
    } catch (err) {
      error = err.message
    }
  }

  async function addSection(e) {
    e.preventDefault()
    try {
      await api.post(`/api/scripts/${script.id}/sections`, { name: newSection })
      newSection = ''
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function addCharacter(e) {
    e.preventDefault()
    try {
      await api.post(`/api/scripts/${script.id}/characters`, { name: newCharacter })
      newCharacter = ''
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function assignActor(character, actorId) {
    try {
      await api.put(`/api/characters/${character.id}`, {
        actor_id: actorId ? Number(actorId) : null,
        name: character.name,
        description: character.description,
      })
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function removeCharacter(c) {
    if (!confirm(`Delete character "${c.name}"?`)) return
    try {
      await api.del(`/api/characters/${c.id}`)
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function addActor(e) {
    e.preventDefault()
    try {
      await api.post('/api/actors', { name: newActor })
      newActor = ''
      await load()
    } catch (err) {
      error = err.message
    }
  }
</script>

<AdminGate>
  <div class="wrap">
    {#if script}
      <p class="cue"><a href="#/admin">Production desk</a> / <a href="#/admin/scripts">Scripts</a> / {script.title}</p>
      <h1>{script.title}</h1>
      {#if error}<p class="error">{error}</p>{/if}

      <div class="panel">
        <p class="cue">Show settings</p>
        <form onsubmit={saveScript}>
          <label for="s-title">Title</label>
          <input id="s-title" type="text" bind:value={script.title} required />
          <label for="s-theme">Theme</label>
          <input id="s-theme" type="text" bind:value={script.theme} />
          <label for="s-desc">Description</label>
          <textarea id="s-desc" rows="2" bind:value={script.description}></textarea>
          <label for="s-secs">Writing time per turn (seconds)</label>
          <input id="s-secs" type="number" min="30" bind:value={script.writing_seconds} />
          <label for="s-mode">Where do people write?</label>
          <select id="s-mode" bind:value={script.station_mode}>
            <option value="station">Backstage writing stations</option>
            <option value="byod">Their own phones</option>
          </select>
          <p class="row">
            <button class="primary" type="submit">Save settings</button>
            {#if saved}<span class="muted">Saved.</span>{/if}
          </p>
        </form>
      </div>

      <div class="panel">
        <p class="cue">Characters</p>
        <table style="margin-top:0.5rem">
          <thead><tr><th>Name</th><th>Played by</th><th></th></tr></thead>
          <tbody>
            {#each characters as c (c.id)}
              <tr>
                <td>{c.name}{#if c.role === 'vosd'}<span class="muted"> (stage directions)</span>{/if}</td>
                <td>
                  {#if c.role !== 'vosd'}
                    <select value={c.actor_id ?? ''} onchange={(e) => assignActor(c, e.target.value)}>
                      <option value="">— uncast —</option>
                      {#each actors as a (a.id)}
                        <option value={a.id}>{a.name}</option>
                      {/each}
                    </select>
                  {/if}
                </td>
                <td style="text-align:right">
                  {#if c.role !== 'vosd'}
                    <button class="danger" onclick={() => removeCharacter(c)}>Delete</button>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
        <form onsubmit={addCharacter} class="row" style="margin-top:0.75rem">
          <input type="text" bind:value={newCharacter} placeholder="New character name" style="flex:1" required />
          <button type="submit">Add character</button>
        </form>
        <form onsubmit={addActor} class="row" style="margin-top:0.5rem">
          <input type="text" bind:value={newActor} placeholder="New actor name" style="flex:1" required />
          <button type="submit">Add actor</button>
        </form>
      </div>

      <div class="spread" style="margin-top:1.5rem">
        <h2 style="margin:0">Sections</h2>
      </div>
      <p class="muted" style="margin-top:0.25rem">
        Sections are the acts of the show, written in order on show night.
      </p>
      {#each sections as sec (sec.id)}
        <SectionCard section={sec} {characters} onDeleted={load} />
      {/each}
      <form onsubmit={addSection} class="row">
        <input type="text" bind:value={newSection} placeholder="Act 1" style="flex:1" required />
        <button type="submit">Add section</button>
      </form>
    {:else}
      <p class="muted">Loading…</p>
    {/if}
  </div>
</AdminGate>
