// main.go
package main

import (
	"Gogame/thread"
)

// Gogame framework version.
const VERSION = "0.0.1"

func main() {
	thread.GetMaster().Wait_thread_over()
}
