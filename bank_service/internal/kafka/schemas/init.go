package schemas

import (
	"github.com/linkedin/goavro/v2"
)

var (
	bankStatementSchema *goavro.Codec
)

func init() {
	var err error

	bankStatementSchema, err = goavro.NewCodec(bankStatementSchemaString)
	if err != nil {
		panic(err)
	}
}
