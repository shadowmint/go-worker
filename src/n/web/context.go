package web

import "n/external/github.com/gorilla/schema"
import "net/http"
import "fmt"
import "n"

/*============================================================================*
 * {{{ Public api
 *============================================================================*/

// Form helpers, for lazy and no view model.
// Otherwise use utils.Decode().
type Form interface {

	// Read a form value as an int64 value.
	Int64(id string) int64

	// Read a form value as an string value.
	String(id string) string
}

// Url helpers
type Url interface {

	// Read a form value as an int64 value.
	Int64(id string) int64

	// Read a form value as an string value.
	String(id string) string
}

// Utils
type Utils interface {

	// Redirect the page to some url.
	Redirect(url string)
	
	// Parse the form and bind it to the given model.
	// Only public attributes are honoured.
  // Pass the &blah to bind to blah.
	Decode(model interface{})
}

// Context object
type Context struct {

  // Http response object
  Writer http.ResponseWriter
  
  // Http request object
  Request *http.Request
  
  // Variables from the url, eg. from /{controller}/{action}/{id}
  vars map[string]string
  
  // Form deserializer 
  decoder *schema.Decoder
  
  // Url handle
  Url Url
  
  // Form handle
  Form Form
  
  // Utils handle
  Utils Utils
}

// Create an instance
func newContext(w http.ResponseWriter, r *http.Request, vars map[string]string) *Context {
	var rtn = &Context {
		Writer : w,
		Request : r,
		vars : vars, 
    decoder : schema.NewDecoder(),
	}
	rtn.Url = &url { rtn }
	rtn.Form = &form { rtn }
	rtn.Utils = &utils { rtn }
	return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*=============================================================================*
 * {{{ Utils impl
 *=============================================================================*/
 
type utils struct {
	c *Context
}

// Decode the things; pass the &blah to bind to blah.
func (self *utils) Decode(model interface{}) {
	self.c.Request.ParseForm()
  self.c.decoder.Decode(model, self.c.Request.Form)
}

// Shortcut for redirect from a context 
func (self *utils) Redirect(url string) {
  http.Redirect(self.c.Writer, self.c.Request, url, http.StatusFound)
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*=============================================================================*
 * {{{ Form impl
 *=============================================================================*/
 
type form struct {
	c *Context
}

// Return an int64 form value or 0
func (self *form) Int64(id string) int64 {
  var rtn int64 = 0
  self.c.Request.ParseForm()
  var raw, found = self.c.Request.Form[id]
  if found {
    var _, err = fmt.Sscanf(raw[0], "%d", &rtn)
    if err != nil {
      n.Log("Failed reading incoming form param: %s is not int64", raw[0])
      n.Log(err.Error())
    }
  }
  return rtn
}

// Return a string form value or ""
func (self *form) String(id string) string {
  var rtn string = ""
  self.c.Request.ParseForm()
  var raw, found = self.c.Request.Form[id]
  if found {
  	rtn = raw[0]
  }
  return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*=============================================================================*
 * {{{ Form impl
 *=============================================================================*/
 
type url struct {
	c *Context
}

// Return an int64 url value
func (self *url) Int64(id string) int64 {
  var rtn int64 = 0
  var raw, found = self.c.vars[id]
  if found {
    var _, err = fmt.Sscanf(raw, "%d", &rtn)
    if err != nil {
      n.Log("Failed reading incoming url param: %s is not int64", raw[0])
      n.Log(err.Error())
    }
  }
  return rtn
}

// Return a string form value or ""
func (self *url) String(id string) string {
  var rtn string = ""
  self.c.Request.ParseForm()
  var raw, found = self.c.vars[id]
  if found {
  	rtn = raw
  }
  return rtn
}

/*=============================================================================*
 * }}}
 *=============================================================================*/