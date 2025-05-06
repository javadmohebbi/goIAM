// Command goiam starts the goIAM identity and access management API server.
// This is the top-level entry point that delegates startup to the internal server package.
package main

import "github.com/javadmohebbi/goIAM/cmd/server"

// main initializes and starts the goIAM server using the cmd/server package.
func main() {
	server.Main()
}
