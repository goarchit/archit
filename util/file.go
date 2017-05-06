//  Filename manipulation routines initially written on 5/6/17
package util

import (
	"github.com/goarchit/archit/log"
	"os/user"
	"path/filepath"
	"strings"
)

func FullPath(path string) string {
	fullpath := path
        if len(path) >= 2 && path[:2] == "~/" {
                usr, _ := user.Current()
                dir := usr.HomeDir + "/"  //  This is unix specific, so "/" is safe
                fullpath = strings.Replace(path, "~/", dir, 1)

        }
	fullpath,err := filepath.Abs(fullpath)
	if err != nil {
		log.Critical("Error getting absolute pathname from",fullpath,":",err)
	}
        log.Trace(path,"expanded to", fullpath)
	return fullpath
}
