<script>
  import { onMount } from 'svelte'
  import { api } from '../../lib/api.js'

  let program = $state(null)
  let error = $state('')

  const cast = $derived(
    program?.open
      ? program.characters
          .filter((c) => c.role !== 'vosd')
          .map((c) => ({ ...c, actor: program.actors.find((a) => a.id === c.actor_id) }))
      : []
  )

  onMount(async () => {
    try {
      program = await api.get('/api/program')
    } catch (err) {
      error = err.message
    }
  })
</script>

<div class="wrap">
  {#if error}<p class="error">{error}</p>{/if}
  {#if program && !program.open}
    <p class="cue">Program</p>
    <h1>The show hasn't opened yet.</h1>
  {:else if program}
    <div class="paper playbill">
      <p class="cue" style="text-align:center">Tonight's program</p>
      <h1 class="pb-title">{program.script.title}</h1>
      {#if program.script.theme}<p class="pb-theme">{program.script.theme}</p>{/if}
      {#if program.script.description}<p class="pb-desc">{program.script.description}</p>{/if}

      {#if program.sections.length > 0}
        <h2 class="pb-heading">The Acts</h2>
        <ol class="pb-acts">
          {#each program.sections as s (s.id)}
            <li>{s.name}{#if s.description} — <span class="pb-dim">{s.description}</span>{/if}</li>
          {/each}
        </ol>
      {/if}

      {#if cast.length > 0}
        <h2 class="pb-heading">The Cast</h2>
        <table class="pb-cast">
          <tbody>
            {#each cast as c (c.id)}
              <tr>
                <td class="pb-char">{c.name}</td>
                <td>{c.actor ? c.actor.name : 'the ensemble'}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}

      <h2 class="pb-heading">How tonight works</h2>
      <p class="pb-how">
        This play does not exist yet. You — the audience — write it tonight,
        one person at a time. Each writer sees only the <em>last line</em> the
        previous writer left behind, adds a few minutes of dialogue and stage
        directions, and passes it on. That's the old parlor game the
        surrealists called an <em>exquisite corpse</em>. When an act is
        finished, we print it, hand it to the cast, and they perform it cold
        — first read, live, no rehearsal. Nobody in the building knows what
        happens next. Including us.
      </p>
      <p class="pb-how">
        Keep your number handy — when it's called, you're the playwright.
      </p>
    </div>
    <p class="no-print" style="margin-top:1rem"><a href="#/audience">Back to my number</a></p>
  {/if}
</div>

<style>
  .playbill {
    font-family: var(--font-serif);
  }
  .pb-title {
    text-align: center;
    font-size: 2.2rem;
    margin: 0.5rem 0 0;
  }
  .pb-theme {
    text-align: center;
    font-style: italic;
    margin: 0.25rem 0 0;
  }
  .pb-desc {
    text-align: center;
    color: var(--ink-dim);
  }
  .pb-heading {
    text-align: center;
    font-size: 0.85rem;
    letter-spacing: 0.22em;
    text-transform: uppercase;
    margin: 2rem 0 0.75rem;
    color: var(--ink-dim);
  }
  .pb-acts {
    max-width: 24rem;
    margin: 0 auto;
    padding-left: 1.5rem;
  }
  .pb-dim {
    color: var(--ink-dim);
  }
  .pb-cast {
    max-width: 26rem;
    margin: 0 auto;
  }
  .pb-cast td {
    border: none;
    padding: 0.15rem 0.75rem;
    color: inherit;
  }
  .pb-char {
    text-align: right;
    font-variant: small-caps;
    font-weight: 600;
  }
  .pb-how {
    max-width: 34rem;
    margin: 0.75rem auto;
  }
</style>
