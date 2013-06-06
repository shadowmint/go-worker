package data

import "n"
import rf "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

type Store interface {

  // Register a type this store can work with.
  // This is always mandatory because go filters code by used objects;
  // failure to use an object at least once (in Register()) will mean
  // objects of that type cannot be deserialized.
  Register(t rf.Type, id string) error

  // Save a record and return a key id for it.
  // Note that only public fields are persisted.
  // ie. Field 'x' no, field 'X' yes.
  //
  // If the 'key' field is passed, this acts as an update, rather
  // than insert operation, and returns the persisted key value.
  Set(typeName string, record interface{}, key ...int64) (int64, error)

  // Fetch a record by key id.
  // Note the type will be a pointer, regardless of the typeName.
  // ie. Use value.(*MyType) to get a local reference to it.
  Get(typeName string, key int64) (interface{}, error)

  // Delete a record by key id.
  Delete(typeName string, key int64) error

  // Return a set of type objects which are registered
  Registered() []string

  // Return a set of objects based on the filter function.
  // 
  // This function returns the 'count' records that the filter function
  // returns true for, after skipping the first 'offset' records that
  // the filter returns true for.
  // 
  // Only the properties specified in the the props property are passed
  // into the filter function, but the entire record is returned as the
  // result. The record will be the appropriate type for the registered
  // typeName.
  Filter(typeName string, offset int, count int, props []string, filter func (record interface{}) bool) (map[int64]interface{}, error)

  // Shortcut for filter with an 'accepts all' filter.
  All(typeName string, offset int, count int) (map[int64]interface{}, error)

  // Count of all records of a given type
  Count(typeName string) (int, error)

  // Drop all records from the given type
  Clear(typeName string) error
}

// 
func newStore(drivers ...interface{} /* Driver */) (out Store) {
  defer func() { 
    var f = recover()
    if f != nil { 
      n.Log("Failed to create Store: %s", f)
      out = nil
    }
  }()

  // Try resolving a driver if one isn't provided
	var raw = n.Resolve(rf.TypeOf((*Driver)(nil)), 0, drivers...)
	if raw == nil {
		n.Log("Failed to create Store: No binding for Driver found")
		return nil
	}
	
  var rtn = &store {
    types : make(map[string] rf.Type),
    driver : raw.(Driver),
  }

  return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/


type store struct {
  
  // The persistence implementation for this store. 
  driver Driver

  // Names and types this store knows about
  types map[string] rf.Type
}

func (self *store) Register(t rf.Type, id string) error {
  for t.Kind() == rf.Ptr {
    t = t.Elem()
  }

  if t.Kind() != rf.Struct {
    return n.Fail(E_BAD_TYPE, "Invalid type; only struct types can be registered")
  }
  
  self.types[id] = t
  return nil
}

func (self *store) Set(typeName string, record interface{}, key ...int64) (int64, error) {

  // Bounds 
  if n.IsNil(record) {
    return -1, n.Fail(E_BAD_RECORD, "Nil record cannot be saved")
  }
  
  // Find type
  var t, found = self.types[typeName]
  if !found {
    return -1, n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  // Invoke the driver
  var rkey, derr = self.driver.Set(t, record, key...)
  return rkey, derr
}

func (self *store) Get(typeName string, key int64) (interface{}, error) {

  // Find type
  var t, found = self.types[typeName]
  if !found {
    return nil, n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  // Return instance
  var rtn, err = self.driver.Get(t, key)
  return rtn, err
}

func (self *store) Delete(typeName string, key int64) error {
  
  // Find type
  var t, found = self.types[typeName]
  if !found {
    return n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  var err = self.driver.Delete(t, key)
  return err
}

func (self *store) Registered() []string {
  var count = len(self.types)
  var rtn = make([]string, count)
  var i = 0
  for k := range self.types {
    rtn[i] = k
    i++
  }
  return rtn
}

func (self *store) Filter(typeName string, offset int, count int, props []string, filter func (record interface{}) bool) (map[int64]interface{}, error) {
  var t, found = self.types[typeName]
  if !found {
    return nil, n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  if filter == nil {
    return nil, n.Fail(E_BAD_FILTER, "Invalid filter (nil)")
  }

  // Entire return space
  var rtn = make(map[int64]interface{})

  // Look for matching records
  var total, _ = self.driver.Count(t)
  var query_offset = 0
  var query_size = count
  var result_offset = 0
  var match_count = 0
  for match_count < count {

    // Find a new block
    var keys, err = self.driver.Keys(t, query_offset, query_size)
    if err != nil {
      return nil, n.Fail(E_DRIVER_FAILURE, err.Error())
    }

    for _, k := range keys {
      var v, ve = self.driver.Get(t, k, props...)
      if ve != nil {
        return nil, n.Fail(E_DRIVER_FAILURE, err.Error())
      }
      
      // Collect matching records
      if filter(v) {
        result_offset += 1
        if result_offset > offset {
          rtn[k] = v
          match_count += 1
          if match_count >= count {
            break;
          }
        } 
      }
    }

    query_offset += query_size
    if query_offset >= total {
      break;
    }
  }

  // Cut result to sub-size
  return rtn, nil
}

func (self *store) All(typeName string, offset int, count int) (map[int64]interface{}, error) {
  var t, found = self.types[typeName]
  if !found {
    return nil, n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  var keys, err = self.driver.Keys(t, offset, count)
  if err != nil {
    return nil, n.Fail(E_DRIVER_FAILURE, err.Error())
  }

  var rtn = make(map[int64]interface{})
  for _, k := range keys {
    var v, ve = self.driver.Get(t, k)
    if ve != nil {
      return nil, n.Fail(E_DRIVER_FAILURE, ve.Error())
    }
    rtn[k] = v
  }

  return rtn, nil
}

func (self *store) Count(typeName string) (int, error) {
  var t, found = self.types[typeName]
  if !found {
    return 0, n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  var rtn, err = self.driver.Count(t)
  return rtn, err
}

func (self *store) Clear(typeName string) error {
  var t, found = self.types[typeName]
  if !found {
    return n.Fail(E_BAD_TYPE, "Invalid type '%s' is not known. Did you register it?", typeName)
  }

  // Delete all keys
  var count, _ = self.driver.Count(t)
  var keys, _ = self.driver.Keys(t, 0, count)
  for _, k := range keys {
    self.driver.Delete(t, k)
  }
  return nil
}

/*============================================================================*
 * }}} 
 *============================================================================*/
