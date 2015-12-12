package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type compilation struct {
	numberOfWorkers []int
}

func NewCompilation(numberOfWorkers []int) Parameter {
	return &compilation{
		numberOfWorkers: numberOfWorkers,
	}
}

func (c *compilation) Apply(input bftinput.Input) bftinput.Input {
	input.CloudConfig.NumberOfCompilationWorkers = c.numberOfWorkers[rand.Intn(len(c.numberOfWorkers))]
	return input
}
