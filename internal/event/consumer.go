/*
MIT License

Copyright (c) 2023 ISSuh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package event

import (
	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	configure *configure.KafkaConfigure
	consumer  *kafka.Consumer
}

func NewConsumer(configure *configure.KafkaConfigure) *Consumer {
	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": configure.BootstrapServer,
			"group.id":          configure.GroupId,
			"auto.offset.reset": "earliest",
		},
	)
	if err != nil {
		panic("can not create kafka consumer")
		return nil
	}

	return &Consumer{
		configure: configure,
		consumer:  c,
	}
}

func (c *Consumer) Init() {
	// streamActiveTopic := "stream-active"
}
