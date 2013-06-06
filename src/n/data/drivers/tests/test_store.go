package tests

import "testing"
import "reflect"

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_store_can_create_instance(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  a.NotNil(c, "Unable to create instance")

  storeTestTeardown(runner, c)
}

func Test_store_can_save_record(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  var record = storeTestTypeA {
    A : 1,
    B : 2,
    C : 3.0,
    D : true,
    E : "Hello World",
  }

  var key, err = c.Set("storeTestTypeA", record)

  a.Nil(err, "Failed to save record")
  a.True(key >= 0, "Invalid key from save operation (%d)", key)

  storeTestTeardown(runner, c)
}

func Test_store_can_get_record(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  var record = storeTestTypeA {
    A : 1,
    B : 2,
    C : 3.0,
    D : true,
    E : "Hello World",
  }

  var key, _ = c.Set("storeTestTypeA", record)
  var raw, err = c.Get("storeTestTypeA", key)
  a.Nil(err, "Failed to read record")

  var rcopy = raw.(*storeTestTypeA)
  a.Equals(rcopy.A, 1, "Invalid key value (int)")
  a.Equals(rcopy.B, int64(2), "Invalid key value (long)")
  a.Equals(rcopy.C, float32(3.0), "Invalid key value (float)")
  a.Equals(rcopy.D, true, "Invalid key value (bool)")
  a.Equals(rcopy.E, "Hello World", "Invalid key value (string)")

  storeTestTeardown(runner, c)
}

func Test_store_can_delete(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  var record = storeTestTypeA {
    A : 1,
    B : 2,
    C : 3.0,
    D : true,
    E : "Hello World",
  }

  var key1, _ = c.Set("storeTestTypeA", record)
  var key2, _ = c.Set("storeTestTypeA", record)
  var key3, _ = c.Set("storeTestTypeA", record)
  var key4, _ = c.Set("storeTestTypeA", record)
  var key5, _ = c.Set("storeTestTypeA", record)
  var count, _ = c.Count("storeTestTypeA")
  a.Equals(count, 5, "Invalid count after insert")

  c.Delete("storeTestTypeA", key2)
  c.Delete("storeTestTypeA", key5)

  var rcopy1, _ = c.Get("storeTestTypeA", key1)
  var rcopy2, _ = c.Get("storeTestTypeA", key2)
  var rcopy3, _ = c.Get("storeTestTypeA", key3)
  var rcopy4, _ = c.Get("storeTestTypeA", key4)
  var rcopy5, _ = c.Get("storeTestTypeA", key5)

  a.NotNil(rcopy1, "Invalid result after delete for key 1")
  a.Nil(rcopy2, "Invalid result after delete for key 2")
  a.NotNil(rcopy3, "Invalid result after delete for key 3")
  a.NotNil(rcopy4, "Invalid result after delete for key 4")
  a.Nil(rcopy5, "Invalid result after delete for key 5")

  storeTestTeardown(runner, c)
}

func Test_store_can_get_types(runner TestRunner, T *testing.T) {
	if !runner.Run() { return; }
  var a, c = storeTestSetup(runner, T)

  c.Register(reflect.TypeOf(T), "Testing")
  c.Register(reflect.TypeOf(a), "Assert")
  tt := storeTestTypeA{}; c.Register(reflect.TypeOf(tt), "storeTestTypeA")

  var types = c.Registered()
  a.Equals(len(types), 3, "Invalid type count after register")

  storeTestTeardown(runner, c)
}

func Test_can_update_record(runner TestRunner, T *testing.T) {
  var a, c = storeTestSetup(runner, T)
  var _ = c
  var _ = a

  var record = storeTestTypeA {
    A : 1,
    B : 2,
    C : 3.0,
    D : true,
    E : "Hello World",
  }

  var key1, _ = c.Set("storeTestTypeA", record)

  var rcopy, _ = c.Get("storeTestTypeA", key1)
  record = *rcopy.(*storeTestTypeA)

  a.Equals(record.E, "Hello World", "Failed to load correctly")
  record.E = "Hello World 2"
  c.Set("storeTestTypeA", record, key1)

  rcopy, _ = c.Get("storeTestTypeA", key1)
  record = *rcopy.(*storeTestTypeA)

  a.Equals(record.A, 1, "Invalid value for record 2")
  a.Equals(record.E, "Hello World 2", "Invalid value for record 2")

  storeTestTeardown(runner, c)
}

/*============================================================================*
 * }}}
 *============================================================================*/

 