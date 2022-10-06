package util

func FilterHandler(filterValues []interface{}) (filter string) {
	if len(filterValues) > 0 {
		filter = ` AND`
	} else {
		filter = ` WHERE`
	}
	return filter
}
