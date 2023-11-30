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
	"fmt"
	"os"

	"github.com/ISSuh/mystream-media_streaming/internal/model"
	"github.com/grafov/m3u8"
)

type PlaylistGenerator struct {
}

func NewPlaylistGenerator() *PlaylistGenerator {
	return &PlaylistGenerator{}
}

func (g *PlaylistGenerator) MakeMediaM3u8(stream *model.Stream) (string, error) {

	playlist, err := m3u8.NewMediaPlaylist(8, 10)
	if err != nil {
		return "", err
	}

	prefix := "/api/v1/segment/"
	testBasePath := "temp/7157824134036780522/20231128162134"
	files, err := os.ReadDir(testBasePath)
	if err != nil {
		return "", err
	}

	var duration float64 = 2.0
	for _, file := range files {
		fmt.Println("[TEST] file : ", file.Name())

		path := prefix + testBasePath + "/" + file.Name()
		if err := playlist.Append(path, duration, "title!!"); err != nil {
			fmt.Println("[TEST] playlist append error. ", err)
		}
	}

	return playlist.String(), nil
}

func (g *PlaylistGenerator) makeMasterPlaylist(stream *model.Stream) (string, error) {
	master := m3u8.NewMasterPlaylist()
	master.SetVersion(3)

	master.Append(
		"chunklist1.m3u8", nil,
		m3u8.VariantParams{
			ProgramId: 123, Resolution: "1280x720", FrameRate: 30.000})

	return "", nil
}

func (g *PlaylistGenerator) makeMediaPlaylist(stream *model.Stream) (string, error) {
	// media := m3u8.NewMasterPlaylist()
	return "", nil
}
