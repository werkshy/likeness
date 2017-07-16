package index

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/Joe-xu/mp4parser"
	"github.com/djherbis/times"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
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
	switch path.Ext(photoPath) {
	case ".jpg":
		maybeTime = exifDate(photoPath)
	case ".mp4":
		maybeTime = mp4Date(photoPath)
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

	// TODO get camModel etc
	//camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	//fmt.Println(camModel.StringVal())

	exifTime, err := exifData.DateTime()
	if err == nil {
		maybeTime.Time = exifTime
		maybeTime.Valid = true
	} else {
		log.Printf("Couldn't read an exif date from from %s\n", path)
	}
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
