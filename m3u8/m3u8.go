package m3u8

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type M3U8 struct {
	outputFile string
	fileChan   chan string
	m3u8Chan   chan string
	tsChan     chan map[string][]string
	endChan    chan string
	errChan    chan string
}

func NewM3U8(outputFile string) *M3U8 {
	if outputFile == "" {
		outputFile = "m3u8"
	}
	return &M3U8{
		outputFile: outputFile,
		fileChan:   make(chan string),
		m3u8Chan:   make(chan string),
		tsChan:     make(chan map[string][]string),
		endChan:    make(chan string),
		errChan:    make(chan string),
	}
}

func (this *M3U8) getM3u8() {
	for {
		url := <-this.fileChan
		var pathinfo = strings.Split(url, "/")
		var filename = fmt.Sprintf("%s/%s", this.outputFile, pathinfo[len(pathinfo)-1])
		res, err := http.Get(url)
		if err != nil {
			this.errChan <- fmt.Sprintf("download error url:%s,err:%s", url, err.Error())
			break
		}
		defer res.Body.Close()
		f, err := os.Create(filename)
		if err != nil {
			this.errChan <- fmt.Sprintf("download create Filename error url:%s,err:%s", url, err.Error())
			break
		}
		defer f.Close()
		io.Copy(f, res.Body)
		this.m3u8Chan <- filename
	}
	close(this.fileChan)
}
func (this *M3U8) readTs() {
	for {
		path := <-this.m3u8Chan
		var tsUrl = []string{}
		tsData := make(map[string][]string, 0)
		handler, err := os.Open(path)
		if err != nil {
			this.errChan <- fmt.Sprintf("readTs open Filename error path:%s,err:%s", path, err.Error())
			break
		}
		defer handler.Close()
		buffer := bufio.NewReader(handler)
		for {
			line, _, err := buffer.ReadLine()
			if err == io.EOF {
				break
			}
			if strings.HasPrefix(string(line), "http") {
				tsUrl = append(tsUrl, string(line))
			}
		}
		tsData["path"] = []string{path}
		tsData["urls"] = tsUrl
		this.tsChan <- tsData
	}
}
func (this *M3U8) makeMp4() {
	for {
		tsMap := <-this.tsChan
		filename := tsMap["path"][0]
		outfile := string(strings.Split(filename, ".")[0]) + ".mp4"

		fd, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			this.errChan <- fmt.Sprintf("makeMp4 open Filename error filename:%s,err:%s", filename, err.Error())
		}
		defer fd.Close()

		for _, v := range tsMap["urls"] {
			fmt.Println("download", v)
			rs, err := http.Get(v)
			//defer rs.Body.Close()
			if err != nil {
				this.errChan <- fmt.Sprintf("makeMp4 get ts url error filename:%s,url:%s,err:%s", filename, v, err.Error())
				break
			}
			io.Copy(fd, rs.Body)
		}
		this.endChan <- ""
	}
}
func (this *M3U8) Download(urls []string) {

	go func() {
		for _, value := range urls {
			this.fileChan <- value
		}
	}()
	num := len(urls)
	for i := 0; i < num; i++ {
		go this.getM3u8()
	}
	for i := 0; i < num; i++ {
		go this.readTs()
	}
	for i := 0; i < num; i++ {
		go this.makeMp4()
	}
	exit := false
	i := 0
	for {
		if exit {
			break
		}
		select {
		case _, ok := <-this.endChan:
			if ok {
				i++
			}
			if i >= num {
				exit = true
				break
			}
		case v, ok := <-this.errChan:
			if ok {
				exit = true
				break
			}
			fmt.Println(v, "err---")
		}
	}
}
