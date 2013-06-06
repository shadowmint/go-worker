package urls 

/*============================================================================*
 * {{{ Public api
 *============================================================================*/

// Fragment urls
func (s *UrlHelper) FragmentIndex() 			   string { return s.url("/fragments/") }
func (s *UrlHelper) FragmentAdd()            string { return s.url("/fragments/add/") }
func (s *UrlHelper) FragmentDelete()         string { return s.url("/fragments/delete/") }
func (s *UrlHelper) FragmentSave()           string { return s.url("/fragments/save/") }
func (s *UrlHelper) FragmentEdit(key int64)  string { return s.url("/fragments/edit/%d", key) }
 
/*============================================================================*
 * }}}
 *============================================================================*/
