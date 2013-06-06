package web

import "net/http"

/*============================================================================*
 * {{{ Internal api 
 *============================================================================*/

type factory struct {
  Context func(w http.ResponseWriter, r *http.Request, vars map[string]string) *Context 
}

var New factory = factory {
  Context : newContext,
}

/*============================================================================*
 * }}}
 *============================================================================*/
