package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type worker struct {
	Conn  *nats.Conn
	CS    *kubernetes.Clientset
	Topic string
}

func (w worker) Do(ctx context.Context, pod v1.Pod) {
	PodLogsConnection := w.CS.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Follow:    true,
		TailLines: &[]int64{int64(10)}[0],
	})

	LogStream, _ := PodLogsConnection.Stream(context.Background())
	defer func(LogStream io.ReadCloser) {
		err := LogStream.Close()
		if err != nil {
			log.Println(fmt.Errorf("failed to close pod=%s stream: %w", pod.Name, err))
		}
	}(LogStream)

	reader := bufio.NewScanner(LogStream)

	var line string

	for {
		select {
		case <-ctx.Done():
			break
		default:
			for reader.Scan() {
				line = reader.Text()

				topic := fmt.Sprintf("%s.logs.%s", w.Topic, strings.ToLower(encodeLog(line)))

				if err := w.Conn.Publish(topic, []byte(line)); err != nil {
					log.Println(fmt.Errorf("failed to publish pod=%s logs: %w", pod.Name, err))
				}
			}
		}
	}
}
