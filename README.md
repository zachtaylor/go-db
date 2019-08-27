# Package db

`import "ztaylor.me/db"`

Package `db` provides database connection helpers

## Package `dbe`

`import "ztaylor.me/db/env"`

Package `dbe` provides bindings to `ztaylor.me/env`

Environment variables

```
DB_USER
DB_PASSWORD
DB_HOST
DB_PORT
DB_NAME
```

## Package `mysql`

`import "ztaylor.me/db/mysql"`

Package `mysql` loads mysql driver using `"github.com/go-sql-driver/mysql"`

# Binary `db-patch`

`go get ztaylor.me/db/cmd/db-patch`

Connect to a database using MySQL, and execute a series of patches

 - Supports runtime options using `"ztaylor.me/db/env"`
 - New runtime option `"PATCH_DIR"` path of dir containing patch files named as `regex("dddd.*\.sql")`. In other words, patch file names begin with 4 numbers, and end with `".sql"`
 - Create table `patches` to record db version
 - Patch files are executed as SQL transactions; each patch will succeed or fail as a whole
