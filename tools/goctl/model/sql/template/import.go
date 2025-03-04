package template

var (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/l306287405/go-zero/core/stores/builder"
	"github.com/l306287405/go-zero/core/stores/cache"
	"github.com/l306287405/go-zero/core/stores/sqlc"
	"github.com/l306287405/go-zero/core/stores/sqlx"
	"github.com/l306287405/go-zero/core/stringx"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/l306287405/go-zero/core/stores/builder"
	"github.com/l306287405/go-zero/core/stores/sqlc"
	"github.com/l306287405/go-zero/core/stores/sqlx"
	"github.com/l306287405/go-zero/core/stringx"
)
`
)
