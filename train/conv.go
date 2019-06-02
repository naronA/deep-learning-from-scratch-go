package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/naronA/zero_deeplearning/mnist"
	"github.com/naronA/zero_deeplearning/network"
	"github.com/naronA/zero_deeplearning/num"
	"github.com/naronA/zero_deeplearning/optimizer"
	"github.com/naronA/zero_deeplearning/vec"
)

// ハイパーパラメタ
const (
	ImageLength  = 784
	ItersNum     = 100000
	BatchSize    = 100
	Hidden       = 50
	LearningRate = 0.0001
	MNIST        = 10
)

func MnistTensor4D(set *mnist.DataSet) (num.Tensor4D, num.Tensor4D) {
	size := len(set.Labels)
	image := num.Tensor4D{}
	label := num.Tensor4D{}
	// label := vec.Vector{}

	for i := 0; i < size; i++ {
		// image = append(image, set.Images[i]...)
		// label = append(label, set.Labels[i]...)
		mat := &num.Matrix{
			Vector:  set.Images[i],
			Rows:    28,
			Columns: 28,
		}
		t3d := num.Tensor3D{}
		t3d = append(t3d, mat)
		image = append(image, t3d)

		labelMat := &num.Matrix{
			Vector:  set.Labels[i],
			Rows:    1,
			Columns: 10,
		}
		labelT3d := num.Tensor3D{}
		labelT3d = append(labelT3d, labelMat)
		label = append(label, labelT3d)
	}

	return image, label
}

func MnistMatrix(set *mnist.DataSet) (*num.Matrix, *num.Matrix) {
	size := len(set.Labels)
	image := vec.Vector{}
	label := vec.Vector{}
	for i := 0; i < size; i++ {
		image = append(image, set.Images[i]...)
		label = append(label, set.Labels[i]...)
	}
	x, _ := num.NewMatrix(size, ImageLength, image)
	t, _ := num.NewMatrix(size, 10, label)
	return x, t
}

func train() {

	train, test := mnist.LoadMnist()

	TrainSize := len(train.Labels)
	// opt := optimizer.NewSGD(LearningRate)
	// opt := optimizer.NewMomentum(LearningRate)
	// opt := optimizer.NewAdaGrad(LearningRate)
	opt := optimizer.NewAdamAny(LearningRate)
	// weightDecayLambda := 0.1
	// net := network.NewMultiLayer(opt, ImageLength, Hidden, MNIST, weightDecayLambda)
	// net := network.NewTwoLayerNet(opt, ImageLength, Hidden, MNIST, weightDecayLambda)

	inputDim := &network.InputDim{
		Channel: 1,
		Height:  28,
		Weidth:  28,
	}

	convParams := &network.ConvParams{
		FilterNum:  30,
		FilterSize: 5,
		Pad:        0,
		Stride:     1,
	}

	net := network.NewSimpleConvNet(opt, inputDim, convParams, 100, 10, 0.01)

	xTrain, tTrain := MnistTensor4D(train)
	xTest, tTest := MnistTensor4D(test)
	iterPerEpoch := func() int {
		// if TrainSize/BatchSize > 1.0 {
		// 	return TrainSize / BatchSize
		// }
		return 1
	}()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < ItersNum; i++ {
		start := time.Now()
		batchIndices := rand.Perm(TrainSize)[:BatchSize]
		xBatch := num.Tensor4D{}
		tBatch := num.Tensor4D{}
		for _, v := range batchIndices {
			mat := &num.Matrix{
				Vector:  train.Images[v],
				Rows:    28,
				Columns: 28,
			}
			t3d := num.Tensor3D{}
			t3d = append(t3d, mat)
			xBatch = append(xBatch, t3d)
			labelMat := &num.Matrix{
				Vector:  train.Labels[v],
				Rows:    1,
				Columns: 10,
			}
			labelT3d := num.Tensor3D{}
			labelT3d = append(labelT3d, labelMat)
			tBatch = append(tBatch, labelT3d)
		}

		grads := net.Gradient(xBatch, tBatch)
		net.UpdateParams(grads)
		loss := net.Loss(xBatch, tBatch)

		if i%iterPerEpoch == 0 && i >= iterPerEpoch {
			trainAcc := net.Accuracy(xTrain, tTrain)
			testAcc := net.Accuracy(xTest, tTest)
			end := time.Now()
			fmt.Printf("elapstime = %v loss = %v\n", end.Sub(start), loss)
			fmt.Printf("train acc / test acc = %v / %v\n", trainAcc, testAcc)
		}
	}
}

func main() {
	train()
}
