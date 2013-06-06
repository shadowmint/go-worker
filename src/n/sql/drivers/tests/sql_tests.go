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
const RUN_DRIVER_TESTS = false

type SqlTests interface {
	Run(runner SqlTestRunner, T *testing.T) bool
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * {{{ Implementation
 *============================================================================*/

type sqlTests struct {
	runner SqlTestRunner
}

func newSqlTests() SqlTests {
	return &sqlTests{}
}

func (self *sqlTests) Run(runner SqlTestRunner, T *testing.T) (rtn bool) {

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
	for i := 0; i < mc; i++ {
		var test_method = tt.Method(i)
		var test_name = test_method.Name
		if strings.HasPrefix(test_name, "Test_") {
			current = test_name
			var test_method_instance = tv.Method(i)
			var test_method_args = []reflect.Value { reflect.ValueOf(runner), reflect.ValueOf(T) }
			test_method_instance.Call(test_method_args)
		}
	}
	
	return true
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
