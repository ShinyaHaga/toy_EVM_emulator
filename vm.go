package main

import (
	"fmt"
	"strconv"

	"github.com/holiman/uint256"
)

//Environment
type Environment struct {
	current_executer [20]byte //直前の実行者（internal transactionの時使用？）
	code_owner       [20]byte //実行するコントラクトのオーナー
	sender           [20]byte //トランザクションの送信者
	gas_price        uint64   //Gas_Price
	calldata         []byte   //トランザクション実行時に渡されるデータ(引数など)
	code             []byte   //実行されるEVMバイトコード
	value            uint64   //送信するETHの量
	accountstate     Account
	//・blockheader
	//・CALL，CREATEの数
	gas_limit uint64 //Gas_Limit
}

func (a *Environment) new(_code_owner [20]byte, _sender [20]byte, _gas_price uint64, _gas_limit uint64, _value uint64, _accountstate Account) Environment {
	_Environment := Environment{_code_owner, _sender, _sender, _gas_price, nil, nil, _value, _accountstate, _gas_limit}
	return _Environment
}

//EVMバイトコードをセット
func (a *Environment) set_code(_code []byte) {
	a.code = _code
}

//インプットデータをセット
func (a *Environment) set_calldata(_calldata []byte) {
	a.calldata = _calldata
}

//EVMの構造体
type EVM struct {
	env     Environment   //トランザクションの構造体
	pc      uint64        //プログラムカウンター
	gas     uint64        //Gasの残量
	sp      uint64        //スタックポインター
	stack   []uint256.Int //スタック領域
	memory  [1000]byte    //メモリー領域
	returns []byte        //リターン領域
	asm     []string      //実行した命令を保存
}

//EVMの初期化
func (b *EVM) new(env Environment) EVM {
	fmt.Println("env: ", env)
	fmt.Println("gas_limit: ", env.gas_limit)
	fmt.Println("gas_price: ", env.gas_price)
	gas := env.gas_limit / env.gas_price
	_EVM := EVM{env, 0, gas, 0, nil, [1000]byte{}, nil, nil}
	return _EVM
}

//EVMのスタックへpush
func (b *EVM) push(value *uint256.Int) {
	b.stack = append(b.stack, *value)
	b.sp++
	fmt.Print("push => stack: ")
	fmt.Println(fmt.Sprintf("%02x", b.stack))
}

//EVMのスタックからpop
func (b *EVM) pop() uint256.Int {
	value := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	b.sp--
	fmt.Print("pop => stack: ")
	fmt.Println(fmt.Sprintf("%02x", b.stack))
	return value
}

//EVMバイトコードの１命令を実行
func (b *EVM) exec(contract Account) bool {
	opcode := b.env.code[b.pc]
	b.pc++

	switch opcode {
	//0x00~: 算術命令
	case 0x00:
		b.op_stop()
	case 0x01:
		b.op_add()
	case 0x02:
		b.op_mul()
	case 0x03:
		b.op_sub()
	case 0x04:
		b.op_div()
	case 0x05:
		b.op_sdiv()
	case 0x06:
		b.op_mod()
	case 0x07:
		b.op_smod()
	case 0x08:
		b.op_addmod()
	case 0x09:
		b.op_mulmod()
	case 0x0a:
		b.op_exp()
	case 0x0b:
		b.op_sig_next_end()

	//0x10~: bit演算
	case 0x10:
		b.op_lt()
	case 0x11:
		b.op_gt()
	case 0x12:
		b.op_slt()
	case 0x13:
		b.op_sgt()
	case 0x14:
		b.op_eq()
	case 0x15:
		b.op_is_zero()
	case 0x16:
		b.op_and()
	case 0x17:
		b.op_or()
	case 0x18:
		b.op_xor()
	case 0x19:
		b.op_not()
	case 0x1a:
		b.op_byte()
	case 0x1b:
		b.op_shl()
	case 0x1c:
		b.op_shr()
	case 0x1d:
		b.op_sar()

	//0x20~: sha3
	case 0x20:
		b.op_sha3()

	//0x30~: Envirionment操作
	case 0x30:
		b.op_address()
	case 0x31:
		b.op_balance()
	case 0x32:
		b.op_origin()
	case 0x33:
		b.op_caller()
	case 0x34:
		b.op_callvalue()
	case 0x35:
		b.op_calldataload()
	case 0x36:
		b.op_calldatasize()
	case 0x37:
		b.op_calldatacopy()
	case 0x38:
		b.op_codesize()
	case 0x39:
		b.op_codecopy()
	case 0x3a:
		b.op_gasprice()
	case 0x3b:
		b.op_extcodesize()
	case 0x3c:
		b.op_extcodecopy()

	//0x50~: EVM内のステート操作
	case 0x50:
		b.op_pop()
	case 0x51:
		b.op_mload()
	case 0x52:
		b.op_mstore()
	case 0x53:
		b.op_mstore8()
	case 0x54:
		b.op_sload()
	case 0x55:
		b.op_sstore()
	case 0x56:
		b.op_jump()
	case 0x57:
		b.op_jumpi()
	case 0x58:
		b.op_pc()
	case 0x59:
		b.op_msize()
	case 0x5a:
		b.op_gas()
	case 0x5b:
		b.op_jumpdest()

	//0x60~: PUSH命令
	case 0x60:
		b.op_push(1)
	case 0x61:
		b.op_push(2)
	case 0x62:
		b.op_push(3)
	case 0x63:
		b.op_push(4)
	case 0x64:
		b.op_push(5)
	case 0x65:
		b.op_push(6)
	case 0x66:
		b.op_push(7)
	case 0x67:
		b.op_push(8)
	case 0x68:
		b.op_push(9)
	case 0x69:
		b.op_push(10)
	case 0x6a:
		b.op_push(11)
	case 0x6b:
		b.op_push(12)
	case 0x6c:
		b.op_push(13)
	case 0x6d:
		b.op_push(14)
	case 0x6e:
		b.op_push(15)
	case 0x6f:
		b.op_push(16)
	case 0x70:
		b.op_push(17)
	case 0x71:
		b.op_push(18)
	case 0x72:
		b.op_push(19)
	case 0x73:
		b.op_push(20)
	case 0x74:
		b.op_push(21)
	case 0x75:
		b.op_push(22)
	case 0x76:
		b.op_push(23)
	case 0x77:
		b.op_push(24)
	case 0x78:
		b.op_push(25)
	case 0x79:
		b.op_push(26)
	case 0x7a:
		b.op_push(27)
	case 0x7b:
		b.op_push(28)
	case 0x7c:
		b.op_push(29)
	case 0x7d:
		b.op_push(30)
	case 0x7e:
		b.op_push(31)
	case 0x7f:
		b.op_push(32)

	//0x80~: DUMP命令
	case 0x80:
		b.op_dump(1)
	case 0x81:
		b.op_dump(2)
	case 0x82:
		b.op_dump(3)
	case 0x83:
		b.op_dump(4)
	case 0x84:
		b.op_dump(5)
	case 0x85:
		b.op_dump(6)
	case 0x86:
		b.op_dump(7)
	case 0x87:
		b.op_dump(8)
	case 0x88:
		b.op_dump(9)
	case 0x89:
		b.op_dump(10)
	case 0x8a:
		b.op_dump(11)
	case 0x8b:
		b.op_dump(12)
	case 0x8c:
		b.op_dump(13)
	case 0x8d:
		b.op_dump(14)
	case 0x8e:
		b.op_dump(15)
	case 0x8f:
		b.op_dump(16)

	//0x90~: SWAP命令
	case 0x90:
		b.op_swap(1)
	case 0x91:
		b.op_swap(2)
	case 0x92:
		b.op_swap(3)
	case 0x93:
		b.op_swap(4)
	case 0x94:
		b.op_swap(5)
	case 0x95:
		b.op_swap(6)
	case 0x96:
		b.op_swap(7)
	case 0x97:
		b.op_swap(8)
	case 0x98:
		b.op_swap(9)
	case 0x99:
		b.op_swap(10)
	case 0x9a:
		b.op_swap(11)
	case 0x9b:
		b.op_swap(12)
	case 0x9c:
		b.op_swap(13)
	case 0x9d:
		b.op_swap(14)
	case 0x9e:
		b.op_swap(15)
	case 0x9f:
		b.op_swap(16)

	//0xa0~: LOG命令
	case 0xa0:
		b.op_log0()
	case 0xa1:
		b.op_log1()
	case 0xa2:
		b.op_log2()
	case 0xa3:
		b.op_log3()
	case 0xa4:
		b.op_log4()

	//0xf0~: SYSTEM命令
	case 0xf0:
		b.op_create()
	case 0xf1:
		b.op_call()
	case 0xf2:
		b.op_callcode()
	case 0xf3:
		b.op_return()
	// case 0xf4:
	// 	b.op_delegatecall()
	// case 0xf5:
	// 	b.op_create2()
	// case 0xfa:
	// 	b.op_staticcall()
	case 0xfd:
		b.op_revert()
		// case 0xff:
		// 	b.op_selfdestruct()
	case 0xfe:
		b.op_invalid()
	default:
		panic("OPCODE: not_implement")

	}

	switch opcode {
	case 0xf3:
		return true
	case 0xfd:
		return true
	default:
		return false
	}

}

//Gas残量の管理．なくなれば終了する
func (b *EVM) consume_gas(_gas uint64) {
	if b.gas >= _gas {
		b.gas = b.gas - _gas
	} else {
		panic("There is a shotage of gas")
	}
}

//トランザクションが終了するまでexecを続ける
func (b *EVM) exec_transaction(contract Account) {
	for int(b.pc) < len(b.env.code) {
		if b.exec(contract) {
			break
		}
	}
}

//逆アセンブル その１
func (b *EVM) disassemble(code string) {

}

//逆アセンブル　その２
func (b *EVM) push_asm(_opcode string) {
	b.asm = append(b.asm, _opcode)
}

//opcode一覧

//0x00: 何もしない
func (b *EVM) op_stop() {
	b.push_asm("STOP")
	fmt.Println("STOP")
}

//0x01: operand1(スタック1番目) + operand2(スタック2番目)
func (b *EVM) op_add() {
	b.consume_gas(3)
	b.push_asm("ADD")
	fmt.Println("ADD")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Add(&operand1, &operand2)
	b.push(result)
}

//0x02: operand1(スタック1番目) * operand2(スタック2番目)
func (b *EVM) op_mul() {
	b.consume_gas(5)
	b.push_asm("MUL")
	fmt.Println("MUL")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Mul(&operand1, &operand2)
	b.push(result)
}

//0x03: operand1(スタック1番目) - operand2(スタック2番目)
func (b *EVM) op_sub() {
	b.consume_gas(3)
	b.push_asm("SUB")
	fmt.Println("SUB")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Sub(&operand1, &operand2)
	b.push(result)
}

//0x04: operand1(スタック1番目) / operand2(スタック2番目)
func (b *EVM) op_div() {
	b.consume_gas(5)
	b.push_asm("DIV")
	fmt.Println("DIV")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Div(&operand1, &operand2)
	b.push(result)
}

//0x05 二の補数付き整数のDIV
func (b *EVM) op_sdiv() {
	b.consume_gas(5)
	b.push_asm("SDIV")
	fmt.Println("SDIV")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.SDiv(&operand1, &operand2)
	b.push(result)
}

//0x06 operand1(スタック1番目) % operand2(スタック2番目)
func (b *EVM) op_mod() {
	b.consume_gas(5)
	b.push_asm("MOD")
	fmt.Println("MOD")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Mod(&operand1, &operand2)
	b.push(result)
}

//0x07　二の補数付き整数のMod
func (b *EVM) op_smod() {
	b.consume_gas(5)
	b.push_asm("SMOD")
	fmt.Println("SMOD")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.SMod(&operand1, &operand2)
	b.push(result)
}

//0x08 {operand1(スタック1番目) + operand2(スタック2番目)} % operand3(スタック3番目)
func (b *EVM) op_addmod() {
	b.consume_gas(8)
	b.push_asm("ADDMOD")
	fmt.Println("ADDMOD")
	operand1 := b.pop()
	operand2 := b.pop()
	operand3 := b.pop()
	result := operand3.AddMod(&operand1, &operand2, &operand3)
	b.push(result)
}

//0x09 {operand1(スタック1番目) * operand2(スタック2番目)} % operand3(スタック3番目)
func (b *EVM) op_mulmod() {
	b.consume_gas(8)
	b.push_asm("MULMOD")
	fmt.Println("MULMOD")
	operand1 := b.pop()
	operand2 := b.pop()
	operand3 := b.pop()
	result := operand3.MulMod(&operand1, &operand2, &operand3)
	b.push(result)
}

//0x0a operand1(スタック1番目) ** operand2(スタック2番目)
func (b *EVM) op_exp() {
	b.consume_gas(10)
	b.push_asm("EXP")
	fmt.Println("EXP")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Exp(&operand1, &operand2)
	b.push(result)
}

//0x0b 二の補数付き整数の長さを拡張する．
//operand2が32byteならそのまま，
//それ以外なら(operand2 * 8 + 7）の符号付き整数と解釈し，32byteとする．
func (b *EVM) op_sig_next_end() {
	b.consume_gas(5)
	b.push_asm("SIGNEXTEND")
	fmt.Println("EXP")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.ExtendSign(&operand1, &operand2)
	b.push(result)
}

//0x10 operand1(スタック1番目) < operand2(スタック2番目)
//ture -> 1    false -> 0
func (b *EVM) op_lt() {
	b.consume_gas(3)
	b.push_asm("LT")
	fmt.Println("LT")
	operand1 := b.pop()
	operand2 := b.pop()
	if operand2.Lt(&operand1) {
		b.push(operand2.SetOne())
	} else {
		b.push(operand2.Clear())
	}
}

//0x11 operand1(スタック1番目) > operand2(スタック2番目)
//ture -> 1    false -> 0
func (b *EVM) op_gt() {
	b.consume_gas(3)
	b.push_asm("GT")
	fmt.Println("GT")
	operand1 := b.pop()
	operand2 := b.pop()
	if operand2.Gt(&operand1) {
		b.push(operand2.SetOne())
	} else {
		b.push(operand2.Clear())
	}
}

//0x12 operand1(スタック1番目) < operand2(スタック2番目) (符号付き整数)
//ture -> 1    false -> 0
func (b *EVM) op_slt() {
	b.consume_gas(3)
	b.push_asm("SLT")
	fmt.Println("SLT")
	operand1 := b.pop()
	operand2 := b.pop()
	if operand2.Slt(&operand1) {
		b.push(operand2.SetOne())
	} else {
		b.push(operand2.Clear())
	}
}

//0x13 operand1(スタック1番目) > operand2(スタック2番目) (符号付き整数)
//ture -> 1    false -> 0
func (b *EVM) op_sgt() {
	b.consume_gas(3)
	b.push_asm("GLT")
	fmt.Println("GLT")
	operand1 := b.pop()
	operand2 := b.pop()
	if operand2.Sgt(&operand1) {
		b.push(operand2.SetOne())
	} else {
		b.push(operand2.Clear())
	}
}

//0x14 operand1(スタック1番目) == operand2(スタック2番目)
//ture -> 1    false -> 0
func (b *EVM) op_eq() {
	b.consume_gas(3)
	b.push_asm("EQ")
	fmt.Println("EQ")
	operand1 := b.pop()
	operand2 := b.pop()
	if operand2.Eq(&operand1) {
		b.push(operand2.SetOne())
	} else {
		b.push(operand2.Clear())
	}

}

//0x15 operand1(スタック1番目) == 0
//ture -> 1    false -> 0
func (b *EVM) op_is_zero() {
	b.consume_gas(3)
	b.push_asm("ISZERO")
	fmt.Println("ISZERO")
	operand1 := b.pop()
	if operand1.IsZero() {
		b.push(operand1.SetOne())
	} else {
		b.push(operand1.SetOne())
	}
}

//0x16 operand1(スタック1番目) & operand2(スタック2番目)
func (b *EVM) op_and() {
	b.consume_gas(3)
	b.push_asm("AND")
	fmt.Println("AND")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.And(&operand1, &operand2)
	b.push(result)
}

//0x17 operand1(スタック1番目) | operand2(スタック2番目)
func (b *EVM) op_or() {
	b.consume_gas(3)
	b.push_asm("OR")
	fmt.Println("OR")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Or(&operand1, &operand2)
	b.push(result)
}

//0x18 operand1(スタック1番目) ~ operand2(スタック2番目)
func (b *EVM) op_xor() {
	b.consume_gas(3)
	b.push_asm("XOR")
	fmt.Println("XOR")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Xor(&operand1, &operand2)
	b.push(result)
}

//0x19 not operand1(スタック1番目)
func (b *EVM) op_not() {
	b.consume_gas(3)
	b.push_asm("NOT")
	fmt.Println("NOT")
	operand1 := b.pop()
	result := operand1.Not(&operand1)
	b.push(result)
}

//0x1a operand2(スタック2番目)のoperand1バイト目を取り出す
func (b *EVM) op_byte() {
	b.consume_gas(3)
	b.push_asm("BYTE")
	fmt.Println("BYTE")
	operand1 := b.pop()
	operand2 := b.pop()
	result := operand2.Byte(&operand1)
	b.push(result)
}

//0x1b shiftビット数だけ左にシフトされたスタック(value)をプッシュする
func (b *EVM) op_shl() {
	b.consume_gas(3)
	b.push_asm("SHL")
	fmt.Println("SHL")
	shift := b.pop()
	value := b.pop()
	if shift.LtUint64(256) {
		value.Lsh(&value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}
	b.push(&value)
}

//0x1c shiftビット数だけ右にシフトされたスタック(value)を0で埋めつつプッシュする．
func (b *EVM) op_shr() {
	b.consume_gas(3)
	b.push_asm("SHR")
	fmt.Println("SHR")
	shift := b.pop()
	value := b.pop()
	if shift.LtUint64(256) {
		value.Rsh(&value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}
	b.push(&value)
}

//0x1d shiftビット数だけ右にシフトされたスタック(value)をプッシュする(符号拡張版)
func (b *EVM) op_sar() {
	b.consume_gas(3)
	b.push_asm("SAR")
	fmt.Println("SAR")
	shift := b.pop()
	value := b.pop()
	if shift.GtUint64(256) {
		if value.Sign() >= 0 {
			value.Clear()
		} else {
			// Max negative shift: all bits set
			value.SetAllOne()
		}
	}
	n := uint(shift.Uint64())
	value.SRsh(&value, n)
	b.push(&value)
}

//0x20 sha3
func (b *EVM) op_sha3() {
	b.push_asm("SHA3")
	b.consume_gas(3)
	panic("SHA3: not implement")
}

//0x30 現在実行者中のアドレスをpush
func (b *EVM) op_address() {
	b.consume_gas(2)
	b.push_asm("ADDRESS")
	fmt.Println("ADDRESS")
	address := new(uint256.Int).SetBytes(b.env.code_owner[:])
	b.push(address)
}

//0x31 与えられたアドレスの残高をpush
func (b *EVM) op_balance() {
	b.consume_gas(400)
	b.push_asm("BALANCE")
	balance := new(uint256.Int).SetUint64(b.env.accountstate.balance)
	b.push(balance)
}

//0x32 コード実行元のアドレスをpush
func (b *EVM) op_origin() {
	b.consume_gas(2)
	b.push_asm("ORIGIN")
	fmt.Println("ORIGIN")
	address := new(uint256.Int).SetBytes(b.env.sender[:])
	b.push(address)
}

//0x33 実行するのアドレスをpush
func (b *EVM) op_caller() {
	b.consume_gas(2)
	b.push_asm("CALLER")
	fmt.Println("CALLER")
	address := new(uint256.Int).SetBytes(b.env.current_executer[:])
	b.push(address)
}

//0x34
func (b *EVM) op_callvalue() {
	b.consume_gas(2)
	b.push_asm("CALLVALUE")
	fmt.Println("CALLVALUE")
	value := b.env.value
	b.push(uint256.NewInt().SetUint64(uint64(value)))
}

//0x35 スタックからpopした値をstartとしてcalldataのstartの位置からstart+32の位置までの32byteのデータをstackにpush
func (b *EVM) op_calldataload() {
	b.consume_gas(3)
	b.push_asm("CALLDATALOAD")
	fmt.Println("CALLDATALOAD")
	x := b.pop()
	if offset, overflow := x.Uint64WithOverflow(); !overflow {
		data := b.env.calldata[offset : offset+32]
		b.push(x.SetBytes(data))
	} else {
		b.push(x.Clear())
	}
}

//0x36 calldataに格納されたデータサイズをstackにpush
func (b *EVM) op_calldatasize() {
	b.consume_gas(2)
	b.push_asm("CALLDATASIZE")
	fmt.Println("CALLDATASIZE")
	size := len(b.env.calldata)
	b.push(new(uint256.Int).SetUint64(uint64(size)))
}

//0x37 calldataに格納されたデータをmemory領域にコピー
func (b *EVM) op_calldatacopy() {
	b.consume_gas(9) //??
	b.push_asm("CALLDATACOPY")
	fmt.Println("CALLDATACOPY")
	memOffset := b.pop()
	dataOffset := b.pop()
	length := b.pop()

	dataOffset64, overflow := dataOffset.Uint64WithOverflow()
	if overflow {
		dataOffset64 = 0xffffffffffffffff
	}

	memOffset64 := memOffset.Uint64()
	length64 := length.Uint64()

	for i := 0; i < int(length64); i++ {
		b.memory[int(memOffset64)+i] = b.env.calldata[int(dataOffset64)+i]
	}

}

//0x38 EVMのコードサイズをスタックにpush
func (b *EVM) op_codesize() {
	b.consume_gas(2)
	b.push_asm("CODESIZE")
	fmt.Println("CODESIZE")
	len := len(b.env.code)
	b.push(new(uint256.Int).SetUint64(uint64(len)))
}

//0x39 コントラクトにデプロイされたコードをmemory領域にコピーする
func (b *EVM) op_codecopy() {
	b.consume_gas(9) //??
	b.push_asm("CODECOPY")
	fmt.Println("CODECOPY")
	memOffset := b.pop()
	codeOffset := b.pop()
	length := b.pop()

	uint64CodeOffset64, overflow := codeOffset.Uint64WithOverflow()
	if overflow {
		uint64CodeOffset64 = 0xffffffffffffffff
	}

	memOffset64 := memOffset.Uint64()
	length64 := length.Uint64()

	for i := 0; i < int(length64); i++ {
		b.memory[int(memOffset64)+i] = b.env.code[int(uint64CodeOffset64)+i]
	}
}

// 0x50: スタックから１命令を取り除く
func (b *EVM) op_pop() {
	b.consume_gas(2)
	b.push_asm("POP")
	fmt.Println("POP")
	_ = b.pop()
}

/// 0x51: スタックからpopしたstartを先頭アドレスしてstart+32までの32byteの値をメモリからロード
func (b *EVM) op_mload() {
	b.consume_gas(3)
	b.push_asm("MLOAD")
	fmt.Println("MLOAD")
	start := b.pop()
	bytes := []byte(b.memory[start.Uint64() : start.Uint64()+32])
	b.push(new(uint256.Int).SetBytes(bytes))
}

//0x52 スタックからstart, valueをpopし、startを先頭アドレスしてstart+32までの32byteのメモリ領域にvalueを格納する
func (b *EVM) op_mstore() {
	b.consume_gas(3)
	b.push_asm("MSTORE")
	fmt.Println("MSTORE")
	address := b.pop()
	value := b.pop()
	for i := 0; i < 32; i++ {
		b.memory[address.Uint64()+uint64(i)] = value.Bytes32()[i]
	}
}

//0x53スタックからstart, valueをpopし、startをアドレスとして1byteのメモリ領域にvalueを格納する
func (b *EVM) op_mstore8() {
	b.consume_gas(3)
	b.push_asm("MSTORE8")
	fmt.Println("MSTORE8")
	address := b.pop()
	value := b.pop()
	b.memory[address.Uint64()] = byte(value.Uint64())
}

//0x54 スタックからpopした値をkeyとしてstorageから対応する値をロード
func (b *EVM) op_sload() {
	b.consume_gas(0) //??
	b.push_asm("SLOAD")
	fmt.Println("SLOAD")
	key := b.pop()
	value := b.env.accountstate.storage[key]
	b.push(&value)
}

// 0x55 storageに書き込みを行う storage[operand1(スタック1番目)] = operand2(スタック2番目)
func (b *EVM) op_sstore() {
	b.consume_gas(0) //??
	b.push_asm("SSTORE")
	fmt.Println("SSTORE")
	key := b.pop()
	value := b.pop()
	b.env.accountstate.storage[key] = value
}

// 0x56 スタックからdestinationをpopしてジャンプ
func (b *EVM) op_jump() {
	b.consume_gas(8)
	b.push_asm("JUMP")
	fmt.Println("JUMP")
	destination := b.pop()
	if b.env.code[destination.Uint64()] != 0x5b {
		panic("op_jumpi: destination must be JUMPDEST")
	}
	b.pc = destination.Uint64()
}

// 0x57: スタックからdestination, conditionをpop
// conditionが0以外ならdestinationにジャンプ
func (b *EVM) op_jumpi() {
	b.consume_gas(10)
	b.push_asm("JUMPI")
	fmt.Println("JUMPI")
	destination := b.pop()
	condition := b.pop()
	if b.env.code[destination.Uint64()] != 0x5b {
		panic("op_jumpi: destination must be JUMPDEST")
	}
	if condition.Uint64() != 0 {
		b.pc = destination.Uint64()
	}
}

//0x58
func (b *EVM) op_pc() {
	b.push_asm("PC")
	fmt.Println("PC")
	pc := new(uint256.Int).SetUint64(b.pc)
	b.push(pc)
}

//0x59
func (b *EVM) op_msize() {
	b.push_asm("MSIZE")
	panic("not implement")
}

//0x5a
func (b *EVM) op_gas() {
	b.push_asm("GAS")
	panic("not implement")
}

//0x5b 動的ジャンプを行う際にスタックからpopした値が示すアドレスにジャンプするが、そのアドレスではこのop_jumpdestがオペコードでなければならない
//そのマーカーになるだけでこのOPCODEは単体で意味を持たない
func (b *EVM) op_jumpdest() {
	b.consume_gas(1)
	b.push_asm("JUMPDEST")
	fmt.Println("JUMPDEST")
}

func (b *EVM) op_gasprice() {
	b.push_asm("GASPRICE")
	panic("not implement")
}

func (b *EVM) op_extcodesize() {
	b.push_asm("EXTCODESIZE")
	panic("not implement")
}

func (b *EVM) op_extcodecopy() {
	b.push_asm("EXTCODECOPY")
	panic("not implement")
}

//0x60~0x7f PUSH命令(与えられるバイト分だけpushする)
func (b *EVM) op_push(len int) {
	var operand *uint256.Int
	var operand_str string
	var x []byte
	for i := 0; i < len; i++ {
		x = append(x, b.env.code[b.pc])
		operand_str += fmt.Sprintf("%02x", b.env.code[b.pc])
		b.pc++
	}
	b.consume_gas(3)
	asm := "PUSH" + strconv.Itoa(len) + " " + operand_str
	operand = new(uint256.Int).SetBytes(x)
	b.push_asm(asm)
	fmt.Println(asm)
	b.push(operand)
}

//0x80 DUMP命令
func (b *EVM) op_dump(index uint64) {
	b.consume_gas(3)
	b.push_asm("DUMP" + strconv.Itoa(int(index)))
	fmt.Println("DUMP" + strconv.Itoa(int(index)))
	operand := b.stack[b.sp-1]
	if b.sp > 1 {
		b.stack[b.sp-index-1] = operand
	} else {
		b.push(&operand)
	}
}

//0x90 SWAP命令
func (b *EVM) op_swap(index uint64) {
	b.consume_gas(3)
	b.push_asm("SWAP" + strconv.Itoa(int(index)))
	fmt.Println("SWAP" + strconv.Itoa(int(index)))
	operand1 := b.stack[b.sp-1]
	operand2 := b.stack[b.sp-index-1]
	b.stack[b.sp-1] = operand2
	b.stack[b.sp-index-1] = operand1
}

//0xa0 LOG命令
func (b *EVM) op_log0() {
	b.push_asm("LOG0")
	panic("not implement")
}

//0xa1 LOG1命令
func (b *EVM) op_log1() {
	b.push_asm("LOG1")
	panic("not implement")
}

//0xa2 LOG2命令
func (b *EVM) op_log2() {
	b.push_asm("LOG2")
	panic("not implement")
}

//0xa3 LOG3命令
func (b *EVM) op_log3() {
	b.push_asm("LOG3")
	panic("not implement")
}

//0xa4 LOG4命令
func (b *EVM) op_log4() {
	b.push_asm("LOG4")
	panic("not implement")
}

//0xf0 CREATE命令
func (b *EVM) op_create() {
	b.push_asm("CREATE")
	panic("not implement")
}

//0xf1 CALL命令
func (b *EVM) op_call() {
	b.push_asm("CALL")
	panic("not implement")
}

//CALLCODE命令
func (b *EVM) op_callcode() {
	b.push_asm("CALLCODE")
	panic("not implement")
}

//RETURN命令
func (b *EVM) op_return() {
	b.push_asm("RETURN")
	fmt.Println("RETURN")
	offset := b.pop()
	length := b.pop()
	return_value := b.memory[offset.Uint64() : offset.Uint64()+length.Uint64()]
	b.returns = return_value
}

//REVERT命令
func (b *EVM) op_revert() {
	b.push_asm("REVERT")
	fmt.Println("REVERT")
	offset := b.pop()
	length := b.pop()
	return_value := b.memory[offset.Uint64() : offset.Uint64()+length.Uint64()]
	b.returns = return_value
}

//無効な命令 INVAILD命令
func (b *EVM) op_invalid() {
	b.push_asm("INVALID")
	fmt.Println("INVALID")
}
