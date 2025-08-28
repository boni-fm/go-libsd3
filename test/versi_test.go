package test

import (
	"testing"

	"github.com/boni-fm/go-libsd3/repo/versi"
)

func TestGetVersiProgramPostgre(t *testing.T) {
	// These are dummy values for testing string logic only
	Constr := "Kamu mau tau ya ehehehe..."
	Kodedc := "6969"
	NamaProgram := "Word.exe"
	Versi := "1.0.0.0"
	IPKomputer := "100.100.100.100"

	result := versi.GetVersiProgramPostgre(Constr, Kodedc, NamaProgram, Versi, IPKomputer)
	t.Logf("Result: %s", result)
}
