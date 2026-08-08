package repository

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
