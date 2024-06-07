package slicez // import "ezpkg.io/slicez"

func First[S ~[]E, E any](s S) (out E) {
	if len(s) > 0 {
		return s[0]
	}
	return out
}

func FirstFunc[S ~[]E, E any](s S, fn func(E) bool) (out E) {
	for _, item := range s {
		if fn(item) {
			return item
		}
	}
	return out
}

func Last[S ~[]E, E any](s S) (out E) {
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return out
}

func LastFunc[S ~[]E, E any](s S, fn func(E) bool) (out E) {
	for i := len(s) - 1; i >= 0; i-- {
		if fn(s[i]) {
			return s[i]
		}
	}
	return out
}

func MapFunc[S ~[]E, E, R any](s S, fn func(E) R) []R {
	result := make([]R, len(s))
	for i, item := range s {
		result[i] = fn(item)
	}
	return result
}

func FilterFunc[S ~[]E, E any](s S, fn func(E) bool) (outs []E) {
	for _, item := range s {
		if fn(item) {
			outs = append(outs, item)
		}
	}
	return outs
}

func MapFilterFunc[S ~[]E, E, R any](s S, fn func(E) (R, bool)) (outs []R) {
	for _, item := range s {
		out, ok := fn(item)
		if ok {
			outs = append(outs, out)
		}
	}
	return outs
}
