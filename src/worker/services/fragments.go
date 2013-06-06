package services

import "worker/model/viewmodel"
import "worker/model"
import "n/data"
import "n"
import rf "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

type Fragments interface {

	// Get a list of fragments to display
	List() *viewmodel.FragmentsListViewModel
	
	// Update a fragment from a view model
	Update(f *viewmodel.FragmentViewModel) error
	
	// Get a fragment by id
	Get(key int64) (*model.Fragment, *model.Tags, error) 
  
  // Add a new fragment
  Add(title string, value string) (*model.Fragment, error) 
  
  // Delete a fragment by id
	Delete(fragment *model.Fragment) error 
	
	// Mark a fragment as recently used
  RecentUse(key int64)
}

func newFragments(args ...interface{} /* Store */) Fragments {

	var raw = n.Resolve(rf.TypeOf((*data.Store)(nil)), 0, args...)
	if raw == nil {
		n.Log("Failed to create fragment service: No Store binding")
		return nil
	}
	
  var rtn = &fragments { 
 	  db : raw.(data.Store),
 	  index : newIndex(),
 	}
  
  // Register types we're going to have to work with.
  {      var t model.Tags; rtn.db.Register(rf.TypeOf(&t), "worker.Tags"); }
  {     var t model.Stars; rtn.db.Register(rf.TypeOf(&t), "worker.Stars"); }
  {  var t model.TagCloud; rtn.db.Register(rf.TypeOf(&t), "worker.TagCloud"); }
  {  var t model.Fragment; rtn.db.Register(rf.TypeOf(&t), "worker.Fragment"); }
  { var t model.StarIndex; rtn.db.Register(rf.TypeOf(&t), "worker.StarIndex"); }
  {     var t model.Index; rtn.db.Register(rf.TypeOf(&t), "worker.Index"); }
  
  return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Api implementation
 *============================================================================*/

type fragments struct {
	db data.Store
	index Index
}

// Update, expand tags
func (self *fragments) Update(f *viewmodel.FragmentViewModel) error {
	var fragment, _, err = self.Get(f.Id)
	if err != nil {
		return n.LogFail(1, "Unable to update fragment; no match for id '%d'", f.Id)
	}
	
	fragment.Load(f)
	err = self.update(fragment)
	return err
}

// Get a list of fragments to display
func (self *fragments) List() *viewmodel.FragmentsListViewModel {
	var count = 20
  var recent, _ = self.findByRecentUse(0, 5) // Last 5 recent items
  var items, _ = self.fetch(0, count)
  var more, _ = self.fetch(count, 1)
  var rtn = viewmodel.FragmentsListViewModel {
  	Items : items,
  	Recent : recent,
  	More : len(more) > 0,
  };
  return &rtn;
}

// Get a fragment and its tags
func (self *fragments) Get(key int64) (*model.Fragment, *model.Tags, error) {

	var rtn, err = self.db.Get("worker.Fragment", key)
	if err != nil {
		return nil, nil, err
	} else {
		rtn.(*model.Fragment).Id = key
	}
	
	var tags, terr = self.db.Filter("worker.Tags", 0, 1, []string{}, func(record interface{}) bool {
		var data = record.(*model.Tags)
		return data.Fragment == key 
	})
	if terr != nil {
		return nil, nil, terr
	}
	var trtn *model.Tags = nil
	if len(tags) > 0 {
		trtn = tags[0].(*model.Tags)
	}

	return rtn.(*model.Fragment), trtn, nil
}

// Add a new fragment
func (self *fragments) Add(title string, value string) (*model.Fragment, error) {
	var f = model.Fragment { 0, title, value, nil }
	var key, err = self.db.Set("worker.Fragment", f) 
	if err != nil {
		return nil, err
	}
	
	// Save own key, because we're going to use it later
	f.Id = key
	self.db.Set("worker.Fragment", f, f.Id) 
	self.RecentUse(f.Id)
	
	return &f, nil
}

// Delete a fragment
func (self *fragments) Delete(fragment *model.Fragment) error {
	if fragment == nil {
		return n.LogFail(1, "Invalid request to delete nil record (bad id?)")
	}
	
	var err = self.db.Delete("worker.Fragment", fragment.Id)
	if err != nil {
		return n.LogFail(1, "Failed to delete record: %s", err.Error())
	}
	
	// If we got here, remove that id from all indexes
	self.index.Delete(fragment.Id)
	
	return nil
}

// Mark a fragment to be recently used
func (self *fragments) RecentUse(key int64) {
	if key != 0 {
		self.index.Add("LastRead", 10, key)
	}
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Internal functions
 *============================================================================*/

// Update a fragment
func (self *fragments) update(fragment *model.Fragment) error {
	var _, err = self.db.Set("worker.Fragment", fragment, fragment.Id) 
	return err
}

// Get a set of fragments and tags
func (self *fragments) fetch(offset int, count int) ([]viewmodel.FragmentSummaryViewModel, error) {
	var all, err = self.db.All("worker.Fragment", offset, count)
	if err != nil {
		return nil, err
	}
	var rtn = make([]viewmodel.FragmentSummaryViewModel, len(all), len(all))
	var index = 0
	for k, v := range all {
		var data = v.(*model.Fragment)
		data.Id = k
		rtn[index] = *data.Summary()
		index++
	}
	return rtn, err
}

// Find by recent use
func (self *fragments) findByRecentUse(offset int, count int) ([]viewmodel.FragmentSummaryViewModel, error) {

	var keys = self.index.Find("LastRead", offset, count)
	
	var rcount = 0
	for _, v := range keys {
		if v != 0 {
			rcount++
		}
	}
	
	var rtn = make([]viewmodel.FragmentSummaryViewModel, rcount, rcount)
	rcount = 0
	for _, v := range keys {
		if v != 0 {
			var fragment, _, err = self.Get(v)
			if err != nil {
				return nil, err
			}
			rtn[rcount] = *fragment.Summary()
			rcount++
		}
	}
	return rtn, nil
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
