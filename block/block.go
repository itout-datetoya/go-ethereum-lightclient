package block

import (
	"encoding/hex"
	"errors"
	"strings"
	"strconv"
	"github.com/tidwall/gjson"
)

type Block struct {
	number uint64 //ブロック番号。 保留中のブロックの場合はnull
	hash [32]byte //32バイト - ブロックのハッシュ。 保留中のブロックの場合はnull
	parentHash [32]byte //32バイト - 親ブロックのハッシュ
	transactionsRoot [32]byte //32バイト - ブロックのトランザクションツリーのルート
	stateRoot [32]byte //32バイト - ブロックの最終状態ツリーのルート
	receiptsRoot [32]byte //32バイト - ブロックのレシートツリーのルート
	timestamp uint64 //ブロックが照合されたときのUNIXタイムスタンプ
}

func Parse(data string) (Block) {
	block := Block{}

	block.number = hexstrToUint64(gjson.Get(data, "result.number").String())
	block.hash = hexstrTo32Bytes(gjson.Get(data, "result.hash").String())
	block.parentHash = hexstrTo32Bytes(gjson.Get(data, "result.parentHash").String())
	block.transactionsRoot = hexstrTo32Bytes(gjson.Get(data, "result.transactionsRoot").String())
	block.stateRoot = hexstrTo32Bytes(gjson.Get(data, "result.stateRoot").String())
	block.receiptsRoot = hexstrTo32Bytes(gjson.Get(data, "result.receiptsRoot").String())
	block.timestamp = hexstrToUint64(gjson.Get(data, "result.timestamp").String())

	return block
}

func hexstrTo32Bytes(hexString string) ([32]byte) {
	hexString = strings.TrimPrefix(hexString, "0x")

	byteArray, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}

	if len(byteArray) != 32 {
		panic(errors.New("hexString is not 32 bytes"))
	} else {
		hash := [32]byte{}
		copy(hash[:], byteArray)
		return hash
	}
}

func hexstrToUint64(hexString string) (uint64) {
	hex, _ := strconv.ParseInt(hexString, 0, 64)
	return uint64(hex)
}