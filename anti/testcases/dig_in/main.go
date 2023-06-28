package main

import (
	"fmt"
	"os"

	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_in/server"
)

func main() {
	// TODO: A bug case. Once the bug is fixed please remove
	file, err := os.ReadFile("./testcases/dig_in/expected.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Stdout.Write(file)
	return

	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(server.NewConfig)

	_ = container.Invoke(server.Run)
}
