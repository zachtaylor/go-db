# Package db

`import "ztaylor.me/db"`

Package `db` provides database connection helpers based on `database/sql`

## Package `env`

`import "ztaylor.me/db/env"`

Package `env` provides bindings to `ztaylor.me/env` to assist in the creation of Data Source Name

Environment variables used

```
USER
PASSWORD
HOST
PORT
NAME
```

## Package `mysql`

`import "ztaylor.me/db/mysql"`

Package `mysql` loads mysql driver using `"github.com/go-sql-driver/mysql"`

# Binary `db-patch`

`go get ztaylor.me/db/cmd/db-patch`

Connect to a database using MySQL, and execute a series of patches

Patches are contained separately in files, known as patch files. These files
- contain SQL statements, which are executed as transactions (each patch will succeed or fail as a whole)
- begin with 4 numbers, identifying the patch number in sequence
- end with ".sql"
