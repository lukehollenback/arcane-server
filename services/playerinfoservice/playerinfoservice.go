package playerinfoservice

import (
	"log"
	"sync"
)

var (
	o    *PlayerInfoService
	once sync.Once
)

//
// PlayerInfoService represents an instance of the Player Info Service, which provides efficient
// access to data (e.g. usernames) about players.
//
type PlayerInfoService struct {
	// TODO: Currently this service is not stateful. However, in the future, it will need to maintain
	//  connections to databases and caches of player info.
}

//
// Instance provides a singleton instance of the service.
//
func Instance() *PlayerInfoService {
	once.Do(func() {
		o = &PlayerInfoService{}
	})

	return o
}

//
// Start implements the method defined by the services.Stop() interface.
//
func (o *PlayerInfoService) Start() (<-chan bool, error) {
	log.Printf("The Player Info Service is starting...")

	ch := make(chan bool, 1)

	ch <- true

	return ch, nil
}

//
// Stop implements the method defined by the services.Stop() interface.
//
func (o *PlayerInfoService) Stop() (<-chan bool, error) {
	log.Printf("The Player Info Service is stopping...")

	ch := make(chan bool, 1)

	ch <- true

	return ch, nil
}

//
// GetUsername retrieves (e.g. from database or cache) the username of the player with the specified
//  player ID.
//
func (o *PlayerInfoService) GetUsername(playerID string) string {
	return "Zaedaux"
}
