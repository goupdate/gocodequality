package testdata

import "fmt"
import "fmt" // duplicate import
import "unsafe"

const Max_Size = 100 // bad naming

var password = "secret123" // hardcoded secret
var unusedVar = 42         // unused var

type MyStruct struct {
	unusedField int
}

func ExportedNoDoc() {} // missing godoc

func unusedPrivateFunc() {} // unused function

func emptyBody() {} // empty body

func tooManyParams(a, b, c, d, e, f int) {} // too many params

func deepNesting() {
	if true {
		if true {
			if true {
				if true {
					if true {
						fmt.Println("deep")
					}
				}
			}
		}
	}
}

func withPanic() {
	panic("oops")
}

func ignoredError() {
	_ = fmt.Println("test")
}

func magicNumbers() {
	x := 42
	_ = x
}

func objectInLoop() {
	for i := 0; i < 10; i++ {
		s := make([]int, 10)
		_ = s
	}
}

func unbufferedChan() {
	ch := make(chan int)
	_ = ch
}

// TODO: implement this
func withTodo() {}

func commentedCode() {
	// fmt.Println("hello")
}

func _unsafe() {
	_ = unsafe.Sizeof(0)
}
