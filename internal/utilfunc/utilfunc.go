package utilfunc

func Retry(n int, fn func() error) error {
	start := 0
	var err error
	for {
		if start >= n {
			return err
		}

		err = fn()
		if err != nil {
			start++
		} else {
			return nil
		}
	}
}
