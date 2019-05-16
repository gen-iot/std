package std

type StringIterator interface {
	Next() (value string, hasNext bool)
}
