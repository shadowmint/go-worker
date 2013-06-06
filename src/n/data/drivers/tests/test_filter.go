package tests

import "testing"

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_store_can_filter_results(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  // Load a bunch of records into the store
  for i := 0; i < 20; i++ {
    var record = storeTestTypeA { A : i }
    c.Set("storeTestTypeA", record)
  }

  // Fetch a few pages
  var records, err = c.All("storeTestTypeA", 5, 5)
  a.Nil(err, "Error fetching records")
  a.Equals(len(records), 5, "Invalid count in record set")

  // Fetch an incomplete page
  records, err = c.All("storeTestTypeA", 15, 10)
  a.Nil(err, "Error fetching records")
  a.Equals(len(records), 5, "Invalid count in record set")

  // Fetch a filtered page
  records, err = c.Filter("storeTestTypeA", 5, 5, []string{}, func(r interface{}) bool {
    var data = r.(*storeTestTypeA)
    return (data.A % 2) == 0
  })
  a.Nil(err, "Error fetching records")
  a.Equals(len(records), 5, "Invalid count in record set")
  for _, v := range records {
    a.True(v.(*storeTestTypeA).A % 2 == 0, "Invalid value on record")
  }

  // Fetch entire filtered set
  var count, _ = c.Count("storeTestTypeA")
  records, err = c.Filter("storeTestTypeA", 0, count, []string{}, func(r interface{}) bool {
    var data = r.(*storeTestTypeA)
    return (data.A % 2) == 0
  })
  a.Nil(err, "Error fetching records")
  a.Equals(len(records), 10, "Invalid count in record set")
  for _, v := range records {
    a.True(v.(*storeTestTypeA).A % 2 == 0, "Invalid value on record")
  }

  // Fetch a filtered page off the end
  records, err = c.Filter("storeTestTypeA", 5, 10, []string{}, func(r interface{}) bool {
    var data = r.(*storeTestTypeA)
    return (data.A % 2) == 0
  })
  a.Nil(err, "Error fetching records")
  a.Equals(len(records), 5, "Invalid count in record set")
  for _, v := range records {
    a.True(v.(*storeTestTypeA).A % 2 == 0, "Invalid value on record")
  }

  storeTestTeardown(runner, c)
}

func Test_store_can_filter_odd_queries(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  // Load a bunch of records into the store
  for i := 0; i < 20; i++ {
    var record = storeTestTypeA { A : i }
    c.Set("storeTestTypeA", record)
  }

  // Fetch a filtered page off the end
  var results, err = c.Filter("storeTestTypeA", 5, 10, []string{}, nil)
  a.NotNil(err, "No error on stupid filter")

  // Fetch a filtered page off the end
  results, err = c.Filter("storeTestTypeA", 5, 10, []string{}, func(record interface{}) bool {
    return false
  })
  a.Nil(err, "No error on stupid filter")
  a.Equals(len(results), 0, "Invalid return on empty request")

  storeTestTeardown(runner, c)
}

func (self *tests) Test_store_can_filter_query_by_properties(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  // Load a bunch of records into the store
  for i := 0; i < 20; i++ {
    var record = storeTestTypeA { A : i + 1, D : i % 2 == 0, E : "Hello World" }
    c.Set("storeTestTypeA", record)
  }

  var results, err = c.Filter("storeTestTypeA", 5, 10, nil, func(record interface{}) bool { return true; })
  a.Nil(err, "Error on property filter with nil properties")

  // Fetch a filtered item
  results, err = c.Filter("storeTestTypeA", 5, 10, []string{ "A", "D" }, func(record interface{}) bool {
  	var data = record.(*storeTestTypeA)	
  	a.True(data.A > 0, "Failed to populate property 'A'")
  	a.Equals(data.E, "", "Populated bad property 'E'")
  	return data.D;
  })
  a.Nil(err, "Error on property filter request")
  a.Equals(len(results), 5, "Invalid return on property filtered query")

  storeTestTeardown(runner, c)
}

/*============================================================================*
 * }}}
 *============================================================================*/
