package supportscolor_test

import (
	"fmt"

	"github.com/jwalton/go-supportscolor"
)

func ExampleSupportsColor() {
	if supportscolor.Stdout().SupportsColor {
		fmt.Println("\u001b[31mThis is Red!\u001b[39m")
	} else {
		fmt.Println("This is not red.")
	}

	// Output: This is not red.
}
