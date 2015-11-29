package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
	var toReplace string = `<![CDATA[<img src="thumbnail.php?237910.jpg" alt="" title="" /><br />`
	var expectedReplace string = `<![CDATA[<img src="http://hdclub.org/thumbnail.php?237910.jpg" alt="" title="" /><br />`

	inBytes := []byte(toReplace)
	var outBytes []byte = imgRE.ReplaceAllFunc(inBytes, replaceBadHrefs)
	if string(outBytes) != expectedReplace {
		t.Fatal(string(outBytes))
	}
}
