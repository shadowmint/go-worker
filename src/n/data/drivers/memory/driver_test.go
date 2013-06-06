package memory

import "n/data"
import "n/test"
import "testing"

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
  rtn.driver = New.Driver()

  return assert, rtn
}

func driverTestTeardown(target *driverTest) {
  // No action required
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
