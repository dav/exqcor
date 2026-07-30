# Exqcor

> **A note from the main developer:**
>
> *The Exquisite Corpse collaborative stage performance system is the
> brainchild of the [late, great Mikl
> Em](https://www.kqed.org/arts/13885565/too-weird-to-live-too-rare-to-die-a-tribute-to-michael-mikl-em-mcelligott).
> I admire Mikl so much and was honored to be asked to implement the first
> version of this system, which was used for three absolutely wonderful
> [production runs](https://www.facebook.com/ExCorpsePlay/) at the Stage
> Werx in San Francisco ~ 2014. I am bringing back this system in hopes that
> other production troupes can make use of it to create amazing
> collaborations in their own communities. I believe it is what Mikl would
> have wanted.*

Software for **exquisite corpse theater**: the audience writes a play one
person at a time — each writer sees only the last line the previous writer
left behind — and improv actors perform each act cold, minutes after it's
written.

## How it works (audience view)

Picture this: the audience arrives at the theater for a play with a specific
theme, e.g. Film Noir. Instead of heading to their seats though, they mingle
for 20 minutes or so with other audience members and the cast of the play on
stage where there is food and drink. The actors (the hard-boiled private
detective, the mysterious "victim", the secretary, etc) are in costume and in
character, improv-ing with each other and the audience.

Each audience member is given a number in a queue. When their number is
called they are given access to a script writing system (either led to a
laptop, or in the modern version, via their phone) where for 5 minutes they
can select characters from a menu and type in dialogue and stage directions.
When they start, they only see the last line from the previous
script-writer's input (this is the exquisite corpse part).

When enough lines have been written for act one, the audience takes their
seat and the play begins. The actors are cold-reading the script, written by
the audience, live on stage. Acts two and three continue to collect lines
while act one is going. In the dozen or so shows we did, it always ended up
as exquisite chaos where everyone has a great time — especially when their
five minutes of contributions hit the stage!

## How it works (production view)

This was 100% Mikl's department, so take my advice here with that grain of
salt. This is from what I recall.

**Pre-production, cast.** Obviously you start with a strong cast of stage
performers comfortable with improv. Pick a theme, assign the characters.
Since obviously you do not have the script, rehearsals are about exploring
these characters and their general relationships to each other, so that you
can use the initial mingling period to guide the audience's imagination. As
I recall, our productions generally had a webpage with photos of the
characters and a little bit of color for them. Also useful to practice cold
reading from random scripts.

**Pre-production, tech.** The Exqcor system is designed to run centralized
on one laptop or PC that is on a Wi-Fi network at the theater. Depending on
how you want to do it, either (A) the audience members' phones or (B) a
laptop/PC for each act are also connected to the same Wi-Fi network. In the
case of using the audience phones, you provide a QR code that will load the
web app for them on their phones. Before the show starts you create a new
production in Exqcor and add in the characters and props that will be
available, set the number of acts and other configurations as needed.

**Show-time tech.**

- *Collecting audience lines* — the old school way was to give each audience
  member a number and call them one at a time to sit down at one of the
  laptops off stage. If using the audience phones, then this is handled
  through the web app on their phone: they just need to enter the queue and
  the UI will let them know when their time has come. An admin view can be
  followed by someone in the production to monitor line collection and in
  general unblock things so lines are being collected continuously. Note
  that if using audience phones, they can continue to contribute from their
  seats.
- *Closing an act* — the admin also can make the call when an act is ready
  to be closed and distributed to the actors.
- The stage scripts can either be printed to paper (one for each actor with
  their lines and stage directions highlighted) or viewed by the actors on
  their own phones/tablets.

## How it runs

One laptop runs a single `exqcor` executable (Windows, macOS, or Linux — no
installer, no runtime). It serves the whole show over the venue's Wi-Fi:

- **Admin** opens the production desk in a browser: set up the show, run the
  queue, print scripts.
- **QR codes** are printed from the admin console: audience, writing-station,
  and cast codes each open the right interface on any phone/tablet/laptop
  browser (Android, iOS, anything).
- **Audience** members scan in, get a writer number, and their phone buzzes
  when they're called — to a backstage writing station, or to write right on
  their phone, per production.
- Timers, the hand-off, and the queue update live (SSE). Every line is saved
  the moment it's written; the server can be power-cycled mid-show without
  losing anything.

## Repository layout

- `server/` — Go server: SQLite storage, show runtime, HTTP/SSE API, and the
  embedded web UI. One static binary per platform.
- `web/` — Svelte SPA (all four surfaces: admin, station, audience, script
  views), built into `server/internal/webui/dist` and embedded.
- `docs/SHOW-NIGHT.md` — the operator runbook. Print it.
- `original/` — the 2012–2014 Rails 3 app (reference).
- `modern/` — a Rails 8 rewrite of the data model/API (reference).
- `Apple/` — an abandoned SwiftUI client experiment (reference).

## Development

```sh
make dev-server   # Go API on :8080 (uses tmp/dev.db)
make dev-web      # Vite dev server on :5173, /api proxied to :8080
make test         # Go tests
make build        # production build: web + single binary at server/exqcor
```

Releases: push a `v*` tag; GitHub Actions builds the frontend, runs tests,
and GoReleaser publishes archives for Linux (amd64/arm64), macOS (universal),
and Windows.
