package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const PORT = ":7878"

func main() {
	httpStart()
}


func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func response(w http.ResponseWriter, r *http.Request) {
	currentPath := getCurrentDirectory()
	file := currentPath + r.URL.Path
	fmt.Println(datetime() + " path: "+file)
	if isFile(file) {
		http.ServeFile(w, r, file)
	} else {
		if !isDir(file){
			http.NotFound(w,r)
			return
		}
		//dir
		fileNames := listDir(file)
		for _, filename := range fileNames {
			if isDir(file + filename){
				fmt.Fprintf(w, "<a href=\".%s/\">.%s</a></br>", string(filename), string(filename))
			}else {
				fmt.Fprintf(w, "<a href=\".%s\">.%s</a></br>", string(filename), string(filename))
			}
		}
	}
	//w.Write([]byte("this dir : " + getCurrentDirectory()))
}

func httpStart(){
	ips := getLocalIP()
	mux := http.NewServeMux()
	mux.HandleFunc("/", response)
	//check args
	port := PORT
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}
	fmt.Println(datetime() + "\n^_^, start http service, listening address: ")
	for _,v := range ips {
		fmt.Println("----> "+ v + port)
	}
	http.ListenAndServe(port, mux)

}

func datetime() string  {
	return time.Now().Format(time.RFC3339)
}

func getLocalIP() []string  {
	address,err := net.InterfaceAddrs()
	var ips []string
	if err != nil {
		os.Exit(1)
	}
	for _,v := range address {
		// 检查ip地址判断是否回环地址 && !ipnet.IP.IsLoopback()
		if ipnet, ok := v.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips,ipnet.IP.String())
			}
		}
	}
	return ips
}


//只遍历一层
func listDir(dirPath string) []string {
	dirInfo, err := ioutil.ReadDir(dirPath)
	var files []string
	if err != nil {
		fmt.Println(err)
	}
	separator := string(os.PathSeparator)
	for _,fileInfo := range dirInfo {
		//println(fileInfo.Name(),string(os.PathSeparator),string(os.PathListSeparator))
		files = append(files, separator + fileInfo.Name())
	}
	return files
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.Mode().IsRegular()
}

