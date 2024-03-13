package fss

import "io/fs"

func MultiFS(fss ...fs.FS) fs.FS {
	return multiFs{fss}
}

type multiFs struct{ fss []fs.FS }

func (m multiFs) Open(name string) (fs.File, error) {
	for _, fs := range m.fss {
		if f, err := fs.Open(name); err == nil {
			return f, nil
		}
	}
	return nil, fs.ErrNotExist
}
