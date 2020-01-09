package dir

import (
	"log"
	"os"
	"path"

	"github.com/axetroy/dvs/internal/fs"
)

var (
	CacheDir string // cache dir
)

func init() {
	var err error

	defer func() {
		if err != nil {
			log.Panicln(err)
		}
	}()

	if c, e := os.UserCacheDir(); e != nil {
		err = e
		return
	} else {
		CacheDir = path.Join(c, "dvm")

		if e := fs.EnsureDir(CacheDir); e != nil {
			err = e
			return
		}
	}
}
