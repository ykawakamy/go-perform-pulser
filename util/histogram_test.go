package util

import (
	"testing"
)

func Test_Print(t *testing.T) {
	hist := CreateHistogram()

	hist.Increament(0)
	hist.Increament(10)
	hist.Increament(54)
	hist.Increament(55)

	hist.Print()

}
