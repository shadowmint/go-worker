package tests

import "testing"
import "strings"
import "n"
import "reflect"
import "runtime/debug"

/*============================================================================*
 * {{{ Public api
 *============================================================================*/
 
// Control if any sql tests are run; turn off in most projects
const RUN_DRIVER_TESTS = true

type Tests interface {
	Run(runner TestRunner, T *testing.T) bool
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * {{{ Implementation
 *============================================================================*/

type tests struct {
	runner TestRunner
}

func newTests() Tests {
	return &tests{}
}

func (self *tests) Run(runner TestRunner, T *testing.T) (rtn bool) {

  // Shortcut; stop if we're not in testing mode
  if !RUN_DRIVER_TESTS { return true; }

	var current = "None"
	defer func() {
		var e = recover()
		if e != nil {
			n.Log("Failed while running test \"%s\": %s", current, e)
			n.Log("%s", debug.Stack())
			rtn = false
		}
	}()
	
	var tt = reflect.TypeOf(self)
	var tv = reflect.ValueOf(self)
	var mc = tt.NumMethod()
	var count = 0
	for i := 0; i < mc; i++ {
		var test_method = tt.Method(i)
		var test_name = test_method.Name
		if strings.HasPrefix(test_name, "Test_") {
			current = test_name
			var test_method_instance = tv.Method(i)
			var test_method_args = []reflect.Value { reflect.ValueOf(runner), reflect.ValueOf(T) }
			test_method_instance.Call(test_method_args)
			count++
		}
	}
	
	if count == 0 {
		n.Log("No tests were found in the driver suite")	
		return false
	}
	
	return true
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
