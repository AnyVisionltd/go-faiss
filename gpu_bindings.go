//go:build gpu
// +build gpu

package faiss

/*
#include <stddef.h>
#include <faiss/c_api/gpu/StandardGpuResources_c.h>
#include <faiss/c_api/gpu/GpuAutoTune_c.h>

size_t get_sizet() {
    const size_t value = 1;
	return value;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func TransferToGpu(index Index) (Index, error) {
	var gpuResources *C.FaissStandardGpuResources
	var gpuIndex *C.FaissGpuIndex
	c := C.faiss_StandardGpuResources_new(&gpuResources)
	if c != 0 {
		return nil, errors.New("error on init gpu %v")
	}

	exitCode := C.faiss_index_cpu_to_gpu(gpuResources, 0, index.cPtr(), &gpuIndex)

	if exitCode != 0 {
		return nil, errors.New("error transferring to gpu")
	}

	return &faissIndex{idx: gpuIndex, resource: gpuResources}, nil
}

func TransferToAllGPUs(index Index, gpuIndexes []int) (Index, error) {
	const amountOfGPUs int = 1
	var gpuResources [amountOfGPUs]*C.FaissStandardGpuResources
	for i := 0; i < len(gpuIndexes); i++ {
		var resourceIndex *C.FaissStandardGpuResources
		gpuResources[i] = resourceIndex
	}

	var gpuIndex *C.FaissGpuIndex
	for i := 0; i < amountOfGPUs; i++ {
		c := C.faiss_StandardGpuResources_new(&gpuResources[i])
		if c != 0 {
			return nil, errors.New("error on init gpu %v")
		}
	}

	var gpuIndexesAsConst [amountOfGPUs]int
	for i, value := range gpuIndexes {
		//const val int  = value
		gpuIndexesAsConst[i] = value
	}
	// todo : handle c.size_t
	//exitCode := C.faiss_index_cpu_to_gpu_multiple((**C.FaissStandardGpuResources)(unsafe.Pointer(&gpuResourcesAsConst[0])), (*C.int)(unsafe.Pointer(&gpuIndexesAsConst[0])),
	//	C.size_t , index.cPtr(), &gpuIndex)
	exitCode := C.faiss_index_cpu_to_gpu_multiple(
		(**C.FaissStandardGpuResources)(unsafe.Pointer(&gpuResources[0])),
		(*C.int)(unsafe.Pointer(&gpuIndexesAsConst[0])),
		C.size_t(1),
		index.cPtr(),
		&gpuIndex)

	if exitCode != 0 {
		return nil, errors.New("error transferring to gpu")
	}
	// todo : pass resouces instead of resource and free all resources
	return &faissIndex{idx: gpuIndex, resource: nil}, nil

	//return &faissIndex{idx: gpuIndex, resource: gpuResources}, nil
}

func TransferToCpu(gpuIndex Index) (Index, error) {
	var cpuIndex *C.FaissIndex

	exitCode := C.faiss_index_gpu_to_cpu(gpuIndex.cPtr(), &cpuIndex)
	if exitCode != 0 {
		return nil, errors.New("error transferring to gpu")
	}

	Free(gpuIndex)

	return &faissIndex{idx: cpuIndex}, nil
}

func Free(index Index) {
	var gpuResource *C.FaissStandardGpuResources
	gpuResource = index.cGpuResource()
	C.faiss_StandardGpuResources_free(gpuResource)
	index.Delete()
}

func CreateGpuIndex() (Index, error) {
	var gpuResource *C.FaissStandardGpuResources
	var gpuIndex *C.FaissGpuIndex
	c := C.faiss_StandardGpuResources_new(&gpuResource)
	if c != 0 {
		return nil, errors.New("error on init gpu %v")
	}

	return &faissIndex{idx: gpuIndex, resource: nil}, nil
}
