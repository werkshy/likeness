package index

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"testing"
	"time"

	"github.com/werkshy/likeness/test_util"
)

func init() {
	log.Println("Running store test init")
	test_util.Init()
	test_util.MigrateDb()
}

var db = test_util.Db // convenience handle to the db creation method

func insertPhoto(t *testing.T, store Store) (photo Photo) {
	photo = Photo{
		Id:       -1,
		Path:     randomString(4) + "/" + randomString(4),
		Md5:      randomBytes(16),
		FileDate: time.Now(),
	}
	log.Printf("Photo: %s\n", photo)

	err := store.InsertPhoto(photo)
	if err != nil {
		t.Error(err)
	}
	return
}

func TestFindPhotoByMd5(t *testing.T) {
	store := NewDbStore(db())
	photo := insertPhoto(t, store)

	found, err := store.FindPhotoByMd5(photo.Md5)

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found photo: %s\n", photo)

	if found.Path != photo.Path {
		t.Errorf("Expected path '%s' to equal '%s'\n", found.Path, photo.Path)
	}
}

func TestFindPhotoByPath(t *testing.T) {
	store := NewDbStore(db())
	photo := insertPhoto(t, store)

	found, err := store.FindPhotoByPath(photo.Path)

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found photo: %s\n", photo)

	if found.Path != photo.Path {
		t.Errorf("Expected path '%s' to equal '%s'\n", found.Path, photo.Path)
	}
}

func TestFindPhotoByPathOrMd5(t *testing.T) {
	store := NewDbStore(db())
	photo := insertPhoto(t, store)

	found, err := store.FindPhotoByPathOrMd5(photo.Path, photo.Md5)

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found photo: %s\n", photo)

	if found.Path != photo.Path {
		t.Errorf("Expected path '%s' to equal '%s'\n", found.Path, photo.Path)
	}
}

func TestInsertDupeChecksum(t *testing.T) {
	store := NewDbStore(db())
	photo := insertPhoto(t, store)

	photo.Path = randomString(4) + "/" + randomString(4)
	err := store.InsertPhoto(photo)
	switch {
	case err == nil:
		t.Error("Expected second insert to return an error")
	case err != UniqueConstraintViolation:
		t.Error("Expected second insert to return a UniqueConstraintViolation error")
	}
}

func TestInsertDupePath(t *testing.T) {
	store := NewDbStore(db())
	photo := insertPhoto(t, store)

	photo.Md5 = randomBytes(16)
	err := store.InsertPhoto(photo)
	switch {
	case err == nil:
		t.Error("Expected second insert to return an error")
	case err == UniqueConstraintViolation:
		t.Error("Expected second insert not to return a UniqueConstraintViolation error")
	}

}

func randomBytes(len int) []byte {
	u := make([]byte, len)
	_, err := rand.Read(u)
	if err != nil {
		log.Fatal(err)
	}

	return u
}

func randomString(len int) string {
	u := randomBytes(len)
	return hex.EncodeToString(u)
}
