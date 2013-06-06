package memory

import "encoding/json"
import "reflect"
import "fmt"
import "sort"

/*============================================================================*
 * {{{ Internal api 
 *============================================================================*/

type RecordSet interface {

  // Set a record
  Set(key int64, record interface{}) error 

  // Unset a record
  Unset(key int64) error

  // Count of records
  Count() int

  // Subset of keys
  Keys(offset int, count int) []int64

  // Generate a new unique key
  Key() int64

  // Get a record
  Get(key int64, props []string) (interface{}, error)
}

func newRecordSet(t reflect.Type) RecordSet {
  return &recordSet {
    Type : t,
    LastKey : 0,
    Data : make(map[int64] map[string] string),
  }
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Implementation  
 *============================================================================*/

// Collection of records and a key repo
type recordSet struct {

  // Type we're looking after here
  Type reflect.Type

  // The last key we generated
  LastKey int64

	// The set of objects we're currently looking after
	Data map[int64] map[string] string
}

func (self *recordSet) Set(key int64, record interface{}) error {
	var count = self.Type.NumField()
	var irecord = reflect.Indirect(reflect.ValueOf(record))
  var content = map[string]string {}
	for i := 0; i < count; i++ {
		var field = self.Type.Field(i)
		var value = irecord.Field(i).Interface()
	  content[field.Name] = self.Serialize(value)
	}
	self.Data[key] = content
  return nil
}

func (self *recordSet) Unset(key int64) error {
  delete(self.Data, key)
  return nil
}

func (self *recordSet) Count() int {
  return len(self.Data)
}

func (self *recordSet) Keys(offset int, count int) []int64 {

  var size = self.Count()

  // No records?
  var lower = offset
  if lower >= size {
    return []int64 {}
  }

  // Get all known keys
  var all = make([]int64, size, size)
  var ai int64arr = all
  var i = 0
  for k, _ := range self.Data {
    all[i] = k
    i++
  }
  sort.Sort(ai)

  // Take a slice from all 
  var upper = offset + count 
  if upper >= size {
    upper = size
  }
  var rtn = all[lower:upper]
    
  return rtn
}

func (self *recordSet) Key() int64 {
  self.LastKey += 1
  return self.LastKey
}

func (self *recordSet) Get(key int64, props []string) (interface{}, error) {

  var content, found = self.Data[key]
  if !found {
    return nil, fmt.Errorf("Invalid key: no such record (%d)", key)
  }
  
  // If not specified, fetch all.
 	if len(props) == 0 {
	  var count = self.Type.NumField()
 		props = make([]string, count, count)
		for i := 0; i < count; i++ {
			var field = self.Type.Field(i)
			props[i] = field.Name
		}
 	}
 	
  var rtn = self.Deserialize(content, props)
  return rtn, nil
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Internal functions  
 *============================================================================*/

// For sorting longs
type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

// Serialize to json
func (self *recordSet) Serialize(record interface{}) string {
  var raw, _ = json.Marshal(record) 
  return string(raw)
}

// Deserialize from json
func (self *recordSet) Deserialize(value map[string]string, props []string) interface{} {
  var instance = reflect.New(self.Type).Interface()
  
  var vi = reflect.Indirect(reflect.ValueOf(instance))
  
 	for i := 0; i < len(props); i++ {
 		var member_type, found = self.Type.FieldByName(props[i])
 		var value_type = vi.FieldByName(props[i])
 		if found {
	 		var member = reflect.New(member_type.Type).Interface()
	 		var vm = reflect.Indirect(reflect.ValueOf(member))
		  var b = []byte(value[props[i]])
		  json.Unmarshal(b, &member) 
		  value_type.Set(vm)
 		}
 	} 
  return instance
}

/*============================================================================*
 * }}}
 *============================================================================*/
