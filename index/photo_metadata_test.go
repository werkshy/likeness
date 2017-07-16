package index

import (
	"testing"
	"time"
)

const (
	sampleJpg = "../samples/20130307_093723.jpg"
	sampleMp4 = "../samples/20130302_104024.mp4"
)

func TestFileDate(t *testing.T) {
	// Don't want to be too rigid, just make sure we get a time
	expectedMinDate := time.Unix(1000000, 0)

	fileDate := FileDate(sampleJpg)

	if fileDate.Before(expectedMinDate) {
		t.Error("File date is too old, doesn't seem right")
	}
}

func TestExifDate(t *testing.T) {
	expectedDate, _ := time.Parse(time.RFC3339, "2013-03-07T09:37:23-05:00")

	metaDate := MetaDate(sampleJpg)

	if !metaDate.Equal(expectedDate) {
		t.Errorf("meta date is wrong: %s != %s\n", metaDate.Format(time.RFC3339), expectedDate)
	}
}

func TestMp4Date(t *testing.T) {
	expectedDate, _ := time.Parse(time.RFC3339, "2013-03-02T08:40:49Z")

	metaDate := MetaDate(sampleMp4)

	if !metaDate.Equal(expectedDate) {
		t.Errorf("meta date is wrong: %s != %s\n", metaDate.Format(time.RFC3339), expectedDate)
	}
}
