package sqlrepo

import "regexp"

//Placeholder is the only placeholder recognized by this package.
//It is transparently parsed in queries into the Dialect specific representation.
const Placeholder = "?"

var placeholderRegexp = regexp.MustCompile(`\` + Placeholder)

//Dialect allows client code to use this package in a DBMS independent way
//by providing methods to abstract the differences betwenn DBMSs.
type Dialect interface {
	//Placeholder should return the DBMS specific placeholder for index.
	//Index is zero-based.
	Placeholder(index int) string
}

//Normalize returns query with every instance of Placeholder replaced by a call
//to d.Placeholder(index) where index is the occurance index of Placeholder in query.
func Normalize(d Dialect, query string) string {
	index := 0
	return placeholderRegexp.ReplaceAllStringFunc(query, func(_ string) string {
		index++
		return d.Placeholder(index - 1)
	})
}
