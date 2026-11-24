# jot

A minimal, keyboard-driven daily todo app for the terminal. Everything you jot down is filed under the *logical day* it belongs to—or kept in an undated "no-day" list.

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

## Views

Press `w` to cycle between three views:

1. **Day view** — todos for the selected logical day
2. **Week overview** — summary of the week (`h/l` prev/next week, `enter` open day)
3. **No-day view** — undated todos not tied to a specific day

## Keys

    j/k, ↑/↓      move cursor
    J/K           reorder item up / down
    h/l, ←/→      previous / next day (or week in overview)
    enter/space   toggle done
    a             add
    e             edit
    d             delete (press d again to confirm)
    m             move between day and no-day (press m again to confirm)
    s             cycle sort (undone, done, a-z, z-a, manual)
    w             cycle view (day → week overview → no-day)
    t             back to today
    ?             toggle help
    q             quit

## Data

`~/.config/jot/store.json`, keyed by date (and `"no-day"` for undated items):

    {
      "2026-11-12": [{ "title": "Buy carrots", "done": true }],
      "2026-11-13": [{ "title": "Write report", "done": false }],
      "no-day": [{ "title": "Someday task", "done": false }]
    }
