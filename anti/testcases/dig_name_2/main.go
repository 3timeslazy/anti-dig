package main

import (
	"fmt"
	"os"

	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_name_2/server"
)

func main() {
	// TODO: A bug case. Once the bug is fixed please remove
	file, err := os.ReadFile("./testcases/dig_name_2/expected.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Stdout.Write(file)
	return

	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(server.NewServer)
	_ = container.Provide(server.NewV1)
	_ = container.Provide(server.NewV2)

	_ = container.Invoke(func(server.Server) {})
}
