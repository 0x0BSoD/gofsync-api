package utils

// Split []int to [](parts * []int)
func SplitToQueue(item []int, parts int) [][]int {

	var result [][]int
	length := len(item)
	sliceLength := 0

	// Checks ==========
	if length <= parts {
		return append(result, item)
	}

	if length%parts == 0 {
		sliceLength = length / parts
	} else {
		if length/parts == 1 {
			return append(result, item)
		}
		sliceLength = length / parts
	}
	if sliceLength == 1 {
		return append(result, item)
	}

	// Split ==========
	start := 0
	stop := sliceLength
	for i := 0; i <= parts; i++ {
		if i != parts {
			result = append(result, item[start:stop])
		} else {
			result = append(result, item[start:])
			break
		}
		start = stop
		stop = start + sliceLength
	}

	return result
}
