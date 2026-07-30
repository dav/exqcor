<script>
  import { onMount } from 'svelte'
  import AdminGate from '../../lib/AdminGate.svelte'
  import { api } from '../../lib/api.js'

  let net = $state(null)
  let scripts = $state([])

  onMount(async () => {
    try {
      ;[net, scripts] = await Promise.all([api.get('/api/netinfo'), api.get('/api/scripts')])
    } catch {
      // panels degrade to empty states
    }
  })
</script>

<AdminGate>
  <div class="wrap">
    <p class="cue">Production desk</p>
    <h1>Exqcor Admin</h1>

    <div class="panel">
      <p class="cue">Tonight</p>
      <div class="row" style="margin-top:0.5rem">
        <a class="btn primary" href="#/admin/show">Run of show</a>
        <a class="btn" href="#/admin/scripts">Scripts</a>
        <a class="btn" href="#/admin/qr">QR codes to print</a>
      </div>
    </div>

    <div class="panel">
      <p class="cue">Network</p>
      {#if net}
        <p style="margin:0.5rem 0 0">
          Phones connect at <strong>{net.base_url}</strong>
          {#if net.candidates?.length > 1}
            <span class="muted">— multiple networks detected; confirm on the QR page.</span>
          {/if}
        </p>
      {:else}
        <p class="muted" style="margin:0.5rem 0 0">Checking the network…</p>
      {/if}
    </div>

    <div class="panel">
      <p class="cue">Scripts</p>
      {#if scripts.length === 0}
        <p class="muted" style="margin:0.5rem 0 0">
          No scripts yet. <a href="#/admin/scripts">Create the first one</a> to set up a show.
        </p>
      {:else}
        <table style="margin-top:0.5rem">
          <tbody>
            {#each scripts as s (s.id)}
              <tr>
                <td><a href="#/admin/scripts/{s.id}">{s.title}</a></td>
                <td class="muted">{s.theme}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>
  </div>
</AdminGate>
