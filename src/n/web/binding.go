package web 

/*============================================================================*
 * {{{ Public api
 *============================================================================*/
 
// Keeps track of bindings between actions and handlers for controllers.
type Binding struct {

  // Attach the controller instance to this
  Instance interface{}

  // Id of this binding
  Id string

  // Map of 
  Map map[string] interface{}
}

/*============================================================================*
 * }}}
 *============================================================================*/