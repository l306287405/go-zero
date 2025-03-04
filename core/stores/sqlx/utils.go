package sqlx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/l306287405/go-zero/core/logx"
	"github.com/l306287405/go-zero/core/mapping"
)

func desensitize(datasource string) string {
	// remove account
	pos := strings.LastIndex(datasource, "@")
	if 0 <= pos && pos+1 < len(datasource) {
		datasource = datasource[pos+1:]
	}

	return datasource
}

func escape(input string) string {
	var b strings.Builder

	for _, ch := range input {
		switch ch {
		case '\x00':
			b.WriteString(`\x00`)
		case '\r':
			b.WriteString(`\r`)
		case '\n':
			b.WriteString(`\n`)
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '"':
			b.WriteString(`\"`)
		case '\x1a':
			b.WriteString(`\x1a`)
		default:
			b.WriteRune(ch)
		}
	}

	return b.String()
}

func format(query string, args ...interface{}) (string, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return query, nil
	}

	var b strings.Builder
	var argIndex int
	bytes := len(query)

	for i := 0; i < bytes; i++ {
		ch := query[i]
		switch ch {
		case '?':
			if argIndex >= numArgs {
				return "", fmt.Errorf("error: %d ? in sql, but less arguments provided", argIndex)
			}

			writeValue(&b, args[argIndex])
			argIndex++
		case '$':
			var j int
			for j = i + 1; j < bytes; j++ {
				char := query[j]
				if char < '0' || '9' < char {
					break
				}
			}
			if j > i+1 {
				index, err := strconv.Atoi(query[i+1 : j])
				if err != nil {
					return "", err
				}

				// index starts from 1 for pg
				if index > argIndex {
					argIndex = index
				}
				index--
				if index < 0 || numArgs <= index {
					return "", fmt.Errorf("error: wrong index %d in sql", index)
				}

				writeValue(&b, args[index])
				i = j - 1
			}
		default:
			b.WriteByte(ch)
		}
	}

	if argIndex < numArgs {
		return "", fmt.Errorf("error: %d arguments provided, not matching sql", argIndex)
	}

	return b.String(), nil
}

func logInstanceError(datasource string, err error) {
	datasource = desensitize(datasource)
	logx.Errorf("Error on getting sql instance of %s: %v", datasource, err)
}

func logSqlError(stmt string, err error) {
	if err != nil && err != ErrNotFound {
		logx.Errorf("stmt: %s, error: %s", stmt, err.Error())
	}
}

func writeValue(buf *strings.Builder, arg interface{}) {
	switch v := arg.(type) {
	case bool:
		if v {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	case string:
		buf.WriteByte('\'')
		buf.WriteString(escape(v))
		buf.WriteByte('\'')
	default:
		buf.WriteString(mapping.Repr(v))
	}
}
