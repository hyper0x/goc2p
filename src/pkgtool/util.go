package pkgtool

func appendIfAbsent(s []string, t ...string) []string {
	for _, t1 := range t {
		var contains bool
		for _, s1 := range s {
			if s1 == t1 {
				contains = true
				break
			}
		}
		if !contains {
			s = append(s, t1)
		}
	}
	return s
}
