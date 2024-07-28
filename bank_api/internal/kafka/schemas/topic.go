package schemas

import _ "embed"

var (
	//go:embed definitions/bank_v1_bank_statement.json
	bankStatementSchemaString string
)

const (
	TopicBankStatement = "bank_v1_bank_statement"
)
