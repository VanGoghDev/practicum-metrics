package testdata

import (
	"fmt"
	"os"
)

//nolint:all // test func calls in osexitcheck_test.go
func main1() {
	defer os.Exit(1)

	fmt.Printf("%d", 1)
}

//nolint:all // test func calls in osexitcheck_test.go
func main2() {
	os.Exit(1) // ignored
}
