<script>
  import { onMount, onDestroy } from 'svelte'
  import { api } from '../../lib/api.js'
  import { subscribe } from '../../lib/sse.js'
  import TurnWriter from '../../lib/TurnWriter.svelte'

  // Station state machine: pick-section → idle → writing → handoff → idle.
  let show = $state(null)
  let sectionId = $state(localStorage.getItem('exqcor_station_section') || '')
  let mode = $state('idle') // idle | writing | handoff
  let writerName = $state('')
  let error = $state('')
  let unsubscribe

  const section = $derived(show?.sections?.find((s) => String(s.id) === String(sectionId)))

  async function loadShow() {
    try {
      show = await api.get('/api/show')
      // Page reload mid-turn: rejoin the active turn.
      if (mode === 'idle' && section?.active_turn_id) mode = 'writing'
      error = ''
    } catch (err) {
      error = err.message
    }
  }

  onMount(() => {
    loadShow()
    unsubscribe = subscribe({
      show_state: loadShow,
      turn_started: loadShow,
      turn_ended: loadShow,
    })
  })
  onDestroy(() => unsubscribe?.())

  function chooseSection(id) {
    sectionId = String(id)
    localStorage.setItem('exqcor_station_section', sectionId)
    mode = 'idle'
    loadShow()
  }

  async function startTurn(e) {
    e.preventDefault()
    error = ''
    try {
      await api.post(`/api/sections/${sectionId}/turns`, { writer_name: writerName })
      mode = 'writing'
    } catch (err) {
      error = err.message
    }
  }
</script>

<div class="wrap">
  {#if !sectionId || (show?.open && !section)}
    <p class="cue">Writing station</p>
    <h1>Which act is this station for?</h1>
    {#if show?.open}
      <div class="row" style="margin-top:1rem">
        {#each show.sections as s (s.id)}
          <button class="primary" onclick={() => chooseSection(s.id)}>{s.name}</button>
        {/each}
      </div>
    {:else}
      <p class="muted">No show is open yet. The admin opens one from the production desk.</p>
    {/if}
  {:else if mode === 'writing'}
    <TurnWriter
      {sectionId}
      sectionName={section?.name}
      onExit={() => {
        mode = 'handoff'
        writerName = ''
        loadShow()
      }}
    />
  {:else if mode === 'handoff'}
    <div class="handoff">
      <p class="cue">{section?.name}</p>
      <h1>Stand up!</h1>
      <p>Thanks for writing. Send the next writer over.</p>
      <p><button class="primary" onclick={() => (mode = 'idle')}>Ready for the next writer</button></p>
    </div>
  {:else}
    <p class="cue">{section?.name} — writing station</p>
    <h1>Take a seat.</h1>
    <p class="muted">
      You'll see the last line of the story so far, and only that. Write what
      happens next — dialogue, stage directions, whatever the story needs.
      When the timer runs out, your last line finishes and you pass it on.
    </p>
    <form onsubmit={startTurn}>
      <label for="wname">Your name (for the program)</label>
      <input id="wname" type="text" bind:value={writerName} placeholder="Optional" />
      {#if error}<p class="error">{error}</p>{/if}
      <p><button class="primary big" type="submit">Start writing</button></p>
    </form>
    <p class="muted no-print" style="margin-top:2rem">
      <button onclick={() => chooseSection('')}>Change act</button>
    </p>
  {/if}
</div>

<style>
  .big {
    font-size: 1.2rem;
    padding: 0.9rem 1.6rem;
  }
  .handoff {
    text-align: center;
    padding-top: 15vh;
  }
  .handoff h1 {
    font-size: 3rem;
    color: var(--tungsten);
  }
</style>
