package anti_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/db"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	httpserver "github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/observability"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Constructor struct {
	Func any
	Opts []dig.ProvideOption
}

func TestAnti(t *testing.T) {
	testCases := []struct {
		Name         string
		Constructors []Constructor
		Opts         []dig.ProvideOption
		Invoker      any
	}{
		{
			Name: "provide_one",
			Constructors: []Constructor{
				{Func: config.NewConfig},
			},
			Invoker: func(config.Config) {},
		},
		{
			Name: "provide_many",
			Constructors: []Constructor{
				{Func: config.NewConfig},
				{Func: db.NewDB},
				{Func: handlerv0.NewHandlerV0},
			},
			Invoker: func(handlers.Handler) {},
		},
		{
			Name: "as_field",
			Constructors: []Constructor{
				{Func: config.NewConfig},
				{Func: observability.NewObservability},
			},
			Invoker: func(observability.Metrics) {},
		},
		{
			Name: "with_dig_name",
			Constructors: []Constructor{
				{
					Func: queue.New1,
					Opts: []dig.ProvideOption{dig.Name("queue_1")},
				},
				{
					Func: queue.New2,
					Opts: []dig.ProvideOption{dig.Name("queue_2")},
				},
				{
					Func: consumer.New,
				},
			},
			Invoker: func(consumer.Consumer) {},
		},
		{
			Name: "with_dig_group",
			Constructors: []Constructor{
				{
					Func: db.NewDB,
				},
				{
					Func: config.NewConfig,
				},
				{
					Func: handlerv0.NewHandlerV0,
					Opts: []dig.ProvideOption{dig.Group("http_handlers")},
				},
				{
					Func: httpserver.NewServer,
				},
			},
			Invoker: func(*httpserver.Server) {},
		},
		{
			Name: "with_tag_group",
			Constructors: []Constructor{
				{Func: db.NewDB},
				{Func: config.NewConfig},
				{Func: handlerv1.NewHandlerV1},
				{Func: grpcserver.NewServer},
			},
			Invoker: func(*grpcserver.Server) {},
		},
		{
			Name: "with_flatten_group",
			Constructors: []Constructor{
				{Func: config.NewConfig},
				{Func: observability.NewObservability},
				{Func: flatten.NewListOfHandlers},
				{Func: httpserver.NewServer},
			},
			Invoker: func(*httpserver.Server) {},
		},
		{
			Name: "with_the_same_package_name",
			Constructors: []Constructor{
				{Func: config.NewConfig},
				{Func: db.NewDB},

				{
					Func: handlerv0.NewHandlerV0,
					Opts: []dig.ProvideOption{dig.Group("http_handlers")},
				},
				{Func: httpserver.NewServer},

				{Func: handlerv1.NewHandlerV1},
				{Func: grpcserver.NewServer},
			},
			Invoker: func(*httpserver.Server, *grpcserver.Server) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			dig.Anti = dig.NewAnti(buf)

			container := dig.New()

			for _, ctor := range testCase.Constructors {
				err := container.Provide(ctor.Func, ctor.Opts...)
				require.NoError(t, err)
			}

			err := container.Invoke(testCase.Invoker)
			require.NoError(t, err)

			expected, err := os.ReadFile(fmt.Sprintf("testdata/%s.go", testCase.Name))
			require.NoError(t, err)

			actual := buf.Bytes()

			assert.Equal(t, string(expected), string(actual))
		})
	}
}
