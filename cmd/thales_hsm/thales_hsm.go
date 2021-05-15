package main

import (
	_ "github.com/hmmftg/libiso/hsm"
	"github.com/hmmftg/libiso/hsm/console"
	"sync"
)

func main() {

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	thalesConsole := console.New()
	go thalesConsole.Show(waitGroup)

	waitGroup.Wait()

}
