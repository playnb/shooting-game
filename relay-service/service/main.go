package main

import (
	"github.com/playnb/util/log"
)

func main() {
	log.InitPanic("../tmp")
	log.Init(log.DefaultLogger("../tmp", "run"))
	defer log.Flush()

	log.Trace("relay service")
}
