package main

import (
	"fmt"
	"yuewen/store/svc/yuewen"
)

func main() {

	fmt.Println(yuewen.SignIn(false))
	// ips := viper.GetStringSlice("ips")
	// slices.
	// fmt.Println(ips)
	// // for i := 0; i < 1; i++ {
	// 	go func ()  {
	// 		yuewen.SignIn(false)
	// 	}()
	// }

	// time.Sleep(time.Minute)
	
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(svc.Chapters(ids[0].ID))
}
