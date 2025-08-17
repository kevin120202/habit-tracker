package migrations

import "embed"

//go:embed *.sql
var FS embed.FS // FS is a variable that holds the embedded file system containing all `.sql` files
