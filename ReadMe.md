# Package db

Provides batteries-included database connection (only MySQL driver)

# db.OpenEnv()

Package db does not have a global database connection. You must pass a database connection to most functions.

`db.OpenEnv` is used to open a `db.DB` connection using runtime options. Several runtime options are available. These may be provided with `.env` file, command line flags, or environment variable

For more information, please see `ztaylor.me/env`

```
const DB_USER = "DB_USER"
const DB_PASSWORD = "DB_PASSWORD"
const DB_HOST = "DB_HOST"
const DB_PORT = "DB_PORT"
const DB_TABLE = "DB_TABLE"
```

# cmd/db-patch

Provides executable that connects to a database, and executes a series of migrations(patches)

db-patch includes additional env variables

```
const PATCH_DIR = "PATCH_DIR"
```

Migrations are interpreted to be in sequential order, and begin with a patch id, which is 4 characters

Migration files are named as `DDDDx+.sql`, where `D` means `digit` and `x+` means any string