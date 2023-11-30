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

package service

import (
	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/ISSuh/mystream-media_streaming/internal/event"
	"github.com/ISSuh/mystream-media_streaming/internal/hls"
	"github.com/ISSuh/mystream-media_streaming/internal/model"
	"github.com/ISSuh/mystream-media_streaming/internal/repository"
	"github.com/ISSuh/mystream-media_streaming/internal/repository/memory"
)

type StreamService struct {
	repository repository.StreamStatusRepository
	generator  *hls.Generator
}

func NewSteamManager(configure *configure.Configure) *StreamService {
	sm := &StreamService{
		repository: memory.NewStreamStatusMemoryRepository(),
		generator:  hls.NewGenerator(),
	}
	return sm
}

func (sm *StreamService) OnActive(status *event.StreamStatus) {

}

func (sm *StreamService) OnDeactive(status *event.StreamStatus) {

}

func (sm *StreamService) FindStream(streamId int) (model.Stream, bool) {
	return sm.repository.Find(streamId)
}

func (sm *StreamService) MakeMasterPlaylist(streamId int) (string, error) {
	return "", nil
}

func (sm *StreamService) MakeMediaPlaylist(streamId int) (string, error) {
	return "", nil
}
