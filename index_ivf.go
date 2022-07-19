package faiss

/*
#include <faiss/c_api/IndexIVF_c.h>
#include <faiss/c_api/gpu/GpuIndexIvf_c.h>

FaissIndexIVF* convert(FaissIndex *idx) {
	return (FaissIndexIVF*)(idx);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func SetNumProbes(index Index, numProbes int) error {
	var ivfidx *C.FaissIndexIVF
	ivfidx = C.convert(index.cPtr())
	if ivfidx == nil {
		return errors.New("Index cannot be converted to ivf")
	}
	C.faiss_IndexIVF_set_nprobe(ivfidx, C.size_t(numProbes))
	return nil
}

func GetProbes(index Index) (uint, error) {
	var ivfidx *C.FaissIndexIVF
	ivfidx = C.convert(index.cPtr())
	if ivfidx == nil {
		return 0, errors.New("Index cannot be converted to ivf")
	}
	var x C.size_t
	x = C.faiss_IndexIVF_nprobe(ivfidx)
	return uint(x), nil
}

func GetNumOfProbesofIndexGpu(index Index) int {
	return int(C.GetNumOfProbesGpu(unsafe.Pointer(index.cPtr())))
}

func SetAmountOfProbesGpu(index Index, probes int) {
	C.SetAmountOfProbes(unsafe.Pointer(index.cPtr()), C.int(probes))
}
