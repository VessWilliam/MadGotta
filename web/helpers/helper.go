package helpers

import "strconv"

func ToggleUrl(id int) string {
	return "/toggle/" + strconv.Itoa(id)
}

func DeleteUrl(id int) string {
	return "/delete/" + strconv.Itoa(id)
}

func ToggleLabel(done bool) string {
	if done {
		return "Undo"
	}
	return "Done"
}

func VisibilityExpr(done bool) string {
	if done {
		return "filter === 'all' || filter === 'done'"
	}
	return "filter === 'all' || filter === 'active'"
}
