package web

import "n"
import "net/http"
import "io/ioutil"
import "html/template"
import _path "path"

/*============================================================================*
 * {{{ Public api
 *============================================================================*/

type Controller interface {

  // Return the url binding for the controller
  Binding() *Binding

  // Set the template path for the controller
  Init(path string)
}

type ControllerBase struct {

  // Config for this controller
  templatePath string

  // Template container
  templates *template.Template
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ ControllerBase api
 *============================================================================*/

// Inject configuration object
func (self *ControllerBase) TemplatePath(path string) {
  self.templatePath = path
}

// Load a template and attach the given name to it.
func (self *ControllerBase) loadTemplate(path string, id string) error {

  var template_path = _path.Join(self.templatePath, path)

  var raw, re = ioutil.ReadFile(template_path)
  if re != nil { 
    // TODO: Controller common error handling.
    n.Log("\nError: failed to load template: %s", path)
    n.Log("%s", re.Error())
    return n.Fail(1, "Invalid template")
  }

  var t *template.Template = self.templates
  if t == nil {
    var t, te = template.New(id).Parse(string(raw))
    if te != nil { 
      // TODO: Controller common error handling.
      n.Log("\nError: invalid template: %s", path)
      n.Log("%s", te.Error())
      return n.Fail(1, "Invalid template")
    } else {
      self.templates = t
    }
  } else {
    var _, te = self.templates.New(id).Parse(string(raw))
    if te != nil { 
      // TODO: Controller common error handling.
      n.Log("\nError: invalid template: %s", path)
      n.Log("%s", te.Error())
      return n.Fail(1, "Invalid template")
    }
  }

  // DEBUG
  // n.Log("Registered template id: %s -> %s", id, template_path)
  return nil
}

// Render a template with a layout file.
// Use {{define ""}} .. {{end}} in the template.
func (self *ControllerBase) View(templates []string, data map[string]interface{}, w http.ResponseWriter) {

  self.templates = nil
  for _, v := range templates {
    var err = self.loadTemplate(v, v)
    if err != nil { return; }
  }

  var ee = self.templates.Execute(w, data)
  if ee != nil { 
    // TODO: Controller common error handling.
    n.Log("\nError: failed executing templates: %s", templates[0])
    n.Log("%s", ee.Error())
  }
}

/*============================================================================*
 * }}}
 *============================================================================*/
