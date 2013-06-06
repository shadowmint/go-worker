package n

import r "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// This is a singleton resolver, it's a little odd, so read the details of
// register before using it.
//
// When using this type remember that reflect.TypeOf() returns the type of
// the VALUE of the interface passed to it, so never do this:
//
// var t MyInterfaceType = Resolver().Get(reflect.TypeOf(t))
type Resolver interface {

  // Register a type to keep track of.
  //
  // Factory is assumed to be a zero argument function that returns a 
  // single value of the appropriate type.
  //
  // Note that the type is only attached to an internal singleton
  // instance after creation is completed, so it is possible to get
  // cyclic resolver chains happening if the factory relies on the
  // resolver itself.
  Register(t r.Type, factory interface{}) error

  // Return a singleton instance of a given type or nil
  Get(t r.Type) (interface{})

  // Return a singleton instance and an error on failure
  Singleton(t r.Type) (interface{}, error)

  // Clear the register.
  // Useful for testing mostly.
  Clear() 
}

// Shortcut to directly resolve an instance from an array of interface{}
func Resolve(t r.Type, offset int, args ...interface{}) interface{} {
	var rtn interface{} = nil
	if len(args) > offset {
		rtn = args[offset]
	}
	if rtn != nil {
		if r.TypeOf(rtn).AssignableTo(t) {
			Log("Invalid param: '%s' != required type '%s'. Autoresolving.", r.TypeOf(rtn), t)
			rtn = nil
		}
	}
	if rtn == nil {
		rtn = New.Resolver().Get(t)
		if rtn == nil {
			Log("Unable to resolve required type: '%s'", t)
		}
	}
	return rtn
}

func newResolver() Resolver {
  if resolverInstance == nil {
    resolverInstance = &resolver {}
    resolverInstance.Clear()
  }
  return resolverInstance
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ IResolver api implementation
 *============================================================================*/

var resolverInstance Resolver = nil

type resolver struct {
  
  // Set of singletons we've already created
  instance map[r.Type] interface{}

  // Set of factory functions
  factory map[r.Type] r.Value
}

func (self *resolver) Register(t r.Type, factory interface{}) error {
  t = self.derefType(t) // Keep only the base type

  if t == nil {
    return Fail(1, "Invalid type (nil); did you mean &t instead of t?")
  }

  var v = r.ValueOf(factory)
  if v.Kind() != r.Func {
    return Fail(1, "Invalid factory; it must be a function.")
  }

  self.factory[t] = v
  var _, found = self.instance[t] 
  if found {
    delete(self.instance, t)
  }

  return nil
}

func (self *resolver) Get(t r.Type) interface{} {
  var rtn, err = self.Singleton(t)
  if err != nil {
    Log("Resolve failure: %s", err)
  }
  return rtn
}

func (self *resolver) Singleton(t r.Type) (rtn interface{}, err error) {

  // Careful, that dangerous interface might not be what we think it is
  defer func() {
    var failure = recover()
    if failure != nil {
      err = Fail(1, "Factory failure: %s", failure.(error).Error())
    }
  }()
  
  t = self.derefType(t) // Keep only the base type
  var factory, found = self.factory[t]
  if !found {
    return nil, Fail(1, "Invalid type '%s'; no binding to that type", t)
  }
  
  var results = factory.Call([]r.Value {})
  if len(results) == 0 {
    return nil, Fail(1, "Factory for type '%s' did not return an instance", t)
  } else if (results[0].IsNil()) {
    return nil, Fail(1, "Factory for type '%s' returned a nil instance", t)
  }

  rtn = results[0].Interface()
  var kind, tvalue = self.derefKind(rtn)
  if ((kind != r.Struct) && (kind != r.Interface)) {
    return nil, Fail(1, "Factory return an invalid object for id '%s': %+v", t, rtn)
  } else if !r.TypeOf(rtn).AssignableTo(t) {
    return nil, Fail(1, "Factory returned invalid type (%s is not assignable to %s)", tvalue, t)
  } else if (rtn == nil) {
    return nil, Fail(1, "Factory returned a nil value for id '%s'", t)
  }
  
  return rtn, nil
}

func (self *resolver) Clear() {
  self.instance = make(map[r.Type] interface{})
  self.factory = make(map[r.Type] r.Value)
}

/*============================================================================*
 * }}} 
 *============================================================================*/

/*============================================================================*
 * {{{ Internal
 *============================================================================*/

// Return a unique key for this type
func (self *resolver) key(t r.Type) string {
	var path = t.PkgPath() + "." + t.Name()
	return path
}

func (self *resolver) derefType(t r.Type) r.Type {
  if t == nil {
    return nil
  }
  for t.Kind() == r.Ptr {
    t = t.Elem()
  }
  return t
}

func (self *resolver) derefKind(value interface{}) (r.Kind, r.Type) {
  var v = r.TypeOf(value)
  for v.Kind() == r.Ptr {
    v = v.Elem()
  }
  return v.Kind(), v
}

/*============================================================================*
 * }}} 
 *============================================================================*/
