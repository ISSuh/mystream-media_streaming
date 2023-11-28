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
	"io/ioutil"

	"github.com/grafov/m3u8"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) MakeMediaM3u8() (string, error) {
	playlist, err := m3u8.NewMediaPlaylist(8, 10)
	if err != nil {
		return "", err
	}

	prefix := "/api/v1/segment/"
	testBasePath := "temp/7157824134036780522/20231128162134"
	files, err := ioutil.ReadDir(testBasePath)
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
