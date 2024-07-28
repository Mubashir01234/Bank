package schemas

import (
	_ "embed"

	"github.com/linkedin/goavro/v2"
)

type BankStatement struct {
	StatementBuf []byte `json:"data"`
}

func (BankStatement) Name() string {
	return TopicBankStatement
}

func (BankStatement) Schema() *goavro.Codec {
	return bankStatementSchema
}
