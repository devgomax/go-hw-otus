package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testdata = "testdata/"

func getMD5Sum(f *os.File) (string, error) {
	sum := md5.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", sum.Sum(nil)), nil
}

func TestCopy(t *testing.T) {
	tests := []struct {
		refTo  string
		offset int64
		limit  int64
	}{
		{refTo: testdata + "out_offset0_limit0.txt", offset: 0, limit: 0},
		{refTo: testdata + "out_offset0_limit10.txt", offset: 0, limit: 10},
		{refTo: testdata + "out_offset0_limit1000.txt", offset: 0, limit: 1000},
		{refTo: testdata + "out_offset0_limit10000.txt", offset: 0, limit: 10000},
		{refTo: testdata + "out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{refTo: testdata + "out_offset6000_limit1000.txt", offset: 6000, limit: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.refTo, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test_result")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			err = Copy(testdata+"input.txt", tempFile.Name(), tt.offset, tt.limit)
			require.NoError(t, err)

			refFile, err := os.OpenFile(tt.refTo, os.O_RDONLY, 0644)
			require.NoError(t, err)
			defer refFile.Close()

			refHashSum, err := getMD5Sum(refFile)
			require.NoError(t, err)

			givenHashSum, err := getMD5Sum(tempFile)
			require.NoError(t, err)

			require.Equal(t, refHashSum, givenHashSum)
		})
	}
}

func TestCopyWithOversizeOffset(t *testing.T) {
	t.Run("offset_gt_len_file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_result")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		fInfo, err := os.Stat(testdata + "input.txt")
		require.NoError(t, err)

		err = Copy(testdata+"input.txt", tempFile.Name(), fInfo.Size()+1, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}

func TestCopyDevNull(t *testing.T) {
	t.Run("unknown_len_file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_result")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		err = Copy("/dev/null", tempFile.Name(), 1, 1)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}

func TestCopyDevUrandom(t *testing.T) {
	t.Run("unknown_len_file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_result")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		err = Copy("/dev/urandom", tempFile.Name(), 1, 1)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}

func TestCopyFromInvalidPath(t *testing.T) {
	t.Run("copy_from_invalid_path", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_result")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		err = Copy(testdata+"dmalsdkasdkasfdmd.txt", tempFile.Name(), 1, 1)
		require.Error(t, err)
	})
}
