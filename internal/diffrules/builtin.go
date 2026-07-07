package diffrules

func BuiltinRegistry() (*Registry, error) {
	return NewRegistry(NewDIF001(), NewDIF002(), NewDIF003())
}
