package webpage

import "testing"

func TestGetVisibleContent(t *testing.T) {
	res, _ := GetVisibleContent("https://www.google.com/")
	t.Log(res)
}
