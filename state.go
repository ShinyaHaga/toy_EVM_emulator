package main

import (
	"fmt"

	"github.com/holiman/uint256"
)

type WorldState struct {
	mapping map[[20]byte]*Account
}

type Account struct {
	nonce   uint
	balance uint64
	//codeHash
	//storageRoot

	storage map[uint256.Int]uint256.Int
	code    []byte
}

//EOAを作成する
func (w *WorldState) CreateEOA(addr [20]byte) {
	w.mapping[addr] = &Account{
		nonce:   0,
		balance: 0,
		storage: map[uint256.Int]uint256.Int{},
		code:    []byte{},
	}
}

//コントラクトアカウントを作成する
func (w *WorldState) CreateCA(addr [20]byte, input []byte) {
	w.mapping[addr] = &Account{
		nonce:   0,
		balance: 0,
		storage: map[uint256.Int]uint256.Int{},
		code:    []byte{},
	}
	env := Environment{}
	env = env.new(addr, [20]byte{0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a, 0x1a}, 10000000, 100000000000000000, 0, *w.mapping[addr])
	fmt.Println("CA env: ", env)
	env.set_code(input)
	vm := EVM{}
	vm = vm.new(env)
	vm.exec_transaction(*w.mapping[addr])
	w.mapping[addr].code = vm.returns
}

//コントラクトアカウントの関数を実行する
func (w *WorldState) run(addr [20]byte, input []byte) {

}
