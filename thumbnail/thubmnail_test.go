package thumbnail

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/disintegration/imaging"
)

type Content struct {
	start   int
	length  int
	ratio   float64
	width   int
	height  int
	fps     float64
	modTime time.Time
}

var suppordedExts = map[string]bool{
	".avi": true,
	".mp4": true,
	".flv": true,
	".mkv": true,
}

func TestMultipleVideos(t *testing.T) {
	baseDir := "/home/bukodi/Downloads/CSENGETETT,MYLORD;;1.2.3.ÉVAD/"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !suppordedExts[ext] {
			return nil
		}
		e := procVideo(path)
		if e != nil {
			t.Error(path+": ", e)
		}
		return nil
	})
	if err != nil {
		t.Error("File walking err: ", err)
	}
}

func TestSingleVideo(t *testing.T) {
	path := "/home/bukodi/Downloads/Rogue.One.2016.RETAiL.BDRiP.x264.HuN-HyperX/rogueone-sd-hyperx.mkv"
	err := procVideo(path)
	if err != nil {
		t.Error(path+": ", err)
	}
}

func procVideo(path string) error {
	var c Content
	err := readVideoInfo(path, &c)
	if err != nil {
		return err
	}
	err = createScreenshots(path, &c, "/tmp/Arduino/")
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	tnpath := path[0:len(path)-len(ext)] + ".jpg"

	out, err := os.Create(tnpath)
	if err != nil {
		return err
	}

	err = makeTile2x2(&c, out, "/tmp/Arduino/")
	if err != nil {
		return err
	}
	out.Close()
	os.Chtimes(tnpath, c.modTime, c.modTime)

	files, _ := filepath.Glob("/tmp/Arduino/00*.jpg")
	for _, f := range files {
		os.Remove(f)
	}
	return nil
}

func readVideoInfo(path string, c *Content) error {
	cmd := exec.Command("/usr/bin/mplayer", "-nosound",
		"-vo", "null", "-frames", "1", "-identify", path)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	outTxt := string(out)
	lines := strings.Split(outTxt, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID_VIDEO_WIDTH=") {
			s := strings.TrimPrefix(line, "ID_VIDEO_WIDTH=")
			c.width, _ = strconv.Atoi(s)
		}
		if strings.HasPrefix(line, "ID_VIDEO_HEIGHT=") {
			c.height, _ = strconv.Atoi(strings.TrimPrefix(line, "ID_VIDEO_HEIGHT="))
		}
		if strings.HasPrefix(line, "ID_LENGTH=") {
			f, _ := strconv.ParseFloat(strings.TrimPrefix(line, "ID_LENGTH="), 64)
			c.length = int(f)
		}
		if strings.HasPrefix(line, "ID_VIDEO_FPS=") {
			c.fps, _ = strconv.ParseFloat(strings.TrimPrefix(line, "ID_VIDEO_LENGTH="), 64)
		}
	}
	if c.width == 0 {
		return errors.New("no ID_VIDEO_WIDTH")
	}
	if c.height == 0 {
		return errors.New("no ID_VIDEO_HEIGHT")
	}
	if c.length == 0 {
		return errors.New("no ID_LENGTH")
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	c.modTime = info.ModTime()
	return nil
}

func createScreenshots(path string, c *Content, tmpDir string) (err error) {
	var sstep int
	sstep = int(c.length / 5)
	cmd := exec.Command("/usr/bin/mplayer", "-noidle", "-nogui", "-nosound",
		"-vo", "jpeg", "-frames", "5", "-sstep", strconv.Itoa(sstep), path)
	cmd.Dir = tmpDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func readImage(path string) (img image.Image, err error) {
	var fimg *os.File
	fimg, err = os.Open(path)
	if err != nil {
		return
	}
	defer fimg.Close()
	// Decode the file to get the image data
	img, _, err = image.Decode(fimg)
	return
}

type quarter int

const (
	topLeft quarter = iota
	topRight
	bottomLeft
	bottomRight
)

func makeTile2x2(c *Content, out io.Writer, tmpDir string) (err error) {
	//var destImg = image.NewRGBA(image.Rect(0, 0, 720, 302))
	var destImg = image.NewRGBA(image.Rect(0, 0, 480, 270))
	err = addQuater(destImg, topLeft, filepath.Join(tmpDir, "00000002.jpg"))
	if err != nil {
		return
	}

	err = addQuater(destImg, topRight, filepath.Join(tmpDir, "00000003.jpg"))
	if err != nil {
		return
	}

	err = addQuater(destImg, bottomLeft, filepath.Join(tmpDir, "00000004.jpg"))
	if err != nil {
		return
	}

	err = addQuater(destImg, bottomRight, filepath.Join(tmpDir, "00000005.jpg"))
	if err != nil {
		return
	}

	option := jpeg.Options{Quality: 50}
	err = jpeg.Encode(out, destImg, &option)
	if err != nil {
		return
	}
	return
}

func addQuater(destImg *image.RGBA, q quarter, srcImpPath string) error {
	var rect image.Rectangle

	switch q {
	case topLeft:
		rect = image.Rect(0, 0, destImg.Rect.Max.X/2, destImg.Rect.Max.Y/2)
	case topRight:
		rect = image.Rect(destImg.Rect.Max.X/2, 0, destImg.Rect.Max.X, destImg.Rect.Max.Y/2)
	case bottomLeft:
		rect = image.Rect(0, destImg.Rect.Max.Y/2, destImg.Rect.Max.X/2, destImg.Rect.Max.Y)
	case bottomRight:
		rect = image.Rect(destImg.Rect.Max.X/2, destImg.Rect.Max.Y/2, destImg.Rect.Max.X, destImg.Rect.Max.Y)
	}

	srcImg, err := readImage(srcImpPath)
	if err != nil {
		return err
	}

	srcImg = imaging.Fit(srcImg, destImg.Rect.Max.X/2, destImg.Rect.Max.Y/2, imaging.NearestNeighbor)

	srcBounds := srcImg.Bounds()
	if srcBounds.Max.X < destImg.Rect.Max.X/2 {
		addX := (destImg.Rect.Max.X/2 - srcBounds.Max.X) / 2
		rect.Min.X += addX
		rect.Max.X += addX
	}
	if srcBounds.Max.Y < destImg.Rect.Max.Y/2 {
		addY := (destImg.Rect.Max.Y/2 - srcBounds.Max.Y) / 2
		rect.Min.Y += addY
		rect.Max.Y += addY
	}

	draw.Draw(destImg, rect, srcImg, image.Point{0, 0}, draw.Over)
	return nil
}

func TestVideoInfo(t *testing.T) {
	_, thisSrcPath, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(thisSrcPath), "sample.mkv")

	cmd := exec.Command("/usr/bin/mplayer", "-noidle", "-nogui", "-nosound", "-vo", "null", "-frames", "1", "-identify", path)
	//cmd := exec.Command("/usr/bin/mplayer", "-noidle", "-nogui", "-nosound", "-vo", "jpeg",
	//	"-frames", "5", "-sstep", "600", path)
	//	cmd.Dir = "/home/bukodi/Downloads/Rogue.One.2016.RETAiL.BDRiP.x264.HuN-HyperX/sample/"
	cmd.Dir = "/tmp/Arduino/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(out))
		//return
	}
	fmt.Println("Output: " + string(out))
}

func TestSMPlayerHash(t *testing.T) {
	path := "sample.mkv"
	hash, _ := smplayerHash(path)
	fmt.Printf("Hash is: %q\n", hash)
}

func TestMultipleSMPlayerHash(t *testing.T) {
	//baseDir := "/home/bukodi/Downloads/CSENGETETT,MYLORD;;1.2.3.ÉVAD/"
	baseDir := "/media/bukodi/DATA/DATA/installData/20170818/"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !suppordedExts[ext] {
			return nil
		}
		hash, e := smplayerHash(path)
		if e != nil {
			t.Error(path+": ", e)
		}
		fmt.Printf("%q\n", hash)
		return nil
	})
	if err != nil {
		t.Error("File walking err: ", err)
	}
}

// Calculates file hash as SMPlayer
// See https://app.assembla.com/spaces/smplayer/subversion/source/6227/smplayer/trunk/src/filehash.cpp
func smplayerHash(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	var a, hash uint64
	hash = uint64(info.Size())

	buff := make([]byte, 65536)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	_, err = f.Read(buff)
	if err != nil {
		return "", err
	}

	for i := 0; i < 8192; i++ {
		a = binary.LittleEndian.Uint64(buff[i*8:])
		hash += a
	}

	_, err = f.Seek(info.Size()-65536, 0)
	if err != nil {
		return "", err
	}
	_, err = f.Read(buff)
	if err != nil {
		return "", err
	}
	for i := 0; i < 8192; i++ {
		a = binary.LittleEndian.Uint64(buff[i*8:])
		hash += a
	}
	return fmt.Sprintf("%x", hash), nil
}
