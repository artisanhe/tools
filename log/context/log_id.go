package context

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

var (
	logIDs = map[int64]string{}
	locker = sync.RWMutex{}
)

//func getGoID() int64 {
//return runtime.GoID()
//}

func getGoID() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine"))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return int64(n)
}

func SetLogID(id string) {
	locker.Lock()
	defer locker.Unlock()

	logIDs[getGoID()] = id
}

func GetLogID() string {
	locker.RLock()
	defer locker.RUnlock()

	goID := getGoID()

	if logID, ok := logIDs[goID]; ok {
		return logID
	}

	return ""
}

func Close() {
	locker.Lock()
	defer locker.Unlock()

	goID := getGoID()

	if _, ok := logIDs[goID]; ok {
		delete(logIDs, goID)
	}
}
