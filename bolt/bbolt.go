package bolt

import (
	"go.etcd.io/bbolt"
)

type (
	BoltDB       = bbolt.DB
	Bucket       = bbolt.Bucket
	BucketStats  = bbolt.BucketStats
	CheckOption  = bbolt.CheckOption
	Cursor       = bbolt.Cursor
	Options      = bbolt.Options
	PageInfo     = bbolt.PageInfo
	Stats        = bbolt.Stats
	TxStats      = bbolt.TxStats
	Tx           = bbolt.Tx
	FreelistType = bbolt.FreelistType
)

const (
	FreelistArrayType = bbolt.FreelistArrayType
	FreelistMapType   = bbolt.FreelistMapType

	DefaultAllocSize     = bbolt.DefaultAllocSize
	DefaultFillPercent   = bbolt.DefaultFillPercent
	DefaultMaxBatchDelay = bbolt.DefaultMaxBatchDelay
	DefaultMaxBatchSize  = bbolt.DefaultMaxBatchSize
)

var (
	Open           = bbolt.Open
	DefaultOptions = bbolt.DefaultOptions
	WithKVStringer = bbolt.WithKVStringer
	HexKVStringer  = bbolt.HexKVStringer
)

var (
	ErrBucketExists       = bbolt.ErrBucketExists
	ErrBucketNameRequired = bbolt.ErrBucketNameRequired
	ErrBucketNotFound     = bbolt.ErrBucketNotFound
	ErrChecksum           = bbolt.ErrChecksum
	ErrDatabaseNotOpen    = bbolt.ErrDatabaseNotOpen
	ErrDatabaseOpen       = bbolt.ErrDatabaseOpen
	ErrDatabaseReadOnly   = bbolt.ErrDatabaseReadOnly
	ErrFreePagesNotLoaded = bbolt.ErrFreePagesNotLoaded
	ErrIncompatibleValue  = bbolt.ErrIncompatibleValue
	ErrInvalid            = bbolt.ErrInvalid
	ErrInvalidMapping     = bbolt.ErrInvalidMapping
	ErrKeyRequired        = bbolt.ErrKeyRequired
	ErrKeyTooLarge        = bbolt.ErrKeyTooLarge
	ErrTimeout            = bbolt.ErrTimeout
	ErrTxClosed           = bbolt.ErrTxClosed
	ErrTxNotWritable      = bbolt.ErrTxNotWritable
	ErrValueTooLarge      = bbolt.ErrValueTooLarge
	ErrVersionMismatch    = bbolt.ErrVersionMismatch
)
