package memory

import "reflect"
import "n/data"

/*============================================================================*
 * {{{ Internal api 
 *============================================================================*/

type factory struct {
  Driver func() data.Driver
  RecordSet func(t reflect.Type) RecordSet
}

var New factory = factory {
  Driver : newDriver,
  RecordSet : newRecordSet,
}

/*============================================================================*
 * }}}
 *============================================================================*/
