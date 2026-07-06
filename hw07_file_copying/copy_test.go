package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func readFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %q: %v", path, err)
	}

	return data
}

func TestCopy(t *testing.T) {
	inputPath := filepath.Join("testdata", "input.txt")

	tests := []struct {
		name      string
		offset    int64
		limit     int64
		wantPath  string
		wantErr   error
		wantErrIs bool
	}{
		{
			name:     "full file",
			offset:   0,
			limit:    0,
			wantPath: filepath.Join("testdata", "out_offset0_limit0.txt"),
		},
		{
			name:     "limit 10",
			offset:   0,
			limit:    10,
			wantPath: filepath.Join("testdata", "out_offset0_limit10.txt"),
		},
		{
			name:     "limit 1000",
			offset:   0,
			limit:    1000,
			wantPath: filepath.Join("testdata", "out_offset0_limit1000.txt"),
		},
		{
			name:     "limit greater than file size",
			offset:   0,
			limit:    10000,
			wantPath: filepath.Join("testdata", "out_offset0_limit10000.txt"),
		},
		{
			name:     "offset 100 limit 1000",
			offset:   100,
			limit:    1000,
			wantPath: filepath.Join("testdata", "out_offset100_limit1000.txt"),
		},
		{
			name:     "offset 6000 limit 1000",
			offset:   6000,
			limit:    1000,
			wantPath: filepath.Join("testdata", "out_offset6000_limit1000.txt"),
		},
		{
			name:      "offset exceeds file size",
			offset:    10000,
			limit:     10,
			wantErr:   ErrOffsetExceedsFileSize,
			wantErrIs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outFile, err := os.CreateTemp(t.TempDir(), "copy-out-*")
			if err != nil {
				t.Fatalf("create temp file: %v", err)
			}
			outPath := outFile.Name()
			outFile.Close()

			err = Copy(inputPath, outPath, tt.offset, tt.limit)
			if tt.wantErrIs {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := readFile(t, outPath)
			want := readFile(t, tt.wantPath)
			if !bytes.Equal(got, want) {
				t.Fatalf("copied content mismatch: got %d bytes, want %d bytes", len(got), len(want))
			}
		})
	}
}

func TestCopyUnsupportedFile(t *testing.T) {
	fifoPath := filepath.Join(t.TempDir(), "fifo")
	if err := syscall.Mkfifo(fifoPath, 0o666); err != nil {
		t.Fatalf("create fifo: %v", err)
	}

	dstFile, err := os.CreateTemp(t.TempDir(), "copy-dst-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	dstPath := dstFile.Name()
	dstFile.Close()

	err = Copy(fifoPath, dstPath, 0, 0)
	if err == nil {
		t.Fatal("expected error for unsupported file, got nil")
	}
	if !errors.Is(err, ErrUnsupportedFile) {
		t.Fatalf("expected error %v, got %v", ErrUnsupportedFile, err)
	}
}
