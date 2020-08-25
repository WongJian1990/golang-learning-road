package monitor

import "golang.org/x/net/context"

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		return
	}
}
