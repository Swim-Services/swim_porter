package internal

import (
	"runtime"
	_ "unsafe"
)

type mapKeyAndValue[T comparable, V any] struct {
	key T
	val V
}

func ParallelMap[T comparable, V any](m map[T]V, f func(key T, val V)) {
	channel := make(chan mapKeyAndValue[T, V])

	procs := runtime.GOMAXPROCS(0)
	for i := 0; i < procs; i++ {
		go func() {
			for keyAndVal := range channel {
				f(keyAndVal.key, keyAndVal.val)
			}
		}()
	}

	for key, val := range m {
		channel <- mapKeyAndValue[T, V]{key: key, val: val}
	}
}

//go:linkname Parallel github.com/disintegration/imaging.parallel
func Parallel(start, stop int, fn func(<-chan int))
