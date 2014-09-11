package web

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coffeehc/logger"
)

const (
	TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	sniffLen   = 512
)

type countingWriter int64

func (w *countingWriter) Write(p []byte) (n int, err error) {
	*w += countingWriter(len(p))
	return len(p), nil
}

type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

func (r httpRange) mimeHeader(contentType string, size int64) textproto.MIMEHeader {
	return textproto.MIMEHeader{
		"Content-Range": {r.contentRange(size)},
		"Content-Type":  {contentType},
	}
}

func FileHandler(req *http.Request, reply *Reply) {
	defer func() {
		if err := recover(); err != nil {
			reply.Error(fmt.Sprintf("%s", err), 500)
		}
	}()
	upath := req.URL.Path
	if upath == "/" || upath == "" {
		if p, ok := reply.GetInterface(Bind_Key_Welcome).(string); ok {
			upath = p
		}
	}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		req.URL.Path = upath
	}
	serveFile(reply, req, path.Clean(upath))
}

func serveFile(reply *Reply, r *http.Request, name string) {
	if fs, ok := reply.GetInterface(Bind_Key_StaticResource).(http.FileSystem); ok {
		f, err := fs.Open(name)
		if err != nil {
			logger.Debug("获取文件失败:%s", name)
			reply.NoFindPage(r)
			return
		}
		d, err := f.Stat()
		if err != nil {
			reply.NoFindPage(r)
			return
		}
		sizeFunc := func() (int64, error) { return d.Size(), nil }
		serveContent(reply, r, d.Name(), d.ModTime(), sizeFunc, f)
	} else {
		reply.NoFindPage(r)
	}
}

func localRedirect(reply *Reply, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	reply.Header().Set("Location", newPath)
	reply.SetStatusCode(http.StatusMovedPermanently)
}

func checkLastModified(reply *Reply, r *http.Request, modtime time.Time) bool {
	if modtime.IsZero() {
		return false
	}
	if t, err := time.Parse(TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := reply.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		reply.SetStatusCode(http.StatusNotModified)
		return true
	}
	reply.Header().Set("Last-Modified", modtime.UTC().Format(TimeFormat))
	return false
}

func checkETag(reply *Reply, r *http.Request) (rangeReq string, done bool) {
	etag := reply.Header().Get("Etag")
	rangeReq = r.Header.Get("Range")
	if ir := r.Header.Get("If-Range"); ir != "" && ir != etag {
		rangeReq = ""
	}
	if inm := r.Header.Get("If-None-Match"); inm != "" {
		if etag == "" {
			return rangeReq, false
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			return rangeReq, false
		}
		if inm == etag || inm == "*" {
			h := reply.Header()
			delete(h, "Content-Type")
			delete(h, "Content-Length")
			reply.SetStatusCode(http.StatusNotModified)
			return "", true
		}
	}
	return rangeReq, false
}

func serveContent(reply *Reply, r *http.Request, name string, modtime time.Time, sizeFunc func() (int64, error), content io.ReadSeeker) {
	if checkLastModified(reply, r, modtime) {
		closeFileReader(content)
		return
	}
	rangeReq, done := checkETag(reply, r)
	if done {
		closeFileReader(content)
		return
	}
	code := http.StatusOK
	// If Content-Type isn't set, use the file's extension to find it, but
	// if the Content-Type is unset explicitly, do not sniff the type.
	ctypes, haveType := reply.Header()["Content-Type"]
	var ctype string
	if !haveType {
		ctype = mime.TypeByExtension(filepath.Ext(name))
		if ctype == "" {
			var buf [sniffLen]byte
			n, _ := io.ReadFull(content, buf[:])
			ctype = http.DetectContentType(buf[:n])
			_, err := content.Seek(0, os.SEEK_SET) // rewind to output whole file
			if err != nil {
				reply.Error("无法找到你需要的资源", http.StatusInternalServerError)
				closeFileReader(content)
				return
			}
		}
		reply.Header().Set("Content-Type", ctype)
	} else if len(ctypes) > 0 {
		ctype = ctypes[0]
	}

	size, err := sizeFunc()
	if err != nil {
		reply.Error(err.Error(), http.StatusInternalServerError)
		closeFileReader(content)
		return
	}

	// handle Content-Range header.
	sendSize := size
	var sendContent io.Reader = content
	if size >= 0 {
		ranges, err := parseRange(rangeReq, size)
		if err != nil {
			reply.Error(err.Error(), http.StatusRequestedRangeNotSatisfiable)
			closeFileReader(content)
			return
		}
		if sumRangesSize(ranges) > size {
			ranges = nil
		}
		switch {
		case len(ranges) == 1:
			ra := ranges[0]
			if _, err := content.Seek(ra.start, os.SEEK_SET); err != nil {
				reply.Error(err.Error(), http.StatusRequestedRangeNotSatisfiable)
				closeFileReader(content)
				return
			}
			sendSize = ra.length
			code = http.StatusPartialContent
			reply.Header().Set("Content-Range", ra.contentRange(size))
		case len(ranges) > 1:
			for _, ra := range ranges {
				if ra.start > size {
					reply.Error(err.Error(), http.StatusRequestedRangeNotSatisfiable)
					closeFileReader(content)
					return
				}
			}
			sendSize = rangesMIMESize(ranges, ctype, size)
			code = http.StatusPartialContent
			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)
			reply.Header().Set("Content-Type", "multipart/byteranges; boundary="+mw.Boundary())
			sendContent = pr
			go func() {
				defer closeFileReader(content)
				for _, ra := range ranges {
					part, err := mw.CreatePart(ra.mimeHeader(ctype, size))
					if err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := content.Seek(ra.start, os.SEEK_SET); err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := io.CopyN(part, content, ra.length); err != nil {
						pw.CloseWithError(err)
						return
					}
				}
				mw.Close()
				pw.Close()
			}()
		}
		reply.Header().Set("Accept-Ranges", "bytes")
		if reply.Header().Get("Content-Encoding") == "" {
			reply.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
		}
	}
	reply.SetStatusCode(code)
	if r.Method != REQUEST_METHOD_HEAD {
		reply.WithReader(sendContent, sendSize)
	}
}

func closeFileReader(content io.Reader) {
	if closer, ok := content.(io.Closer); ok {
		closer.Close()
	}
}

func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i > size || i < 0 {
				return nil, errors.New("invalid range")
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func sumRangesSize(ranges []httpRange) (size int64) {
	for _, ra := range ranges {
		size += ra.length
	}
	return
}

func rangesMIMESize(ranges []httpRange, contentType string, contentSize int64) (encSize int64) {
	var w countingWriter
	mw := multipart.NewWriter(&w)
	for _, ra := range ranges {
		mw.CreatePart(ra.mimeHeader(contentType, contentSize))
		encSize += ra.length
	}
	mw.Close()
	encSize += int64(w)
	return
}
