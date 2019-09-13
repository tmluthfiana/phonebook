package knot

import "net/http"

var _ http.FileSystem = CustomFileSystem{}

type CustomFileSystem struct {
	source http.FileSystem
}

func (f CustomFileSystem) Open(path string) (http.File, error) {
	file, err := f.source.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, err
	}

	return file, err
}
