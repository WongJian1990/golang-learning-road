package monitor

import "golang.org/x/net/context"

var pods []*KafkaProducer

func AddPod(pod *KafkaProducer) {
	pods = append(pods, pod)
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for _, pod := range pods {
			pod.Run()
		}
	}()
	select {
	case <-ctx.Done():
		return
	}
}
