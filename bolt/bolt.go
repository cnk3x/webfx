package bolt

type Process = func(tx *Tx) error

type DB struct {
	BoltDB
}

func (db *DB) Update(processes ...Process) (err error) {
	err = db.BoltDB.Update(func(tx *Tx) (err error) {
		for _, process := range processes {
			if err = process(tx); err != nil {
				return
			}
		}
		return
	})
	return
}

func (db *DB) View(processes ...Process) (err error) {
	err = db.BoltDB.View(func(tx *Tx) (err error) {
		for _, process := range processes {
			if err = process(tx); err != nil {
				return
			}
		}
		return
	})
	return
}

func (db *DB) Get(bucketName, key string) (val string, err error) {
	err = db.View(func(tx *Tx) error {
		bucket, err := BucketOpen(tx, bucketName)
		if err != nil {
			return err
		}
		val = string(bucket.Get([]byte(key)))
		return nil
	})

	return
}

type BoltTx interface {
	Bucket(name []byte) *Bucket
	CreateBucketIfNotExists(name []byte) (*Bucket, error)
	DeleteBucket(key []byte) error
}

type BoltBucket interface {
	BoltTx
	Put(key, value []byte) error
	Get(key []byte) []byte
	Delete(key []byte) error
	Cursor() *Cursor
	NextSequence() (uint64, error)
	ForEach(fn func(k, v []byte) error) error
}

var _ = BoltTx((*Tx)(nil))
var _ = BoltTx((*Bucket)(nil))
