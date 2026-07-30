<script>
  // The live writing surface, shared by backstage stations and (in
  // bring-your-own-device shows) audience phones. Self-contained: polls the
  // turn context, listens to SSE, and calls onExit when the turn ends.
  import { onMount, onDestroy } from 'svelte'
  import { api } from './api.js'
  import { subscribe } from './sse.js'
  import { clockOffset, secondsLeft, formatClock } from './timer.js'

  let { sectionId, sectionName = '', onExit } = $props()

  let ctx = $state(null)
  let selectedCharacter = $state(null)
  let text = $state('')
  let error = $state('')
  let offset = $state(0)
  let now = $state(Date.now())
  let connection = $state('connected')
  let textarea
  let pollHandle, tickHandle, unsubscribe

  const turn = $derived(ctx?.turn ?? null)
  const grace = $derived(ctx?.grace_seconds ?? 15)
  const remaining = $derived(turn?.ends_at ? secondsLeft(turn.ends_at, offset, now) : 0)
  const inGrace = $derived(turn && remaining <= 0 && remaining > -grace)
  const over = $derived(turn && remaining <= -grace)

  async function poll() {
    try {
      const fresh = await api.get(`/api/sections/${sectionId}/turn`)
      offset = clockOffset(fresh.server_now)
      if (ctx?.turn && !fresh.turn) {
        onExit?.()
        return
      }
      ctx = fresh
      error = ''
    } catch (err) {
      error = err.message
    }
  }

  onMount(() => {
    poll()
    pollHandle = setInterval(poll, 3000)
    tickHandle = setInterval(() => (now = Date.now()), 250)
    unsubscribe = subscribe(
      {
        turn_ended: (d) => {
          if (d?.section_id === Number(sectionId)) poll()
        },
      },
      (state) => (connection = state)
    )
  })
  onDestroy(() => {
    clearInterval(pollHandle)
    clearInterval(tickHandle)
    unsubscribe?.()
  })

  async function addLine(e) {
    e.preventDefault()
    if (!selectedCharacter || !text.trim()) return
    try {
      await api.post(`/api/turns/${turn.id}/lines`, {
        character_id: selectedCharacter.id,
        text: text.trim(),
      })
      text = ''
      error = ''
      await poll()
      textarea?.focus()
    } catch (err) {
      error = err.message
    }
  }

  async function finish() {
    try {
      await api.post(`/api/turns/${turn.id}/done`)
      onExit?.()
    } catch (err) {
      error = err.message
    }
  }

  function onKeydown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      addLine(e)
    }
  }
</script>

{#if turn}
  <div class="spread no-print">
    <p class="cue">{sectionName || ctx.section?.name}</p>
    <p class="timer" class:low={remaining <= 60 && remaining > 0} class:grace={inGrace || over} aria-live="polite">
      {#if inGrace}
        Time! Finish your sentence…
      {:else if over}
        0:00
      {:else}
        {formatClock(remaining)}
      {/if}
    </p>
  </div>

  {#if connection === 'reconnecting'}
    <p class="error">Connection wobbling — your saved lines are safe. Reconnecting…</p>
  {/if}

  {#if ctx.last_line}
    <p class="cue" style="margin-top:1rem">The story so far ends with…</p>
    <div class="fold" style="margin-top:0.4rem">
      <div class="script-line" class:vosd={ctx.last_line.character_role === 'vosd'}>
        {#if ctx.last_line.character_role !== 'vosd'}
          <span class="speaker">{ctx.last_line.character_name}</span>
        {/if}
        <span class="dialogue">{ctx.last_line.text}</span>
      </div>
    </div>
  {/if}

  <p class="cue" style="margin-top:1.5rem">Who speaks next?</p>
  <div class="row" style="margin-top:0.4rem">
    {#each ctx.characters as c (c.id)}
      <button
        class="chip"
        class:selected={selectedCharacter?.id === c.id}
        class:vosd-chip={c.role === 'vosd'}
        onclick={() => (selectedCharacter = c)}
      >
        {c.role === 'vosd' ? 'Stage directions' : c.name}
      </button>
    {/each}
  </div>

  <form onsubmit={addLine} style="margin-top:1rem">
    <textarea
      rows="3"
      bind:value={text}
      bind:this={textarea}
      onkeydown={onKeydown}
      disabled={over}
      placeholder={selectedCharacter
        ? selectedCharacter.role === 'vosd'
          ? 'Describe what happens on stage…'
          : `What does ${selectedCharacter.name} say?`
        : 'Pick a character first'}
    ></textarea>
    <p class="row">
      <button class="primary" type="submit" disabled={!selectedCharacter || !text.trim() || over}>
        Add line
      </button>
      <button type="button" onclick={finish}>I'm done — pass it on</button>
    </p>
  </form>

  {#if error}<p class="error">{error}</p>{/if}

  {#if ctx.my_lines?.length}
    <p class="cue" style="margin-top:1rem">Your lines this turn</p>
    <div class="paper" style="margin-top:0.4rem">
      {#each ctx.my_lines as l (l.id)}
        <div class="script-line" class:vosd={l.character_role === 'vosd'}>
          {#if l.character_role !== 'vosd'}<span class="speaker">{l.character_name}</span>{/if}
          <span class="dialogue">{l.text}</span>
        </div>
      {/each}
    </div>
  {/if}
{:else}
  <p class="muted">Waiting for the turn to begin…</p>
{/if}

<style>
  .timer {
    font-family: var(--font-script);
    font-size: 2rem;
    font-weight: 700;
    color: var(--tungsten);
    margin: 0;
  }
  .timer.low {
    color: var(--curtain);
  }
  .timer.grace {
    color: var(--curtain);
    font-size: 1.25rem;
  }
  .chip {
    border-radius: 999px;
    padding: 0.5rem 1rem;
  }
  .chip.selected {
    background: var(--tungsten);
    border-color: var(--tungsten);
    color: #241a08;
  }
  .chip.vosd-chip {
    font-style: italic;
  }
</style>
