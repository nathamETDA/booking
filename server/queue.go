package main

import (
	"time"
)

func wait2door() {
	for {
		queueMutex.Lock()
		now := time.Now()
		more := maxExitDoor - len(exitDoor)
		if more > len(waitQueue) {
			more = len(waitQueue)
		}
		for i := 0; i < more; i++ {
			exitDoor = append(exitDoor, clientState{waitQueue[i].id, now})
		}
		copy(waitQueue, waitQueue[more:])
		waitQueue = waitQueue[:len(waitQueue)-more]
		queueMutex.Unlock()

		time.Sleep(refreshInterfal)
	}
}

func purgeWaitQueue() {
	for {
		queueMutex.Lock()
		now := time.Now()
		for i := 0; i < len(waitQueue); i++ {
			if waitQueue[i].lastUpdate.Add(waitQueueTimeout).Before(now) {
				copy(waitQueue[i:], waitQueue[i+1:])
				waitQueue = waitQueue[:len(waitQueue)-1]
			}
		}
		queueMutex.Unlock()
		time.Sleep(refreshInterfal)
	}
}

func purgeExitDoor() {
	for {
		queueMutex.Lock()
		now := time.Now()
		for i := 0; i < len(exitDoor); i++ {
			if exitDoor[i].lastUpdate.Add(exitDooorTimeout).Before(now) {
				copy(exitDoor[i:], exitDoor[i+1:])
				exitDoor = exitDoor[:len(exitDoor)-1]
			}
		}
		queueMutex.Unlock()
		time.Sleep(refreshInterfal)
	}
}

func moveOut(id string) bool {
	for i := 0; i < len(exitDoor); i++ {
		if id == exitDoor[i].id {
			queueMutex.Lock()
			//recheck again in case of race condision
			if id == exitDoor[i].id {
				copy(exitDoor[i:], exitDoor[i+1:])
				exitDoor = exitDoor[:len(exitDoor)-1]
			}
			queueMutex.Unlock()
			return true
		}
	}
	return false
}
