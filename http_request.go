package giglet

import (
	"errors"
	"giglet/safe"
	"giglet/specs"
	"io"
	"mime/multipart"
	"net"
)

type HttpRequest struct {
	_ safe.NoCopy

	server *Server
	conn net.Conn
	hijacker HijackHandler
	extras map[string]any

	protoMajor, protoMinor uint16
	method specs.HttpMethod
	url *specs.Url
	header *specs.ReadOnlyHeader

	body io.Reader
	bodyReaded bool
	cachedBody []byte
	cachedMultipart *multipart.Form
	cachedForm specs.Query
}

func (req *HttpRequest) ProtoAtLeast(major, minor uint16) bool {
	return req.protoMajor > major ||
		req.protoMajor == major && req.protoMinor >= minor
}

func (req *HttpRequest) ProtoNoHigher(major, minor uint16) bool {
	return req.protoMajor < major ||
		req.protoMajor == major && req.protoMinor <= minor
}

func (req *HttpRequest) Read(buf []byte) (n int, err error) {
	if req.body == nil || req.bodyReaded {
		return 0, io.EOF
	}
	if req.server != nil {
		req.server.applyReadTimeout(req.conn)
		defer req.conn.SetDeadline(zeroTime)
	}
	
	n, err = req.body.Read(buf)
	if err == io.EOF {
		req.bodyReaded = true
	}
	return n, err
}

func (req *HttpRequest) PostBody() (buf []byte, err error) {
	if req.body == nil || (req.bodyReaded && req.cachedBody == nil) {
		return nil, io.EOF
	} else if req.cachedBody != nil {
		return req.cachedBody, nil
	}

	buf, err = io.ReadAll(req)
	if err == nil {
		req.cachedBody = buf
	}
	
	return
}

func (req *HttpRequest) RemoteAddr() net.Addr {
	return req.conn.RemoteAddr()
}

func (req *HttpRequest) Hijack(handler HijackHandler) {
	req.hijacker = handler
}

func (req *HttpRequest) Method() specs.HttpMethod {
	return req.method
}

func (req *HttpRequest) Url() *specs.Url {
	return req.url
}

func (req *HttpRequest) Header() *specs.ReadOnlyHeader {
	return req.header
}

func (req *HttpRequest) GetExtra(key string) any {
	if req.extras == nil {
		return nil
	}
	return req.extras[key]
}

func (req *HttpRequest) SetExtra(key string, value any) {
	if req.extras == nil {
		req.extras = map[string]any{}
	}
	req.extras[key] = value
}

func (req *HttpRequest) PostForm() (specs.Query, error) {
	if req.body == nil {
		return nil, io.EOF
	} else if req.cachedForm != nil {
		return req.cachedForm, nil
	} else if req.Header().ContentType() != specs.ContentTypeMultipart {
		form, err := req.MultipartForm()
		if err !=  nil {
			return nil, err
		}
		return form.Value, nil
	} else if req.Header().ContentType() != specs.ContentTypeForm {
		return nil, errors.New("this Content-Type is not a urlencoded-form")
	} else if req.bodyReaded {
		return nil, nil
	}
	req.bodyReaded = true

	buf, err := io.ReadAll(req)
	if err != nil {
		return nil, err
	}
	req.cachedForm, err = specs.ParseQuery(string(buf))
	if err != nil {
		return nil, err
	}
	return req.cachedForm, nil
}

func (req *HttpRequest) MultipartForm() (*multipart.Form, error) {
	if req.body == nil {
		return nil, io.EOF
	} else if req.Header().ContentType() != specs.ContentTypeMultipart {
		return nil, errors.New("this Content-Type is not a multipart-form")
	} else if req.cachedMultipart != nil {
		return req.cachedMultipart, nil
	} else if req.bodyReaded {
		return nil, nil
	}
	req.bodyReaded = true

	boundary := req.header.GetMediaParams("boundary")
	if len(boundary) == 0 {
		return nil, errors.New("this request Content-Type does not contains boundary")
	}

	reader := multipart.NewReader(req, boundary)
	form, err := reader.ReadForm(0)
	if err != nil {
		return nil, err
	}
	req.cachedMultipart = form
	req.cachedForm = req.cachedMultipart.Value
	return req.cachedMultipart, nil
}
