//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../proto/bin ../../proto/upload.proto

package main

import (
	"log"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func basic(binId string) []string {
	extractedStrings, err := Strings("localhost:1234", binId)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return extractedStrings
}

func main() {
	js := mewn.String("./frontend/dist/my-app/main-es2015.js")
	css := mewn.String("./frontend/dist/my-app/styles.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "Pippi",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)
	app.Run()
}
