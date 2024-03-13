package bolt

import (
	"bytes"
	"errors"
)

type S interface{ ~string | ~[]byte }

func BucketOpen[T S](tx BoltTx, bucketNames ...T) (bucket *Bucket, err error) {
	if len(bucketNames) == 0 {
		err = ErrBucketNameRequired
		return
	}

	for _, bucketName := range bucketNames {
		if bucket = tx.Bucket([]byte(bucketName)); bucket == nil {
			err = ErrBucketNotFound
			return
		}
		tx = bucket
	}

	return
}

func BucketCreate[T S](tx BoltTx, bucketNames ...T) (bucket *Bucket, err error) {
	if len(bucketNames) == 0 {
		err = ErrBucketNameRequired
		return
	}

	for _, bucketName := range bucketNames {
		if bucket, err = tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return
		}
		tx = bucket
	}

	bucket = tx.(*Bucket)
	return
}

func Set[B S, K S, V S](bucketName B, key K, value V) Process {
	return func(tx *Tx) error { return TxSet(tx, bucketName, key, value) }
}

func Append[T S](bucketName, key, val T) Process {
	return func(tx *Tx) error { return TxAppend(tx, bucketName, key, val) }
}

func Del[B S, K S](bucketName B, keys ...K) Process {
	return func(tx *Tx) error { return TxDel(tx, bucketName, keys...) }
}

func TxSet[B S, K S, V S](tx BoltTx, bucketName B, key K, value V) error {
	bucket, err := BucketCreate(tx, bucketName)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), []byte(value))
}

func TxAppend[T S](tx BoltTx, bucketName, key, val T) error {
	bucket, err := BucketCreate(tx, bucketName)
	if err != nil {
		return err
	}
	v := []byte(val)
	if exist := bucket.Get([]byte(key)); len(exist) > 0 {
		v = append(append(exist, '\n'), v...)
	}
	return bucket.Put([]byte(key), v)
}

func TxDel[B S, K S](tx BoltTx, bucketName B, keys ...K) error {
	bucket, err := BucketOpen(tx, bucketName)
	if err != nil {
		if errors.Is(err, ErrBucketNotFound) {
			return nil
		}
		return err
	}
	for _, key := range keys {
		if err := bucket.Delete([]byte(key)); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[T S](bucketName T, fn func(k, v []byte) error) Process {
	return func(tx *Tx) error {
		bucket, err := BucketOpen(tx, bucketName)
		if err != nil {
			return err
		}
		return bucket.ForEach(fn)
	}
}

func ForEachWithPrefix[T S, P S](bucketName T, prefix P, fn func(k, v []byte) error) Process {
	return func(tx *Tx) error {
		bucket, err := BucketOpen(tx, bucketName)
		if err != nil {
			return err
		}
		c := bucket.Cursor()
		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}
		return nil
	}
}
