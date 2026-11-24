# inventory
Home inventory app to track items initially for insurance purposes.

## Photo coverage filter

The inventory page includes a **Photos** filter: **All Items**, **Without Photos**, or **With Photos**. It can be combined with category and search filters. The `GET /api/items` response includes `photo_count` for each item; counts update after a photo is uploaded or deleted.

Photo counts are computed from the existing `item_photos` table, so no additional migration is needed.

## Tests

Run all Go tests with `go test ./...`. The MySQL repository tests use the test database configured in `config/config.json` and apply/reset Goose migrations; use a dedicated disposable test database, not your inventory database.

## Room/location, reference links, and bulk editing

Apply Goose migration `00005_item_locations_and_links.sql` **before starting the updated server**.
Existing items receive empty values for the new metadata fields; existing inventory values are preserved.

- Each item has an editable **Room / Location** input with suggestions taken from rooms already used in the inventory. Free-form values such as `Garage > Cabinet` are supported; consistent spelling matters.
- **Room** in the toolbar filters items by their assigned location.
- **Website URL** and **Manual URL** are separate optional http(s) links. When populated, the item detail view offers an **Open link** action.
- **Condition**, **Purchase Source**, and **Purchase Price** are also available in the item dialog.
- Use the checkboxes beside items (or the header checkbox to select all **currently visible** items), then **Bulk Edit**. The confirmation dialog shows the number of selected items. Only the chosen field is changed. Supported bulk fields: room, category, condition, purchase source, website URL, and manual URL. Empty text clears the chosen field.
- Bulk edits validate all selected IDs before writing. If a database error occurs during writing, the API reports how many succeeded; **bulk edits are not currently transactional**.
- CSV export includes the new metadata columns. CSV import accepts older spreadsheets without those columns, treating omitted values as blank/zero. The existing importer **skips duplicate name + short-description entries** rather than updating them, so re-importing a spreadsheet does not update existing records.

Run `go test ./...` with a disposable MySQL test database before deployment, and verify a few add/edit, bulk edit, filtering, and CSV import/export flows in the browser.

## Production layout and configuration

The application supports a persistent `--config` flag on the root command:

```sh
inventory --config /etc/inventory/config.json server --port 8080
```

The configuration file supports optional `photos_dir` and `web_dir` properties. When omitted, they default to `photos` and `./web/public`, respectively, to preserve development behavior. Paths are resolved relative to the process working directory unless absolute paths are supplied.

Recommended production paths:

- Executable: `/usr/local/inventory/bin/inventory`
- Static web files: `/usr/local/inventory/web/public`
- Config (including database credentials): `/etc/inventory/config.json`
- Persistent photos: `/var/lib/inventory/photos`
- Logs: `/var/log/inventory`

Example additional properties in the production config JSON:

```json
{
  "photos_dir": "/var/lib/inventory/photos",
  "web_dir": "/usr/local/inventory/web/public"
}
```

Keep the existing `database` and `test_database` properties when adding these fields. Do not commit production credentials or photographs to Git.

### Safe deployment policy

Run `go test ./...` against a **dedicated disposable MySQL test database**, never against production. The repository tests can reset Goose migrations and modify test tables. A deployment must abort before installation if tests or build fail. After tests, stage a new binary and web assets, retain a previous release for rollback, restart the service, and check HTTP health. Coordinate production path changes with the backup script so it continues to include the actual configuration and photographs. **Do not migrate the running service until these steps have been verified.**
