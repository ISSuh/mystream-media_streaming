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
	"context"
	"encoding/json"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/segmentio/kafka-go"
)

type Consumers struct {
	consumerFactory *ConsumerFactory

	streamActiveConsumer   *kafka.Reader
	streamDeactiveConsumer *kafka.Reader
	litener                StreamListener

	activeChan   chan kafka.Message
	deactiveChan chan kafka.Message

	running bool
	wg      sync.WaitGroup
}

func NewConsumers(configure *configure.KafkaConfigure, litener StreamListener) *Consumers {
	consumerFactory := NewConsumerFactory(configure)

	return &Consumers{
		consumerFactory:        consumerFactory,
		streamActiveConsumer:   consumerFactory.streamActiveConsumer(),
		streamDeactiveConsumer: consumerFactory.streamDeactiveConsumer(),
		litener:                litener,
		activeChan:             make(chan kafka.Message),
		deactiveChan:           make(chan kafka.Message),
		running:                false,
		wg:                     sync.WaitGroup{},
	}
}

func (c *Consumers) RunBackground() {
	c.running = true

	c.wg.Add(1)
	go func() {
		log.Info("[Consumers][RunBackground] run StreamActiveConsumer")
		defer c.wg.Done()

		for c.running {
			m, err := c.streamActiveConsumer.ReadMessage(context.Background())
			if err != nil {
				continue
			}

			c.activeChan <- m
		}
	}()

	c.wg.Add(1)
	go func() {
		log.Info("[Consumers][RunBackground] run StreamDeactiveConsumer")
		defer c.wg.Done()

		for c.running {
			m, err := c.streamDeactiveConsumer.ReadMessage(context.Background())
			if err != nil {
				continue
			}

			c.activeChan <- m
		}
	}()

	go c.work()
}

func (c *Consumers) work() {
	for c.running {
		select {
		case message := <-c.activeChan:
			log.Info("[Consumers][work] active event : ", string(message.Value))
			s, err := c.parseStreamStatus(message.Value)
			if err != nil {
				continue
			}

			c.litener.OnActive(s)
		case message := <-c.deactiveChan:
			log.Info("[Consumers][work] deactive event : ", string(message.Value))
			s, err := c.parseStreamStatus(message.Value)
			if err != nil {
				continue
			}

			c.litener.OnDeactive(s)
		}
	}
}

func (c *Consumers) Close() {
	c.running = false

	c.streamActiveConsumer.Close()
	c.streamDeactiveConsumer.Close()

	c.wg.Wait()
}

func (c *Consumers) parseStreamStatus(event []byte) (*StreamStatus, error) {
	input, err := strconv.Unquote(string(event))
	if err != nil {
		log.Error("[Consumers][parseStreamStatus] invalid event format. ", err.Error())
		return nil, err
	}

	status := &StreamStatus{}
	err = json.Unmarshal([]byte(input), status)
	if err != nil {
		log.Error("[Consumers][parseStreamStatus] parse event error. ", err.Error())
		return nil, err
	}
	return status, nil
}
