package main

import "os"

type FileLog struct {
	file *os.File
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func OpenLog(appName string) FileLog {

	if _, err := os.Stat("log"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		err := os.Mkdir("log", os.ModePerm)
		check(err)
	}

	f, err := os.Create("log/" + appName + ".log")
	check(err)

	return FileLog{file: f}
}

func (fileLog *FileLog) LogOnce(content string) {
	_, err := fileLog.file.WriteString(content + "\n")
	check(err)
	err = fileLog.file.Close()
	check(err)
}
