package file

import "n/data"
import "n/test"
import "testing"
import "path/filepath"
import "os"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type driverTest struct {
  assert test.Assert
  driver data.Driver
}

func driverTestSetup(T *testing.T) (test.Assert, *driverTest) {
  var assert = test.New.Assert(T)
  var rtn = &driverTest{}
  rtn.driver = newDriver()

  return assert, rtn
}

func driverTestTeardown(target *driverTest) {
  var path, err = filepath.Abs("data")
  if err == nil {
    os.RemoveAll(path)
  } 
}

/*============================================================================*

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_driver_can_create_instance(T *testing.T) {
  var a, c = driverTestSetup(T)

  a.NotNil(c.driver, "Unable to create instance")

  driverTestTeardown(c)
}

/*============================================================================*
 * }}}
 *============================================================================*/
