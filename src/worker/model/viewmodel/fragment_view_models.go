package viewmodel 

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

type FragmentsListViewModel struct {
	Recent []FragmentSummaryViewModel
	Items []FragmentSummaryViewModel
	More bool
}

type FragmentSummaryViewModel struct {
	Id int64
	Title string
	Tags []string
}

type FragmentViewModel struct {
	Id int64
  Title string
  Content string
  Tags []string
  AllTags string
}

/*============================================================================*
 * }}}
 *============================================================================*/
