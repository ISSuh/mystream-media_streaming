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
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/ISSuh/mystream-media_streaming/internal/event"
	"github.com/ISSuh/mystream-media_streaming/internal/hls"
	"github.com/ISSuh/mystream-media_streaming/internal/model"
	"github.com/ISSuh/mystream-media_streaming/internal/repository"
	"github.com/ISSuh/mystream-media_streaming/internal/repository/memory"
)

type StreamService struct {
	repository repository.StreamStatusRepository
	generator  *hls.PlaylistGenerator
}

func NewSteamManager(configure *configure.Configure) *StreamService {
	sm := &StreamService{
		repository: memory.NewStreamStatusMemoryRepository(),
		generator:  hls.NewPlaylistGenerator(),
	}
	return sm
}

func (sm *StreamService) OnActive(status *event.StreamStatus) {
	log.Info("[StreamService][OnActive] status : ", status)
	if !status.Active {
		log.Warn("[StreamService][OnActive] invalid stream status.")
		return
	}

	finded, exist := sm.repository.Find(status.StreamId)
	if exist && finded.Active {
		log.Warn("[StreamService][OnActive] stream already activated. ", finded)
		return
	}

	stream := sm.convertToStream(*status)
	sm.repository.Save(stream)
}

func (sm *StreamService) OnDeactive(status *event.StreamStatus) {
	log.Info("[StreamService][OnDeactive] status : ", status)
	if status.Active {
		log.Warn("[StreamService][OnActive] invalid stream status.")
		return
	}

	finded, exist := sm.repository.Find(status.StreamId)
	if exist && !finded.Active {
		log.Warn("[StreamService][OnActive] stream already deactivated. ", finded)
		return
	}

	sm.repository.Delete(status.StreamId)
}

func (sm *StreamService) FindMasterPlaylist(streamId int) (string, error) {
	stream, exist := sm.repository.Find(streamId)
	if !exist {
		return "", errors.New("[StreamService][FindMasterPlaylist] not found stream. " + strconv.Itoa(streamId))
	}

	if !stream.Active {
		return "", errors.New("[StreamService][FindMasterPlaylist] stream is not activated. " + strconv.Itoa(streamId))
	}

	return "", nil
}

func (sm *StreamService) FindMediaPlaylist(streamId int) (string, error) {
	stream, exist := sm.repository.Find(streamId)
	if !exist {
		return "", errors.New("[StreamService][FindMediaPlaylist] not found stream. " + strconv.Itoa(streamId))
	}

	if !stream.Active {
		return "", errors.New("[StreamService][FindMediaPlaylist] stream is not activated. " + strconv.Itoa(streamId))
	}

	return "", nil
}

func (sm *StreamService) convertToStream(status event.StreamStatus) model.Stream {
	return model.Stream{
		StreamId:       status.StreamId,
		Active:         status.Active,
		Url:            status.Url,
		ActiveAt:       status.ActiveAt,
		DeactiveAt:     status.DeactiveAt,
		MasterPlaylist: "",
		MediaPlaylist:  "",
	}
}
