package std

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnZip(zipFilePath, output string) error {
	zrc, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	for _, f := range zrc.File {
		info := f.FileInfo()
		outputTarget := filepath.Join(output, f.Name)
		if info.IsDir() {
			err := os.MkdirAll(outputTarget, info.Mode())
			if err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(outputTarget, os.O_CREATE|os.O_TRUNC|os.O_RDWR, info.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(file, rc)
		_ = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
