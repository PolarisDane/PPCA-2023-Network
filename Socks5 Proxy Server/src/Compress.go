package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"

	"github.com/andybalholm/brotli"
)

func Decompressbr(data []byte) []byte {
	brReader := brotli.NewReader(bytes.NewReader(data))
	ret, err := ioutil.ReadAll(brReader)
	if err != nil {
		log.Panic(err)
	}
	return ret
}

func Compressbr(data []byte) []byte {
	buffer := bytes.Buffer{}
	brWriter := brotli.NewWriter(&buffer)
	brWriter.Write(data)
	brWriter.Close()
	return buffer.Bytes()
}

func Decompressgzip(data []byte) []byte {
	buffer := bytes.Buffer{}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
		return nil
	}
	io.Copy(&buffer, gzipReader)
	gzipReader.Close()
	return buffer.Bytes()
}

func Compressgzip(data []byte) []byte {
	buffer := bytes.Buffer{}
	gzipWriter := gzip.NewWriter(&buffer)
	gzipWriter.Write(data)
	gzipWriter.Close()
	return buffer.Bytes()
}

func Decompressflate(data []byte) []byte {
	flateReader := flate.NewReader(bytes.NewReader(data))
	ret, err := ioutil.ReadAll(flateReader)
	if err != nil {
		log.Panic(err)
	}
	return ret
}

func Compressflate(data []byte) []byte {
	buffer := bytes.Buffer{}
	flateWriter, _ := flate.NewWriter(&buffer, flate.DefaultCompression)
	flateWriter.Write(data)
	flateWriter.Close()
	return buffer.Bytes()
}
