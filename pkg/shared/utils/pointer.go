package utils

// ToStringPointer convierte un string a un puntero a string
func ToStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ToIntPointer convierte un int a un puntero a int
func ToIntPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// PointerToString convierte un puntero a string a un string
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
