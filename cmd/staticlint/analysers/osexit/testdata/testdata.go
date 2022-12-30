package testdata

import (
	"os"
)

func main() {
	os.Exit(0) // want "os.Exit found in main function"
}
