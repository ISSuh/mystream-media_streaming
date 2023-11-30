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

package memory

import (
	"sync"

	"github.com/ISSuh/mystream-media_streaming/internal/model"
)

type StreamStatusMemoryRepository struct {
	engin map[int]model.Stream
	mtx   sync.Mutex
}

func NewStreamStatusMemoryRepository() *StreamStatusMemoryRepository {
	return &StreamStatusMemoryRepository{
		engin: make(map[int]model.Stream),
	}
}

func (m *StreamStatusMemoryRepository) Find(streamId int) (model.Stream, bool) {
	data, exist := m.engin[streamId]
	return data, exist
}

func (m *StreamStatusMemoryRepository) Save(stream model.Stream) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.engin[stream.StreamId] = stream
}

func (m *StreamStatusMemoryRepository) Delete(streamId int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.engin, streamId)
}
