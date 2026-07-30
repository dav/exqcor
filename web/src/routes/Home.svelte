<script>
  import { onMount } from 'svelte'
  import { replace } from 'svelte-spa-router'
  import { api } from '../lib/api.js'

  // Scanning a QR lands on /j/<token> which redirects to a role home; hitting
  // the bare origin sends known roles onward and greets everyone else.
  onMount(async () => {
    try {
      const { role } = await api.get('/api/session')
      if (role === 'admin') replace('/admin')
      else if (role === 'station') replace('/write')
      else if (role === 'actor') replace('/program')
      else if (role === 'audience') replace('/audience')
    } catch {
      // stay on the landing page
    }
  })
</script>

<div class="wrap">
  <p class="cue">Exquisite Corpse Theater</p>
  <h1>Exqcor</h1>
  <p class="muted">
    To join tonight's show, scan the QR code posted in the lobby. If you're
    running the show, open this page on the server laptop to reach the admin
    console.
  </p>
</div>
