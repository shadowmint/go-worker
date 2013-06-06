package controllers

import "n"
import "net/http"
import "n/web"
import "worker/model/viewmodel"
import "worker/services"
import "worker/urls"
import "n/external/github.com/gorilla/schema"
import "html/template"

/*==========================================================================*
 * {{{ Public api 
 *==========================================================================*/

func newFragments() web.Controller {
  _ = n.Log
  var rtn = &fragments{
    service : services.New.Fragments(),
    decoder : schema.NewDecoder(),
    urls : urls.New.UrlHelper(),
  }
  return rtn
}

/*==========================================================================*
 * }}}
 *==========================================================================*/

/*==========================================================================*
 * {{{ IController api implementation
 *==========================================================================*/

// Public binding of this controller
func (self *fragments) Init(path string) {
  self.TemplatePath(path)
}

// Public binding of this controller
func (self *fragments) Binding() *web.Binding { 
  var rtn = web.Binding {
    Id : "fragments",
    Instance : self,
    Map : map[string] interface{} {
       "index" : (*fragments).index,
        "edit" : (*fragments).edit,
         "add" : (*fragments).add,
        "save" : (*fragments).save,
      "delete" : (*fragments).delete,
    },
  }
  return &rtn
}

/*==========================================================================*
 * }}}
 *==========================================================================*/

/*==========================================================================*
 * {{{ Actions
 *==========================================================================*/

type fragments struct {
  web.ControllerBase
  service services.Fragments
  decoder *schema.Decoder 
  urls *urls.UrlHelper
}

func (self *fragments) edit(c *web.Context) {

	var id = c.Url.Int64("id")
	var fragment, _, _ = self.service.Get(id)
  var m = self.model(fragment.Model())
  	
	// Update the last read index
	self.service.RecentUse(id)
	
  m["Raw"] = template.HTML(fragment.Content)
  self.draw(c.Writer, m, "fragments/edit.html")
}

func (self *fragments) add(c *web.Context) {
	self.service.Add("New fragment", "")
  http.Redirect(c.Writer, c.Request, self.urls.FragmentIndex(), http.StatusFound)
}

func (self *fragments) index(c *web.Context) {
	var m = self.model(self.service.List());
  self.draw(c.Writer, m, "fragments/index.html")
}

func (self *fragments) delete(c *web.Context) {

	// Parse input
	var id = c.Form.Int64("Id")
	
	// Delete record
	var f, _, _ = self.service.Get(id)
	self.service.Delete(f)
	
  http.Redirect(c.Writer, c.Request, self.urls.FragmentIndex(), http.StatusFound)
}

func (self *fragments) save(c *web.Context) {

	var model viewmodel.FragmentViewModel
	
	// Parse input
	c.Request.ParseForm()
	self.decoder.Decode(&model, c.Request.Form)
	
	// Update 
	self.service.Update(&model)
	
  http.Redirect(c.Writer, c.Request, self.urls.FragmentEdit(model.Id), http.StatusFound)
}

/*==========================================================================*
 * }}}
 *==========================================================================*/

/*==========================================================================*
 * {{{ Internal
 *==========================================================================*/
 
// Return a model for the given object
func (self *fragments) model(m ...interface{}) map[string]interface{} {

	// Add a model if we have one
	var rtn = make(map[string]interface{})
	if len(m) > 0 {
		rtn["Model"] = m[0]
	}
	
	// Add url helper
  rtn["Urls"] = self.urls
  
	return rtn
}

// Render a view, including the required layout file
func (self *fragments) draw(w http.ResponseWriter, model map[string]interface{}, template string) {

	// This is html
  w.Header().Add("Content-Type", "text/html")
  
  // ++layout files
  var templates = []string { 
  	template, 
  	"fragments/_layout.html", 
  	"fragments/_styles.html", 
	}
  
  self.View(templates, model, w)
}

/*==========================================================================*
 * }}}
 *==========================================================================*/
