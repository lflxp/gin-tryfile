package tryfile

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
)

/* dist =>
   //go:embed dist
   var dist embed.FS

   http.FileSystem => http.FS(dist)
*/
func RegisterTryFile(router *gin.Engine, hfs http.FileSystem, staticFileDir ...string) {
	if len(staticFileDir) == 1 {
		router.Any(fmt.Sprintf("%s/*any", staticFileDir[0]), WrapHandler(staticFileDir[0], hfs))
	} else if len(staticFileDir) == 2 {
		router.Any(fmt.Sprintf("%s/*any", staticFileDir[0]), WrapHandler(staticFileDir[1], hfs))
	} else if len(staticFileDir) == 3 {
		router.Any(fmt.Sprintf("%s/*any", staticFileDir[0]), func(c *gin.Context) {
			server := TryFileCustom(staticFileDir[1], staticFileDir[2], hfs)
			server.ServeHTTP(c.Writer, c.Request)
		})
	} else {
		log.Error("staticFileDir Params is invaild")
	}
}

func WrapHandler(staticFileDir string, hfs http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		server := TryFileCustom(staticFileDir, "", hfs)
		server.ServeHTTP(c.Writer, c.Request)
	}
}

func TryFileCustom(staticFileDir, staticFile string, hfs http.FileSystem) http.HandlerFunc {
	if staticFile == "" {
		staticFile = "index.html"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		hserver := http.FileServer(hfs)
		// catch 404 error
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		hserver.ServeHTTP(nfrw, r)
		if nfrw.status == 404 {
			log.Debugf("Redirecting %s to index.html", r.RequestURI)
			targetFile := fmt.Sprintf("%s/%s", staticFileDir, staticFile)

			f, err := hfs.Open(targetFile)
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}
			defer f.Close()

			d, err := f.Stat()
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}

			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Header().Add("Accept-Ranges", "bytes")
			w.Header().Add("Content-Length", string(d.Size()))
			http.ServeContent(w, r, targetFile, d.ModTime(), f)
		}
	}
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
