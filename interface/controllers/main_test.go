package controllers_test

import (
	"goilerplate/infrastructure/testutil"
	"os"
	"testing"
)

var suittest *testutil.T

func TestMain(m *testing.M) {
	suittest = testutil.New()

	run := m.Run()

	suittest.Close()

	os.Exit(run)
}
