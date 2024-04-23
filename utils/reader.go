package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/hpcloud/tail"
)

type FileReader struct {
	lineChan chan string
	done     chan bool
	file     *os.File
	*sync.Mutex
}

func NewFileReader() *FileReader {
	return &FileReader{
		lineChan: nil,
		file:     nil,
	}
}

func (f *FileReader) tail() {

	f.printOffset()
	t, err := tail.TailFile(args.File, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: io.SeekEnd,
		},
		ReOpen: true,
		Follow: true,
		Logger: tail.DiscardingLogger,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	for line := range t.Lines {
		f.lineChan <- line.Text
	}
}

func getLine(f *os.File, offset *int64) string {
	b := make([]byte, 0)

	for {
		*offset -= int64(1)

		if *offset <= 0 {
			break
		}

		t := make([]byte, 1)
		if _, err := f.ReadAt(t, *offset); err != nil {
			panic(err)
		}

		if t[0] == 10 {
			break
		}
		b = append(b, t[0])
	}

	line := ""
	for i := len(b) - 1; i >= 0; i-- {
		line += string(b[i])
	}

	return line
}

func (fr *FileReader) printOffset() {

	f, _ := os.Open(args.File)

	st, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := st.Size()
	offset := fileSize
	for i := 0; i < args.LinesOffset; i++ {
		line := getLine(f, &offset)
		if line == "" {
			i--
			continue
		}
		fr.lineChan <- line
	}
}

func (f *FileReader) Paint() {
	f.lineChan = make(chan string)
	f.done = make(chan bool)

	go f.controller()
	if args.IsFile {
		if args.Tail {
			f.tail()
		} else {
			// TODO: Open Custom File Explorer
		}
	} else {
		f.printStdin()
	}

	f.done <- true
}

// INFO: Print Log Line with Color based on a level
func (f FileReader) paintLine(level string, line string) string {
	color := config.colors[level]

	if color != "" {
		def := config.defaultColor
		return fmt.Sprintf("%s%s%s", color, line, def)
	}

	return line
}

func (f *FileReader) printStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		f.lineChan <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}

func (fr *FileReader) linesHadler(line string, level *string, r *regexp.Regexp) {

	if args.Grep != "" {
		r2 := regexp.MustCompile(args.Grep)
		if !r2.MatchString(line) {
			return
		}
	}

	if args.IgnoreGrep != "" {
		r2 := regexp.MustCompile(args.IgnoreGrep)
		if r2.MatchString(line) {
			return
		}
	}

	str := r.FindStringSubmatch(line)
	if len(str) > 1 {
		*level = strings.ToLower(str[1])
	}

	fmt.Println(fr.paintLine(*level, line))
}

func (f *FileReader) controller() {
	level := ""
	r := regexp.MustCompile(args.Pattern)

	for {
		select {
		case line := <-f.lineChan:
			f.linesHadler(line, &level, r)

		case <-f.done:
			if !args.IsFile {
				close(f.lineChan)
				f.file.Close()
				os.Remove(args.File)
			}
			return
		}
	}
}
