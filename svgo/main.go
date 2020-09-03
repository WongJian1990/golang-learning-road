package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
)

var width = 800
var height = 400
var startTime = time.Now().UnixNano()

func drawPoint(osvg *svg.SVG, pnt int, process int) {
	sec := time.Now().UnixNano()
	diff := (int64(sec) - int64(startTime)) / 300000
	pointLocation := int(diff)
	// if pointLocation > 800 {
	// 	pointLocation = 800
	// }
	pointLocationV := 0
	color := "#000000"
	switch {
	case process == 1:
		pointLocationV = 60
		color = "#cc6666"
	default:
		pointLocationV = 180
		color = "#66cc66"
	}
	fmt.Println("pointLocation: ", pointLocation)
	osvg.Rect(pointLocation, pointLocationV, 3, 5, "fill:"+color+";stroke:none;")
	time.Sleep(1 * time.Millisecond)
}

func visualize(rw http.ResponseWriter, req *http.Request) {
	startTime = time.Now().UnixNano()
	fmt.Println("Request to /visualize")
	rw.Header().Set("Content-Type", "image/svg+xml")
	outputSVG := svg.New(rw)
	outputSVG.Start(width, height)
	outputSVG.Rect(10, 10, 78, 100, "fill:#eeeeee;stroke:none;")
	outputSVG.Text(20, 30, "Process 1 Timeline", "text-anchor:start;font-size:12px;fill:#333333")
	outputSVG.Rect(10, 130, 780, 100, "fill:#eeeeee;stroke:none;")
	outputSVG.Text(20, 150, "Process 2 Timeline", "text-anchor:start;font-size:12px;fill:#333333")
	// outputSVG.Text(650, 360, "Run without goroutines", "text-anchor:start;font-size:12px;fill:#333333")
	outputSVG.Rect(10, 340, 780, 100, "fill:#eeeeee;stroke:none;")
	outputSVG.Text(650, 360, "Run without goroutines", "text-anchor:start;font-size:12px;fill:#333333")

	for i := 0; i < 801; i++ {
		timeText := strconv.FormatInt(int64(i), 10)
		if i%100 == 0 {
			outputSVG.Text(i, 380, timeText, "text-anchor:middle;font-size:10px;fill:#000000")
		} else if i%4 == 0 {
			outputSVG.Circle(i, 377, 1, "fill:#cccccc;stroke:none")
		}

		if i%10 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#ddddd")
		} else if i%50 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#cccccc")
		}
	}
	for i := 0; i < 100; i++ {
		go drawPoint(outputSVG, i, 1)
		drawPoint(outputSVG, i, 2)
	}

	outputSVG.End()

}
func main() {
	http.Handle("/visualize", http.HandlerFunc(visualize))

	err := http.ListenAndServe("localhost:1900", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
