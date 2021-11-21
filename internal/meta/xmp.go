package meta

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"path/filepath"
	"runtime/debug"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"

	"github.com/photoprism/photoprism/pkg/txt"
	"github.com/photoprism/photoprism/pkg/fs"
)

// XMP parses an XMP file and returns a Data struct.
func XMP(fileName string) (data Data, err error) {
	err = data.XMP(fileName)

	return data, err
}

// Parses XMP data out of a media file
func (data *Data) XMPMedia(fileName string, fileType fs.FileType) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("metadata: %s in %s (xmp panic)\nstack: %s", e, txt.Quote(filepath.Base(fileName)), debug.Stack())
		}
	}()
	f, err := os.Open(fileName)

	defer f.Close()

	stat, err := f.Stat()

	size := stat.Size()

	sl, err := Parse(f, int(size))
	
	var seg *jpegstructure.Segment
	_, seg, err = sl.FindXmp()
	doc := XmpDocument{}

	if err := doc.FromBytes(seg.Data); err != nil {
		return fmt.Errorf("metadata: can't read %s (xmp)", txt.Quote(filepath.Base(fileName)))
	}

	return data.xmpCommon(doc)
}

// Parse parses a JPEG uses an `io.ReadSeeker`. Even if it fails, it will return
// the list of segments encountered prior to the failure.
func Parse(rs io.ReadSeeker, size int) (ec *jpegstructure.SegmentList, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = fmt.Errorf("metadata: %s in %s (xmp panic)\nstack: %s", e, txt.Quote(filepath.Base(fileName)), debug.Stack())
		}
	}()

	s := bufio.NewScanner(rs)

	// Since each segment can be any size, our buffer must allowed to grow as
	// large as the file.
	buffer := []byte{}
	s.Buffer(buffer, size)

	js := jpegstructure.NewJpegSplitter(nil)
	s.Split(js.Split)

	for s.Scan() != false {
	}

	// Always return the segments that were parsed, at least until there was an
	// error.
	ec = js.Segments()

	return ec, nil
}

// XMP parses an XMP file and returns a Data struct.
func (data *Data) XMP(fileName string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("metadata: %s in %s (xmp panic)\nstack: %s", e, txt.Quote(filepath.Base(fileName)), debug.Stack())
		}
	}()

	doc := XmpDocument{}

	if err := doc.Load(fileName); err != nil {
		return fmt.Errorf("metadata: can't read %s (xmp)", txt.Quote(filepath.Base(fileName)))
	}

	return data.xmpCommon(doc)

}

func (data *Data) xmpCommon(doc XmpDocument) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("metadata: %s in %s (xmp panic)\nstack: %s", e, txt.Quote(filepath.Base(fileName)), debug.Stack())
		}
	}()


	if doc.Title() != "" {
		data.Title = doc.Title()
	}

	if doc.Artist() != "" {
		data.Artist = doc.Artist()
	}

	if doc.Description() != "" {
		data.Description = doc.Description()
	}

	if doc.Copyright() != "" {
		data.Copyright = doc.Copyright()
	}

	if doc.CameraMake() != "" {
		data.CameraMake = doc.CameraMake()
	}

	if doc.CameraModel() != "" {
		data.CameraModel = doc.CameraModel()
	}

	if doc.LensModel() != "" {
		data.LensModel = doc.LensModel()
	}

	if takenAt := doc.TakenAt(); !takenAt.IsZero() {
		data.TakenAt = takenAt
	}

	if len(doc.Keywords()) != 0 {
		data.AddKeywords(doc.Keywords())
	}

	return nil
}
