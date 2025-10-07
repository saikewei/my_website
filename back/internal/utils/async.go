package utils

func RunTaskAsync(task func() error) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- task()
	}()
	return <-errChan
}

func RunTaskAsyncWithResult[T any](task func() (T, error)) (T, error) {
	type result struct {
		Value T
		Err   error
	}

	resultChan := make(chan result, 1)

	go func() {
		defer close(resultChan)

		val, err := task()

		resultChan <- result{Value: val, Err: err}
	}()

	res := <-resultChan

	return res.Value, res.Err
}
