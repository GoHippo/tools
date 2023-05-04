package main

import (
	"fmt"
	"github.com/GoHippo/tools/m_file"
)

// only test func
func main() {
	err := m_file.Zip_Extract("/test.zip", "")
	if err != nil {
		fmt.Println(err)
	}
}
