package web

import "io/fs"

var templateFs fs.FS

func SetTemplateFs(fsys fs.FS) {
	templateFs = fsys
}
