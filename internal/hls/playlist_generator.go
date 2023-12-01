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

package hls

import (
	"github.com/ISSuh/mystream-media_streaming/internal/model"
	"github.com/grafov/m3u8"
)

type PlaylistGenerator struct {
	reader SegmentsReader
}

func NewPlaylistGenerator(reader SegmentsReader) *PlaylistGenerator {
	return &PlaylistGenerator{
		reader: reader,
	}
}

func (g *PlaylistGenerator) MakeMasterPlaylist(stream *model.Stream) (string, error) {
	master := m3u8.NewMasterPlaylist()
	master.SetVersion(3)

	master.Append(
		stream.MediaPlaylist.Uri,
		nil,
		m3u8.VariantParams{
			ProgramId: 50, Resolution: "1280x720", FrameRate: 30.000, Video: "720p30"})

	return master.String(), nil
}

func (g *PlaylistGenerator) MakeMediaPlaylist(stream *model.Stream, streamPath string) (string, error) {
	totalSegment := 10

	list, err := g.reader.SegmentsList(streamPath, 0, totalSegment)
	if err != nil {
		return "", err
	}

	itemLen := len(list)
	playlist, err := m3u8.NewMediaPlaylist(uint(itemLen), uint(totalSegment))
	if err != nil {
		return "", err
	}

	playlist.SeqNo = stream.MediaPlaylist.Sequence

	var duration float64 = 2
	title := "live"
	for i := 0; i < itemLen; i++ {
		path := stream.MediaPlaylist.Uri + "/" + list[i]
		playlist.Append(path, duration, title)
	}

	return playlist.String(), nil
}
