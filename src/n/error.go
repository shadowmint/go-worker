package n

import "fmt"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// Standard error type
type Error struct {
  
  // Error code associated with this response.
  Code int

  // Error message
  Message string
}

// Generate a standard error
func Fail(code int, msg string, args ...interface{}) *Error {
  return &Error {
    Code : code,
    Message : fmt.Sprintf(msg, args...),
  }
}

// Generate a standard error and log it to std out
func LogFail(code int, msg string, args ...interface{}) *Error {
	Log(msg, args...)
	return Fail(code, msg, args...)
}

/*============================================================================*
 * }}} 
 *============================================================================*/

/*============================================================================*
 * {{{ Error implementation
 *============================================================================*/

func (self *Error) Error() string {
  return fmt.Sprintf("%s (%d)", self.Message, self.Code)
}

/*============================================================================*
 * }}} 
 *============================================================================*/
