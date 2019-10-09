// schemes: https
// host: p1payapi.tst.protocol.one
// version: 1.0.0
// swagger:meta
package main

import (
	"github.com/paysuper/paysuper-management-api/cmd/http"
	"github.com/paysuper/paysuper-management-api/cmd/micro"
	"github.com/paysuper/paysuper-management-api/cmd/root"
)

func main() {
	root.Execute(http.Cmd, micro.Cmd)
}
