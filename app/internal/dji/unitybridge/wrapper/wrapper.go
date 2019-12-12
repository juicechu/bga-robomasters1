// +build !windows !amd64

package wrapper

type wrapper struct{}

func Instance() *wrapper {
	panic("Only windows amd64 is supported.")
}

