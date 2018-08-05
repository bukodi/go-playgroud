package thumbnail

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type Content struct {
	start  int
	length int
	ratio  float64
	width  int
	height int
	fps    float64
}

var suppordedExts = []string{"mp4", "mkv", "avi"}

func TestMultipleVideos(t *testing.T) {
	baseDir := "/home/bukodi/Downloads/CSENGETETT,MYLORD;;1.2.3.ÉVAD/"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".avi") {
			return nil
		}
		e := procVideo(path)
		if e != nil {
			t.Error(path+": ", e)
		}
		return nil
	})
	if err != nil {
		t.Error("Final err: ", err)
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
	err = makeTile2x2(&c, "/tmp/Arduino/")
	if err != nil {
		return err
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

func makeTile2x2(c *Content, tmpDir string) error {
	var destImg = image.NewRGBA(image.Rect(0, 0, c.width*2, c.height*2))
	{ // Top left
		img02, err := readImage(filepath.Join(tmpDir, "00000002.jpg"))
		if err != nil {
			return err
		}
		draw.Draw(destImg, image.Rect(0, 0, c.width, c.height), img02, image.Point{0, 0}, draw.Over)
	}

	{ // Top right
		img03, err := readImage(filepath.Join(tmpDir, "00000003.jpg"))
		if err != nil {
			return err
		}
		draw.Draw(destImg, image.Rect(c.width, 0, c.width*2, c.height), img03, image.Point{0, 0}, draw.Over)
	}

	{ // Bottom left
		img04, err := readImage(filepath.Join(tmpDir, "00000004.jpg"))
		if err != nil {
			return err
		}
		draw.Draw(destImg, image.Rect(0, c.height, c.width, c.height*2), img04, image.Point{0, 0}, draw.Over)
	}

	{ // Bottom right
		img05, err := readImage(filepath.Join(tmpDir, "00000005.jpg"))
		if err != nil {
			return err
		}
		draw.Draw(destImg, image.Rect(c.width, c.height, c.width*2, c.height*2), img05, image.Point{0, 0}, draw.Over)
	}

	out, err := os.Create(filepath.Join(tmpDir, "tn.jpg"))
	if err != nil {
		return err
	}
	defer out.Close()

	option := jpeg.Options{Quality: 50}
	err = jpeg.Encode(out, destImg, &option)
	if err != nil {
		return err
	}
	return nil
}

func TestVideoInfo(t *testing.T) {
	path := "/home/bukodi/Downloads/Rogue.One.2016.RETAiL.BDRiP.x264.HuN-HyperX/rogueone-sd-hyperx.mkv"
	//path := "/home/bukodi/Downloads/Rogue.One.2016.RETAiL.BDRiP.x264.HuN-HyperX/sample/sample.mkv"
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
