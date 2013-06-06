package n

import "reflect"
import "os"
import "os/exec"
import "strings"
import "path"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/
 
// Shortcut for nil testing without panic
func IsNil(value interface{}) bool {
	var t = reflect.TypeOf(value)
	var rtn = t == nil
	if !rtn {
		if t.Kind() == reflect.Chan || t.Kind() == reflect.Func || t.Kind() == reflect.Interface || t.Kind() == reflect.Map || t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		  var v = reflect.ValueOf(value)
			rtn = v.IsNil()
		}
	}
	return rtn
}

// Shortcut to get the path to the current executable
func ExecPath() string {
  var here = os.Args[0]
  if !strings.HasPrefix(here, "/") {
    here, _ = exec.LookPath(os.Args[0])
    if !strings.HasPrefix(here, "/") {
      var wd, _ = os.Getwd()
      here = path.Join(wd, here)
    }
  }
  here = path.Dir(here)
  return here
}

/*============================================================================*
 * }}} 
 *============================================================================*/
