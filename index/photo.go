package index

import (
	"encoding/hex"
	"time"
)

type Photo struct {
	Id       int64
	Path     string
	Md5      []byte    `db:"checksum_value"`
	FileDate time.Time `db:"file_date"`
	MetaDate NullTime  `db:"meta_date"`
}

func (photo Photo) String() string {
	return photo.Path + " [" + hex.EncodeToString(photo.Md5) + "]"
}
