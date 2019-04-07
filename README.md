# Package db

`import "ztaylor.me/db"`

Package db provides database connection helpers for MySQL

## `db.OpenEnv()`

`db.OpenEnv` creates `sql.DB` using runtime options. For more information, please see `ztaylor.me/env`.

Environment variables
```
DB_USER
DB_PASSWORD
DB_HOST
DB_PORT
DB_NAME
```

## `cmd/db-patch`

Executable: connect to database, and execute a series of patches; supports options defined by `db.OpenEnv`

Additional environment variables
```
PATCH_DIR
```

Patches are executed in sequential order. They must be named as `"dddd.*\.sql"`. In other words, patch file names begin with 4 numbers, and end with `".sql"`

Patch files are executed as SQL transactions; each patch file will succeed or fail as a whole
