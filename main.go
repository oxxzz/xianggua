package main

import (
	"fmt"
	"runtime"
	"yuewen/store/svc/origin"
	"yuewen/store/svc/yuewen"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	// s := make(chan os.Signal, 1)
	// defer close(s)
	// signal.Notify(s, os.Interrupt, syscall.SIGINT)
	// go func() { svc.Process(s) }()
	// if viper.GetBool("enableAPI") {
	// 	go func() {
	// 		if err := api.Setup(true).Run(viper.GetString("listen")); err != nil {
	// 			logrus.Errorf("api.Setup.Run failed: %v", err)
	// 			s <- syscall.SIGINT
	// 		}
	// 	}()
	// }

	// <-s
	// time.Sleep(time.Second * 5)
	// logrus.Debugf("shutdown")

	books, err := origin.Books()
	if err != nil {
		panic(err)
	}

	for _, v := range books {
		book, err := origin.Info(v.ID)
		if err != nil {
			panic(err)
		}

		id, err := yuewen.PushBookInfo(book)
		fmt.Println(id, err)
	}
}
