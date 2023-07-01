package main

import (
	"fmt"
	"os"

	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti"
	rename "github.com/3timeslazy/anti-dig/anti/gorename"
	"github.com/3timeslazy/anti-dig/anti/optimiser/testcases/rename/consumer"
	"github.com/3timeslazy/anti-dig/anti/optimiser/testcases/rename/registry"
)

func main() {
	expectedtxt, err := os.Create("expected.txt")
	if err != nil {
		fmt.Println("os.Create failed:", err)
		return
	}

	dig.Anti = anti.New(expectedtxt).
		Optimise(true).
		Rename(true)
	rename.Diff = true

	container := dig.New()

	err = registry.Invoke(container, func(c consumer.Consumer) {})
	if err != nil {
		fmt.Println("registry.Invoke failed:", err)
	}
}
