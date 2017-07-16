package index

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var UniqueConstraintViolation = errors.New("Unique constraint violation")

const PqUniqueConstraint = "23505"

type Store interface {
	FindPhotoByMd5(md5 []byte) (Photo, error)
	FindPhotoByPath(path string) (Photo, error)
	FindPhotoByPathOrMd5(path string, md5 []byte) (Photo, error)
	InsertPhoto(photo Photo) error
}

type DbStore struct {
	*sqlx.DB
}

func NewDbStore(db *sqlx.DB) Store {
	return DbStore{db}
}

func (store DbStore) FindPhotoByMd5(md5 []byte) (photo Photo, err error) {
	selectSql := `
		SELECT id, path, checksum_value, file_date, meta_date
		FROM photos
		WHERE checksum_type = 'md5' AND checksum_value = $1`

	err = store.Get(&photo, selectSql, md5)
	logPqError(err)
	return
}

func (store DbStore) FindPhotoByPath(path string) (photo Photo, err error) {
	selectSql := `
		SELECT id, path, checksum_value, file_date, meta_date
		FROM photos
		WHERE path = $1`

	err = store.Get(&photo, selectSql, path)
	logPqError(err)
	return
}

func (store DbStore) FindPhotoByPathOrMd5(path string, md5 []byte) (photo Photo, err error) {
	selectSql := `
		SELECT id, path, checksum_value, file_date, meta_date
		FROM photos
		WHERE path = $1
		OR checksum_type = 'md5' AND checksum_value = $2`

	err = store.Get(&photo, selectSql, path, md5)
	logPqError(err)
	return
}

func (store DbStore) InsertPhoto(photo Photo) (err error) {
	// ON CONFLICT clause requires Postgres >= 9.5
	insertSql := `
		INSERT INTO photos (path, meta_date, file_date, checksum_type, checksum_value)
		VALUES ($1, $2, $3, 'md5', $4)
		ON CONFLICT (checksum_type, checksum_value) DO NOTHING
		RETURNING id`
	row := store.QueryRow(insertSql, photo.Path, photo.MetaDate, photo.FileDate, photo.Md5)

	var insertedId int64
	err = row.Scan(&insertedId)
	logPqError(err)

	// If no row was inserted it means the ON CONFLICT clause was triggered
	// i.e. we hit a constraint violation on the checksum unique index
	if err == sql.ErrNoRows {
		log.Printf("WARN: constraint violation on md5 '%s'\n", hex.EncodeToString(photo.Md5))
		err = UniqueConstraintViolation
	}
	if err == nil {
		log.Printf("Inserted photo #%d\n", insertedId)
	}

	return
}

// Log the PQ error code and reason if this error can be cast as a pq.Error
func logPqError(err error) {
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Printf("pq error: %s - %s\n", pqErr.Code, pqErr.Message)
		}
	}
}
