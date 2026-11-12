# Home Insurance Inventory — Project Plan

The plan for the inventory app: what it holds, what it can do, and what
is left. Tick boxes off as items land. Current state: sheet holds 92
items, ~$26,346 in known values; repo is `inventory`; four slices are
done, green, and pushed; Phase 6 scaffolding is done (database, user,
CLI root + migrate, migration, `migrate up` green, committed and
pushed). Next: the MySQL provider methods (6.3).

## Phase 1 — Core: an item goes in, the list gives it back (done)

- [X] 1.1 Item type with the V1 fields: name, manufacturer, model,
  value, description, quantity, serial number, category, purchase date
- [X] 1.2 DatabaseProvider with Repository, Save and List (in-memory)
- [X] 1.3 DatabaseProviderInterface: Save, List
- [X] 1.4 First test green on the real provider, not the mock; pushed

## Phase 2 — Search: list by category (done)

- [X] 2.1 ListByCategory on the real provider
- [X] 2.2 Test seeds Repository directly (Taylor 610CE in Guitars,
  Paradym Loud in Entertainment) and checks the filter; green, pushed

## Phase 3 — Totals by category (done)

- [X] 3.1 Totals: a map of per-category totals (Value × Quantity), no
  parameter, no filter
- [X] 3.2 Test: Music Gear 358, Computers 2070; green, pushed

## Phase 4 — CSV export (done, one follow-up)

- [X] 4.1 Export method on the provider, headers in sheet column order
- [X] 4.2 Test writes to a file in t.TempDir() and reads it back;
  green, pushed
- [X] 4.3 Three CSV-reading helpers live in the test file (not needed
  in prod code; can be moved and exported if needed)
- [ ] 4.4 Put Description back in the export (left out deliberately
  for length; decided for completeness it should be there)

## Phase 5 — Photos, serial numbers, URLs (skipped for now)

Skipped 2026-11-12 — the database backend comes first instead. Photo
storage is settled (files on disk, filename pointer in the DB, see
Phase 9); photo upload is a web feature (Phase 9). Serial numbers are
already in the struct. These items are parked, not dropped.

- [ ] 5.1 Photo support per item (deferred)
- [ ] 5.2 Website URL per item (deferred)

## Phase 6 — Database backend (MySQL)

- [X] 6.1 Schema: items table, auto-increment PK, fields matching the
  Item struct. No unique keys at this point (duplicates are not a
  concern now; revisit if others use it). Migration ran, `migrate up`
  green.
- [X] 6.2 CLI: root + migrate commands (copied from register's
  migrate), pkg config, models, repository. Database and user created.
  Build green, tests green, committed, pushed.
- [ ] 6.3 MySQL provider implementing the same interface (Save, List,
  ListByCategory, Totals, Export)
- [ ] 6.4 Swap the in-memory provider for MySQL behind the interface;
  all tests still green
- [ ] 6.5 The DB is the reference data source from here on; the sheet
  is capture only

## Phase 7 — One-shot import from the sheet

- [ ] 7.1 Finish capture in the sheet first: 6 guitars, a guitar
  synth, 2 big TVs, other guitar effects, plus furniture (the import
  waits on this)
- [ ] 7.2 Import is a CLI command reading a CSV export of the sheet,
  once into the DB (no Sheets API, no live link to the sheet)
- [ ] 7.3 Import reads by header name, not row position (columns were
  moved and rows re-sorted)
- [ ] 7.4 Import tolerates blanks: 16 rows with no purchase date, Gamo
  rifle with no value
- [ ] 7.5 Import writes whole dollars (Value is whole dollars, no
  cents)
- [ ] 7.6 No tests for the import (Rob's call — the code is mostly
  written already)
- [ ] 7.7 No unique keys; duplicates not worried about now

## Phase 8 — Web UI, basic

Web-based from the start: pick the framework up front and build basic
on it, no converting later.

- [ ] 8.1 Choose the JS framework (decide once, up front)
- [ ] 8.2 Go server serves the web app
- [ ] 8.3 List view: all items, filterable by category
- [ ] 8.4 Item detail view
- [ ] 8.5 Add and edit items manually via the web UI (after the
  import, this replaces the sheet as the entry point)
- [ ] 8.6 Color scheme: blues and natural colors

## Phase 9 — Web: photo upload

- [ ] 9.1 Upload photos to an item from the web UI
- [ ] 9.2 Photos stored on disk in local storage; the DB column holds
  a pointer to the file (a file name)
- [ ] 9.3 Item view shows its photos

## Phase 10 — Web: search and totals (later)

- [ ] 10.1 Search items by name, manufacturer, or model
- [ ] 10.2 Totals view (per-category totals, Value × Quantity)

## Phase 11 — Multi-user (later, only if others get interested)

- [ ] 11.1 Password-protected, multi-user support (separate data per
  user, possibly encrypted)

## Phase 12 — Wallet inventory (later)

- [ ] 12.1 Wallet inventory (card names and last four digits only,
  never full numbers)

## Standing decisions

- Value is the per-item (unit) price, in whole dollars, no cents. A
  row's total is Value × Quantity.
- The spreadsheet is one-time capture only, NOT a live source. Read it
  once into the DB, then code works off the DB. No Sheets API.
- The DB is the reference data source.
- Finish capture first, then one import. No re-runnable import. No
  tests for the import. No unique keys; duplicates deferred until
  others use it.
- The CLI covers migrations and the CSV sheet import.
- After the import, new items go in manually via the web app; the sheet
  retires.
- Photos: files on disk in local storage; the DB column holds a pointer
  to the file (a file name). Upload is a web feature.
- Web-based from the start: pick the JS framework up front, build basic
  on it, no converting later. Colors: blues and natural colors.
- Plaintext config with secrets is never committed; config.json stays
  gitignored, encrypted config.json.age like register.
- Rob writes the code; Muse guides and second pair of eyes. Hints
  first, full answer only on "just tell me".
