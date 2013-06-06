package services

import "n/test"
import "testing"
import "n"
import "n/data"
import "n/data/drivers/memory"
import rf "reflect"

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_index_can_create_instance(T *testing.T) {
  var a, i = indexTestSetup(T)

  a.NotNil(i, "Unable to create instance")

  indexTestTeardown(i)
}

func Test_index_can_add_to_index(T *testing.T) {
  var a, i = indexTestSetup(T)

	i.Add("my_index_1", 10, 1)
	i.Add("my_index_1", 10, 2)
	i.Add("my_index_1", 10, 3)
	i.Add("my_index_1", 10, 4)
	
	_ = a
  indexTestTeardown(i)
}

func Test_index_can_read_index(T *testing.T) {
  var a, i = indexTestSetup(T)

	i.Add("my_index_2", 5, 1)
	i.Add("my_index_2", 5, 7)
	i.Add("my_index_2", 5, 6)
	i.Add("my_index_2", 5, 5)
	i.Add("my_index_2", 5, 4)
	i.Add("my_index_2", 5, 3)
	i.Add("my_index_2", 5, 1)
	
	var set = i.Find("my_index_2", 0, 15)
	
	a.Equals(len(set), 5, "Invalid index length")
	a.Equals(set[0], int64(1), "Invalid index value")
	a.Equals(set[1], int64(3), "Invalid index value")
	a.Equals(set[2], int64(4), "Invalid index value")
	a.Equals(set[3], int64(5), "Invalid index value")
	a.Equals(set[4], int64(6), "Invalid index value")
	
	_ = a
  indexTestTeardown(i)
}

func Test_index_can_delete_index_item(T *testing.T) {
  var a, i = indexTestSetup(T)

	i.Add("my_index_2", 15, 1)
	i.Add("my_index_2", 15, 7)
	i.Add("my_index_2", 15, 6)
	i.Add("my_index_2", 15, 5)
	i.Add("my_index_2", 15, 4)
	i.Add("my_index_2", 15, 3)
	i.Add("my_index_2", 15, 1)
	
	i.Add("my_index_1", 15, 99)
	i.Add("my_index_1", 15, 107)
	i.Add("my_index_1", 15, 103)
	i.Add("my_index_1", 15, 1)
	
	var set1 = i.Find("my_index_1", 0, 15)
	a.Equals(len(set1), 4, "Invalid index length")
	
	var set2 = i.Find("my_index_2", 0, 15)
	a.Equals(len(set2), 6, "Invalid index length")
	
	i.Delete(1)
	
	set1 = i.Find("my_index_1", 0, 15)
	a.Equals(len(set1), 3, "Invalid index length")
	
	set2 = i.Find("my_index_2", 0, 15)
	a.Equals(len(set2), 5, "Invalid index length")
	
  indexTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
 
/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

func indexTestSetup(T *testing.T) (test.Assert, Index) {
	
  // Setup dependencies
	var r = n.New.Resolver()
	r.Register(rf.TypeOf((*data.Driver)(nil)), memory.New.Driver)
	r.Register(rf.TypeOf((*data.Store)(nil)), data.New.Store)
  
  var assert = test.New.Assert(T)
  var instance = New.Index()
  
  return assert, instance
}

func indexTestTeardown(target Index) {

  // Delete working indexes
  target.DeleteIndex("my_index_1")
  target.DeleteIndex("my_index_2")
  
	n.New.Resolver().Clear()
}

/*============================================================================*
 * }}}
 *============================================================================*/
