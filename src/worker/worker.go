package worker 

import "fmt"
import "net/http"
import "worker/controllers"
import "n"
import "n/web"
import "n/resources"
import "time"
import rf "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// Config keys
const (
  Static = "StaticFolder"
  Views = "ViewsFolder"
  Port = "ListenPort"
)

// Load controllers
func Load(config resources.Config) {
  bind(controllers.New.Fragments, config)
}

// Invoke action on controller
func Resolve(controller string, action string, vars map[string]string, w http.ResponseWriter, r *http.Request) {
  var binding, found = actions[controller]
  if found {
    var handler, found = binding.Map[action]
    if (found) {
      var c = web.New.Context(w, r, vars)
      var v = rf.ValueOf(handler)
      var args = []rf.Value { rf.ValueOf(binding.Instance), rf.ValueOf(c) }
      v.Call(args) 
      return
    }
  }
  fmt.Fprintf(w, "Not found")
  n.Log(": No match for %s::%s", controller, action)
}

// Run the tasks system
func Tasks(config resources.Config) {
  go func() {
    for {
      // TODO: run update tasks here.
      time.Sleep(5000000000)
    }
  }()
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Internal
 *============================================================================*/

// Set of bindings for controller ids and params 
var actions = map [string] *web.Binding {}

// Attach a single binding 
func bind(src func() web.Controller, config resources.Config) {
  var c = src()
  var b = c.Binding()

  // Attach template path
  var path = config.Get(Views)
  c.Init(path)

  var _, found = actions[b.Id]
  if found {
    panic(fmt.Sprintf("Conflict! Duplicate binding controller name '%s'", b.Id))
  }
  actions[b.Id] = b
}

/*============================================================================*
 * }}}
 *============================================================================*/
