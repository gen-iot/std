package std

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileGroup struct {
	rawDir          string
	dirName         string
	childFiles      []string
	childFileGroups []*FileGroup
}

func (this *FileGroup) GetDirName() string {
	return this.dirName
}

func (this *FileGroup) GetChildFiles() []string {
	return this.childFiles
}

func (this *FileGroup) SetCustomDirName(dirName string) *FileGroup {
	this.dirName = dirName + separatorStr
	return this
}

const separatorStr = string(filepath.Separator)

func NewFileGroup(dir string, childFiles ...string) *FileGroup {
	out := &FileGroup{
		rawDir:     dir,
		dirName:    "",
		childFiles: nil,
	}
	if len(dir) != 0 {
		out.dirName = fmt.Sprintf("%s%c", filepath.Base(dir), filepath.Separator)
	}
	if len(childFiles) == 0 {
		return out
	}
	out.childFiles = make([]string, 0, len(childFiles))
	for idx := range childFiles {
		out.childFiles = append(out.childFiles, filepath.Base(childFiles[idx]))
	}
	return out
}

func (this *FileGroup) addChildFileGroup(fg FileGroup) {
	fg.rawDir = filepath.Join(this.rawDir, filepath.Base(fg.rawDir))
	fg.dirName = filepath.Join(this.dirName, fg.dirName) + separatorStr
	this.childFileGroups = append(this.childFileGroups, &fg)
}

func (this *FileGroup) AddChildFileGroup(fgs ...*FileGroup) *FileGroup {
	if this.childFileGroups == nil {
		this.childFileGroups = make([]*FileGroup, 0, len(fgs))
	}
	for _, fg := range fgs {
		this.addChildFileGroup(*fg)
	}
	return this
}

func zipFilesDetails(zipWriter *zip.Writer, fgs []*FileGroup) error {
	curDir := "."
	for _, fg := range fgs {
		if len(fg.dirName) != 0 {
			_, err := zipWriter.Create(fg.dirName)
			if err != nil {
				return err
			}
			curDir = fg.dirName
		} else {
			curDir = "."
		}
		for _, fName := range fg.childFiles {
			writer, err := zipWriter.Create(filepath.Join(curDir, fName))
			if err != nil {
				return err
			}
			file, err := os.Open(filepath.Join(fg.rawDir, fName))
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, file)
			if err != nil {
				_ = file.Close()
				return err
			}
			_ = file.Close()
		}
		if len(fg.childFileGroups) == 0 {
			continue
		}
		err := zipFilesDetails(zipWriter, fg.childFileGroups)
		if err != nil {
			return err
		}
	}
	return nil
}

func Zip(output io.Writer, fgs []*FileGroup) error {
	zipWriter := zip.NewWriter(output)
	defer CloseIgnoreErr(zipWriter) // implicit flush to file
	return zipFilesDetails(zipWriter, fgs)
}

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
