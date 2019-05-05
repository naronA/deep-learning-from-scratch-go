package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/naronA/zero_deeplearning/mat"
	"github.com/naronA/zero_deeplearning/mnist"
	"github.com/naronA/zero_deeplearning/network"
	"github.com/naronA/zero_deeplearning/vec"
)

// ハイパーパラメタ
const (
	ImageLength  = 784
	ItersNum     = 100000
	BatchSize    = 100
	Hidden       = 50
	LearningRate = 0.0001
)

func MnistMatrix(set *mnist.DataSet) (*mat.Matrix, *mat.Matrix) {
	size := len(set.Labels)
	image := vec.Vector{}
	label := vec.Vector{}
	for i := 0; i < size; i++ {
		image = append(image, set.Images[i]...)
		label = append(label, set.Labels[i]...)
	}
	x, _ := mat.NewMatrix(size, ImageLength, image)
	t, _ := mat.NewMatrix(size, 10, label)
	return x, t
}

func train() {

	train, test := mnist.LoadMnist()

	TrainSize := len(train.Labels)
	net := network.NewTwoLayerNet(ImageLength, Hidden, 10, 0.01)

	trainLossList := []float64{}
	trainAccList := []float64{}
	testAccList := []float64{}
	xTrain, tTrain := MnistMatrix(train)
	xTest, tTest := MnistMatrix(test)
	iterPerEpoch := func() int {
		if TrainSize/BatchSize > 1.0 {
			return TrainSize / BatchSize
		}
		return 1
	}()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < ItersNum; i++ {
		start := time.Now()
		batchIndices := rand.Perm(TrainSize)[:BatchSize]
		image := vec.Vector{}
		label := vec.Vector{}
		for _, v := range batchIndices {
			image = append(image, train.Images[v]...)
			label = append(label, train.Labels[v]...)
		}

		xBatch, _ := mat.NewMatrix(BatchSize, ImageLength, image)
		tBatch, _ := mat.NewMatrix(BatchSize, 10, label)
		grads := net.Gradient(xBatch, tBatch)
		newParams := map[string]*mat.Matrix{}
		keys := []string{"W1", "b1", "W2", "b2"}
		for _, k := range keys {
			mullr := grads[k].Mul(LearningRate)
			newParams[k] = net.Params[k].Sub(mullr)
		}
		net.UpdateParams(newParams)

		loss := net.Loss(xBatch, tBatch)

		if i%iterPerEpoch == 0 && i >= iterPerEpoch {
			trainAcc := net.Accuracy(xTrain, tTrain)
			testAcc := net.Accuracy(xTest, tTest)
			// trainAccList = append(trainAccList, trainAcc)
			// testAccList = append(testAccList, testAcc)
			end := time.Now()
			fmt.Printf("elapstime = %v loss = %v\n", end.Sub(start), loss)
			fmt.Printf("train acc / test acc = %v / %v\n", trainAcc, testAcc)
		}
	}
	fmt.Println(trainLossList)
	fmt.Println(trainAccList)
	fmt.Println(testAccList)
}

func main() {
	train()
}