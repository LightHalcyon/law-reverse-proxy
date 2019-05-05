package main

import (
	"io"

	"log"
	"net/http"

	// "mime/multipart"
	// "time"
	"os"
	// "strconv"
	"encoding/csv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type appError struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}

var files map[string]string

// CSVToMap reads csv to map
func CSVToMap(reader io.Reader, arr map[string]string) (map[string]string, error) {
	r := csv.NewReader(reader)
	var err error
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		arr[record[0]] = record[1]
	}
	log.Println(arr)
	return arr, err
}

func download(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")

	id := c.Param("id")

	csvf := "dl/list.csv"
	f, err := os.Open(csvf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:    http.StatusInternalServerError,
			Message: "Index file corrupted",
		})
		return
	}
	defer f.Close()
	log.Println("File Opened")
	files, err = CSVToMap(f, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if _, ok := files[id]; !ok {
		c.JSON(http.StatusNotFound, appError{
			Code:    http.StatusNotFound,
			Message: "File not found, have you uploaded it here?",
		})
		return
	}

	targetFile := files[id]
	fileName := strings.Replace(targetFile, "dl/", "", -1)

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, appError{
			Code:    http.StatusNotFound,
			Message: "File not found, file may perhaps have been moved or removed",
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(targetFile)
	return
}

func main() {
	files = make(map[string]string)
	log.Println(files)
	if _, err := os.Stat("/dl"); os.IsNotExist(err) {
		os.Mkdir("/dl", 0644)
	}

	r := gin.Default()

	r.GET("/dl/:id", download)
	// r.POST("/ul", upload)
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
