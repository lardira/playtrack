package testutil

import (
	"math/rand/v2"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

var (
	faker     *gofakeit.Faker
	fakerOnce sync.Once

	seed uint64
)

func init() {
	initFaker()
}

func initFaker() {
	fakerOnce.Do(func() {
		seed = rand.Uint64()
		faker = gofakeit.NewFaker(source.NewJSF(seed), true)
	})
}

func Faker() *gofakeit.Faker {
	return faker
}

func GetSeed() uint64 {
	return seed
}
