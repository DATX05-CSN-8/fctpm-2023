package tpminstantiator

type tpmInstantiatorService struct {
	swtpmPath string
}

type TpmInstance struct {
	Id         string
	SocketPath string
	pid        int
}

func NewTpmInstantiatorService(swtpmPath string) *tpmInstantiatorService {
	return &tpmInstantiatorService{
		swtpmPath: swtpmPath,
	}
}

func (s *tpmInstantiatorService) Create() (*TpmInstance, error) {
	// TODO instantiate TPM
	return &TpmInstance{
		Id:         "demo",
		SocketPath: "/tmp/aaa",
	}, nil
}

func (s *tpmInstantiatorService) Destroy(instance *TpmInstance) error {
	// TODO destroy running TPM
	return nil
}
