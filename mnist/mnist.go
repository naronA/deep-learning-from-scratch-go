package mnist

import (
	"compress/gzip"
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/naronA/zero_deeplearning/vec"
)

const (
	imageMagic = 0x00000803
	labelMagic = 0x00000801
	Width      = 28
	Height     = 28
)

type DataSet struct {
	Rows      int
	Cols      int
	RawImages []RawImage
	Images    []vec.Vector
	Labels    []vec.Vector
}

func LoadMnist() (*DataSet, *DataSet) {

	trainImagesFile, err := os.Open("./mnist/train-images-idx3-ubyte.gz")
	if err != nil {
		return nil, nil
	}
	defer trainImagesFile.Close()
	trainLabelsFile, err := os.Open("./mnist/train-labels-idx1-ubyte.gz")
	if err != nil {
		return nil, nil
	}
	defer trainLabelsFile.Close()
	testImagesFile, err := os.Open("./mnist/t10k-images-idx3-ubyte.gz")
	if err != nil {
		return nil, nil
	}
	defer testImagesFile.Close()
	testLabelsFile, err := os.Open("./mnist/t10k-labels-idx1-ubyte.gz")
	if err != nil {
		return nil, nil
	}
	defer testLabelsFile.Close()

	trainRows, trainColumns, trainImages, trainFImages := readImages(trainImagesFile)
	trainLabels := readLabels(trainLabelsFile)
	train := &DataSet{
		Rows:      trainRows,
		Cols:      trainColumns,
		RawImages: trainImages,
		Images:    trainFImages,
		Labels:    trainLabels,
	}
	testRows, testColumns, testImages, testFImages := readImages(testImagesFile)
	testLabels := readLabels(testLabelsFile)
	test := &DataSet{
		Rows:      testRows,
		Cols:      testColumns,
		RawImages: testImages,
		Images:    testFImages,
		Labels:    testLabels,
	}
	return train, test
}

func oneHot(n uint8) vec.Vector {
	oneHot := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	oneHot[n] = 1
	return oneHot
}

func readLabels(file io.Reader) []vec.Vector {
	r, err := gzip.NewReader(file)
	if err != nil {
		return nil
	}
	defer r.Close()

	var (
		magic int32
		n     int32
	)
	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		return nil
	}
	if magic != labelMagic {
		return nil
	}
	if err := binary.Read(r, binary.BigEndian, &n); err != nil {
		return nil
	}
	// N個のラベルデータが含まれているのでN要素の配列をつくる
	labels := make([]vec.Vector, n)
	for i := 0; i < int(n); i++ {
		var num uint8
		if err := binary.Read(r, binary.BigEndian, &num); err != nil {
			return nil
		}
		// labels[i].number = num
		// labels[i].oneHot = oneHot(num)
		labels[i] = oneHot(num)
	}
	return labels
}

type RawImage []byte

func (img RawImage) ColorModel() color.Model {
	return color.GrayModel
}

func (img RawImage) At(x, y int) color.Color {
	return color.Gray{img[y*Width+x]}
}

func (img RawImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{Width, Height},
	}
}

func readImages(file io.Reader) (int, int, []RawImage, []vec.Vector) {
	r, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	var (
		magic int32
		n     int32
		nrow  int32
		ncol  int32
	)

	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		panic(err)
	}
	if magic != imageMagic {
		panic(err)
	}
	if err := binary.Read(r, binary.BigEndian, &n); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.BigEndian, &nrow); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.BigEndian, &ncol); err != nil {
		panic(err)
	}
	// N個のラベルデータが含まれているのでN要素の配列をつくる
	imgs := make([]RawImage, n)
	fimgs := make([]vec.Vector, n)
	m := int(nrow * ncol)
	for i := 0; i < int(n); i++ {
		imgs[i] = make(RawImage, m)
		fimgs[i] = make(vec.Vector, m)
		mfull, err := io.ReadFull(r, imgs[i])

		if err != nil {
			panic(err)
		}
		if mfull != m {
			return 0, 0, nil, nil
		}

		for j, b := range imgs[i] {
			fimgs[i][j] = float64(b)
		}
		// fmt.Println(imgs[i])
		// fmt.Println(fimgs[i])
		// panic(nil)
	}
	return int(nrow), int(ncol), imgs, fimgs
}
