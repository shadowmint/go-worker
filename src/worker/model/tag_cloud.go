package model 

/*============================================================================*
 * {{{ Public api 
 *============================================================================*/

// Tag cloud marker
type TagCloudMarker struct {
	Tag string
	Hits int
}

// A set of tags associated with a fragment
type TagCloud struct {
	Tags []TagCloudMarker
}

/*============================================================================*
 * }}}
 *============================================================================*/
