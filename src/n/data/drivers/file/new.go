package file

import "reflect"
import "n/data"

/*============================================================================*
 * {{{ Internal api 
 *============================================================================*/

type factory struct {
  RecordSet func(t reflect.Type) RecordSet
  Driver func() data.Driver
}

var New factory = factory {
  RecordSet : newRecordSet,
  Driver : newDriver,
}

/*============================================================================*
 * }}}
 *============================================================================*/
