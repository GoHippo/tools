package main

import (
	"fmt"
	"github.com/GoHippo/tools/m_file"
)

func main() {
	a, _ := m_file.File_open("go.mod")
	fmt.Println(string(a))
}
