<script>
  import { onMount, onDestroy } from 'svelte'
  import { api } from '../../lib/api.js'
  import { subscribe } from '../../lib/sse.js'
  import TurnWriter from '../../lib/TurnWriter.svelte'

  let me = $state(null) // /api/audience/me payload
  let error = $state('')
  let connection = $state('connected')
  let unsubscribe

  const member = $derived(me?.member)
  const byod = $derived(me?.station_mode === 'byod')

  async function refresh() {
    try {
      const current = await api.get('/api/audience/me')
      me = current.joined ? current : await api.post('/api/audience/join', {})
      error = ''
    } catch (err) {
      error = err.message
    }
  }

  onMount(() => {
    refresh()
    unsubscribe = subscribe(
      {
        queue_changed: refresh,
        show_state: refresh,
        your_turn: () => {
          refresh()
          navigator.vibrate?.([200, 100, 200])
        },
      },
      (state) => (connection = state)
    )
  })
  onDestroy(() => unsubscribe?.())

  async function startWriting() {
    error = ''
    try {
      await api.post(`/api/sections/${member.called_section_id}/turns`, {})
      await refresh()
    } catch (err) {
      error = err.message
    }
  }
</script>

<div class="wrap">
  {#if error}<p class="error">{error}</p>{/if}
  {#if connection === 'reconnecting'}
    <p class="error">Reconnecting… keep this page open.</p>
  {/if}

  {#if !me}
    <p class="cue">Tonight's show</p>
    <h1>Joining the queue…</h1>
  {:else if !member}
    <p class="cue">Tonight's show</p>
    <h1>The house isn't open yet.</h1>
    <p class="muted">Keep this page open — you'll get your writer number as soon as the show opens.</p>
  {:else if member.status === 'writing' && byod}
    <TurnWriter
      sectionId={member.called_section_id}
      sectionName={me.called_section?.name}
      onExit={refresh}
    />
  {:else if member.status === 'called'}
    <div class="ticket called">
      <p class="cue">{me.script_title}</p>
      <p class="number">{member.number}</p>
      <h1>You're up!</h1>
      {#if byod}
        <p>Find a quiet corner — you're writing <strong>{me.called_section?.name}</strong> right here on your phone.</p>
        <p><button class="primary big" onclick={startWriting}>Start writing</button></p>
      {:else}
        <p>Head backstage to the <strong>{me.called_section?.name}</strong> writing station and take a seat.</p>
      {/if}
    </div>
  {:else if member.status === 'done' || member.status === 'writing'}
    <div class="ticket">
      <p class="cue">{me.script_title}</p>
      <h1>Thanks for writing!</h1>
      <p class="muted">Your words are in the script. Enjoy the show — you'll hear them on stage soon.</p>
      <p><a href="#/program">See the program</a></p>
    </div>
  {:else if member.status === 'skipped'}
    <div class="ticket">
      <p class="number">{member.number}</p>
      <h1>We missed you.</h1>
      <p class="muted">Your number was called but we couldn't find you. Flag down a producer to get back in line.</p>
    </div>
  {:else}
    <div class="ticket">
      <p class="cue">{me.script_title} — you're writer number</p>
      <p class="number">{member.number}</p>
      {#if member.position === 0}
        <h2>You're next.</h2>
        <p class="muted">Stay close — your call is coming any moment.</p>
      {:else}
        <h2>{member.position} {member.position === 1 ? 'writer' : 'writers'} ahead of you.</h2>
        <p class="muted">
          Keep this page open. Your phone will buzz when it's your turn to
          {me.station_mode === 'byod' ? 'write' : 'head backstage'}.
        </p>
      {/if}
      <p><a href="#/program">Read the program while you wait</a></p>
    </div>
  {/if}
</div>

<style>
  .ticket {
    text-align: center;
    padding-top: 8vh;
  }
  /* The deli-ticket number: the one thing an usher can read from across the lobby. */
  .number {
    font-family: var(--font-script);
    font-size: clamp(6rem, 30vw, 11rem);
    font-weight: 700;
    line-height: 1;
    color: var(--tungsten);
    margin: 1rem 0;
    text-shadow: 0 0 60px rgba(232, 176, 84, 0.35);
  }
  .called h1 {
    color: var(--tungsten);
    font-size: 2.5rem;
  }
  .big {
    font-size: 1.2rem;
    padding: 0.9rem 1.6rem;
  }
</style>
