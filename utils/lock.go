package utils

type Locker interface {
	Lock(string) (unLock func() error, err error)
}

func WithLock(locker Locker, lockKey string, f func() error) error {
	var err error
	unlock, err := locker.Lock(lockKey)
	if err != nil {
		return err
	}

	defer unlock()
	err = f()
	if err != nil {
		return err
	}
	return nil
}
