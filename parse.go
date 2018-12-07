package unbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
)

func Test() {
	print("test")
}

func ReadJson(filename string) *Pack {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	p := &Pack{}
	err = json.Unmarshal(data, p)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("out,%#v\n", p.Frames)
	for k, v := range p.Frames {
		fd := &FrameData{}
		mapstructure.Decode(v, fd)
		fd.Name = k
		p.FrameList = append(p.FrameList, fd)
		// fmt.Printf("out,%#v\n", fd)
	}

	return p
}

func readImage(filename string) (image.Image, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func UnpackDir(dirName, outPath, ext string) error {
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			fmt.Printf("visited file or dir: %q, %v\n", path, info.IsDir())
			ext2 := filepath.Ext(path)
			if ext2 == ext {
				outPath2 := filepath.Join(outPath, filepath.Dir(path)[len(dirName):])
				fmt.Printf("unpack: %s\n", path)
				UnpackFile(path, outPath2)
				// fmt.Printf("out dir: %s\n", filepath.Join(outPath, filepath.Dir(path)[len(dirName):]))
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dirName, err)
		return err
	}
	return nil
}

type SubImageer interface {
	SubImage(r image.Rectangle) image.Image
}

func UnpackFile(filename, outPath string) error {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	var baseName = path.Base(name)
	var imgName = name + ".png"
	data := ReadJson(filename)
	img, err := readImage(imgName)
	if err != nil {
		return err
	}

	mSize1 := data.Meta.Size
	mSize2 := img.Bounds().Max

	// fmt.Printf("img,data,%v\n", data.Meta)
	// fmt.Printf("img,b,%v\n", img.Bounds())
	outPath = path.Join(outPath, baseName)
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err = os.MkdirAll(outPath, 0777)
		if err != nil {
			return err
		}
	}

	if mSize1.Width != mSize2.X || mSize1.Height != mSize2.Y {
		return errors.New("配置, 资源大小不一致")
	}

	// fmt.Printf("img,b,%v\n", img.ColorModel())

	subImageer, ok := img.(SubImageer)
	if !ok {
		return errors.New("不支持SubImage")
	}
	for _, f := range data.FrameList {
		frame := f.Frame
		pngName := f.Name
		rect := image.Rect(frame.X, frame.Y, frame.X+frame.W, frame.Y+frame.H)
		subImg := subImageer.SubImage(rect)
		outFileName := path.Join(outPath, pngName)
		err = Save(outFileName, &subImg)
		if err != nil {
			return err
		}
	}
	return nil
}

func Save(filename string, img *image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, *img)
}
