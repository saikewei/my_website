package utils

func RunTaskAsync(task func() error) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- task()
	}()
	return <-errChan
}
