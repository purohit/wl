package oauth

import (
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type LoggableRequest struct {
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             io.ReadCloser
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
}

func newLoggableRequest(req http.Request) LoggableRequest {
	var form, postForm url.Values
	if req.Form != nil {
		form = sanitizeCredentialsFromForm(req.Form)
	}

	if req.PostForm != nil {
		postForm = sanitizeCredentialsFromForm(req.PostForm)
	}

	req.Header["Authorization"] = nil

	return LoggableRequest{
		Method:           req.Method,
		URL:              req.URL,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Header:           req.Header,
		Body:             req.Body,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Close:            req.Close,
		Host:             req.Host,
		Form:             form,
		PostForm:         postForm,
		MultipartForm:    req.MultipartForm,
		Trailer:          req.Trailer,
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
		TLS:              req.TLS,
	}
}

func sanitizeCredentialsFromForm(form url.Values) url.Values {
	form.Set("password", "***")
	return form
}
