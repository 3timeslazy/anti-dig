package antidig

import "math/rand"

func SetRand(r *rand.Rand) Option {
	return setRand(r)
}
