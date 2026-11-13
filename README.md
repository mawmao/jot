# jot

A minimal, keyboard-driven daily todo app for the terminal. Everything you jot down is filed under the *logical day* it belongs to.

## Install

    go install github.com/mawmao/jot@latest

Make sure `~/go/bin` is on your `PATH`.

## Day boundary

By default a day starts at midnight. If you work past midnight and
want a 2am task to count as "yesterday", set the boundary in
`~/.config/jot/config.json`:

    { "day_start": "04:00" }

Or override per-run:

    jot --day-start=04:00

## Keys

    j/k, ↑/↓      move
    h/l, ←/→      previous / next day
    enter/space   toggle done
    a             add
    e             edit
    d             delete (press d again to confirm)
    t             back to today
    ?             toggle help
    q             quit

## Data

`~/.config/jot/store.json`, keyed by date:

    {
      "2026-11-12": [{ "title": "Buy carrots", "done": true }],
      "2026-11-13": [{ "title": "Write report", "done": false }]
    }
