package memory

import "n/test"
import "testing"
import "reflect"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type recordSetTestTypeB struct {
  Id string
  Value string
  T1 *recordSetTestTypeA 
  T2 recordSetTestTypeA 
}

type recordSetTestTypeA struct {
  X int
  Y int
}

type recordSetTest struct {
  assert test.Assert
  set RecordSet
}

func recordSetTestSetup(T *testing.T) (test.Assert, *recordSetTest) {
  var assert = test.New.Assert(T)
  var rtn = &recordSetTest{}

  var r recordSetTestTypeA
  rtn.set = New.RecordSet(reflect.TypeOf(r))

  return assert, rtn
}

func recordSetTestSetupB(T *testing.T) (test.Assert, *recordSetTest) {
  var assert = test.New.Assert(T)
  var rtn = &recordSetTest{}

  var r recordSetTestTypeB
  rtn.set = New.RecordSet(reflect.TypeOf(r))

  return assert, rtn
}

func recordSetTestTeardown(target *recordSetTest) {
  // No action required
}

/*============================================================================*

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_recordSet_can_create_instance(T *testing.T) {
  var a, c = recordSetTestSetup(T)

  a.NotNil(c.set, "Unable to create instance")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_set(T *testing.T) {
  var a, c = recordSetTestSetup(T)

  var record = recordSetTestTypeA {
    X : 10,
    Y : 10,
  }

  var err = c.set.Set(45, record)
  a.Nil(err, "Set failed")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_count(T *testing.T) {
  var a, c = recordSetTestSetup(T)

  var record = recordSetTestTypeA {
    X : 10,
    Y : 10,
  }

  c.set.Set(1, record)
  c.set.Set(2, record)
  c.set.Set(3, record)
  c.set.Set(4, record)
  c.set.Set(5, record)

  var count = c.set.Count()
  a.Equals(count, 5, "Invalid record count after insert")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_get(T *testing.T) {
  var a, c = recordSetTestSetup(T)

  var record = recordSetTestTypeA {
    X : 10,
    Y : 10,
  }

  c.set.Set(45, record)
  var nraw, err = c.set.Get(45, []string{})
  a.Nil(err, "Failed to get record")

  var nrecord = nraw.(*recordSetTestTypeA)
  a.Equals(nrecord.X, 10, "Invalid X value")
  a.Equals(nrecord.Y, 10, "Invalid Y value")
  
  nraw, err = c.set.Get(45, []string{"X"})
  a.Nil(err, "Failed to get record")

  nrecord = nraw.(*recordSetTestTypeA)
  a.Equals(nrecord.X, 10, "Invalid X value with restricted props")
  a.Equals(nrecord.Y, 0, "Invalid Y value with restricted props")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_get_new_keys(T *testing.T) {
  var a, c = recordSetTestSetup(T)

  var prev_key = c.set.Key()
  for i := 0; i < 25; i++ {
    var key = c.set.Key()
    a.True(prev_key != key, "Invalid key")
  }

  recordSetTestTeardown(c)
}

func Test_recordSet_can_delete_records(T *testing.T) {
  var a, c = recordSetTestSetup(T)
  
  var record = recordSetTestTypeA { X : 10, Y : 10, }

  for i := 0; i < 25; i++ {
    record.X = i
    c.set.Set(int64(i), record)
  }

  var count = c.set.Count()
  a.Equals(count, 25, "Invalid initial count")

  c.set.Unset(5)
  c.set.Unset(10)
  c.set.Unset(15)

  count = c.set.Count()
  a.Equals(count, 22, "Invalid count after delete")

  var nraw, err = c.set.Get(5, []string { "X", "Y" })
  a.Nil(nraw, "Invalid return on invalid key")
  a.NotNil(err, "Invalid error on invalid key")

  nraw, err = c.set.Get(6, []string{})
  a.NotNil(nraw, "Invalid return on valid key")
  a.Nil(err, "Invalid error on valid key")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_get_key_range(T *testing.T) {
  var a, c = recordSetTestSetup(T)
  
  var record = recordSetTestTypeA { X : 10, Y : 10, }

  for i := 0; i < 25; i++ {
    record.X = i
    c.set.Set(int64(i), record)
  }

  c.set.Unset(5)
  c.set.Unset(10)
  c.set.Unset(15)
  c.set.Unset(20)

  var keys1 = c.set.Keys( 0, 5)
  var keys2 = c.set.Keys( 5, 5)
  var keys3 = c.set.Keys(10, 5)
  var keys4 = c.set.Keys(15, 5)
  var keys5 = c.set.Keys(20, 5)

  var sample = make(map[int64]bool)
  for i := 0; i < len(keys1); i++ { sample[keys1[i]] = true; }
  for i := 0; i < len(keys2); i++ { sample[keys2[i]] = true; }
  for i := 0; i < len(keys3); i++ { sample[keys3[i]] = true; }
  for i := 0; i < len(keys4); i++ { sample[keys4[i]] = true; }
  for i := 0; i < len(keys5); i++ { sample[keys5[i]] = true; }

  a.Equals(len(sample), 21, "Invalid number of unique keys after delete")

  recordSetTestTeardown(c)
}

func Test_recordSet_can_persist_complex_types(T *testing.T) {
  var a, c = recordSetTestSetupB(T)

  var record = recordSetTestTypeB {
    Id : "ID HERE",
    Value : "VALUE HERE",
    T1 : &recordSetTestTypeA { X : 10, Y : 20 },
    T2 : recordSetTestTypeA { X : 30, Y : 40 },
  }

  // Notice we pass in a pointer this time
  c.set.Set(10, &record)

  var nraw, err = c.set.Get(10, []string{})

  a.Nil(err, "Error after valid get call")
  a.NotNil(nraw, "Invalid return after valid get call")

  var nrecord = nraw.(*recordSetTestTypeB)
  a.Equals(nrecord.T1.X, 10, "Invalid T1 X value")
  a.Equals(nrecord.T1.Y, 20, "Invalid T1 Y value")
  a.Equals(nrecord.T2.X, 30, "Invalid T2 X value")
  a.Equals(nrecord.T2.Y, 40, "Invalid T2 Y value")

  recordSetTestTeardown(c)
}


/*============================================================================*
 * }}}
 *============================================================================*/
