package n

import "n/test"
import "testing"
import r "reflect"

/*============================================================================*
 * {{{ Helpers
 *============================================================================*/

type resolverTestInfA interface {
  Value() int
}

type resolverTestTypeA struct {
  child resolverTestInfB
}

func (self *resolverTestTypeA) Value() int {
  return 100
}

func resolverTestFactoryA() resolverTestInfA {
  var rtn = &resolverTestTypeA{}

  // Resolve dependencies
  var rs = New.Resolver()
  rtn.child = rs.Get(r.TypeOf(&rtn.child)).(resolverTestInfB)

  return rtn
}

type resolverTestInfB interface {
  Value() int
  Value2() int
}

type resolverTestTypeB struct {
}

func (self *resolverTestTypeB) Value() int {
  return 200
}

func (self *resolverTestTypeB) Value2() int {
  return 400
}

func resolverTestFactoryB() resolverTestInfB {
  return &resolverTestTypeB{}
}

func resolverTestBadFactoryB1() int {
  return 300
}

func resolverTestBadFactoryB2() resolverTestInfB {
  return nil
}

func resolverTestBadFactoryB3() resolverTestInfA {
  return &resolverTestTypeA{}
}

func resolverTestSetup(T *testing.T) (test.Assert, Resolver) {
  var assert = test.New.Assert(T)
  var instance = New.Resolver()
  return assert, instance
}

func resolverTestTeardown(target Resolver) {
  target.Clear()
}

/*============================================================================*
 * }}}
 *============================================================================*/

/*============================================================================*
 * {{{ Tests
 *============================================================================*/

func Test_resolver_can_create_instance(T *testing.T) {
  var a, i = resolverTestSetup(T)

  a.NotNil(i, "Unable to create instance")

  resolverTestTeardown(i)
}

func Test_resolver_cannot_manufacture_from_bad_factory(T *testing.T) {
  var a, i = resolverTestSetup(T)
  var b resolverTestInfB
  var bi interface{}
  var err error
  
  i.Register(r.TypeOf(&b), resolverTestBadFactoryB1)
  bi, err = i.Singleton(r.TypeOf(&b))
  a.NotNil(err, "Didn't fail on invalid factory B1")
  a.Nil(bi, "Created instance from invalid factory B1")

  i.Register(r.TypeOf(&b), resolverTestBadFactoryB2)
  bi, err = i.Singleton(r.TypeOf(&b))
  a.NotNil(err, "Didn't fail on invalid factory B2")
  a.Nil(bi, "Created instance from invalid factory B2")

  i.Register(r.TypeOf(&b), resolverTestBadFactoryB3)
  bi, err = i.Singleton(r.TypeOf(&b))
  a.NotNil(err, "Didn't fail on invalid factory B3")
  a.Nil(bi, "Created instance from invalid factory B3")

  resolverTestTeardown(i)
}

func Test_resolver_cannot_register_stupid_types(T *testing.T) {
  var a, i = resolverTestSetup(T)
  
  var b resolverTestInfB; var err = i.Register(r.TypeOf(b), resolverTestFactoryB)
  a.NotNil(err, "Didn't fail on invalid register request")

  err = i.Register(r.TypeOf(&b), resolverTestFactoryB)
  a.Nil(err, "Failed on valid register request")

  resolverTestTeardown(i)
}

func Test_resolver_can_manufacture_objects(T *testing.T) {
  var a, i = resolverTestSetup(T)
  
  var ta resolverTestInfA; i.Register(r.TypeOf(&ta), resolverTestFactoryA)
  var tb resolverTestInfB; i.Register(r.TypeOf(&tb), resolverTestFactoryB)

  var bi, errb = i.Singleton(r.TypeOf(&tb))
  a.Nil(errb, "Fail on factory B")
  a.NotNil(bi, "Failed to create instance from valid factory B")
  tb = bi.(resolverTestInfB)

  var ai, erra = i.Singleton(r.TypeOf(&ta))
  a.Nil(erra, "Fail on factory A")
  a.NotNil(ai, "Failed to create instance from valid factory A")
  ta = ai.(resolverTestInfA)

  var inst resolverTestInfA; inst = i.Get(r.TypeOf(&inst)).(resolverTestInfA)
  a.NotNil(inst, "Created instance from valid factory A")

  resolverTestTeardown(i)
}

/*============================================================================*
 * }}}
 *============================================================================*/
