lk# Home Insurance Inventory — Project Plan

The plan for the inventory app: what it holds, what it can do, and what
is left. Tick boxes off as slices land. Current state: sheet holds 92
items, ~$26,346 in known values; repo is `inventory`; four slices are
done, green, and pushed.

## Phase 1 — Core: an item goes in, the list gives it back (done)

- [x] Item type with the V1 fields: name, manufacturer, model, value,
  description, quantity, serial number, category, purchase date
- [x] DatabaseProvider with Repository, Save and List (in-memory)
- [x] DatabaseProviderInterface: Save, List
- [x] First test green on the real provider, not the mock; pushed

## Phase 2 — Search: list by category (done)

- [x] ListByCategory on the real provider
- [x] Test seeds Repository directly (Taylor 610CE in Guitars, Paradym
  Loud in Entertainment) and checks the filter; green, pushed

## Phase 3 — Totals by category (done)

- [x] Totals: a map of per-category totals (Value × Quantity), no
  parameter, no filter
- [x] Test: Music Gear 358, Computers 2070; green, pushed

## Phase 4 — CSV export (done, one follow-up)

- [x] Export method on the provider, headers in sheet column order
- [x] Test writes to a file in t.TempDir() and reads it back; green,
  pushed
- [x] Three CSV-reading helpers live in the test file (not needed in
  prod code; can be moved and exported if needed)
- [ ] Put Description back in the export (left out deliberately for
  length; decided for completeness it should be there)

## Phase 5 — Photos, serial numbers, URLs

- [ ] Photo support per item (storage decision still open: files on
  disk with paths in the record)
- [ ] Website URL per item
- [ ] Serial numbers are already in the struct; record them from
  photos and screenshots as they come in

## Phase 6 — MySQL store

- [ ] MySQL provider implementing the same interface (Save, List,
  ListByCategory, Totals, Export)
- [X] Schema: items table with an auto-increment PK and fields matching
  the Item struct
- [ ] Swap the in-memory provider for MySQL behind the interface; all
  tests still green

## Phase 7 — Finish capture, then the one-shot import

- [ ] Finish capture in the sheet: 6 guitars, a guitar synth, 2 big
  screen TVs, other guitar effects, plus furniture
- [ ] Resolve the Costco TV (bought in-store; no order email exists)
- [ ] Gamo air rifle value (its order total covers several items; sort
  out or send the order-details view)
- [ ] Confirm the Swift shock mount (confirmation email lists it at
  $89; Sweetwater's site invoice does not)
- [ ] One-shot import: read the sheet once into the DB (no Sheets API,
  no live link to the sheet)
- [ ] Import tolerates blanks: 16 rows with no purchase date, Gamo
  rifle with no value
- [ ] Import reads by header name, not row position (columns were moved
  and rows re-sorted)
- [ ] Import writes whole dollars: Value is whole dollars, no cents;
  the sheet itself goes to ints
- [ ] Tests after the fact for the import (Rob's call)

## Phase 8 — Manual entry via the app

- [ ] After the import, new items are entered manually via the app;
  the sheet retires as an entry point

## Later / maybe

- [ ] Web UI (responsive, works on iPhone and iPad; password-protected
  multi-user SPA only if family or friends want in)
- [ ] Wallet inventory (card names and last four digits only, never
  full numbers)
- [ ] Multiple photos per item

## Standing decisions

- Value is the per-item (unit) price, in whole dollars, no cents. A
  row's total is Value × Quantity.
- The spreadsheet is one-time capture only, NOT a live source. Read it
  once into the DB, then code works off the DB. No Sheets API.
- Finish capture first, then one import. No re-runnable import.
- After the import, new items go in manually via the app; the sheet
  retires.
- Tests after the fact for the import are Rob's call, not TDD-first.
- Rob writes the code; Muse guides and second pair of eyes. Hints
  first, full answer only on "just tell me".
