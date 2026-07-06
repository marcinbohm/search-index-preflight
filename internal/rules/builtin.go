package rules

func BuiltinRegistry() (*Registry, error) {
	return NewRegistry(NewSIL001(), NewSIL002())
}
