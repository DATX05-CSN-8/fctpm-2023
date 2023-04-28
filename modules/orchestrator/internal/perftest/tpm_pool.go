package perftest

import "github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"

type TpmPoolTouplpe struct {
	alloc    tpmallocator
	instance *tpminstantiator.TpmInstance
}
type TpmPool struct {
	tpmq []*TpmPoolTouplpe
	//readyq []tpmallocator
}

func NewTpmPool(size int, basepath string) (*TpmPool, error) {

	s := TpmPool{
		tpmq: make([]*TpmPoolTouplpe, 0),
		//readyq: make([]tpmallocator, 0),
	}

	for i := 0; i < size; i++ {

		// AAA todo basepath + strconv.Itoa(*inum)
		tpmalloc := tpminstantiator.NewTpmInstantiatorServiceWithBasePath(basepath)

		tpminstance, err := tpmalloc.Allocate()
		if err != nil {
			return nil, err
		}

		s.tpmq = append(s.tpmq, &TpmPoolTouplpe{alloc: tpmalloc, instance: tpminstance})
	}

	return &s, nil
}
