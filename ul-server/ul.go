package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/csv"
	"io"
	"os"
	"math/rand"
	"fmt"
	"bytes"
	"path/filepath"
	"time"
	"crypto/md5"
	"strconv"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
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
	return arr, err
}

// MapToCSV writes map to csv
func MapToCSV(writer io.Writer, arr map[string]string) error {
	w := csv.NewWriter(writer)
	var err error
	r := make([][]string, len(arr))
	i := 0
	for k, v := range arr {
		r[i] = make([]string, 2)
		r[i][0] = k
		r[i][1] = v
		log.Println(r)
		if err != nil {
			break
		}
		i++
	}
	err = w.WriteAll(r)
	log.Println(err)
	return err
}

// GetMD5Hash hashes string with md5
func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

// TokenGenerator generates token for as key for downloading
func TokenGenerator() string {
	b := make([]byte, 18)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func upload(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")

	clientIP := c.ClientIP()
	// clientIP := "127.0.0.1"
	log.Println(clientIP)
	timestamp := strconv.FormatInt(time.Now().UTC().Add(time.Hour * time.Duration(2)).Unix(), 10)
	hash := GetMD5Hash(clientIP + " secret")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		// log.Println(err)
		c.JSON(http.StatusBadRequest, appError{
			Code:    http.StatusBadRequest,
			Message: "File get error, did you upload a file?",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, appError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	sep := string(filepath.Separator)
	filename := "/dl/" + filepath.Base(sep+"dl"+sep+fileHeader.Filename)

	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)

	// key := TokenGenerator()
	// files[key] = filename
	// log.Println(files[key])

	// csv, err := os.OpenFile("/dl/list.csv", os.O_WRONLY, 0660)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, appError{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: "Index file corrupted",
	// 	})
	// 	return
	// }
	// defer csv.Close()
	// err = MapToCSV(csv, files)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, appError{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, appError{
		Code: http.StatusOK,
		// Message:	os.Getenv("PROXYURL") + "/download/" + fileHeader.Filename +"?md5=" + hash + "&expires=" + timestamp,
		Message: "http://localhost:22106/download/" + fileHeader.Filename +"?md5=" + hash + "&expires=" + timestamp,
	})
	return
}

func main() {
	if _, err := os.Stat("dl"); os.IsNotExist(err) {
		os.Mkdir("dl", 0755)
	}
	// files = make(map[string]string)

	// if _, err := os.Stat("/dl/list.csv"); os.IsNotExist(err) {
	// 	f, _ := os.Create("/dl/list.csv")
	// 	defer f.Close()
	// } else {
	// 	f, err := os.Open("/dl/list.csv")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	files, err = CSVToMap(f, files)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// log.Println(files)

	router := gin.Default()
	router.Use(static.Serve("/static", static.LocalFile("static", true)))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Upload Page",
		})
	})
	router.POST("/", upload)
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	conf.AddAllowHeaders("X-ROUTING-KEY")
	conf.AddAllowHeaders("Content-Type")
	conf.AddAllowHeaders("Access-Control-Allow-Origin")
	conf.AddAllowHeaders("Access-Control-Allow-Headers")
	conf.AddAllowHeaders("Access-Control-Allow-Methods")
	conf.AddAllowHeaders("Host")
	router.Use(cors.New(conf))
	router.Run("0.0.0.0:23063")
}