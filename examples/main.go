package main

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	tryfile "github.com/lflxp/gin-tryfile"
)

//go:embed dist
var distFile embed.FS

func main() {
	r := gin.Default()
	// As Default, Gin Router equal StaticPath
	tryfile.RegisterTryFile(r, http.FS(distFile), "/dist")

	// Custom Gin Router and StaticPath
	// tryfile.RegisterTryFile(r, http.FS(distFile), "/static", "/dist")

	// Fully custom parameters
	// tryfile.RegisterTryFile(r, http.FS(distFile), "/static", "/dist/custom", "index.html")

	// As you wish
	r.Any("/try/file/*any", tryfile.WrapHandler("/dist", http.FS(distFile)))
	r.Run()
}
