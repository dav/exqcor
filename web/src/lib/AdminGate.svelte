<script>
  import { onMount } from 'svelte'
  import { api } from './api.js'

  let { children } = $props()

  let state = $state('checking') // checking | locked | open
  let passphrase = $state('')
  let error = $state('')

  onMount(async () => {
    try {
      const { role } = await api.get('/api/session')
      state = role === 'admin' ? 'open' : 'locked'
    } catch {
      state = 'locked'
    }
  })

  async function unlock(e) {
    e.preventDefault()
    error = ''
    try {
      await api.post('/api/session/admin', { passphrase })
      state = 'open'
    } catch (err) {
      error = err.message
    }
  }
</script>

{#if state === 'open'}
  {@render children()}
{:else if state === 'locked'}
  <div class="wrap">
    <p class="cue">Stage door</p>
    <h1>Admin access</h1>
    <p class="muted">
      Enter the admin passphrase. It was shown in the server window when the
      show was first set up.
    </p>
    <form onsubmit={unlock}>
      <label for="pass">Passphrase</label>
      <input id="pass" type="password" bind:value={passphrase} autocomplete="current-password" />
      {#if error}<p class="error">{error}</p>{/if}
      <p><button class="primary" type="submit">Unlock</button></p>
    </form>
  </div>
{/if}
