package common

func MergeDistinctByName[T INamed](mainSlice, dbSlice []T) []T {
	existing := make(map[string]struct{})

	for _, item := range mainSlice {
		existing[item.GetName()] = struct{}{}
	}

	for _, item := range dbSlice {
		if _, ok := existing[item.GetName()]; !ok {
			mainSlice = append(mainSlice, item)
		}
	}

	return mainSlice
}
