package services

import "n/test"
import "testing"
import "n"
import "n/data"
import "n/data/drivers/memory"
import rf "reflect"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

func fragmentTestSetup(T *testing.T) (test.Assert, Fragments) {

  // Setup dependencies
	var r = n.New.Resolver()
	r.Register(rf.TypeOf((*data.Driver)(nil)), memory.New.Driver)
	r.Register(rf.TypeOf((*data.Store)(nil)), data.New.Store)
  
  var assert = test.New.Assert(T)
  var instance = New.Fragments()
  return assert, instance
}

func fragmentTestTeardown(target Fragments) {
	n.New.Resolver().Clear()
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_fragment_can_create_instance(T *testing.T) {
  var a, i = fragmentTestSetup(T)

  a.NotNil(i, "Unable to create instance")

  fragmentTestTeardown(i)
}

func Test_fragment_can_insert(T *testing.T) {
  var a, i = fragmentTestSetup(T)

	var f, err = i.Add("Title", "Value")
	a.Nil(err, "Failed to add instance")
	a.NotNil(f, "Invalid return on add call")

  fragmentTestTeardown(i)
}

func Test_fragment_can_get(T *testing.T) {
  var a, i = fragmentTestSetup(T)

	var f, _ = i.Add("Title", "Value")
	var key = f.Id
	
	var fn, ft, ferr = i.Get(key)
	a.Nil(ferr, "Failed to get instance")
	a.Nil(ft, "Tags on new fragment")
	a.NotNil(fn, "Invalid fragment on get call")
	a.Equals(fn.Title, "Title", "Fetch on record returns invalid data")
	a.Equals(fn.Content, "Value", "Fetch on record returns invalid data")

  fragmentTestTeardown(i)
}

func Test_fragment_can_delete(T *testing.T) {
  var a, i = fragmentTestSetup(T)

	var f, _ = i.Add("Title", "Value")
	
	var err = i.Delete(f)
	a.Nil(err, "Failed to delete instance")
	
	var fn, _, ferr = i.Get(f.Id)
	a.NotNil(ferr, "Missing error on call")
	a.Nil(fn, "Read deleted record")

  fragmentTestTeardown(i)
}

// TODO: Test deleting deletes tags
// TODO: Test deleting deletes star entries

func Test_fragment_can_update(T *testing.T) {
  var a, i = fragmentTestSetup(T)

	var f, _ = i.Add("Title", "Value")
	f.Title = "Hello"
	f.Content = "World"
  var fm = f.Model()
	
	var err = i.Update(fm)
	a.Nil(err, "Invalid error on update")
	
	var fn, _, _ = i.Get(f.Id)
	a.Equals(fn.Title, "Hello", "Fetch on record returns invalid data")
	a.Equals(fn.Content, "World", "Fetch on record returns invalid data")

  fragmentTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
