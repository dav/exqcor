<script>
  import { onMount } from 'svelte'
  import { api } from '../../lib/api.js'

  // One component serves three views:
  //   /script/:id                    full script
  //   /script/:id/actor/:actorId     that actor's lines highlighted
  //   /script/:id/character/:characterId  that character's lines highlighted
  let { params = {} } = $props()

  let full = $state(null)
  let stats = $state(null)
  let error = $state('')

  const actorId = $derived(params.actorId ? Number(params.actorId) : null)
  const characterId = $derived(params.characterId ? Number(params.characterId) : null)

  const highlightCharacterIds = $derived.by(() => {
    if (!full) return new Set()
    if (characterId) return new Set([characterId])
    if (actorId) return new Set(full.characters.filter((c) => c.actor_id === actorId).map((c) => c.id))
    return new Set()
  })

  const focusName = $derived.by(() => {
    if (!full) return ''
    if (characterId) return full.characters.find((c) => c.id === characterId)?.name ?? ''
    if (actorId) return full.actors.find((a) => a.id === actorId)?.name ?? ''
    return ''
  })

  const cast = $derived(
    full
      ? full.characters
          .filter((c) => c.role !== 'vosd')
          .map((c) => ({ ...c, actor: full.actors.find((a) => a.id === c.actor_id) }))
      : []
  )

  onMount(async () => {
    try {
      full = await api.get(`/api/scripts/${params.id}/full`)
    } catch (err) {
      error = err.message
      return
    }
    stats = await api.get(`/api/scripts/${params.id}/stats`).catch(() => null)
  })
</script>

<div class="wrap script-page">
  {#if error}
    <p class="error no-print">{error}</p>
  {:else if !full}
    <p class="muted">Loading the script…</p>
  {:else}
    <div class="row no-print" style="margin-bottom:1rem">
      <button class="primary" onclick={() => window.print()}>Print</button>
      {#if actorId || characterId}
        <a class="btn" href="#/script/{params.id}">Full script</a>
      {/if}
      <details style="flex-basis:100%">
        <summary class="muted">Print a single actor's or character's script</summary>
        <p class="row" style="margin-top:0.5rem">
          {#each full.actors as a (a.id)}
            <a class="btn" href="#/script/{params.id}/actor/{a.id}">{a.name}</a>
          {/each}
          {#each full.characters.filter((c) => c.role !== 'vosd') as c (c.id)}
            <a class="btn" href="#/script/{params.id}/character/{c.id}">{c.name}</a>
          {/each}
        </p>
      </details>
    </div>

    <article class="paper print-sheet">
      <header class="title-page">
        <p class="cue">An exquisite corpse in {full.sections.length} {full.sections.length === 1 ? 'part' : 'parts'}</p>
        <h1>{full.script.title}</h1>
        {#if full.script.theme}<p class="theme">{full.script.theme}</p>{/if}
        {#if focusName}<p class="focus">Sides for {focusName}</p>{/if}
        <p class="byline">Written by the audience · Performed cold by the cast</p>
      </header>

      {#if cast.length > 0}
        <section class="cast">
          <h2 class="section-heading">Cast</h2>
          <table class="cast-table">
            <tbody>
              {#each cast as c (c.id)}
                <tr>
                  <td class="cast-name">{c.name}</td>
                  <td class="cast-actor">{c.actor ? c.actor.name : ''}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      {#each full.sections as sec (sec.id)}
        <section class="act">
          <h2 class="section-heading">{sec.name}</h2>
          {#if sec.props.length > 0}
            <p class="props">Props: {sec.props.map((p) => p.name).join(', ')}</p>
          {/if}
          {#if sec.lines.length === 0}
            <p class="unwritten">(not yet written)</p>
          {/if}
          {#each sec.lines as l (l.id)}
            <div
              class="script-line"
              class:vosd={l.character_role === 'vosd'}
              class:highlight={highlightCharacterIds.has(l.character_id)}
            >
              {#if l.character_role !== 'vosd'}<span class="speaker">{l.character_name}</span>{/if}
              <span class="dialogue">{l.text}</span>
            </div>
          {/each}
          {#if sec.writers.length > 0}
            <p class="writers">Written by {sec.writers.join(', ')}</p>
          {/if}
        </section>
      {/each}
    </article>

    {#if stats?.length}
      <div class="panel no-print" style="margin-top:1.5rem">
        <p class="cue">Writing stats</p>
        <table style="margin-top:0.5rem">
          <thead><tr><th>Section</th><th>Turns</th><th>Lines</th><th>Avg turn</th></tr></thead>
          <tbody>
            {#each stats as st (st.section_id)}
              <tr>
                <td>{st.section_name}</td>
                <td>{st.turns}</td>
                <td>{st.lines}</td>
                <td>{st.avg_turn_seconds ? Math.round(st.avg_turn_seconds) + 's' : '—'}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  {/if}
</div>

<style>
  .script-page :global(.paper) {
    padding: 3rem 2.5rem;
  }
  .title-page {
    text-align: center;
    margin-bottom: 2.5rem;
  }
  .title-page h1 {
    font-family: var(--font-script);
    font-size: 2.2rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
  }
  .theme {
    font-style: italic;
    margin: 0.25rem 0;
  }
  .focus {
    font-family: var(--font-script);
    font-weight: 700;
    margin: 0.75rem 0 0;
  }
  .byline {
    color: var(--ink-dim);
    font-size: 0.9rem;
  }
  .section-heading {
    font-family: var(--font-script);
    text-align: center;
    text-transform: uppercase;
    letter-spacing: 0.15em;
    font-size: 1.1rem;
    margin: 2.5rem 0 1.25rem;
  }
  .cast-table {
    max-width: 28rem;
    margin: 0 auto;
    font-family: var(--font-script);
  }
  .cast-table td {
    border: none;
    padding: 0.15rem 0.75rem;
    color: inherit;
  }
  .cast-name {
    text-align: right;
    text-transform: uppercase;
  }
  .cast-actor {
    text-align: left;
  }
  .props,
  .writers,
  .unwritten {
    font-family: var(--font-script);
    font-style: italic;
    font-size: 0.9rem;
    color: var(--ink-dim);
    text-align: center;
  }
  .writers {
    margin-top: 1.25rem;
  }
  .script-line.highlight .speaker,
  .script-line.highlight .dialogue {
    font-weight: 700;
  }
  .script-line.highlight {
    background: rgba(232, 176, 84, 0.18);
    border-radius: 2px;
  }

  @media print {
    .act {
      break-inside: auto;
    }
    .script-line {
      break-inside: avoid;
    }
    .section-heading {
      break-after: avoid;
    }
    .title-page {
      break-after: page;
    }
    .cast {
      break-after: page;
    }
    .script-line.highlight {
      background: #f3e2b8;
    }
  }
</style>
