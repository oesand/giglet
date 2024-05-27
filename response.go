package giglet

import (
	"giglet/safe"
	"giglet/specs"
	"io"
	"strconv"
	"sync/atomic"
)

type Response interface {
	StatusCode() *specs.StatusCode
	SetStatusCode(*specs.StatusCode)

	Header() *specs.Header
}

type PreparableResponse interface {
	Response
	Prepare()
}

type WritableResponse interface {
	Response
	Write(io.Writer)
}

type HeaderResponse struct {
	_ safe.NoCopy
	statusCode *specs.StatusCode
	header *specs.Header
}

func (resp *HeaderResponse) StatusCode() *specs.StatusCode {
	if resp.statusCode == nil {
		resp.statusCode = specs.StatusCodeOK
	}
	return resp.statusCode
}

func (resp *HeaderResponse) SetStatusCode(code *specs.StatusCode) {
	resp.statusCode = code
}

func (resp *HeaderResponse) Header() *specs.Header {
	if resp.header == nil {
		resp.header = &specs.Header{}
	}
	return resp.header
}


type OncePrepareResponse struct {
	HeaderResponse
	once atomic.Bool
}

func (resp *OncePrepareResponse) MarkOnce() bool {
	val := resp.once.Load()
	if !val {
		resp.once.Store(true)
	}
	return val
}


type TextResponse struct {
	OncePrepareResponse

	Content string
	ContentType specs.ContentType
}

func (resp *TextResponse) Prepare() {
	if resp.MarkOnce() { return }

	if resp.ContentType == specs.ContentTypeUndefined {
		resp.ContentType = specs.ContentTypePlain
	}
	resp.Header().Set("Content-Length", strconv.Itoa(len(resp.ContentType)))
	resp.Header().Set("Content-Type", string(resp.ContentType))
}

func (resp *TextResponse) Write(writer io.Writer) {
	writer.Write(safe.StringToBuffer(resp.Content))
}


type StreamResponse struct {
	OncePrepareResponse

	Stream io.Reader
	Size int64
	ContentType specs.ContentType
}

func (resp *StreamResponse) Prepare() {
	if resp.MarkOnce() { return }

	if resp.ContentType == specs.ContentTypeUndefined {
		resp.ContentType = specs.ContentTypeRaw
	}
	if resp.Size > 0 {
		resp.Header().Set("Content-Length", strconv.FormatInt(resp.Size, 10))
	}
	resp.Header().Set("Content-Type", string(resp.ContentType))
}

func (resp *StreamResponse) Write(writer io.Writer) {
	io.Copy(writer, resp.Stream)
}

