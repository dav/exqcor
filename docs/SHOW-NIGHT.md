# Exqcor Show-Night Runbook

Print this. Tape it to the production desk.

## Before the house opens (30+ minutes)

1. **Wi-Fi.** Put the server laptop and a test phone on the same Wi-Fi
   network. Venue Wi-Fi with client isolation ("guest" networks often have
   it) will silently block phones — if the test phone can't load the page, use
   a cheap travel router or a phone hotspot instead, and put everything on it.
2. **Start the server.** Double-click `exqcor` (macOS/Linux) or `exqcor.exe`
   (Windows). A browser opens the admin page. **If the OS firewall asks,
   click Allow** — otherwise phones cannot connect.
3. **First run only:** the terminal window shows the admin passphrase. Write
   it down. You need it to reach the admin pages from any device that isn't
   the server laptop.
4. **Open the show.** Admin → Run of show → open tonight's script. (Set up
   the script, characters, casting, acts, and opening lines in advance;
   duplicate a past show to reuse its structure.)
5. **Print QR codes.** Admin → QR codes to print. If the page warns about
   multiple networks, pick the one the phones are on. Print; post the
   audience sheet in the lobby, keep station/cast sheets backstage.
6. **Set up writing stations** (station mode): on each station laptop or
   tablet, scan the station QR, then pick which act that station serves.
7. **Test end-to-end** with one phone: scan the audience QR → get number →
   call it from Run of show → write a line at the station → end the turn.
   Then requeue or ignore that test entry.

## During the show

- **Audience arrives:** they scan the lobby QR and get a number. The Run of
  show page shows the queue live.
- **Calling writers:** tap **Call next writer here** on the act that has a
  free seat. The person's phone buzzes and tells them where to go (or, in
  phones-only mode, lets them start writing right where they sit).
- **No-shows:** if a called number doesn't appear, mark **No-show** and call
  the next. You can put no-shows **Back in line** later.
- **Cutting a turn short:** **End turn now** on the act's card. The writer's
  saved lines are kept; their unfinished sentence is not.
- **When an act is done:** stop calling writers to it and **Mark section
  complete**. Open the script view, hit **Print**, and hand it to the cast.
  Per-actor sides are on the same page under "Print a single actor's script."

## If things go wrong

- **Server laptop crashes or reboots:** start `exqcor` again. Everything —
  queue, numbers, scripts, the turn that was in progress and its timer —
  comes back from the database. Phones reconnect on their own.
- **A phone loses Wi-Fi:** the page shows "Reconnecting…" and recovers by
  itself; their number is tied to the device, so even a full reload keeps it.
- **Wrong network on the QR page:** pick the right one on the QR page and
  reprint. Old sheets point at the wrong address, so bin them.
- **A QR sheet leaks somewhere it shouldn't** (e.g. the station QR posted in
  the lobby): QR page → **Invalidate & regenerate** that sheet, reprint.
- **Lost admin passphrase:** on the server laptop itself the admin pages
  always open without one (localhost is trusted). To reset it properly, stop
  the server, delete the `admin_pass_hash` row: 
  `sqlite3 exqcor.db "DELETE FROM settings WHERE key='admin_pass_hash'"`,
  start again, and a fresh passphrase is printed.
- **Nuclear option:** the whole show lives in `exqcor.db` next to the
  program. Copy that file and you've backed up the production; move it to
  another laptop and start `exqcor` there to migrate mid-show (re-print QR
  codes — the address changes).

## Numbers to know

- Writing timer: set per script (default 5:00). Writers get 15 grace seconds
  after zero to finish a sentence.
- Every line is saved the moment the writer adds it. A dead battery loses at
  most the sentence being typed.
