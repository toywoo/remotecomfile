package handler

import (
    "crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

func makeCodeFile(code, lang) {

}

// comfile
func comfile(res http.ResponseWriter, req *http.Request) {
	lang := req.FormValue("lang")
	code := req.FormValue("code")

	hash := sha256.New()

	hash.Write([]byte(code))

	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	codeFileName := "/codes/" + mdStr[:10] + "." + ""
    
	codeFile, err := os.Create(codeFileName)
    defer codeFile.Close()
	if err != nil {
        errCookie := http.Cookie{
			Name:     "errorServer",
			Value:    "err: create file",
			SameSite: http.SameSiteLaxMode,
		}
		res.Header().Set("Set-Cookie", errCookie.String())
		http.Redirect(res, req, MAIN_URL, http.StatusSeeOther)
	}
    fmt.Fprintf(codeFile, code)

	switch lang {
	case "py":
		// connect py compiler
	case "c":
		// connect c compiler
	case "cpp":
		//connect cpp compiler
	}
}

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