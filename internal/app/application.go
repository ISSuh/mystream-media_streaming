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

package app

import (
	"github.com/ISSuh/mystream-media_streaming/internal/api/router"
	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/ISSuh/mystream-media_streaming/internal/event"
	"github.com/ISSuh/mystream-media_streaming/internal/service"
	"github.com/gin-gonic/gin"
)

type MainApplication struct {
	configure *configure.Configure
	engine    *gin.Engine
	consumer  *event.Consumers
}

func NewApplication(configure *configure.Configure) *MainApplication {
	streamService := service.NewSteamManager(configure)
	consumer := event.NewConsumers(&configure.Kafka, streamService)

	return &MainApplication{
		configure: configure,
		engine:    router.Setup(streamService),
		consumer:  consumer,
	}
}

func (a *MainApplication) Run() error {
	if err := a.engine.Run(":" + a.configure.Server.Port); err != nil {
		return err
	}
	return nil
}
