package std

type UniverseFreeNoexcept interface {
	Free() error
}

type UniverseFree interface {
	Free()
}

func UniverseFreeProxy(ele interface{}) {
	if uf := ele.(UniverseFree); uf != nil {
		uf.Free()
	}
}

func UniverseErrFreeProxy(ele interface{}) error {
	if uf := ele.(UniverseFreeNoexcept); uf != nil {
		return uf.Free()
	}
	return nil
}
