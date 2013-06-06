package web

import "n/test"
import "testing"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type controllerTestTypeA struct {
  ControllerBase
}

func (self controllerTestTypeA) Binding() *Binding {
  return nil
}

func (self controllerTestTypeA) Init(path string) {
}

func controllerTestSetup(T *testing.T) (test.Assert, Controller) {
  var assert = test.New.Assert(T)
  var instance = controllerTestTypeA{}
  return assert, instance
}

func controllerTestTeardown(target Controller) {
  // No action required
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_controller_can_create_instance(T *testing.T) {
  var a, i = controllerTestSetup(T)
  a.NotNil(i, "Unable to create instance")
  controllerTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
