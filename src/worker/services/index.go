package services

import "worker/model"
import "n/data"
import "n"
import rf "reflect"

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// Useful indexing apis
type Index interface {

	// Index a key by some property and put this at the top of the index
  Add(key string, limit int, id int64) 
  
  // Remove a key from all indexes
  Delete(id int64)
  
  // Remove an entire index
	DeleteIndex(key string)
  
  // Find a subset of the keys in an index.
  //
  // The returned set is only valid keys; for a raw set of 
  // [0, 1, 34, 3, 33] the output is [1, 34, 3, 33] 
  // 
  // If count was 3, the output would be [1, 34, 3].
  //
  // @param offset The offset into the index records to scan
  // @param count The count of results to return, at most.
  Find(key string, offset int, count int) []int64
}

func newIndex(args ...interface{} /* Store */) Index {

	var raw = n.Resolve(rf.TypeOf((*data.Store)(nil)), 0, args...)
	if raw == nil {
		n.Log("Failed to create index service: No Store binding")
		return nil
	}
	
  var rtn = &index { 
 	  db : raw.(data.Store),
 	}
  
  // Register types we're going to have to work with.
  { var t model.Index; rtn.db.Register(rf.TypeOf(&t), "worker.Index"); }
  
  return rtn
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Api implementation
 *============================================================================*/

// Maximum index count; arbitrary large number
const index_MAX = 100

type index struct {
	db data.Store
}

func (self *index) Find(key string, offset int, count int) []int64 {

	var index = self.get(key, index_MAX)
	
	var items = 0
	var skip = 0
	for i := 0; i < index.Size && items <= count; i++  {
		if index.Ids[i] != 0 {
			skip++
			if skip > offset {
				items++
			}
		}
	}
	
	var rtn = make([]int64, items, items)
	
	items = 0
	skip = 0
	for i := 0; i < index.Size && items < count; i++  {
		if index.Ids[i] != 0 {
			skip++
			if skip > offset {
				rtn[items] = index.Ids[i]
				items++
			}
		}
	}
	return rtn
}

func (self *index) Add(key string, limit int, id int64) {

	// Attempt to find a record, and create one if we can't
	var index = self.get(key, limit)
	
	// Add this id to the set, keep at most N items, if not already present
	var found = false
	for i := 0; i < index.Size; i++ {
		if index.Ids[i] == id {
			found = true
			for j := i; j >= 1; j-- {
				index.Ids[j] = index.Ids[j - 1]
			}
			break
		}
	}
	if !found {	
		for i := (index.Size - 1); i >= 1; i-- {
			index.Ids[i] = index.Ids[i - 1]
		}
	}
	index.Ids[0] = id
	
	// Save index
	self.set(index)
}

// Arbitrary limit; max of 100 index or we're doing something weird.
func (self *index) Delete(id int64) {
	var items, err = self.db.All("worker.Index", 0, 100)
	if err != nil {
		n.Log("Failed to remove index: %s", err.Error())
	}
	
	for _, v := range(items) {
		var index = v.(*model.Index)
		var tmp = make([]int64, index.Size, index.Size)
		var update = false
		var offset = 0
		for j := 0; j < index.Size; j++ {
			if index.Ids[j] != id {
				tmp[offset] = index.Ids[j]
				offset++
			} else {
				update = true;
			}
		}
		if update {
			index.Ids = tmp
			self.set(index)
		}
	}
}

// Arbitrary limit; max of 100 index or we're doing something weird.
func (self *index) DeleteIndex(key string) {
	var index = self.get(key, index_MAX)
	self.db.Delete("worker.Index", index.Id)
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * Private apis
 *============================================================================*/
 
// Save record
func (self *index) set(index *model.Index) {
	if (index.Id == 0) {
		self.db.Set("worker.Index", index)
	} else {
		self.db.Set("worker.Index", index, index.Id)
	}
}

// Get or create record
func (self *index) get(key string, size int) *model.Index {
	var index *model.Index
	var set, ierr = self.db.Filter("worker.Index", 0, 1, []string{}, func(record interface{}) bool {
		var data = record.(*model.Index)
		return data.Name == key
	})
	if ierr != nil {
		n.Log("Index error: %s", ierr.Error())
	}
	if len(set) == 0 {
		index = &model.Index {
			Name : key,
			Size : size,
			Ids : make([]int64, size, size),
		}
	} else {
		for k, v := range set {
		  index = v.(*model.Index)
		  index.Id = k
		  break
		}
	}
	return index
}
 
/*============================================================================*
 * }}}
 *============================================================================*/
