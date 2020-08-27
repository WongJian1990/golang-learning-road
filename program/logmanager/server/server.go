package server

import (
	"context"
)

var pods []Pod

func AddPod(pod Pod) {
	pods = append(pods, pod)
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, pod := range pods {
		go func(pod Pod) {
			pod.Run()
		}(pod)
	}
	select {
	case <-ctx.Done():
		return
	}

}
