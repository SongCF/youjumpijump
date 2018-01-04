package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	jump "github.com/songcf/youjumpijump"
)

var similar *jump.Similar

func screenshot(filename string) image.Image {
	_, err := exec.Command("/system/bin/screencap", "-p", filename).Output()
	if err != nil {
		panic("screenshot failed")
	}

	inFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	src, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}
	inFile.Close()
	return src
}

func main() {
	defer func() {
		jump.Debugger()
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("the program has crashed, press any key to exit")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var inputRatio float64
	var similarSleep int
	var err error
	{
		fmt.Print("input jump ratio (recommend 2.25):")
		_, err = fmt.Scanln(&inputRatio)
		if err != nil {
			log.Printf("input is empty, will use 2.25 as default ratio")
			inputRatio = 2.25
		}
		fmt.Print("input similarSleep (recommend 300):")
		_, err = fmt.Scanln(&similarSleep)
		if err != nil {
			log.Printf("input is empty, will use 300 as default ratio")
			similarSleep = 300
		}
	}

	similar = jump.NewSimilar(inputRatio)

	for {
		src := screenshot("jump.png")

		start, end := jump.Find(src)
		if start == nil {
			log.Print("can't find the starting point，please export the debugger directory")
			break
		} else if end == nil {
			log.Print("can't find the end point，please export the debugger directory")
			break
		}

		scale := float64(src.Bounds().Max.X) / 720
		nowDistance := jump.Distance(start, end)
		// similarDistance, nowRatio := similar.Find(nowDistance)
		similarDistance, nowRatio := 0.0, inputRatio

		{
			if nowDistance < 200 {
				nowRatio = -0.0134*nowDistance + 4.952
			} else if nowDistance > 300 {
				nowRatio = -0.00075*nowDistance + 2.475
			} else {
				nowRatio = 2.25
			}

		}

		log.Printf("from:%v to:%v distance:%.2f similar:%.2f ratio:%v press:%.2fms ", start, end, nowDistance, similarDistance, nowRatio, nowDistance*nowRatio)

		_, err = exec.Command("/system/bin/sh", "/system/bin/input", "swipe",
			strconv.FormatFloat(float64(start[0])*scale, 'f', 0, 32),
			strconv.FormatFloat(float64(start[1])*scale, 'f', 0, 32),
			strconv.FormatFloat(float64(end[0])*scale, 'f', 0, 32),
			strconv.FormatFloat(float64(end[1])*scale, 'f', 0, 32),
			strconv.Itoa(int(nowDistance*nowRatio))).Output()
		if err != nil {
			panic("touch failed")
		}
		jump.Debugger2(nowRatio, nowDistance)



		time.Sleep(time.Millisecond * 1000)
	}
}
