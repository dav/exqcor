<script>
  import { onMount } from 'svelte'
  import AdminGate from '../../lib/AdminGate.svelte'
  import { api } from '../../lib/api.js'

  let net = $state(null)
  let tokens = $state(null)
  let error = $state('')

  const sheets = $derived(
    tokens
      ? [
          {
            role: 'audience',
            title: 'Join tonight’s show',
            instructions:
              'Scan with your phone camera to get your writer number. Keep the page open — it buzzes when you’re up.',
            url: tokens.audience,
          },
          {
            role: 'station',
            title: 'Writing station',
            instructions:
              'Backstage only. Scan on each writing-station laptop or tablet, then pick which act it serves.',
            url: tokens.station,
          },
          {
            role: 'actor',
            title: 'Cast & crew',
            instructions:
              'For actors: opens the program now, and your script pages once acts are written.',
            url: tokens.actor,
          },
        ]
      : []
  )

  async function load() {
    try {
      ;[net, tokens] = await Promise.all([api.get('/api/netinfo'), api.get('/api/tokens')])
    } catch (err) {
      error = err.message
    }
  }
  onMount(load)

  async function pickIP(ip) {
    error = ''
    try {
      await api.post('/api/netinfo', { ip })
      await load()
    } catch (err) {
      error = err.message
    }
  }

  async function regenerate(role) {
    if (!confirm(`Print new ${role} QR codes? Every ${role} QR already printed stops working.`)) return
    try {
      tokens = await api.post('/api/tokens/regenerate', { role })
    } catch (err) {
      error = err.message
    }
  }
</script>

<AdminGate>
  <div class="wrap">
    <div class="no-print">
      <p class="cue"><a href="#/admin">Production desk</a> / QR codes</p>
      <h1>QR codes to print</h1>
      {#if error}<p class="error">{error}</p>{/if}

      {#if net?.candidates?.length > 1}
        <div class="panel">
          <p class="cue">This laptop is on more than one network</p>
          <p class="muted" style="margin:0.5rem 0">
            Pick the Wi-Fi network the audience's phones will be on. The QR
            codes below point at it.
          </p>
          <div class="row">
            {#each net.candidates as c (c.ip)}
              <button
                class:primary={net.base_url.includes(c.ip)}
                onclick={() => pickIP(c.ip)}
              >
                {c.interface} — {c.ip}{c.default ? ' (likely)' : ''}
              </button>
            {/each}
          </div>
        </div>
      {/if}

      <p class="muted">
        Phones must be on the same Wi-Fi as this laptop
        {#if net}(<strong>{net.base_url}</strong>){/if}. Print this page —
        each QR gets its own sheet. Post the audience sheet in the lobby;
        keep the station and cast sheets backstage.
      </p>
      <p class="row">
        <button class="primary" onclick={() => window.print()}>Print all three sheets</button>
      </p>
    </div>

    {#each sheets as sheet (sheet.role)}
      <section class="sheet paper">
        <p class="cue">Exquisite corpse theater</p>
        <h2 class="sheet-title">{sheet.title}</h2>
        <img
          class="qr"
          src={`/api/qr.png?url=${encodeURIComponent(sheet.url)}`}
          alt={`QR code: ${sheet.url}`}
        />
        <p class="sheet-url">{sheet.url}</p>
        <p class="sheet-instructions">{sheet.instructions}</p>
        <p class="no-print">
          <button class="danger" onclick={() => regenerate(sheet.role)}>
            Invalidate & regenerate this QR
          </button>
        </p>
      </section>
    {/each}
  </div>
</AdminGate>

<style>
  .sheet {
    text-align: center;
    margin: 1.5rem 0;
    padding: 2.5rem 1.5rem;
  }
  .sheet-title {
    font-size: 2rem;
    margin: 0.5rem 0 1rem;
  }
  .qr {
    width: min(70vw, 320px);
    height: auto;
    image-rendering: pixelated;
  }
  .sheet-url {
    font-family: var(--font-script);
    font-size: 0.95rem;
    word-break: break-all;
    color: var(--ink-dim);
  }
  .sheet-instructions {
    max-width: 28rem;
    margin: 0.75rem auto 0;
  }
  @media print {
    .sheet {
      break-after: page;
      margin: 0;
    }
    .qr {
      width: 9cm;
    }
  }
</style>
