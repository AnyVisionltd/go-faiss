package faiss

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestFlatIndexOnGpu(t *testing.T) {
	index, err := NewIndexFlatL2(1)
	if err != nil {
		log.Fatal(err)
	}

	idx, err := TransferToGpu(index)
	if err != nil {
		log.Fatal(err)
	}
	vectorsToAdd := []float32{1,2,3,4,5}
	err = idx.Add(vectorsToAdd)
	if err != nil {
		fmt.Println(err.Error())
	}
	distances, resultIds, err := idx.Search(vectorsToAdd, 5)
	fmt.Println(distances, resultIds, err)
	for i := range vectorsToAdd {
		require.Equal(t, int64(i), resultIds[len(vectorsToAdd)*i])
		require.Equal(t, float32(0), distances[len(vectorsToAdd)*i])
	}
}

func TestIndexIDMapOnGPU(t *testing.T) {
	index, err := NewIndexFlatL2(1)
	if err != nil {
		log.Fatal(err)
	}

	indexMap, err := NewIndexIDMap(index)
	if err != nil {
		fmt.Println(err.Error())
	}
	idx, err := TransferToGpu(indexMap)
	if err != nil {
		log.Fatal(err)
	}
	vectorsToAdd := []float32{1,2,3,4,5}
	ids := make([]int64, len(vectorsToAdd))
	for i := 0; i < len(vectorsToAdd); i++ {
		ids[i] = int64(i)
	}

	err = idx.AddWithIDs(vectorsToAdd, ids)
	if err != nil {
		fmt.Println(err.Error())
	}
	distances, resultIds, err := idx.Search(vectorsToAdd, 5)
	fmt.Println(idx.D(), idx.Ntotal())
	fmt.Println(distances, resultIds, err)
	for i := range vectorsToAdd {
		require.Equal(t, ids[i], resultIds[len(vectorsToAdd)*i])
		require.Equal(t, float32(0), distances[len(vectorsToAdd)*i])
	}
}
