package rpc

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"encoding/hex"
	"io"
	"strconv"
)

type ReqBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Method string `json:"method"`
	Params []interface{} `json:"params"`
	Id string `json:"id"`
}

const URL_DEAULT = "https://mainnet.infura.io/v3/cdeb7402eca247e0a054717f350b4e50"

func GetBlockByHash(hash [32]byte) (data string) {
	reqBody := ReqBody{}
	reqBody.Id = "1"
	reqBody.Jsonrpc = "2.0"
	reqBody.Method = "eth_getBlockByHash"
	reqBody.Params = []interface{}{"0x" + hex.EncodeToString(hash[:]), false}

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String:", reqBody)
	fmt.Println("JSON:", string(reqBodyJson))
	
	res, err := http.Post(URL_DEAULT, "application/json", bytes.NewBuffer(reqBodyJson))

	if err != nil {
		fmt.Println("[!] " + err.Error())
	} else {
		fmt.Println("[*] " + res.Status)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return string(body)
}

func GetBlockByNumber(number uint64) (data string) {
	reqBody := ReqBody{}
	reqBody.Id = "1"
	reqBody.Jsonrpc = "2.0"
	reqBody.Method = "eth_getBlockByNumber"
	reqBody.Params = []interface{}{"0x" + fmt.Sprintf("%x", number), false}

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String:", reqBody)
	fmt.Println("JSON:", string(reqBodyJson))
	
	res, err := http.Post(URL_DEAULT, "application/json", bytes.NewBuffer(reqBodyJson))

	if err != nil {
		fmt.Println("[!] " + err.Error())
	} else {
		fmt.Println("[*] " + res.Status)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return string(body)
}

func GetBeaconBlockHeader(slot uint64) (data string) {
	res, err := http.Get(URL_DEAULT + strconv.FormatUint(slot, 10))

	if err != nil {
		fmt.Println("[!] " + err.Error())
	} else {
		fmt.Println("[*] " + res.Status)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return string(body)
}