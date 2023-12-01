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

package file

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type SegmentsReaderOnFile struct {
}

func NewSegmentsReaderOnFile() *SegmentsReaderOnFile {
	return &SegmentsReaderOnFile{}
}

func (s *SegmentsReaderOnFile) SegmentsList(uri string, offset int, limit int) ([]string, error) {
	result := make([]string, 0)

	files, err := os.ReadDir(uri)
	if err != nil {
		log.Error("[SegmentsReaderOnFile][SegmentsList] can not open dir. ", err.Error())
		return nil, err
	}

	filesLen := len(files)
	if filesLen > limit {
		files = files[filesLen-limit:]
	}

	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, err
}

func (s *SegmentsReaderOnFile) ReadSegment(uri, segmentName string) ([]byte, error) {
	path := uri + "/" + segmentName
	segment, err := os.ReadFile(path)
	if err != nil {
		log.Info("[SegmentsReaderOnFile][SegmentsList] can not open segment. ", err.Error())
		return nil, err
	}
	return segment, nil
}
