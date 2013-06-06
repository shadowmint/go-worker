package test

import "testing"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type assertTest struct {
  assert Assert
}

func assertTestSetup(T *testing.T) (Assert, *assertTest) {
  var assert = New.Assert(T)
  var rtn = &assertTest{}
  rtn.assert = assert

  return assert, rtn
}

func assertTestTeardown(target *assertTest) {
  // No action required
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_assert_can_create_instance(T *testing.T) {
  var a, c = assertTestSetup(T)

  a.NotNil(c.assert, "Unable to create instance")

  assertTestTeardown(c)
}

func Test_assert_can_pass_own_tests(T *testing.T) {
  var a, c = assertTestSetup(T)

  a.Nil(nil, "Nil was not nil")
  a.NotNil(a, "NotNil was nil")
  a.True(true, "True is not true")
  a.False(false, "False is not false")
  a.Equals(1, 1, "Equals is not equal")

  assertTestTeardown(c)
}

/*============================================================================*
 * }}}
 *============================================================================*/
