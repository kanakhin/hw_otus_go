package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !srcInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := srcInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, err := src.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.CopyN(dst, src, bytesToCopy); err != nil {
		return err
	}

	return nil
}
