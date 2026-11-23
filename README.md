# inventory
Home inventory app to track items initially for insurance purposes.

## Photo coverage filter

The inventory page includes a **Photos** filter: **All Items**, **Without Photos**, or **With Photos**. It can be combined with category and search filters. The `GET /api/items` response includes `photo_count` for each item; counts update after a photo is uploaded or deleted.

Photo counts are computed from the existing `item_photos` table, so no additional migration is needed.

## Tests

Run all Go tests with `go test ./...`. The MySQL repository tests use the test database configured in `config/config.json` and apply/reset Goose migrations; use a dedicated disposable test database, not your inventory database.
