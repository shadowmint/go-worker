package services

/*============================================================================*
 * {{{ Internal api 
 *============================================================================*/

type factory struct {
  Fragments func(...interface{}) Fragments
  Index func(...interface{}) Index
}

var New factory = factory {
  Fragments : newFragments,
  Index : newIndex,
}

/*============================================================================*
 * }}}
 *============================================================================*/
