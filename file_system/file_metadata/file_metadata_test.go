package file_metadata

import (
	"fmt"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	fmdt := NewFileMetaData()
	err := fmdt.Upload("/foo/path/somedir/file.txt")
	if err != nil {
		fmt.Println("error uploading")
	}
	err = fmdt.Upload("/foo/path/file.txt")
	if err != nil {
		fmt.Println("error uploading")
	}
	err = fmdt.Upload("/foo/path/somedir/someotherdir/file2.txt")
	if err != nil {
		fmt.Println("error uploading")
	}

	s1, _ := fmdt.Ls("/")
	fmt.Println(s1)
	if !strings.Contains(s1, "foo") {
		t.Fatalf("incorrect ls result / %s", s1)
	}

	s2, _ := fmdt.Ls("/foo/")
	fmt.Println(s2)
	if !strings.Contains(s2, "path") {
		t.Fatalf("incorrect ls result foo %s", s2)
	}

	s3, _ := fmdt.Ls("/foo/path/somedir/someotherdir/")
	fmt.Println(s3)
	if !strings.Contains(s3, "file2.txt") {
		t.Fatalf("incorrect ls result somedir %s", s3)
	}
}
