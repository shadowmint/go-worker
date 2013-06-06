package data

import "n/test"
import "testing"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

func storeTestSetup(T *testing.T) (test.Assert, Store) {
  var assert = test.New.Assert(T)
  var instance = New.Store()
  return assert, instance
}

func storeTestTeardown(target Store) {
  // No action required
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_store_cannot_create_empty_instance(T *testing.T) {
  var a, i = storeTestSetup(T)

  a.Nil(i, "Able to create invalid instance")

  storeTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
