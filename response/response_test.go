package response

import (
	"fmt"
	"testing"
	"time"
)

func Test_Stream(t *testing.T) {
	channel := make(chan *Response, 100)

	go func() {
		timer := time.NewTimer(time.Second * 10)
		i := 0
		for {
			select {
			case <-timer.C:
				close(channel)
				fmt.Println("开始退出生产", "i=", i)
				return
			default:
				i += 1
				data := NewResponse()
				data.Created = i
				channel <- data
			}
		}
	}()

	i := 0
	for v := range channel {
		i++
		fmt.Println("i=", i, "time=", time.Now().String(), "v", v.Created)
		time.Sleep(time.Second * 1)
	}
}
