package file

import "n/data"
import "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

func newDriver() data.Driver {
  return &driver {
    data : make(map[reflect.Type] RecordSet),
  }
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Implementation
 *============================================================================*/

type driver struct {

  // Map of collection of objects to record sets
  data map[reflect.Type] RecordSet
}

// Save a record and return a key id for it.
func (self *driver) Set(t reflect.Type, record interface{}, args ...int64) (int64, error) {
  var key int64
  var rs = self.Records(t)
  if len(args) > 0 {
    key = args[0]
  } else {
    key = rs.Key()
  }
  var err = rs.Set(key, record)
  return key, err
}

// Update a record 
func (self *driver) Update(t reflect.Type, key int64, record interface{}) error {
  var rs = self.Records(t)
  var err = rs.Set(key, record)
  return err
}

// Fetch a record by key id.
func (self *driver) Get(t reflect.Type, key int64, props ...string) (interface{}, error) {
  var rs = self.Records(t)
  var rtn, err = rs.Get(key, props)
  return rtn, err
}

// Delete a record by key id.
func (self *driver) Delete(t reflect.Type, key int64) error {
  var rs = self.Records(t)
  var err = rs.Unset(key)
  return err
}

// Return a count of this type of record.
func (self *driver) Count(t reflect.Type) (int, error) {
  var rs = self.Records(t)
  var rtn = rs.Count()
  return rtn, nil
}

// Return the keys for this record type, from offset to offset + count.
// If < count records are available after offset, return all found.
func (self *driver) Keys(t reflect.Type, offset int, count int) ([]int64, error) {
  var rs = self.Records(t)
  var rtn = rs.Keys(offset, count)
  return rtn, nil
}

/*============================================================================*
 * }}} 
 *============================================================================*/

/*============================================================================*
 * {{{ Internal
 *============================================================================*/

// Return the recordset for a type, or create it
func (self *driver) Records(t reflect.Type) RecordSet {
  var rtn, found = self.data[t]
  if !found {
    self.data[t] = New.RecordSet(t)
    rtn = self.data[t]
  }
  return rtn
}

/*============================================================================*
 * }}} 
 *============================================================================*/
