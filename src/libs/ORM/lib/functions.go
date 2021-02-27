package lib

import "fmt"

func PrintPre(structure interface{}) {
	fmt.Printf("%+v\n", structure)
}
