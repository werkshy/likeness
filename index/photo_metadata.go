package index

import (
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Joe-xu/mp4parser"
	"github.com/djherbis/times"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	mediainfo "github.com/zhulik/go_mediainfo"
)

// Retrieve the file modifcation (or birth time if available)
// Uses https://github.com/djherbis/times
func FileDate(path string) (t time.Time) {
	timespec, err := times.Stat(path)
	if err != nil {
		log.Printf("Can't read file %s\n", path)
		return time.Unix(0, 0)
	}

	t = timespec.ModTime()
	if timespec.HasBirthTime() {
		t = timespec.BirthTime()
	}
	return
}

// TODO get 'best guess' date for photo
// 1. Exif date for jpg
// 2.Mp4 metatdata for mp4
// 3. Parse date from filename (20XXXXX or 20XX-XX-XX that matches a real date)
// 4. File modification time
// Aaaand move this to a different class

func MetaDate(photoPath string) (maybeTime NullTime) {
	switch strings.ToLower(path.Ext(photoPath)) {
	case ".jpg":
		maybeTime = exifDate(photoPath)
	case ".mp4":
		maybeTime = mp4Date(photoPath)
	case ".mov":
		maybeTime = mediaInfoDate(photoPath)
	default:
		log.Printf("Don't know how to read metadata for %s\n", photoPath)
	}
	return
}

func exifDate(path string) (maybeTime NullTime) {
	maybeTime.Valid = false
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	exifData, err := exif.Decode(f)
	if err != nil {
		log.Printf("Couldn't read exif data for %s\n", path)
		return
	}

	// NOTE: exif dates don't contain timezone. This code converts to the local timezone
	//       which is probably what we want anyway.
	exifTime, err := exifData.DateTime()
	if err == nil {
		maybeTime.Time = exifTime
		maybeTime.Valid = true
	} else {
		log.Printf("Couldn't read an exif date from from %s\n", path)
	}
	return
}

func mediaInfoDate(path string) (maybeTime NullTime) {
	mi := mediainfo.NewMediaInfo()
	err := mi.OpenFile(path)
	defer mi.Close()
	if err != nil {
		log.Printf("Mediainfo: failed to parse the bytes")
		return
	}

	encodedDateString := mi.Get("Encoded_Date")
	if encodedDateString == "" {
		log.Printf("Mediainfo: No Encoded_Date found\n")
		return
	}

	format := "MST 2006-01-02 15:04:05"

	encodedDate, err := time.Parse(format, encodedDateString)
	if err != nil {
		log.Printf("Mediainfo: couldn't parse %s : %s\n", encodedDateString, err)
		return
	}

	maybeTime.Time = encodedDate
	maybeTime.Valid = true

	return
}

func mp4Date(path string) (maybeTime NullTime) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Printf("Couldn't open %s\n", path)
		return
	}

	p := mp4parser.NewParser(f)
	info, err := p.Parse()
	if err == nil {
		maybeTime.Time = info.CreationTime()
		maybeTime.Valid = true
	} else {
		log.Printf("Couldn't read an mp4 date from from %s\n", path)
	}

	return
}
