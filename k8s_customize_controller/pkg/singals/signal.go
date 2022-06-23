package singals

import "os"

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	//这里主协程就直接返回了
	return stop
}
