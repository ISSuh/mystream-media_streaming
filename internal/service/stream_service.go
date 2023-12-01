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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_streaming/internal/configure"
	"github.com/ISSuh/mystream-media_streaming/internal/event"
	"github.com/ISSuh/mystream-media_streaming/internal/hls"
	"github.com/ISSuh/mystream-media_streaming/internal/hls/file"
	"github.com/ISSuh/mystream-media_streaming/internal/model"
	"github.com/ISSuh/mystream-media_streaming/internal/repository"
	"github.com/ISSuh/mystream-media_streaming/internal/repository/memory"
)

type StreamService struct {
	configure  *configure.Configure
	repository repository.StreamStatusRepository
	generator  *hls.PlaylistGenerator
	reader     hls.SegmentsReader
}

func NewSteamManager(configure *configure.Configure) *StreamService {
	reader := file.NewSegmentsReaderOnFile()

	sm := &StreamService{
		configure:  configure,
		repository: memory.NewStreamStatusMemoryRepository(),
		generator:  hls.NewPlaylistGenerator(reader),
		reader:     reader,
	}

	// sm.insertTestData()

	return sm
}

// func (sm *StreamService) insertTestData() {
// 	stream := model.Stream{
// 		StreamId:          1,
// 		Active:            true,
// 		Uri:               "/1/20231130215747",
// 		MasterPlaylistUri: "/api/v1/playlist/master/1/20231130215747",
// 		MediaPlaylistUri:  "/api/v1/playlist/media/1/20231130215747",
// 	}

// 	sm.repository.Save(stream)
// }

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

	stream := sm.newStreamFromStatus(*status)

	playlist, err := sm.generator.MakeMasterPlaylist(&stream)
	if err != nil {
		log.Error("[StreamService][OnActive] cat not generate mast playlist", err.Error())
		return
	}

	stream.MasterPlaylist.Playlist = playlist

	log.Info("[StreamService][OnActive] stream : ", stream)
	sm.repository.Save(stream)
}

func (sm *StreamService) OnDeactive(status *event.StreamStatus) {
	log.Info("[StreamService][OnDeactive] status : ", status)
	if status.Active {
		log.Warn("[StreamService][OnDeactive] invalid stream status.")
		return
	}

	finded, exist := sm.repository.Find(status.StreamId)
	if exist && !finded.Active {
		log.Warn("[StreamService][OnDeactive] stream already deactivated. ", finded)
		return
	}

	sm.repository.Delete(status.StreamId)
}

func (sm *StreamService) FindMasterPlaylist(streamId int, streamPath string) (string, error) {
	stream, exist := sm.repository.Find(streamId)
	if !exist {
		return "", errors.New("[StreamService][FindMasterPlaylist] not found stream. " + strconv.Itoa(streamId))
	}

	if !stream.Active {
		return "", errors.New("[StreamService][FindMasterPlaylist] stream is not activated. " + strconv.Itoa(streamId))
	}

	path := "/" + strconv.Itoa(streamId) + "/" + streamPath
	if stream.Uri != path {
		return "", errors.New("[StreamService][FindMasterPlaylist] invalid stream path " + strconv.Itoa(streamId))
	}

	return stream.MasterPlaylist.Playlist, nil
}

func (sm *StreamService) FindMediaPlaylist(streamId int, streamPath string) (string, error) {
	stream, exist := sm.repository.Find(streamId)
	if !exist {
		return "", errors.New("[StreamService][FindMediaPlaylist] not found stream. " + strconv.Itoa(streamId))
	}

	if !stream.Active {
		return "", errors.New("[StreamService][FindMediaPlaylist] stream is not activated. " + strconv.Itoa(streamId))
	}

	path := "/" + strconv.Itoa(streamId) + "/" + streamPath
	if stream.Uri != path {
		return "", errors.New("[StreamService][FindMasterPlaylist] invalid stream path " + strconv.Itoa(streamId))
	}

	// returns saved media playlist if it has been updated less than 2 seconds
	if time.Since(stream.MediaPlaylist.UpdateTime) < 2*time.Second {
		return stream.MediaPlaylist.Playlist, nil
	}

	segmentDir := sm.configure.Segment.BasePath + "/" + path
	playlist, err := sm.generator.MakeMediaPlaylist(&stream, segmentDir)
	if err != nil {
		return "", err
	}

	stream.MediaPlaylist.Playlist = playlist
	stream.MediaPlaylist.UpdateTime = time.Now()
	stream.MediaPlaylist.Sequence++
	sm.repository.Save(stream)

	log.Info("[StreamService][FindMediaPlaylist] stream : ", stream)
	return playlist, nil
}

func (sm *StreamService) FindSegment(streamId int, streamPath string, segmentName string) ([]byte, error) {
	path := sm.configure.Segment.BasePath + "/" + strconv.Itoa(streamId) + "/" + streamPath
	return sm.reader.ReadSegment(path, segmentName)
}

func (sm *StreamService) newStreamFromStatus(status event.StreamStatus) model.Stream {
	baseUri := "/api/v1/playlist/"

	masterPlaylistUri := baseUri + "master" + status.Uri
	masterPlayList := model.MasterPlaylist{
		Playlist: "",
		Uri:      masterPlaylistUri,
	}

	mediaPlaylistUri := baseUri + "media" + status.Uri
	mediaPlayList := model.MediaPlaylist{
		Playlist:   "",
		Uri:        mediaPlaylistUri,
		Sequence:   0,
		UpdateTime: time.Now(),
	}

	return model.Stream{
		StreamId:       status.StreamId,
		Active:         status.Active,
		Uri:            status.Uri,
		ActiveAt:       status.ActiveAt,
		DeactiveAt:     status.DeactiveAt,
		MasterPlaylist: masterPlayList,
		MediaPlaylist:  mediaPlayList,
	}
}
