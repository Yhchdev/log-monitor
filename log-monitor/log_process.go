package main

/*面向接口编程3步：
1、抽象出上层接口
	1.1 从文件中读/从标准输出读  读接口


2. 结合具体的业务逻辑，做出具体的接口实现

3.将接口的实现，以接口的形式注入到接口调用方

*/
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogProcess struct {
	rchan chan []byte
	wchan chan string
	read  Reader
	write Writer
}

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Writer(wc chan string)
}

type readFormFile struct {
	path string
}

type writeToInfluxDB struct {
	influxDBDsn string
}

func (r *readFormFile) Read(rc chan []byte) {

	// 打开文件
	f, err := os.Open("./access.log")
	if err != nil {
		panic(fmt.Sprintf("open file fiead:%s", err.Error()))
	}

	// 移动到文件末尾，就的日志，没必要监听
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		data, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				panic(fmt.Sprintf("read file err:%s", err.Error()))
			}
		}
		rc <- data[:len(data)-1]
	}
}

func (w *writeToInfluxDB) Writer(wc chan string) {
	for v := range wc {
		fmt.Println(v)
	}

}

func (l *LogProcess) process() {
	// 处理消息模块

	for v := range l.rchan {
		l.wchan <- strings.ToUpper(string(v))
	}

}

func main() {
	rfile := &readFormFile{path: "/tmp/access.log"}
	wdb := &writeToInfluxDB{influxDBDsn: "username&password.."}

	lp := &LogProcess{
		rchan: make(chan []byte),
		wchan: make(chan string),
		read:  rfile,
		write: wdb,
	}

	go lp.read.Read(lp.rchan)
	go lp.process()
	go lp.write.Writer(lp.wchan)

	time.Sleep(3600 * time.Second)

}
