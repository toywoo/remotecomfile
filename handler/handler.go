package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	service "classNote/service"
)

const MAIN_URL = "http://localhost:3002"

// index
func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")

	htmlFile, readError := ioutil.ReadFile("./Client/main.html")
	if readError != nil {
		http.NotFound(res, req)
	} else {
		res.Write(htmlFile)
	}
}

// comfile
func comfile(res http.ResponseWriter)

func (handler *Handler) PathNav(res http.ResponseWriter, req *http.Request) {
	var path = req.URL.Path
	splitedPath := strings.Split(path, "/")

	switch splitedPath[1] {
	case "":
		index(res, req)

	case "comfile":
	    comfile()

	default:
		http.NotFound(res, req)
	}

}