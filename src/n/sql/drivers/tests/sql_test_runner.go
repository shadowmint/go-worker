package tests

import nsql "n/sql"
import "n/test"
import "testing"

/*============================================================================*
 * {{{ Public api
 *============================================================================*/
 
type SqlTestRunner interface {
	Run() bool
  Table() string
	Setup(T *testing.T) (test.Assert, nsql.Sql)
	Teardown(driver nsql.Sql)	 
}
 
/*============================================================================*
 * }}}
 *============================================================================*/
 
