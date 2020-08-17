package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//创建tar存档
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	files := []struct {
		Name,
		Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher name:\nGeorge\nGeoffrey\nGonza"},
		{"todo.txt", "Get animal handling license."},
	}

	//写tar存档
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	//读tar存档
	r := bytes.NewReader(buf.Bytes())
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			//tar 归档结束标记
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s: \n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatalln(err)
		}
		fmt.Println()
	}
	file, err := os.OpenFile("test.tar", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	buf.WriteTo(file)

	//创建zip
	buf = new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, file := range files {
		f, err := zw.Create(file.Name)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatalln(err)
		}
	}
	err = zw.Close()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println()
	fmt.Println()
	//读zip
	r = bytes.NewReader(buf.Bytes())
	zr, err := zip.NewReader(r, int64(buf.Len()))
	if err != nil {
		log.Fatalln(err)
	}
	for _, f := range zr.File {
		fmt.Printf("Contents of %s: \n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatalln(err)
		}
		_, err = io.Copy(os.Stdout, rc)
		if err != nil {
			log.Fatalln(err)
		}
		rc.Close()
		fmt.Println()
	}

}
