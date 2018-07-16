package sqlrepo

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
)

//IdsPlaceholdersArgs returns a list of placeholders (from List(Placeholder, len(ids))
//and a slice of arguments that is a copy of ids.
func IdsPlaceholdersArgs(ids []data.Id) (string, []interface{}) {
	if len(ids) == 0 {
		return "", []interface{}{}
	}

	placeholders := List(Placeholder, len(ids))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	return placeholders, args
}
