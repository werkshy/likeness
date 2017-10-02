package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/werkshy/likeness/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
