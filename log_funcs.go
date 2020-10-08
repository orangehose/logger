package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

var maxLength int = 3
var logFileName string = "log_file"

func writeToLocalLog(payload string) {
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("Unable to open file: ", err)
		os.Exit(1)
	}

	fileLength, err := lineCounter(file)
	if err != nil {
		fmt.Println("Unable to read file len: ", err)
		os.Exit(1)
	}

	// Cyclic overwriting logs on disk
	if fileLength < maxLength {
		appendToFile(file, payload)
	} else {
		err = rewriteLogFile(file, payload)
		if err != nil {
			fmt.Println("Unable to create new file: ", err)
			os.Exit(1)
		}
	}

	defer file.Close()
}

func appendToFile(extFile *os.File, payload string) {
	writer := bufio.NewWriter(extFile)
	writer.WriteString(payload)
	writer.WriteString("\n")
	writer.Flush()
}

func rewriteLogFile(extFile *os.File, payload string) error {
	tempFileName := "dest_file"

	// Go to beginning of file
	extFile.Seek(0, 0)
	reader := bufio.NewReader(extFile)

	line, err := reader.ReadSlice('\n')
	if err != nil {
		return err
	}

	firstLineLen := int64(len(line))

	// Creating a temporary log file to replace the current one
	destFile, err := os.OpenFile(tempFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Unable to open file: ", err)
		os.Exit(1)
	}

	// Getting rid of the oldest recording
	_, err = extFile.Seek(firstLineLen, io.SeekStart)
	if err != nil {
		panic(err)
	}

	// Overwriting logs on disk
	_, err = io.Copy(destFile, extFile)
	if err != nil {
		panic(err)
	}

	appendToFile(destFile, payload)

	if err := os.Remove(logFileName); err != nil {
		panic(err)
	}

	if err := os.Rename(tempFileName, logFileName); err != nil {
		panic(err)
	}

	defer destFile.Close()

	return nil
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)

		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
