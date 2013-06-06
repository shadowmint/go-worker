package memory

import tests "n/data/drivers/tests"
import "n/data"
import "n/test"
import "testing"

/*============================================================================*
 * {{{ Tests
 *============================================================================*/
 
func Test_driver_tests(T *testing.T) {
  var assert = test.New.Assert(T)
	var runner = &testRunner{}
	var tests = tests.New.Tests()
	var result = tests.Run(runner, T)
	assert.True(result, "Driver tests failed")
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type testRunner struct {
}

func (self *testRunner) Run() bool {
	return true
}

func (self *testRunner) Setup(T *testing.T) (test.Assert, data.Store) {
  var assert = test.New.Assert(T)
  var instance = data.New.Store(New.Driver())
  return assert, instance
}

func (self *testRunner) Teardown(instance data.Store) {
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
