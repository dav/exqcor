<script>
  import { onMount } from 'svelte'
  import { api } from '../../lib/api.js'

  let { section, characters, onDeleted } = $props()

  let charSections = $state([])
  let props = $state([])
  let primingText = $state('')
  let propName = $state('')
  let error = $state('')
  let open = $state(false)

  const vosd = $derived(characters.find((c) => c.role === 'vosd'))

  async function load() {
    let priming
    ;[charSections, props, priming] = await Promise.all([
      api.get(`/api/sections/${section.id}/character-sections`),
      api.get(`/api/sections/${section.id}/props`),
      api.get(`/api/sections/${section.id}/priming-line`),
    ])
    if (priming.line) primingText = priming.line.text
  }
  onMount(load)

  function linkFor(charId) {
    return charSections.find((cs) => cs.character_id === charId)
  }

  async function setLink(char, attached, onStage) {
    error = ''
    try {
      await api.post(`/api/sections/${section.id}/character-sections`, {
        character_id: char.id,
        attached,
        on_stage: onStage,
      })
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function savePriming(e) {
    e.preventDefault()
    error = ''
    try {
      await api.post(`/api/sections/${section.id}/priming-line`, {
        character_id: vosd.id,
        text: primingText,
      })
    } catch (err) {
      error = err.message
    }
  }

  async function addProp(e) {
    e.preventDefault()
    error = ''
    try {
      await api.post(`/api/sections/${section.id}/props`, { name: propName })
      propName = ''
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function removeProp(id) {
    await api.del(`/api/props/${id}`)
    await load()
  }

  async function removeSection() {
    if (!confirm(`Delete section "${section.name}" and all its writing?`)) return
    try {
      await api.del(`/api/sections/${section.id}`)
      onDeleted?.()
    } catch (err) {
      error = err.message
    }
  }
</script>

<div class="panel">
  <div class="spread">
    <h3 style="margin:0">{section.name}</h3>
    <button onclick={() => (open = !open)}>{open ? 'Close' : 'Edit'}</button>
  </div>

  {#if open}
    {#if error}<p class="error">{error}</p>{/if}

    <p class="cue" style="margin-top:1rem">Opening line (what the first writer sees)</p>
    <form onsubmit={savePriming} class="row" style="margin-top:0.4rem">
      <input
        type="text"
        bind:value={primingText}
        placeholder="It was a dark and stormy night…"
        style="flex:1"
      />
      <button class="primary" type="submit">Save line</button>
    </form>

    <p class="cue" style="margin-top:1.25rem">Who's in this section</p>
    <table style="margin-top:0.4rem">
      <thead><tr><th>Character</th><th>In section</th><th>On stage</th></tr></thead>
      <tbody>
        {#each characters as c (c.id)}
          {@const link = linkFor(c.id)}
          <tr>
            <td>{c.name}{#if c.role === 'vosd'}<span class="muted"> (stage directions)</span>{/if}</td>
            <td>
              <input
                type="checkbox"
                checked={!!link}
                disabled={c.role === 'vosd'}
                onchange={(e) => setLink(c, e.target.checked, link?.on_stage ?? false)}
              />
            </td>
            <td>
              {#if c.role !== 'vosd'}
                <input
                  type="checkbox"
                  checked={link?.on_stage ?? false}
                  disabled={!link}
                  onchange={(e) => setLink(c, true, e.target.checked)}
                />
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <p class="cue" style="margin-top:1.25rem">Props</p>
    {#if props.length > 0}
      <ul style="margin:0.4rem 0">
        {#each props as p (p.id)}
          <li class="row">{p.name} <button onclick={() => removeProp(p.id)}>Remove</button></li>
        {/each}
      </ul>
    {/if}
    <form onsubmit={addProp} class="row" style="margin-top:0.4rem">
      <input type="text" bind:value={propName} placeholder="A suspicious envelope" style="flex:1" required />
      <button type="submit">Add prop</button>
    </form>

    <p style="margin-top:1.25rem"><button class="danger" onclick={removeSection}>Delete section</button></p>
  {/if}
</div>
