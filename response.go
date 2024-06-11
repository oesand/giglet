package giglet

import (
	"fmt"
	"giglet/safe"
	"giglet/specs"
	"io"
	"strconv"
	"sync/atomic"
)

type Response interface {
	StatusCode() specs.StatusCode
	Header() *specs.Header
}

type PreparableResponse interface {
	Response
	Prepare()
}

type WritableResponse interface {
	Response
	WriteBody(io.Writer)
}

type HeaderResponse struct {
	_ safe.NoCopy

	statusCode *specs.StatusCode
	header *specs.Header
}

func (resp *HeaderResponse) StatusCode() specs.StatusCode {
	if resp.statusCode == nil {
		resp.statusCode = &specs.StatusCodeOK
	}
	return *resp.statusCode
}

func (resp *HeaderResponse) SetStatusCode(code specs.StatusCode) Response {
	resp.statusCode = &code
	return resp
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


type RedirectResponse struct {
	header *specs.Header

	Url string
	Permanent bool
}

func (resp *RedirectResponse) StatusCode() *specs.StatusCode {
	if resp.Permanent {
		return &specs.StatusCodePermanentRedirect
	}
	return &specs.StatusCodeTemporaryRedirect
}

func (resp *RedirectResponse) Header() *specs.Header {
	if resp.header == nil {
		resp.header = &specs.Header{}
		resp.header.Set("Location", resp.Url)
	}
	return resp.header
}

func (resp *RedirectResponse) Error() string {
	output := "redirect to: " + resp.Url
	if resp.Permanent {
		output = "permanent " + output 
	}
	return output
}


type TextResponse struct {
	OncePrepareResponse

	Text string
	ContentType specs.ContentType
}

func (resp *TextResponse) Prepare() {
	if resp.MarkOnce() { return }

	if resp.ContentType == specs.ContentTypeUndefined {
		resp.ContentType = specs.ContentTypePlain
	}
	resp.Header().Set("Content-Length", strconv.Itoa(len(resp.Text)))
	resp.Header().Set("Content-Type", string(resp.ContentType))
}

func (resp *TextResponse) WriteBody(writer io.Writer) {
	writer.Write(safe.StringToBuffer(resp.Text))
}

func (resp *TextResponse) Error() string {
	return fmt.Sprintf("text response<%d>: %s", resp.statusCode.Code, resp.Text)
}

type StreamResponse struct {
	OncePrepareResponse

	Stream io.Reader
	Size uint64
	ContentType specs.ContentType
}

func (resp *StreamResponse) Prepare() {
	if resp.MarkOnce() { return }

	if resp.ContentType == specs.ContentTypeUndefined {
		resp.ContentType = specs.ContentTypeRaw
	}
	if resp.Size > 0 {
		resp.Header().Set("Content-Length", strconv.FormatUint(resp.Size, 10))
	}
	resp.Header().Set("Content-Type", string(resp.ContentType))
}

func (resp *StreamResponse) WriteBody(writer io.Writer) {
	io.Copy(writer, resp.Stream)
}

func (resp *StreamResponse) Error() string {
	return fmt.Sprintf("stream response<%d> %s", resp.statusCode.Code, strconv.FormatUint(resp.Size, 10))
}

