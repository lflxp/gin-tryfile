# Gin Middleware TryFile

This project is to solve the problem that the gin framework processes the dynamic routing file in the front-end compilation file and imitates nginx try_file function

本项目是解决gin框架处理前端编译文件中动态路由文件，模仿nginx try_file功能

```
location /images/ {
    root /opt/html/;
    try_files $uri   index.html; 
}
```

# Usage

`func RegisterTryFile(router *gin.Engine, hfs http.FileSystem, staticFileDir ...string)`

Parameter analysis
* router *gin.Engine
* hfs http.FileSystem
    * http.FS
    * http.Dir
    * http.FS(embed.FS)
* staticFileDir ...string
    * only one word as default,which means gin router and static path
    * include two word,first is gin router,second is static path, default try file is index.html
    * include three word,first is gin router,second is static path,three is try file name eg: index.html
# Demo

```
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
```

## Verification results

```
➜  gin-tryfile git:(main) ✗ curl http://127.0.0.1:8080/dist/this/is/not/exist
<html>
    <body>
        <h3>Hello Gin-TryFile</h3>
    </body>
</html>
➜  gin-tryfile git:(main) ✗ curl http://127.0.0.1:8080/try/file/this/is/not/exist
<html>
    <body>
        <h3>Hello Gin-TryFile</h3>
    </body>
</html>
```