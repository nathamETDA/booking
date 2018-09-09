package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	speed            = 1
	maxWaitQueue     = 100000
	maxExitDoor      = 100
	waitQueueTimeout = time.Minute * 4
	exitDooorTimeout = time.Minute * 5
	refreshInterfal  = time.Second * 5
)

var open = false
var (
	waitQueue  = make([]clientState, 0, maxWaitQueue)
	exitDoor   = make([]clientState, 0, maxExitDoor)
	queueMutex = &sync.Mutex{}
)

var (
	waitSecret  = "unicornsAreAwesome"
	loginSecret = "unicornsAreAwesome2"
)

type clientState struct {
	id         string
	lastUpdate time.Time
}

var id = 1
var idMutex = &sync.Mutex{}

func newClientState() clientState {
	idMutex.Lock()
	r := fmt.Sprint(id)
	id++
	idMutex.Unlock()
	return clientState{r, time.Now()}
}
