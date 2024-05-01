package sql

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/artie-labs/transfer/lib/config/constants"
)

type Dialect interface {
	NeedsEscaping(identifier string) bool // TODO: Remove this when we escape everything
	QuoteIdentifier(identifier string) string
}

type DefaultDialect struct{}

func (DefaultDialect) NeedsEscaping(_ string) bool { return true }

func (DefaultDialect) QuoteIdentifier(identifier string) string {
	return fmt.Sprintf(`"%s"`, identifier)
}

type BigQueryDialect struct{}

func (BigQueryDialect) NeedsEscaping(_ string) bool { return true }

func (BigQueryDialect) QuoteIdentifier(identifier string) string {
	// BigQuery needs backticks to quote.
	return fmt.Sprintf("`%s`", identifier)
}

type RedshiftDialect struct{}

func (RedshiftDialect) NeedsEscaping(_ string) bool { return true }

func (rd RedshiftDialect) QuoteIdentifier(identifier string) string {
	// Preserve the existing behavior of Redshift identifiers being lowercased due to not being quoted.
	return fmt.Sprintf(`"%s"`, strings.ToLower(identifier))
}

type SnowflakeDialect struct {
	UppercaseEscNames bool
}

// symbolsToEscape are additional keywords that we need to escape
var symbolsToEscape = []string{":"}

func (sd SnowflakeDialect) NeedsEscaping(name string) bool {
	if sd.UppercaseEscNames {
		// If uppercaseEscNames is true then we will escape all identifiers that do not start with the Artie priefix.
		// Since they will be uppercased afer they are escaped then they will result in the same value as if we
		// we were to use them in a query without any escaping at all.
		return true
	} else {
		if slices.Contains(constants.ReservedKeywords, name) {
			return true
		}
		// If it does not contain any reserved words, does it contain any symbols that need to be escaped?
		for _, symbol := range symbolsToEscape {
			if strings.Contains(name, symbol) {
				return true
			}
		}
		// If it still doesn't need to be escaped, we should check if it's a number.
		if _, err := strconv.Atoi(name); err == nil {
			return true
		}
		return false
	}
}

func (sd SnowflakeDialect) QuoteIdentifier(identifier string) string {
	if sd.UppercaseEscNames {
		identifier = strings.ToUpper(identifier)
	} else {
		slog.Warn("Escaped Snowflake identifier is not being uppercased",
			slog.String("name", identifier),
			slog.Bool("uppercaseEscapedNames", sd.UppercaseEscNames),
		)
	}

	return fmt.Sprintf(`"%s"`, identifier)
}