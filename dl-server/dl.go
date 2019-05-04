package main

import (
	"bytes"
	"io"
	// "log"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"math/rand"
	// "mime/multipart"
	// "time"
	// "os"
	// "strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

type appError struct {
	Code	int    `json:"status"`
	Message	string `json:"message"`
}

var files map[string]string

// TokenGenerator generates token for as key for downloading
func TokenGenerator() string {
	b := make([]byte, 18)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func download(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin","*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")

	id := c.Param("id")

	if _, ok := files[id]; !ok {
		c.JSON(http.StatusNotFound, appError{
			Code:		http.StatusNotFound,
			Message:	"File not found, have you uploaded it here?",
		})
		return
	}

	targetFile := files[id]
	fileName := strings.Replace(targetFile, "dl/" , "", -1)

    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", "attachment; filename=" + fileName)
    c.Header("Content-Type", "application/octet-stream")
	c.File(targetFile)
	return
}

func upload(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin","*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		// log.Println(err)
		c.JSON(http.StatusBadRequest, appError{
			Code:		http.StatusBadRequest,
			Message:	"File get error, did you upload a file?",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:		http.StatusInternalServerError,
			Message:	err.Error(),
		})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:		http.StatusInternalServerError,
			Message:	err.Error(),
		})
		return
	}

	sep := string(filepath.Separator)
	filename := "dl/" + filepath.Base(sep + "dl" + sep + fileHeader.Filename)

	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)

	key := TokenGenerator()
	files[key] = filename

	c.JSON(http.StatusOK, appError{
		Code:		http.StatusOK,
		// Message:	os.getenv("PROXYURL") + os.getenv("DLROUTE") + key,
		Message: "localhost/download/" + key,
	})
}

func main() {
	files = make(map[string]string)

	r := gin.Default()

	r.GET("/:id", download)
	r.POST("/ul", upload)
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	conf.AddAllowHeaders("X-ROUTING-KEY")
	conf.AddAllowHeaders("Content-Type")
	conf.AddAllowHeaders("Access-Control-Allow-Origin")
	conf.AddAllowHeaders("Access-Control-Allow-Headers")
	conf.AddAllowHeaders("Access-Control-Allow-Methods")
	conf.AddAllowHeaders("Host")
	r.Use(cors.New(conf))

	r.Run("0.0.0.0:23062")
}