package pkg

import (
	"strconv"
	"sync"
)

func GeneratedId() func() string {
	var id int
	var mu sync.Mutex

	return func() string {
		mu.Lock()
		defer mu.Unlock()
		id++
		return strconv.Itoa(id)

	}

}
