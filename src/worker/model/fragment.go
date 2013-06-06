package model 

import "worker/model/viewmodel"
import "strings"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// A single content fragment 
type Fragment struct {
	Id int64
  Title string
  Content string
  Tags []string
}

// Import data from a view model
func (self *Fragment) Load(f *viewmodel.FragmentViewModel) {
	self.Id = f.Id
	self.Title = f.Title
	self.Content = f.Content;
	
	// Expand all tags
	var tags = strings.Split(f.AllTags, ",")
	for k, v := range tags {
		tags[k] = strings.Trim(v, " \n\r\t")
	}
	self.Tags = tags
}


// Convert this into a summary model 
func (self *Fragment) Summary() *viewmodel.FragmentSummaryViewModel {
	var rtn = &viewmodel.FragmentSummaryViewModel {
		Id : self.Id,
		Title : self.Title,
		Tags : self.Tags,
	}
	return rtn
}

// Convert this into a view model
func (self *Fragment) Model() *viewmodel.FragmentViewModel {
	var rtn = &viewmodel.FragmentViewModel {
		Id : self.Id,
		Title : self.Title,
		Content : self.Content,
		Tags : []string{},
		AllTags : "",
	}
	if (self.Tags != nil) && (len(self.Tags) > 0) {
		rtn.AllTags = strings.Join(self.Tags, ", ")
	}
	return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/
