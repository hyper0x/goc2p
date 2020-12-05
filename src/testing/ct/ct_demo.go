package ct

func TypeCategoryOf(v interface{}) string {
	switch v.(type) {
	case bool:
		return "boolean"
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return "integer"
	case float32, float64:
		return "float"
	case complex64, complex128:
		return "complex"
	case string:
		a := "string"
		return a
	}
	return "unknown"
}
