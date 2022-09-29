package file_metadata

import (
	"fmt"
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

	s1 := fmdt.Ls("/")
	fmt.Println(s1)
	if s1 != "foo" {
		t.Fatalf("incorrect ls result / %s", s1)
	}

	s2 := fmdt.Ls("/foo/")
	fmt.Println(s2)
	if s2 != "path" {
		t.Fatalf("incorrect ls result foo %s", s2)
	}

	s3 := fmdt.Ls("/foo/path/somedir/someotherdir/")
	fmt.Println(s3)
	if s3 != "file2.txt" {
		t.Fatalf("incorrect ls result somedir %s", s3)
	}
}
