package main

import (
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"io"
	"os/exec"
	"path/filepath"
	"log"
	"flag"
	"io/ioutil"
	"errors"
)

type MediaInfo map[string]map[string]string

var recursive bool
var output string

func init() {
	const (
		recursiveDefault = false
		recursiveInfo = "Recurse files in directory"
		outputDefault = ""
		outputInfo = "Output file"
	)
	flag.BoolVar(&recursive, "recursive", recursiveDefault, recursiveInfo)
	flag.BoolVar(&recursive, "r", recursiveDefault, recursiveInfo + " (shorthand)")
	flag.StringVar(&output, "output", outputDefault, outputInfo)
	flag.StringVar(&output, "o", outputDefault, outputInfo + " (shorthand)")
	flag.Parse()
}

func checkFatal(e error) {
	if e != nil {
		log.Fatalf("Fatal: %s", e)
	}
}

func check(e error) {
	if e != nil {
		log.Printf("Error: %s", e)
	}
}

func main() {
	args := flag.Args()

	if len(args) != 1 {
		checkFatal(errors.New("Specify a single directory path or file as argument"))
	}

	dir, err := os.Stat(args[0])
	checkFatal(err)

	var files []string
	if !dir.IsDir() {
		files = args[:1]
	}else if recursive {
		files = walkFiles(args[0])
	} else {
		files = listFiles(args[0])
	}

	var outfile *os.File

	if output != "" {
		outfile, err = os.Create(output)
		checkFatal(err)
	} else {
		outfile = os.Stdout
	}

	w := bufio.NewWriter(outfile)
	
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetIndent("", "  ")

	var mediaDataSlice []MediaInfo

	for _, file := range files {
		data := getFileMediaInfo(file)
		mediaDataSlice = append(mediaDataSlice, data)
	}

	jsonEncoder.Encode(mediaDataSlice)

	
	w.Flush()
	outfile.Close()
}

func listFiles(path string) []string{
		var files []string
		fileinfos, err := ioutil.ReadDir(path)
    checkFatal(err)
    for _, f := range fileinfos {
			if !f.IsDir() {
				files = append(files, f.Name())
			}
		}
		return files
}

func walkFiles(path string) []string{
	var files []string
	err := filepath.Walk(path, func (fpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error: %s", err)
		}
		if !info.IsDir() {
			files = append(files, fpath)
		}
		return nil
	})
	checkFatal(err)
	return files
}

func getFileMediaInfo(filepath string) MediaInfo {
	cmd := exec.Command("mediainfo", filepath)
	stdout, err := cmd.StdoutPipe()
	check(err)
	err = cmd.Start()
	check(err)
	data := getMediaInfo(stdout)
	err = cmd.Wait()
	check(err)
	return data
}

func getMediaInfo(reader io.Reader) MediaInfo {
	data := make(map[string]map[string]string)
	input := bufio.NewScanner(reader)
	var header string
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if line == "" {
			continue
		}
		colonidx := strings.Index(line, ":")
		if colonidx == -1 {
			header = strings.TrimSpace(line)
			data[header] = make(map[string]string)
		} else {
			key := strings.TrimSpace(line[:colonidx])
			value := strings.TrimSpace(line[colonidx+1:])
			data[header][key] = value
		}
	}
	return MediaInfo(data)
}