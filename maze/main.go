package main

func main() {
	setupLogging()
	done := make(chan bool, 1)
	//go newHTTPProxy()
	//go newUDPProxy()
	go winRace()

	// prevent program from closing since
	// both proxies are in separate goroutines
	<-done
}

// capture 2: three keypresses in all four directions, three jumps
