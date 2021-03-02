package main

import "fmt"

func mainP() {
	var a string
	fmt.Print("deploy または transaction と入力してください: ")
	fmt.Scan(&a)

	switch a {
	case "deploy":
		fmt.Println("deploy")
	case "transaction":
		fmt.Println("transaction")
	default:
		panic("Error:指定されていないコマンド")
	}

}

func main() {
	w := WorldState{
		mapping: map[[20]byte]*Account{},
	}
	w.CreateEOA([20]byte{0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a})
	w.CreateCA([20]byte{}, str_to_bytes("6080604052600460005534801561001557600080fd5b5060c7806100246000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea2646970667358221220daf43c1797a846e0f0e14bdc656fcde79ac0a05a083e3ceabb8ef8dc713c360a64736f6c63430007000033"))
	fmt.Println(w)
	fmt.Println(w.mapping[[20]byte{0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a}])
	fmt.Println(w.mapping[[20]byte{}])
}
