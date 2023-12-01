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

package controller

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_streaming/internal/api/response"
	"github.com/ISSuh/mystream-media_streaming/internal/service"
	"github.com/gin-gonic/gin"
)

type StraemController struct {
	streamService *service.StreamService
}

func NewStraemController(streamService *service.StreamService) *StraemController {
	return &StraemController{
		streamService: streamService,
	}
}

func (c *StraemController) Test(context *gin.Context) {
	log.Info("[StraemController][Test]")
	context.Status(http.StatusOK)
}

func (c *StraemController) MasterPlaylistOptions(context *gin.Context) {
	context.Status(http.StatusOK)
}

func (c *StraemController) MasterPlaylist(context *gin.Context) {
	streamPath := context.Param("streamPath")
	streamId, err := strconv.Atoi(context.Param("streamId"))
	if err != nil {
		resp := response.Error(http.StatusBadRequest, "invalid streamId. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Info("[StraemController][MasterPlaylist] id : ", streamId, ", streamPath : ", streamPath)

	playlist, err := c.streamService.FindMasterPlaylist(streamId, streamPath)
	if err != nil {
		resp := response.Error(http.StatusBadRequest, "generate playlist fail. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	context.Data(http.StatusOK, "application/x-mpegURL", []byte(playlist))
}

func (c *StraemController) MediaPlaylistOptions(context *gin.Context) {
	context.Status(http.StatusOK)
}

func (c *StraemController) MediaPlaylist(context *gin.Context) {
	streamPath := context.Param("streamPath")
	streamId, err := strconv.Atoi(context.Param("streamId"))
	if err != nil {
		resp := response.Error(http.StatusBadRequest, "invalid streamId. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Info("[StraemController][MediaPlaylist] id : ", streamId, ", streamPath : ", streamPath)
	playlist, err := c.streamService.FindMediaPlaylist(streamId, streamPath)
	if err != nil {
		resp := response.Error(http.StatusBadRequest, "generate playlist fail. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	context.Data(http.StatusOK, "application/x-mpegURL", []byte(playlist))
}

func (c *StraemController) Segment(context *gin.Context) {
	streamPath := context.Param("streamPath")
	segmentName := context.Param("segment")
	streamId, err := strconv.Atoi(context.Param("streamId"))
	if err != nil {
		resp := response.Error(http.StatusBadRequest, "invalid streamId. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Info("[StraemController][Segment] id : ", streamId, ", streamPath : ", streamPath, ", segmentName : ", segmentName)
	segment, err := c.streamService.FindSegment(streamId, streamPath, segmentName)
	if err != nil {
		resp := response.Error(http.StatusNotFound, "not found segment. "+err.Error())
		context.JSON(http.StatusInternalServerError, resp)
		return
	}

	context.Data(http.StatusOK, "video/mp2t", segment)
}
