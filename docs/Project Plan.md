# Home Insurance Inventory — Project Plan

The plan for the inventory app: what it holds, what it can do, and what
is left. Tick boxes off as items land. Keep tasks small enough to use as
git commit messages; phases and tasks may be reordered as priorities change.

Current confirmed status: MySQL is the source of truth; Goose handles
migrations; the CSV import is implemented, skips duplicates, and works.
The uploaded plan records 93 imported items as of 2026-11-14. The web
list/category filter and first PDF export are implemented. The existing
checklist states are retained below; newer work is listed separately.

## Phase 1 — Core: an itemPhoto goes in, the list gives it back (done)

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

## Phase 4 — CSV export (done)

- [X] 4.1 Export method on the provider, headers in sheet column order
- [X] 4.2 Test writes to a file in t.TempDir() and reads it back;
  green, pushed. Rewritten 2026-11-14 as Test_ItemsExportToCSVFile
  in cmd/export_test.go: header check via stdlib csv, dates,
  TotalValue; green, pushed.
- [X] 4.3 Three CSV-reading helpers live in the test file (not needed
  in prod code; can be moved and exported if needed)
- [X] 4.4 Description back in the export (was left out deliberately
  for length; added 2026-11-12 in the import/export rework). Export
  now writes Value and Total Cost as separate columns (2026-11-14),
  fixing the old row-total-under-Value-header trap.

## Phase 5 — Photos, serial numbers, URLs (skipped for now)

Skipped 2026-11-12 — the database backend comes first instead. Photo
storage is settled (files on disk, filename pointer in the DB, see
Phase 9); itemPhoto upload is a web feature (Phase 9). Serial numbers are
already in the struct. These items are parked, not dropped.

- [ ] 5.1 Photo support per itemPhoto (deferred)
- [ ] 5.2 Website URL per itemPhoto (deferred)

## Phase 6 — Database backend (MySQL)

- [X] 6.1 Schema: items table, auto-increment PK, fields matching the
  Item struct. Unique key is (name, description) (was name alone
  2026-11-12, then (name, model); changed 2026-11-14 because generic
  Names collide and Details is TEXT, which MySQL can't index whole).
  Field lengths tightened 2026-11-14: name varchar(30), description
  varchar(60), model varchar(60), details TEXT. Baseline rewritten
  (not a 00002) because the DB isn't live yet — drop and rebuild is
  fine until then. Migration ran, `migrate up` green; DB holds all
  93 rows.
- [X] 6.2 CLI: root + migrate commands (copied from register's
  migrate), pkg config, models, repository. Database and user created.
  Build green, tests green, committed, pushed.
- [X] 6.3 MySQL provider implementing the same interface (Save,
  SaveItems, List, ListByCategory, TotalByCategory,
  SortedListForOutput — renamed from SortedListForCSV). Methods exist
  and work against real MySQL (import loaded 93 rows via SaveItems;
  export reads via SortedListForOutput). Ticked 2026-11-14: the
  methods implement the interface and work; the automated real-DB
  test is separate work (the inventory_test database, after the
  duplicate skip).
- [X] 6.4 Swap the in-memory provider for MySQL behind the interface;
  all tests still green. Commands now build the MySQL provider via
  NewMySQLProvider (import and export); the in-memory provider is no
  longer used by the commands.
- [X] 6.5 The DB is the reference data source from here on; the sheet
  is capture only. DB holds all 93 rows and is what export/PDF read
  from.
- [X] 6.6 Automated tests against a real MySQL database
  (inventory_test). Tests live in pkg/repository/mysql/mysql_test.go
  (the mysql package's test file); the fake-based provider tests
  stay in database_provider_test.go. One-time setup: CREATE DATABASE
  inventory_test + grant. Test setup builds a DBConfig with Name
  "inventory_test", connects, runs goose migrations up, truncates
  items; then seed out of order and assert order — proves the
  provider methods against real MySQL instead of the fake.
  (Password: env var, or read config.json and swap the Name.)

## Phase 7 — CSV import from the sheet

- [ ] 7.1 Finish capture in the sheet first: 6 guitars, a guitar
  synth, 2 big TVs, other guitar effects, plus furniture (the import
  waits on this). Note: the import already ran on the 93 captured
  rows (2026-11-14) before capture finished; the remaining items
  still need a path into the DB (open question — web UI 8.5, or
  another import).
- [X] 7.2 Import is a CLI command reading a CSV export of the sheet,
  once into the DB (no Sheets API, no live link to the sheet).
  Command exists and ran 2026-11-14 (93 rows landed).
- [X] 7.3 Import reads by header name, not row position (columns were
  moved and rows re-sorted) — proven by the 2026-11-14 import, which
  handled the Value column move.
- [X] 7.4 Import tolerates blanks: 16 rows with no purchase date, Gamo
  rifle with no value — proven, all 93 rows landed.
- [X] 7.5 Import writes whole dollars (Value is whole dollars, no
  cents) — proven; CurrencyInt parses the accounting-formatted
  values ($ 199) and the DB holds whole dollars.
- [X] 7.6 No tests for the import (Rob's call — the code is mostly
  written already)
- [X] 7.7 (name, description) is the unique key (see 6.1); a duplicate
  on import is skipped — the skip must be visible (a count or an
  error), not silent. The silent skip bit us 2026-11-14 (two Desks
  and two Sheet Sets with empty models collided under the old (name,
  model) key; the (name, description) key fixed it).

## Phase 8 — Web UI, basic

Web-based from the start: pick the framework up front and build basic
on it, no converting later.

- [X] 8.1 Choose the JS framework (decide once, up front). Decided
  2026-11-14: vanilla JS for now, to refresh his memory on JS
  first; can upgrade to a framework later.
- [X] 8.2 Go server serves the web app
- [X] 8.3 List view: all items, filterable by category
- [X] 8.4 Item detail view
- [X] 8.5 Add and edit items manually via the web UI (after the
  import, this replaces the sheet as the entry point)
- [X] 8.6 Color scheme: blues and natural colors

## Phase 9 — Web: itemPhoto upload

- [X] 9.1 Upload photos to an itemPhoto from the web UI
- [X] 9.2 Photos stored on disk in local storage; the DB column holds
  a pointer to the file (a file name)
- [X] 9.3 Item view shows its photos

## Phase 10 — Web: search and totals (later)

- [X] 10.1 Search items by name, manufacturer, or model
- [X] 10.2 Totals view (per-category totals, Value × Quantity)

## Phase 11 — Multi-user (later, only if others get interested)

- [ ] 11.1 Password-protected, multi-user support (separate data per
  user, possibly encrypted)

## Phase 12 — Wallet inventory (later)

- [ ] 12.1 Wallet inventory (card names and last four digits only,
  never full numbers)

## Phase 13 — PDF export (later)

Rob's idea 2026-11-12: an additional export format so output looks
cleaner and more professional, handling wrapping of long column
values. PDF is his lean (a learning challenge; he has done PDF before
and remembers inches-to-points conversion); doc or Google Sheets also
possible.

- [X] 13.1 Export items to PDF, wrapping long column values (first
  version built 2026-11-12 with gpdf: landscape, header with
  name/address/date, Confidential footer; the $ left-alignment is
  still being tweaked, and Purchase Date / Serial Number were dropped
  from this version). Description column added 2026-11-14 with
  adjusted widths. Known gpdf bug: when a column value wraps, the
  stripe background isn't replicated across the other columns on the
  wrapped line (Rob may PR a fix). Idea: a second, roomier PDF for
  Details, which is too long for the dense table.

## Phase 14 — HTTP handlers: functional tests (implemented; coverage reconciliation pending)

Existing automated tests in `server_test.go` confirmed by Rob on 2026-11-22:
SaveItem handler — invalid JSON, success, duplicate, error, and bad request;
DeleteItem handler — success and error; ListItem handler — test present.
These tests were previously omitted from the checklist status. The specific
scenarios below remain unchecked until the tests are matched against each
requirement; unchecked does not mean no tests exist.

- [ ] 14.1 Inventory existing routes, handlers, request/response formats, and error paths
- [ ] 14.2 Add `httptest` coverage for listing all items and filtering by category (list handler test exists; category case to verify)
- [ ] 14.3 Test itemPhoto detail endpoint: found, not found, and invalid ID (verify whether a separate detail route exists)
- [ ] 14.4 Test create-itemPhoto endpoint: valid payload, missing fields, invalid values, and duplicate (save handler tests cover success, invalid JSON, bad request, duplicate; check create-specific cases)
- [ ] 14.5 Test update-itemPhoto endpoint: valid changes, not found, invalid payload, and duplicate (save handler tests exist; check update-specific cases)
- [ ] 14.6 Test handler responses for database failures and expected HTTP status codes (save and delete error tests exist; verify complete coverage)
- [ ] 14.7 Ensure handler tests are isolated and deterministic; document how to run them

## Phase 15 — Database and migration test coverage (reported complete; verification pending)

Rob reports on 2026-11-22 that Phase 15's repository coverage, empty-result and duplicate-key tests, totals/sorting tests, test database isolation, migration checks, and test instructions have been done. Retain these as reported complete pending a quick source review; this is not a new test run.

- [X] 15.1 Audit MySQL repository test coverage against every public repository method
- [X] 15.2 Add tests for empty results, duplicate-key behavior, and invalid data
- [X] 15.3 Verify totals and sorting with multiple quantities and tied values
- [X] 15.4 Make `inventory_test` setup and cleanup repeatable without touching production data
- [X] 15.5 Verify Goose migrations against a clean test database
- [X] 15.6 Document test prerequisites and a single command for running the suite

## Phase 16 — Frontend code organization and usability (partially complete; review later)

Rob reports that portions of Phase 16 are complete, including whole-dollar/thousands-separator formatting. Responsive narrow-screen table work (16.5) is not currently a priority; revisit later.

- [ ] 16.1 Review JavaScript responsibilities: API requests, state, rendering, and events
- [ ] 16.2 Separate reusable formatting/helpers from DOM rendering without changing behavior
- [ ] 16.3 Consolidate CSS rules and variables; remove obsolete/duplicate overrides
- [X] 16.4 Restore consistent thousands separators and whole-dollar formatting throughout UI (reported complete 2026-11-22)
- [ ] 16.5 Improve responsive table behavior on narrower screens (deferred by choice)
- [ ] 16.6 Add clear loading, empty-list, success, and error states
- [ ] 16.7 Add lightweight frontend tests for formatting and other pure functions

## Phase 17 — Inventory details and valuation

- [ ] 17.1 Add room/location field to the schema, model, API, and UI
- [ ] 17.2 Add purchase price separately from current estimated replacement value
- [ ] 17.3 Clarify existing Value field semantics and migrate data without losing values
- [ ] 17.4 Add condition field with a small, documented set of values
- [ ] 17.5 Display purchase date, serial number, details, and website URL in itemPhoto detail view
- [ ] 17.6 Allow filtering by room/location as well as category
- [ ] 17.7 Document how replacement estimates are entered and reviewed over time

## Phase 18 — Receipts and attachments

- [ ] 18.1 Define receipt attachment metadata and on-disk storage layout
- [ ] 18.2 Upload and associate a receipt with an itemPhoto
- [ ] 18.3 View/download attached receipts from itemPhoto details
- [ ] 18.4 Validate file type/size and prevent unsafe file paths
- [ ] 18.5 Include itemPhoto and receipt files in backup/restore procedure

## Phase 19 — Dashboard, search, and reporting

- [ ] 19.1 Build dashboard showing total replacement value, itemPhoto count, and category totals
- [ ] 19.2 Add room/location totals and identify high-value items
- [ ] 19.3 Add search across itemPhoto name, description, manufacturer, and model
- [ ] 19.4 Add an insurance-oriented PDF summary by category and location
- [ ] 19.5 Improve PDF money formatting, wrapping, and row background consistency
- [ ] 19.6 Consider a second, roomier PDF report for full itemPhoto details and attachments
- [ ] 19.7 Document that insurance policies may have special limits for selected categories

## Phase 20 — Backup, recovery, and release readiness

- [ ] 20.1 Document MySQL backup and restore commands
- [ ] 20.2 Back up itemPhoto photos and receipts alongside database data
- [ ] 20.3 Test restoring a backup into a separate database and file directory
- [ ] 20.4 Add a simple, repeatable local backup workflow
- [ ] 20.5 Review configuration, secrets handling, and ignored files before release
- [ ] 20.6 Write setup/run/test instructions for a fresh machine

## Standing decisions

- Value is the per-itemPhoto (unit) amount in whole dollars; TotalValue is Value × Quantity. Clarify how Value relates to replacement value before adding a separate purchase-price field.
- Name is generic; Description is the short specific name; Details holds longer text.
- Export sort: TotalValue descending, then name and description ascending.
- Unique itemPhoto key: (name, description). CSV import works and skips duplicates.
- MySQL is the source of truth; Goose manages schema migrations. The spreadsheet is not a live data source; there is no Sheets API integration.
- The CLI provides migrations and CSV import; the web app is the intended entry point for new items.
- Photos are stored on disk with a database reference to the file; uploads belong in the web UI.
- The web UI uses vanilla JavaScript for now; a framework is not required.
- Secrets must not be committed; plaintext `config.json` stays gitignored.
- Rob writes the code; the assistant provides guidance and review, with hints first and full implementations on request.

## Plan maintenance notes

- Existing checklist marks above were preserved from the uploaded plan, including historical items whose descriptions mention earlier intentions. Phase 7's original unchecked tasks should be reconciled individually rather than silently marked complete.
- The user confirmed the CSV import is complete and skips duplicates. The older 'one-shot only / no re-runnable import' assumptions were removed from Standing decisions.
- New phases are proposals, not claims that features or tests are currently missing; confirm against the latest repository before implementation.

## Recent UI and data notes (2026-11-22)

- Item Details is now a multiline textarea; it inherits the form font and follows view/edit disabled state.
- Purchase Date, Serial Number, and Details display in the itemPhoto dialog; 17.5 remains partial because website URL is outstanding.
- Reimporting the spreadsheet created an unintended row with Name `1`; correcting the spreadsheet and removing the extra row resolved the unique (name, description) conflict. Review import data quality on future imports.
- The Category dropdown appearance is deferred for later styling review.

## Next session — Phase 9 itemPhoto uploads

Priority for the next session: take more itemPhoto photographs and begin implementing itemPhoto uploads. Start with the file storage layout and database filename reference (9.2), then the upload route and UI (9.1), and displaying photos in itemPhoto details (9.3). Review existing repository code before implementation.
