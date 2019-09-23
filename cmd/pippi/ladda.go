package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func logger(handler http.Handler) http.Handler {
	logger := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logger)
}

func recvUploads() {
	handler := logger(http.HandlerFunc(upload))
	if err := http.ListenAndServe(uploadAddr, handler); err != nil {
		log.Fatal(err)
	}
}

const (
	uploadAddr = ":2000"
	ext        = ".bin"
)

// upload handles upload requests.
func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if _, err := fmt.Fprintf(w, uploadPage[1:]); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		return
	case "OPTIONS":
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,content-type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Allow", "POST")
		if _, err := w.Write([]byte("POST")); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		return
	case "POST":
		if err := r.ParseMultipartForm(1024 * 1024 * 1024); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		defer file.Close()
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, file); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		cacheDir = filepath.Join(cacheDir, "pippi")
		rawHash := sha256.Sum256(buf.Bytes())
		hash := hex.EncodeToString(rawHash[:])
		dirName := filepath.Join(cacheDir, hash)
		if err := os.Mkdir(dirName, 0755); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
		dstPath := filepath.Join(dirName, hash+ext)
		log.Printf("creating %q", dstPath)
		if err := ioutil.WriteFile(dstPath, buf.Bytes(), 0644); err != nil {
			log.Printf("%+v", errors.WithStack(err))
			return
		}
	}
}

const uploadPage = `
<!doctype html>
<html>
	<head>
		<title>upload file</title>
	</head>
	<body>
		<form enctype="multipart/form-data" action="/" method="POST">
			<input type="file" name="file">
			<input type="submit" value="upload">
		</form>
	</body>
</html>
`
