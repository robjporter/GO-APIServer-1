package main

import (
	"fmt"
	conf "github.com/roporter/go-libs/go-config"
)

func main() {
	config := conf.NewConfig()
	config,_ = conf.ReadFromFile("../config/config.json")
	tmp := config.Get("security.blocked")
	fmt.Println(tmp)
	tmp2 := config.GetKeyAsStringArray("security.blocked")
	fmt.Println(tmp2)
	fmt.Println(len(tmp2))
	
	
	//tmp2 := tmp.(map[string]interface{})
	//tmp2 := config.GetKey("security.blocked")
	//fmt.Println(tmp2)
	//tmp3 := config.GetStringArray("security.blocked")
	//fmt.Println(tmp3)
}