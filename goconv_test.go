package goconv

import (
	"os"
	"testing"
)

func TestConv(t *testing.T) {
	s := NewService()
	f, err := os.Open("./test.pptx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	out, err := os.Create("./test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	err = s.Convert(f, out)
	if err != nil {
		t.Fatal(err)
	}
}
