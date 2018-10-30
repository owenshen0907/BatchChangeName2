// BatchChangeName2 project main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	//"github.com/axgle/mahonia"

	"github.com/larspensjo/config"
	//"github.com/qiniu/iconv"
)

var (
	configFile = flag.String("configfile", "Sender_config.ini", "General configuration file")
	Version    = "BatchChangeName V1.0.20180814 "
	Auther     = "Owen Shen"
	tmp        []string
)

func getArgs() {
	version := flag.Bool("v", false, "version")
	flag.Parse()
	if *version {
		fmt.Println("Version：", Version)
		fmt.Println("Auther:", Auther)
		return
	}
}
func main() {

	if len(os.Args) > 1 {
		getArgs()
	} else {
		body()
	}

}

func body() {
	var TOPIC = readconfigfile()
	fmt.Println("Hello World!")
	var pwd, filename string

	filename = TOPIC["filename"]

	pwd = getCurrentDirectory()
	filename = pwd + "\\" + filename
	fmt.Println("pwd:" + pwd)
	fmt.Println(filename)
	ReadLine(filename, Print)
	fmt.Println(tmp)
	//	cd, err := iconv.Open("gbk", "utf-8")
	//	erro(err)
	//	defer cd.Close()

	for _, v := range tmp {
		t := strings.Split(v, ",")
		if len(t) == 2 && t[1] != "" {
			fmt.Println(t[0] + "-->" + t[1])
			fmt.Println(pwd + "\\" + TOPIC["source"] + "\\" + t[0])
			//CopyFile(pwd+"\\"+TOPIC["destination"]+"\\"+t[1], pwd+"\\"+TOPIC["source"]+"\\"+t[0])
			if checkFileIsExist(pwd + "\\" + TOPIC["source"] + "\\" + t[0]) {
				fmt.Println("文件存在：" + t[0])

				CopyFile(pwd+"\\"+TOPIC["destination"]+"\\"+t[1], pwd+"\\"+TOPIC["source"]+"\\"+t[0])
			} else {
				fmt.Println("文件不存在，无法复制" + t[0])
			}
		}

	}
	fmt.Println("按任意键退出哦！！！")
	fmt.Scanln(&filename)
}

func ReadLine(fileName string, handler func(string)) error {
	//conv := mahonia.NewEncoder("utf-8")
	//cd, err := iconv.Open("gbk", "utf-8")
	//erro(err)
	//defer cd.Close()

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		//line = cd.ConvString(line)
		//a := conv.ConvertString(line)
		//fmt.Println(a)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func Print(line string) {
	fmt.Println(line)
	tmp = append(tmp, line)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "\\", -1)
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func readconfigfile() (TOPIC map[string]string) {
	TOPIC = make(map[string]string)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("topicArr") {
		section, err := cfg.SectionOptions("topicArr")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("topicArr", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END
	return TOPIC
}
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
func erro(err error) {
	if err != nil {
		panic(err)
	}
}
