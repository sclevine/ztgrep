package ztgrep_test

import (
	"strings"
	"testing"

	"github.com/sclevine/ztgrep"
)

func TestZTgrep(t *testing.T) {
	zt, err := ztgrep.New("test")
	if err != nil {
		t.Fatal(err)
	}
	tt := []string{
		"testdata/test-l2.tar.gz:test-l1.tar",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test-l1.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test.tar.bz2",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile1",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile2",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test.tar.bz2:testfile3",
		"testdata/test-l2.tar.gz:test.tar.gz",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:test.tar.xz",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile1",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile2",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test.tar.xz:testfile3",
		"testdata/test-l2.tar.gz:test.tar.zst",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile1",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile2",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:test.tar.zst:testfile3",
		"testdata/test-l2.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:testfile1",
		"testdata/test-l2.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:testfile2",
		"testdata/test-l2.tar.gz:testfile3",
		"testdata/test-l2.tar.gz:testfile3",
	}
	i := 0
	for res := range zt.Start([]string{"testdata/test-l2.tar.gz"}) {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
		if p := strings.Join(res.Path, ":"); p != tt[i] {
			t.Errorf("%s != %s", p, tt[i])
		}
		i++
	}
	if i != len(tt) {
		t.Error("Too few results")
	}
}


func TestZTgrepZip(t *testing.T) {
	zt, err := ztgrep.New("test")
	if err != nil {
		t.Fatal(err)
	}
	tt := []string{
		"testdata/test-l2.zip:test-l1.zip",
		"testdata/test-l2.zip:test-l1.zip:test.tgz",
		"testdata/test-l2.zip:test-l1.zip:test.tgz:testfile1",
		"testdata/test-l2.zip:test-l1.zip:test.tgz:testfile1",
		"testdata/test-l2.zip:test-l1.zip:test.tgz:testfile2",
		"testdata/test-l2.zip:test-l1.zip:test.tgz:testfile2",
		"testdata/test-l2.zip:test-l1.zip:test.zip",
		"testdata/test-l2.zip:test-l1.zip:test.zip:testfile1",
		"testdata/test-l2.zip:test-l1.zip:test.zip:testfile1",
		"testdata/test-l2.zip:test-l1.zip:test.zip:testfile2",
		"testdata/test-l2.zip:test-l1.zip:test.zip:testfile2",
		"testdata/test-l2.zip:test-l1.zip:testfile1",
		"testdata/test-l2.zip:test-l1.zip:testfile1",
		"testdata/test-l2.zip:test-l1.zip:testfile2",
		"testdata/test-l2.zip:test-l1.zip:testfile2",
		"testdata/test-l2.zip:test.tgz",
		"testdata/test-l2.zip:test.tgz:testfile1",
		"testdata/test-l2.zip:test.tgz:testfile1",
		"testdata/test-l2.zip:test.tgz:testfile2",
		"testdata/test-l2.zip:test.tgz:testfile2",
	}
	i := 0
	for res := range zt.Start([]string{"testdata/test-l2.zip"}) {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
		if p := strings.Join(res.Path, ":"); p != tt[i] {
			t.Errorf("%s != %s", p, tt[i])
		}
		i++
	}
	if i != len(tt) {
		t.Error("Too few results")
	}
}