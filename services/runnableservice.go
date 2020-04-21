package services

//
// RunnableService provides a generic interface for interacting with services that require explicit
// start-up and shut-down procedures.
//
type RunnableService interface {
	//
	// Start initiates the start-up process for the service.
	//
	// It returns a channel that can be blocked on until a "true" signal is recieved over it –
	// indicating that the service has completed its start-up process. It is up to the caller to
	// ensure that subsequent calls are not performed until said signal is recieved. Failing to do so
	// may cause a corrupt program state.
	//
	Start() (<-chan bool, error)

	//
	// Stop initiates the shut-down process for the service.
	//
	// It returns a channel that can be blocked on until a "true" signal is recieved over it –
	// indicating that the service has completed its shut-down process. It is up to the caller to
	// ensure that subsequent calls are not performed until said signal is recieved. Failing to do so
	// may cause a corrupt program state.
	//
	Stop() (<-chan bool, error)
}
