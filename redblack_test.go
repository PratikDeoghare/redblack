package redblack

import (
	"fmt"
	"testing"
)


func TestRedblack(t *testing.T) {
	a := New[int, string](func(x, y int) bool { return x < y })
	a.Set(12, "twelve")
	b := New[string, int](func(x, y string) bool { return x < y })
	b.Set("twelve", 12)
	b.Set("a", 12)
	b.Set("b", 12)
	b.(*rbmap[string, int]).Preorder()
	fmt.Println(a, b)
}

