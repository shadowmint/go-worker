package tests

import "n/test"
import "n/data"
import "testing"
import "reflect"

/*============================================================================*
 * {{{ Public api
 *============================================================================*/
 
type TestRunner interface {
	Run() bool
	Setup(T *testing.T) (test.Assert, data.Store)
	Teardown(store data.Store)	 
}
 
/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type storeTestTypeA struct {
  A int
  B int64
  C float32
  D bool
  E string
}

type storeTest struct {
  assert test.Assert
  store data.Store
  driver data.Driver
}

func storeTestSetup(runner TestRunner, T *testing.T) (test.Assert, data.Store) {
  var assert, rtn = runner.Setup(T)
  
 	// Register
 	rtn.Register(reflect.TypeOf((*storeTestTypeA)(nil)), "storeTestTypeA")
 
  return assert, rtn
}

func storeTestTeardown(runner TestRunner, target data.Store) {
  target.Clear("storeTestTypeA")
  runner.Teardown(target)
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
