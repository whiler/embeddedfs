package embeddedfs

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	ErrInvalidOffset = errors.New("invalid offset")
	ErrInvalidWhence = errors.New("invalid whence")
	ErrInvalidDir    = errors.New("invalid dir")
	ErrInvalidCount  = errors.New("invalid count")
)

type FileSystem interface {
	Open(name string) (http.File, error)
	Stat(name string) (os.FileInfo, error)
}

type DefaultFileSystem struct{}

func (DefaultFileSystem) Open(name string) (http.File, error) {
	return os.Open(name)
}

func (DefaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

type EmbeddedFileSystem map[string]*EmbeddedFile

func (fs EmbeddedFileSystem) Open(name string) (http.File, error) {
	if basic, exists := fs[name]; !exists {
		return nil, os.ErrNotExist
	} else {
		return &embeddedFile{info: basic.Info, content: basic.Content, children: basic.Children}, nil
	}
}

func (fs EmbeddedFileSystem) Stat(name string) (os.FileInfo, error) {
	if basic, exists := fs[name]; !exists {
		return nil, os.ErrNotExist
	} else {
		return basic.Info, nil
	}
}

type EmbeddedFile struct {
	Info     *FileInfo
	Content  []byte
	Children []*FileInfo
}

type embeddedFile struct {
	info     *FileInfo
	content  []byte
	children []*FileInfo
	offset   int64
	closed   bool
}

func (f *embeddedFile) Close() error {
	f.closed = true
	return nil
}

func (f *embeddedFile) Read(p []byte) (n int, err error) {
	if f.closed {
		return 0, os.ErrClosed
	}
	if f.offset >= int64(len(f.content)) {
		return 0, io.EOF
	}
	n = copy(p, f.content[f.offset:])
	f.offset += int64(n)
	return
}

func (f *embeddedFile) ReadAt(p []byte, offset int64) (n int, err error) {
	if f.closed {
		return 0, os.ErrClosed
	}
	if offset < 0 {
		return 0, ErrInvalidOffset
	}
	if offset >= int64(len(f.content)) {
		return 0, io.EOF
	}
	n = copy(p, f.content[offset:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (f *embeddedFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed {
		return 0, os.ErrClosed
	}
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = f.offset + offset
	case io.SeekEnd:
		abs = int64(len(f.content)) + offset
	default:
		return 0, ErrInvalidWhence
	}
	if abs < 0 {
		return 0, ErrInvalidOffset
	}
	f.offset = abs
	return abs, nil
}

func (f *embeddedFile) Stat() (os.FileInfo, error) {
	if f.closed {
		return nil, os.ErrClosed
	}
	return f.info, nil
}

func (f *embeddedFile) Readdir(count int) ([]os.FileInfo, error) {
	switch {
	case f.closed:
		return nil, os.ErrClosed
	case !f.info.RawMode.IsDir():
		return nil, ErrInvalidDir
	case 0 > count:
		ret := make([]os.FileInfo, len(f.children))
		for i, info := range f.children {
			ret[i] = info
		}
		return ret, nil
	case len(f.children) < count:
		return nil, ErrInvalidCount
	default:
		ret := make([]os.FileInfo, count)
		for i, info := range f.children[0:count] {
			ret[i] = info
		}
		return ret, nil
	}
}

type FileInfo struct {
	RawName    string
	RawSize    int64
	RawMode    os.FileMode
	RawModTime time.Time
}

func (i *FileInfo) Name() string {
	return i.RawName
}

func (i *FileInfo) Size() int64 {
	return i.RawSize
}

func (i *FileInfo) Mode() os.FileMode {
	return i.RawMode
}

func (i *FileInfo) ModTime() time.Time {
	return i.RawModTime
}

func (i *FileInfo) IsDir() bool {
	return i.RawMode.IsDir()
}

func (i *FileInfo) Sys() interface{} {
	return nil
}
