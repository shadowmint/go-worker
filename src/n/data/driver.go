package data

import "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

type Driver interface {
  
  // Save a record and return a key id for it.
  // Or, if a key is passed, use that instead.
  Set(t reflect.Type, record interface{}, key ...int64) (int64, error)

  // Fetch a record by key id, populating only the fields with the given names.
  // If no properties are supplied, return all properties on the object.
  Get(t reflect.Type, key int64, props ...string) (interface{}, error)

  // Delete a record by key id.
  Delete(t reflect.Type, key int64) error

  // Return a count of this type of record.
  Count(t reflect.Type) (int, error)

  // Return the keys for this record type, from offset to offset + count.
  // If < count records are available after offset, return all found.
  Keys(t reflect.Type, offset int, count int) ([]int64, error)
}

/*============================================================================*
 * }}} 
 *============================================================================*/
