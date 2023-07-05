package registry

import (
	dig "github.com/3timeslazy/anti-dig"
)

var providers = []any{}

func Provide(fn any) {
	providers = append(providers, fn)
}

func Invoke(container *dig.Container, fn any) error {
	for _, provider := range providers {
		err := container.Provide(provider)
		if err != nil {
			return err
		}
	}
	return container.Invoke(fn)
}
