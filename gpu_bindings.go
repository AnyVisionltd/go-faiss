package faiss

/*
#include <faiss/c_api/gpu/StandardGpuResources_c.h>
#include <faiss/c_api/gpu/GpuAutoTune_c.h>
*/
import "C"
import (
	"errors"
)

func TransferToGpu(index Index) (Index, error) {
	var gpuResources *C.FaissStandardGpuResources
	var gpuIndex *C.FaissGpuIndex
	c := C.faiss_StandardGpuResources_new(&gpuResources)
	if c != 0 {
		return nil, errors.New("error on init gpu")
	}

	exitCode := C.faiss_index_cpu_to_gpu(gpuResources, 0, index.cPtr(), &gpuIndex)
	if exitCode != 0 {
		return nil, errors.New("error gpu blabla")
	}

	return &faissIndex{idx: gpuIndex}, nil
}
