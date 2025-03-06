package internal

import (
	"runtime"
	"sync"
	_ "unsafe"
)

type mapKeyAndValue[T comparable, V any] struct {
	key T
	val V
}

func ParallelMap[T comparable, V any](m map[T]V, f func(key T, val V)) {
	channel := make(chan mapKeyAndValue[T, V])

	procs := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func() {
			for keyAndVal := range channel {
				f(keyAndVal.key, keyAndVal.val)
			}
			wg.Done()
		}()
	}

	for key, val := range m {
		channel <- mapKeyAndValue[T, V]{key: key, val: val}
	}
	close(channel)
	wg.Wait()
}

//go:linkname Parallel github.com/disintegration/imaging.parallel
func Parallel(start, stop int, fn func(<-chan int))
