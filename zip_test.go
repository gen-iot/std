package std

import (
	"archive/zip"
	"fmt"
	"os"
	"testing"
)

func TestUnZip(t *testing.T) {
	err := os.Mkdir("zip_output", 0755)
	AssertError(err, "mkdir zip_output")
	defer func() {
		_ = os.RemoveAll("zip_output")
	}()
	err = UnZip("sample.zip", "zip_output")
	AssertError(err, "unzip")
}

func TestZip(t *testing.T) {
	wd, err := os.Getwd()
	AssertError(err, "getwd")
	fmt.Println("wd=", wd)
	file, err := os.OpenFile("output.zip", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	AssertError(err, "create output")
	defer CloseIgnoreErr(file)
	err = Zip(file, []*FileGroup{
		NewFileGroup(".", "go.mod").SetCustomDirName("c").AddFileItemWithAlias("go.sum", "goSum"),
		NewFileGroup(".idea", "misc.xml", "modules.xml", "workspace.xml").SetCustomDirName(".").
			AddChildFileGroup(NewFileGroup("inspectionProfiles", "Project_Default.xml").SetCustomDirName("a")),
	})
	AssertError(err, "zip failed")
}

func TestZip2(t *testing.T) {
	file, err := os.OpenFile("output.zip", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	AssertError(err, "create output")

	writer := zip.NewWriter(file)
	defer func() {
		CloseIgnoreErr(writer)
		CloseIgnoreErr(file)
	}()
	_, err = writer.Create("a/")
	AssertError(err, "create a/")
	_, err = writer.Create("a/b/")
	AssertError(err, "create a/b/")
	_, err = writer.Create("a/")
	AssertError(err, "create a/ again")
}
