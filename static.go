// static
package web

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const sniffLen = 512

type StaticService struct {
	root      http.FileSystem
	urlPrefix string
}

func NewStaticFilter(root http.FileSystem, urlPrefix string) ActionFilter {
	service := &StaticService{root, urlPrefix}
	return service.StaticFilter
}

func (this *StaticService) StaticFilter(request *http.Request, reply *Reply, chain FilterChain) {
	upath := strings.TrimPrefix(request.URL.Path, this.urlPrefix)
	if len(upath) < len(request.URL.Path) {
		request.URL.Path = upath
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			request.URL.Path = upath
		}
		serveFile(reply, request, this.root, path.Clean(upath), true)
	} else {
		chain(request, reply)
	}
}

func localRedirect(reply *Reply, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	reply.SetHeader("Location", newPath)
	reply.SetCode(http.StatusMovedPermanently)
}

func serveFile(reply *Reply, r *http.Request, fs http.FileSystem, name string, redirect bool) {
	const indexPage = "/index.html"
	if strings.HasSuffix(r.URL.Path, indexPage) {
		localRedirect(reply, r, "./")
		return
	}
	f, err := fs.Open(name)
	if err != nil {
		reply.SetCode(http.StatusNotFound)
		return
	}
	d, err1 := f.Stat()
	if err1 != nil {
		reply.SetCode(http.StatusNotFound)
		return
	}
	if redirect {
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(reply, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(reply, r, "../"+path.Base(url))
				return
			}
		}
	}
	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(index)
		if err == nil {
			dd, err := ff.Stat()
			if err == nil {
				name = index
				d = dd
				f = ff
			}
		}
	}
	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		if checkLastModified(reply, r, d.ModTime()) {
			return
		}
		dirList(reply, f)
		return
	}
	// serveContent will check modification time
	sizeFunc := func() (int64, error) { return d.Size(), nil }
	serveContent(reply, r, d.Name(), d.ModTime(), sizeFunc, f)
}

func serveContent(reply *Reply, r *http.Request, name string, modtime time.Time, sizeFunc func() (int64, error), content io.ReadSeeker) {
	if checkLastModified(reply, r, modtime) {
		return
	}
	rangeReq, done := checkETag(reply, r, modtime)
	if done {
		return
	}
	code := http.StatusOK
	ctype := mime.TypeByExtension(filepath.Ext(name))
	if ctype == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [sniffLen]byte
		n, _ := io.ReadFull(content, buf[:])
		ctype = http.DetectContentType(buf[:n])
		_, err := content.Seek(0, os.SEEK_SET) // rewind to output whole file
		if err != nil {
			reply.SetCode(http.StatusInternalServerError)
			reply.With("seeker can't seek")
			return
		}
	}
	reply.SetContentType(ctype)

	size, err := sizeFunc()
	if err != nil {
		reply.SetCode(http.StatusInternalServerError)
		reply.With(err.Error())
		return
	}
	// handle Content-Range header.
	sendSize := size
	var sendContent io.Reader = content
	if size >= 0 {
		ranges, err := parseRange(rangeReq, size)
		if err != nil {
			reply.SetCode(http.StatusRequestedRangeNotSatisfiable).With(err.Error())
		}
		if sumRangesSize(ranges) > size {
			ranges = nil
		}
		switch {
		case len(ranges) == 1:
			ra := ranges[0]
			if _, err := content.Seek(ra.start, os.SEEK_SET); err != nil {
				reply.SetCode(http.StatusRequestedRangeNotSatisfiable).With(err.Error())
				return
			}
			sendSize = ra.length
			code = http.StatusPartialContent
			reply.SetHeader("Content-Range", ra.contentRange(size))
		case len(ranges) > 1:
			sendSize = rangesMIMESize(ranges, ctype, size)
			code = http.StatusPartialContent
			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)
			reply.SetHeader("Content-Type", "multipart/byteranges; boundary="+mw.Boundary())
			sendContent = pr
			// cause writing goroutine to fail and exit if CopyN doesn't finish.
			go func() {
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

		reply.SetHeader("Accept-Ranges", "bytes")
		if reply.GetHeader("Content-Encoding") == "" {
			reply.SetHeader("Content-Length", strconv.FormatInt(sendSize, 10))
		}
	}
	reply.SetCode(code)
	if r.Method != "HEAD" {
		reply.With(io.LimitReader(sendContent, sendSize))
	}
}

func checkLastModified(reply *Reply, r *http.Request, modtime time.Time) bool {
	if modtime.IsZero() {
		return false
	}

	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(time.Second)) {
		reply.DelHeader("Content-Type")
		reply.DelHeader("Content-Length")
		reply.SetCode(http.StatusNotModified)
		return true
	}
	reply.SetHeader("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return false
}

func dirList(reply *Reply, f http.File) {
	reply.SetHeader("Content-Type", "text/html; charset=utf-8")
	buf := bytes.NewBuffer(nil)
	buf.WriteString("<pre>\n")
	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			// name may contain '?' or '#', which must be escaped to remain
			// part of the URL path, and not indicate the start of a query
			// string or fragment.
			url := url.URL{Path: name}
			fmt.Fprintf(buf, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
		}
	}
	fmt.Fprintf(buf, "</pre>\n")
	reply.With(buf)
}

func checkETag(reply *Reply, r *http.Request, modtime time.Time) (rangeReq string, done bool) {
	etag := reply.GetHeader("Etag")
	rangeReq = r.Header.Get("Range")

	if ir := r.Header.Get("If-Range"); ir != "" && ir != etag {
		// The If-Range value is typically the ETag value, but it may also be
		// the modtime date. See golang.org/issue/8367.
		timeMatches := false
		if !modtime.IsZero() {
			if t, err := http.ParseTime(ir); err == nil && t.Unix() == modtime.Unix() {
				timeMatches = true
			}
		}
		if !timeMatches {
			rangeReq = ""
		}
	}

	if inm := r.Header.Get("If-None-Match"); inm != "" {
		// Must know ETag.
		if etag == "" {
			return rangeReq, false
		}

		// TODO(bradfitz): non-GET/HEAD requests require more work:
		// sending a different status code on matches, and
		// also can't use weak cache validators (those with a "W/
		// prefix).  But most users of ServeContent will be using
		// it on GET or HEAD, so only support those for now.
		if r.Method != "GET" && r.Method != "HEAD" {
			return rangeReq, false
		}

		// TODO(bradfitz): deal with comma-separated or multiple-valued
		// list of If-None-match values.  For now just handle the common
		// case of a single item.
		if inm == etag || inm == "*" {
			//			h := w.Header()
			reply.DelHeader("Content-Type")
			reply.DelHeader("Content-Length")
			reply.SetCode(http.StatusNotModified)
			return "", true
		}
	}
	return rangeReq, false
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

func sumRangesSize(ranges []httpRange) (size int64) {
	for _, ra := range ranges {
		size += ra.length
	}
	return
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

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

type countingWriter int64

func (w *countingWriter) Write(p []byte) (n int, err error) {
	*w += countingWriter(len(p))
	return len(p), nil
}
