package resources 

import "n"
import "path"

/*================================================================================*
 * {{{ Public api
 *================================================================================*/

// For common configuration settings
type Config interface {
  Parse(raw string) error
  Load(uri string) error
  Set(key string, value string) 
  Get(key string) string
  Has(key string) bool
  Default(key string, value string)
  Require(key string, msg string)

  // Validate the config and log errors
  Validate() bool

  // Convert a value to an absolute path value using cwd
  FixPath(key string, prefix string)
}

/*================================================================================*
 * }}}
 *================================================================================*/

/*================================================================================*
 * {{{ Default Config implementation
 *================================================================================*/

type configBase struct {
  values map[string] string
  defaults map[string] string
  requires map[string] string
}

func (self *configBase) Default(key string, value string) {
  self.defaults[key] = value
}

func (self *configBase) Require(key string, msg string) {
  self.requires[key] = msg
}

func (self *configBase) Validate() bool {
  var failed = false
  for k, _ := range self.requires {
    if !self.Has(k) {
      var _, found = self.defaults[k]
      if !found {
        n.Log("Missing config key: %s", k)
        failed = true
      }
    }
  }
  return !failed
}

func (self *configBase) Set(key string, value string) {
  self.values[key] = value
}

func (self *configBase) Get(key string) string {
  var rtn, found = self.values[key]
  if (!found) {
    rtn = ""
  }
  return rtn
}

// Check we have a value
func (self *configBase) Has(key string) bool {
  var _, found = self.values[key]
  return found
}

// Convert 
func (self *configBase) FixPath(key string, prefix string) {
  if self.Has(key) {
    var value = self.Get(key)
    if !path.IsAbs(value) {
      var new_value = path.Join(prefix, value)
      self.Set(key, new_value)
    }
  }
}

/*================================================================================*
 * }}}
 *================================================================================*/
