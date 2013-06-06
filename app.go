package main 

import "n/external/github.com/gorilla/mux"
import "github.com/tebeka/desktop"
import "net/http"
import "worker/build"
import "worker"
import "path"
import "fmt"
import "n"
import "time"
import "n/data"
import "n/data/drivers/file"
import "n/resources"
import rf "reflect"

/*============================================================================*
 * {{{ Main
 *============================================================================*/

func main() {
  
  // Resolver setup
  bind();

  // Load config
  var config = loadConfig()
  worker.Load(config)

  // Setup background worker
  worker.Tasks(config)

  // Controllers
  var r = mux.NewRouter()
  r.HandleFunc("/{controller}/{action}/{id}/", dispatch)
  r.HandleFunc("/{controller}/{action}/{id}", dispatch)
  r.HandleFunc("/{controller}/{action}/", dispatch)
  r.HandleFunc("/{controller}/{action}", dispatch)
  r.HandleFunc("/{controller}/", dispatch)
  r.HandleFunc("/{controller}", dispatch)
  r.HandleFunc("/", dispatch)
  http.Handle("/", r);

  // Static content 
  var prefix = "/assets/"
  var path = config.Get(worker.Static)
  var h = http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
  http.Handle(prefix, h)
  fmt.Printf(" Templates: '%s'\n", config.Get(worker.Views))
  fmt.Printf("    Static: '%s'\n", path)
  fmt.Printf("Static URL: '%s'\n", prefix)

  // Launch the browser after a moment
  go func() {
    time.Sleep(1000 * 1000 * 5)
    var url = "http://127.0.0.1" + config.Get(worker.Port)
    desktop.Open(url)
  }()

  // Start server~
  var uri = "localhost" + config.Get(worker.Port)
  fmt.Printf(" Listening: %s\n", uri)
  var err = http.ListenAndServe(uri, nil)
  if err != nil {
    fmt.Printf("Failed: %s\n", err.Error())
  }
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Internal
 *============================================================================*/

// Resolve an incoming request to an appropriate action.
func dispatch(w http.ResponseWriter, r *http.Request) {
  var vars = mux.Vars(r) 

  var controller, cfound = vars["controller"]
  if !cfound { controller = "fragments" }

  var action, afound = vars["action"]
  if !afound { action = "index" }

  worker.Resolve(controller, action, vars, w, r)
}

// Load config file and validate all required keys are present
func loadConfig() resources.Config {

  var config = resources.New.FileConfig()
  if build.DEBUG {
    config.Load(path.Join("content", "app.config"))
  } else {
    config.Load(path.Join(n.ExecPath(), "content", "app.config"))
  }

  config.Require(worker.Static, "Missing static assets folder path")
  config.Require(worker.Views, "Missing templates folder path")
  config.Require(worker.Port, "Missing HTTP binding port")
  if !config.Validate() {
    panic("Unable to start application. Bad config.")
  }

  // Correct paths to be absolute
  if !build.DEBUG {
    config.FixPath(worker.Static, n.ExecPath())
    config.FixPath(worker.Views, n.ExecPath())
  }

  return config
}

// Bind any injectable types that need to be configured 
func bind() {
  var r = n.New.Resolver()
  { var t data.Store; r.Register(rf.TypeOf(&t), data.New.Store); }
  { var t data.Driver; r.Register(rf.TypeOf(&t), file.New.Driver); }
}

/*============================================================================*
 * }}}
 *============================================================================*/
