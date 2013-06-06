package file

import "n"
import "math/rand"
import "encoding/json"
import "reflect"
import "path/filepath"
import "strings"
import "io/ioutil"
import "os"
import "fmt"

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
}

func (self *recordSet) Set(key int64, record interface{}) error {
  var err = self.Init()
  if err != nil { return err }
  
  var kpath, fpaths = self.KeyPath(key)
  err = ioutil.WriteFile(kpath, self.Serialize(record), 0755)
  if err != nil { return err }
  
  for field, fpath := range(fpaths) {
  	var vrecord = reflect.ValueOf(record)
  	if vrecord.Kind() == reflect.Ptr || vrecord.Kind() == reflect.Interface {
  		vrecord = vrecord.Elem()
  	}
  	var fvalue = vrecord.FieldByName(field).Interface()
	  err = ioutil.WriteFile(fpath, self.Serialize(fvalue), 0755)
	  if err != nil { return err }
  }
  
  return nil
}

func (self *recordSet) Count() int {
  _ = self.Init()
  var count = 0
  filepath.Walk(self.RecordsPath(), func(path string, info os.FileInfo, err error) error {
    if err != nil {
      n.Log("Error walking record set: %s", err)
    }
    if info != nil {
      if !info.IsDir() {
      	if strings.Count(info.Name(), ".") == 1 {
	        count++
        }
      }
    }
    return nil
  })
  return count
}

// Deserialize entire record if no specific props were requested
func (self *recordSet) Get(key int64, props []string) (interface{}, error) {
	var kpath, ppaths = self.KeyPath(key)
	
	// Entire record
	if len(props) == 0 {
	  var content, err = ioutil.ReadFile(kpath)
	  if err != nil { return nil, err; }
	  var rtn = self.Deserialize(content, self.Type)
	  return rtn, nil
	}
	
	// Only parts
	var rtn = reflect.New(self.Type).Interface()
	var vrtn = reflect.ValueOf(rtn)
	if vrtn.Kind() == reflect.Ptr || vrtn.Kind() == reflect.Interface {
		vrtn = vrtn.Elem()
  }
	
	for i := 0; i < len(props); i++ {
		path, found := ppaths[props[i]]
		if found {
		
			// Special case; if the field really actually doesn't exist, 
			// load the original record instead.
			var field_content, ferr = ioutil.ReadFile(path)
			if ferr != nil { return self.Get(key, []string{}) }
			
			var field, fterr = self.Type.FieldByName(props[i])
			if !fterr { return nil, n.Fail(1, "The given type does not have a property '%s'", props[i]) }
			
			var field_type = field.Type
			var field_value = self.Deserialize(field_content, field_type)
			
			var field_value_as_value = reflect.ValueOf(field_value)
			var field_pointer = vrtn.FieldByName(props[i])
			
			if field_value_as_value.Kind() == reflect.Ptr || field_value_as_value.Kind() == reflect.Interface {
				field_value_as_value = field_value_as_value.Elem()
			}
				
			field_pointer.Set(field_value_as_value)
		} else {
			return nil, n.Fail(1, "Invalid property list: '%s' was not an object property", props[i])
		}
	}
		
  return rtn, nil
}

func (self *recordSet) Unset(key int64) error {
	
	var kpath, fpaths = self.KeyPath(key)
  var err = os.Remove(kpath)
  if err != nil { return err }
  
  for _, fpath := range(fpaths) {
	  err = os.Remove(fpath)
	  if err != nil { return err }
  }
  
  return nil
}

func (self *recordSet) Keys(offset int, count int) []int64 {
  _ = self.Init()
  var sofar = 0
  var hits = 0
  var all = make([]int64, count)
  filepath.Walk(self.RecordsPath(), func(path string, info os.FileInfo, err error) error {
    if err != nil {
      n.Log("Error walking record set: %s", err)
    }
    if info != nil {
      if !info.IsDir() {
      	if strings.Count(info.Name(), ".") == 1 {
	        if sofar >= offset {
	          var _, file = filepath.Split(path)
	          var key int64
	          var _, err = fmt.Sscanf(file, "R%d.json", &key)
	          if err == nil && hits < count {
	            all[hits] = key
	            hits++
	          }
	        }
	        sofar++
        }
      }
    }
    return nil
  })
  all = all[:hits]
  return all
}

func (self *recordSet) Key() int64 {
  var rtn int64 = 0
  var found = false
  var tries = 0
  for !found {
    var key = rand.Int63()
    if self.ValidKey(key) {
      rtn = key
      found = true
    } else {
      tries++
    }
    if tries > 100 {
      n.Log("Unable to find unused key")
      break
    }
  }
  return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Serialization helpers  
 *============================================================================*/

// Is new key ok?
func (self *recordSet) ValidKey(key int64) bool {
  var rtn = true
  var path, _ = self.KeyPath(key)
  _, err := os.Stat(path)
  if err == nil {
    rtn = false
  }
  return rtn
}

// Init stuff
func (self *recordSet) Init() error {
  var err = os.MkdirAll(self.RecordsPath(), 0755)
  if err != nil {
    n.Log("Error creating record set: %s", err)
  }
  return err
}

// Return path for a records
func (self *recordSet) RecordsPath() string {
  var root, err = filepath.Abs("data")
  if err != nil {
    n.Log("Failed to resolve base path: %s", err)
  }
  var rtn = filepath.Join(root, self.Type.Name())
  return rtn
}

// Serialize to json
func (self *recordSet) Serialize(record interface{}) []byte {
  var raw, _ = json.Marshal(record) 
  return raw
}

// Deserialize from json
func (self *recordSet) Deserialize(value []byte, t reflect.Type) interface{} {
  var instance = reflect.New(t).Interface()
  json.Unmarshal(value, &instance) 
  return instance
}

// Decomponse a type into a set of key / value pairs
func (self *recordSet) Decompose(key string) map[string]string {
	var rtn = map[string]string {}
	var count = self.Type.NumField()
	for i := 0; i < count; i++ {
		var field = self.Type.Field(i)
  	var path = fmt.Sprintf("%s.%s.json", key, field.Name)
  	rtn[field.Name] = filepath.Join(self.RecordsPath(), path)
	}
	return rtn
}

// Return path for all items
func (self *recordSet) KeyPath(key int64) (string, map[string]string) {
  var path = fmt.Sprintf("R%020d", key)
	var items = self.Decompose(path)
  path = fmt.Sprintf("%s.json", path)
  return filepath.Join(self.RecordsPath(), path), items
}
 
/*============================================================================*
 * }}}
 *============================================================================*/