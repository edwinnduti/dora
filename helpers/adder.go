package helpers

func Add(box *[]int, choice *string) *[]int {

	var POOR string = "Poor"
	var FAIR string = "Fair"
	var GOOD string = "Good"
	var VERYGOOD string = "Very Good"
	var EXCELLENT string = "Excellent"

	// switch elements
	switch *choice {
	case POOR:
		*box = append(*box, 1)
		return box
	case FAIR:
		*box = append(*box, 2)
		return box
	case GOOD:
		*box = append(*box, 3)
		return box
	case VERYGOOD:
		*box = append(*box, 4)
		return box
	case EXCELLENT:
		*box = append(*box, 5)
		return box
	default:
		*box = append(*box, 0)
		return box
	}
}
