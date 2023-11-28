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
	"net/http"

	"github.com/ISSuh/mystream-media_streaming/internal/hls"
	"github.com/ISSuh/mystream-media_streaming/internal/router/response"
	"github.com/gin-gonic/gin"
)

type HlsRouter struct {
	generator *hls.Generator
}

func NewHlsRouter() *HlsRouter {
	return &HlsRouter{
		generator: hls.NewGenerator(),
	}
}

func (c *HlsRouter) View(context *gin.Context) {
	id := context.Param("streamId")
	streamPath := context.Param("streamPath")

	fmt.Println("[TEST][View] id : ", id, ", streamPath : ", streamPath)

	playlist, err := c.generator.MakeMediaM3u8()
	if err != nil {
		resp := response.Error(http.StatusInternalServerError, "generate playlist fail. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	context.Data(http.StatusOK, "application/x-mpegURL", []byte(playlist))
}

func (c *HlsRouter) Segment(context *gin.Context) {
	// id := context.Param("streamId")
	baseDir := context.Param("baseDir")
	sessionId := context.Param("sessionId")
	time := context.Param("time")
	segmentName := context.Param("segment")

	path := baseDir + "/" + sessionId + "/" + time + "/" + segmentName
	fmt.Println("[TEST][Segment] path : ", path)

	segment, err := ioutil.ReadFile(path)
	if err != nil {
		context.String(http.StatusNotFound, "")
		return
	}

	context.Data(http.StatusOK, "video/mp2t", segment)
}
