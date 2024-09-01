package block

import (
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/api"
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

func ParseBlock(data string) (Block) {
	block := Block{}

	block.number = util.HexstrToUint64(gjson.Get(data, "result.number").String())
	block.hash = util.HexstrTo32Bytes(gjson.Get(data, "result.hash").String())
	block.parentHash = util.HexstrTo32Bytes(gjson.Get(data, "result.parentHash").String())
	block.transactionsRoot = util.HexstrTo32Bytes(gjson.Get(data, "result.transactionsRoot").String())
	block.stateRoot = util.HexstrTo32Bytes(gjson.Get(data, "result.stateRoot").String())
	block.receiptsRoot = util.HexstrTo32Bytes(gjson.Get(data, "result.receiptsRoot").String())
	block.timestamp = util.HexstrToUint64(gjson.Get(data, "result.timestamp").String())

	return block
}

func GetBlockByHash(hash [32]byte, url string) (Block) {
	data := api.GetBlockByHash(hash, url)
	return ParseBlock(data)
}

func GetBlockByNumber(number uint64, url string) (Block) {
	data := api.GetBlockByNumber(number, url)
	return ParseBlock(data)
}