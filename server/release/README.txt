EXQCOR — exquisite corpse theater, in a single program
=======================================================

WHAT THIS IS
  The whole show server: audience queue, writing stations, timers, and
  script printing. It runs on this laptop; everyone else joins from their
  phone's browser by scanning a QR code you print.

RUNNING IT
  macOS/Linux:  double-click `exqcor` (or run ./exqcor in a terminal)
  Windows:      double-click `exqcor.exe`

  A browser window opens with the admin console. The first run prints an
  admin passphrase in the terminal window — write it down. Your show data
  lives in `exqcor.db` next to the program; copy that file to back up the
  production.

FIRST-RUN WARNINGS (once per laptop)
  macOS:    "Apple could not verify..." — right-click the file, choose Open,
            then Open again. Or in a terminal:
            xattr -d com.apple.quarantine ./exqcor
  Windows:  SmartScreen: click "More info", then "Run anyway".
  Firewall: when asked whether to allow incoming connections, click ALLOW.
            If you decline, phones cannot reach the show.

WI-FI
  The laptop and every phone must be on the same Wi-Fi network. Guest
  networks that isolate clients will not work — a cheap travel router or a
  phone hotspot is the reliable choice. Test with one phone before the
  house opens.

MORE
  The full show-night runbook: https://github.com/dav/exqcor/blob/master/docs/SHOW-NIGHT.md
