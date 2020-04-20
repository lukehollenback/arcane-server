package handlers

import "log"

//
// PkgInit initializes this package so that it can be subsequently used by the software. Should be
// called at the very beginning of program execution.
//
func PkgInit() {
	// NOTE: For this particular package, no initialization is necessary. This is simply called in
	//  order to cause the various "init()" functions in this package to be fired by the runtime.

	log.Print("The \"handler\" package has been initialized.")
}
