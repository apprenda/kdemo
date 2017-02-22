// The MIT License (MIT)
//
// Copyright (c) 2016 Apprenda Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Grid defines a grid of iframes
type Grid struct {
	Cols int
	Rows int
}

// Range returns an array of float arrays to use for iterating inside the HTML template.
// The inner-array contents are the column widths (in %) to be used by the iframes.
func (grid Grid) Range() [][]float32 {
	width := 100.00 / float32(grid.Cols)
	cols := make([]float32, grid.Cols)
	for i := range cols {
		cols[i] = width
	}
	rows := make([][]float32, grid.Rows)
	for i := range rows {
		rows[i] = cols
	}
	return rows
}

// Super simple Kubernetes demo app that shows the $HOSTNAME of the Pod running it.
// The `/grid` path can generate a grid of iframes to show many Pods at once.
// It accepts two parameters, `cols` and `rows` to define how the iframe grid is generated.
func main() {

	bgColor := "green"

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bgColor = "red"
	}()

	const pod = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="Cache-control" content="No-Cache">
   	    <link type="image/x-icon" href="http://apprenda.com/favicon.ico" rel="shortcut icon" />
		<title>K8s Simple Demo</title>
	</head>
	<body style="font-family: Helvetica; font-size: .7em; text-align: center">
		<h1>I'm a Kubernetes Pod!</h1>
        <h3 style="padding: 1em; background-color: %s; color: white;">ID: %q</h3>
	</body>
</html>`

	const gridTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="Cache-control" content="No-Cache">
    	<link type="image/x-icon" href="http://apprenda.com/favicon.ico" rel="shortcut icon" />
		<title>K8s Simple Demo</title>
	</head>
	<body style="margin: 0; padding: 0; font-family: Helvetica; text-align: center">
        {{range .}}
        <div>
            {{range .}}
            <iframe id="ifrm" src="/" scrolling="no" style="border: 0; margin: 0; padding: 0; float: left; width: {{.}}%;"></iframe>
            {{end}}
        </div>
        {{end}}
	</body>
</html>`

	hostname, _ := os.Hostname()

	tmpl := template.Must(template.New("gridTmpl").Parse(gridTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		fmt.Fprintf(w, pod, bgColor, hostname)
	})

	http.HandleFunc("/grid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		var cols, rows int
		cols, err := strconv.Atoi(r.FormValue("cols"))
		if err != nil {
			cols = 1
		}
		rows, err = strconv.Atoi(r.FormValue("rows"))
		if err != nil {
			rows = 1
		}
		grid := Grid{cols, rows}
		tmpl.Execute(w, grid.Range())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
