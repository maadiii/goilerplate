package admin_test

import (
	"goldfish/app/testcase"
	"os"
	"testing"
)

var suittest *testcase.T

func TestMain(m *testing.M) {
	suittest = testcase.New()

	run := m.Run()

	suittest.Close()

	os.Exit(run)
}
