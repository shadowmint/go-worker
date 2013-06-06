package resources 

import "n/test"
import "testing"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

func fileConfigTestSetup(T *testing.T) (test.Assert, Config) {
  var assert = test.New.Assert(T)
  var instance = New.FileConfig()
  return assert, instance
}

func fileConfigTestTeardown(target Config) {
  // No action required
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_file_config_can_create_instance(T *testing.T) {
  var a, i = fileConfigTestSetup(T)
  a.NotNil(i, "Failed to create instance")
  fileConfigTestTeardown(i)
}

func Test_file_config_can_do_basic_ops(T *testing.T) {
  var a, i = fileConfigTestSetup(T)

  var data = `
    CONFIG0 = VALUE1
    CONFIG2 = VALUE2
    CONFIG3 = VALUE3

    # This is a comment block
    # CONFIG4 = VALUE4
    CONFIG5 = VALUE5
  `

  i.Parse(data)

  a.True(i.Has("CONFIG0"), "Missing key 'CONFIG0'")
  a.True(i.Has("CONFIG2"), "Missing key 'CONFIG2'")
  a.True(i.Has("CONFIG3"), "Missing key 'CONFIG3'")
  a.False(i.Has("CONFIG4"), "Found non-existent key 'CONFIG4'")
  a.True(i.Has("CONFIG5"), "Missing key 'CONFIG5'")

  a.Equals(i.Get("CONFIG0"), "VALUE1", "Invalid CONFIG0 value read")
  a.Equals(i.Get("CONFIG2"), "VALUE2", "Invalid CONFIG2 value read")
  a.Equals(i.Get("CONFIG3"), "VALUE3", "Invalid CONFIG3 value read")
  a.Equals(i.Get("CONFIG5"), "VALUE5", "Invalid CONFIG5 value read")

  fileConfigTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
