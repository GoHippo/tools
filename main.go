package main

import (
	"fmt"
	"github.com/GoHippo/tools/m_file"
)

// only test func
func main() {
	err := m_file.Zip_Extract("putty-0.73-ru-17.zip", "")
	if err != nil {
		fmt.Println(err)
	}
}
