package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy копирует limit байт из файла по пути fromPath в файл по пути toPath с отступом offset.
func Copy(fromPath, toPath string, offset, limit int64) error {
	fInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if fInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > fInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || offset+limit > fInfo.Size() {
		limit = fInfo.Size() - offset
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	if offset != 0 {
		_, err = fileFrom.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fileFrom)

	fileTo, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0o777)
	if err != nil {
		return err
	}

	defer fileTo.Close()

	_, err = io.CopyN(fileTo, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
