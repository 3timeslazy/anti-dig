package main

import (
	dig "github.com/3timeslazy/anti-dig"
	grpcserver "github.com/3timeslazy/anti-dig/anti/testcases/same_pkg_name/grpc/server"
	httpserver "github.com/3timeslazy/anti-dig/anti/testcases/same_pkg_name/http/server"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(httpserver.New)
	_ = container.Provide(grpcserver.New)
	_ = container.Invoke(func(httpserver.Server, grpcserver.Server) {})
}
