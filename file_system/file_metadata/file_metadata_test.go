package file_metadata

import (
	"fmt"
	"testing"
)

func TestUpload(t *testing.T) {
	fmdt := newFileMetaData()
	fmdt.upload("/foo/path/somedir/file.txt")
	fmdt.upload("/foo/path/file.txt")
	fmdt.upload("/foo/path/somedir/file2.txt")

	fmt.Println(fmdt.ls("/"))
	if fmdt.ls("/") != "foo" {
		t.Fatalf("incorrect ls result /")
	}

	fmt.Println(fmdt.ls("/foo/"))
	if fmdt.ls("/foo/") != "path" {
		t.Fatalf("incorrect ls result foo")
	}

	fmt.Println(fmdt.ls("/foo/path/somedir/"))
	if fmdt.ls("/foo/path/somedir/") != "file.txt file2.txt" {
		t.Fatalf("incorrect ls result somedir")
	}
}
