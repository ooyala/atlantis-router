package zk

func ArrayDiff(curr, prev []string) []string {
	prevMap := make(map[string]bool)
	for _, s := range prev {
		prevMap[s] = true
	}

	created := make([]string, 0)
	for _, s := range curr {
		if !prevMap[s] {
			created = append(created, s)
		}
	}

	return created
}
