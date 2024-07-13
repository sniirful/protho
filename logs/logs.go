package logs

import "fmt"

func PrintE(a ...any) {
	fmt.Println(a...)
}

func PrintI(a ...any) {
	fmt.Println(a...)
}

func PrintV(verbose bool, a ...any) {
	if verbose {
		fmt.Println(a...)
	}
}
