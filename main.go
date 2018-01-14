package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

/*
	区块结构
*/
type Block struct {
	Timestamp     int64  //时间戳,当前 区块 的创建时间
	Data          []byte //当前区块 存放到信息(如果是比特货,就是账单信息喽)
	PrevBlockHash []byte //上一个区块的 加密hash值
	Hash          []byte //当前区块的 加密hash值
}

/*
	取了 Block 结构的一些字段（Timestamp, Data 和 PrevBlockHash），
	并将它们相互连接起来，
	然后在连接后的结果上计算一个 SHA-256 的哈希.
	在 SetHash 方法中完成这个任务
*/
func (this *Block) SetHash() {
	//将this.Timestamp中毫秒不分去掉,并且得到[]byte二进制形式
	timestamp := []byte(strconv.FormatInt(this.Timestamp, 10))
	//将this.PrevBlockHash + this.Data 和 timestamp 二进制数据进行拼接
	//中间的拼接以空二进制数据[]byte{}链接
	//例如 this.PrevBlockHash = "ABCDEFG"
	//     this.Data = "55kai"
	//     timestamp = 1712345670
	//最后拼接到headers 应该是 "ABCDEFG55kai1712345670"
	headers := bytes.Join([][]byte{this.PrevBlockHash, this.Data, timestamp}, []byte{})

	//将headers作 SHA256加密
	hash := sha256.Sum256(headers)

	this.Hash = hash[:]
}

/*
	创建一个区块
*/
func NewBlock(data string, prevBlockHash []byte) *Block {
	//time.Now() 会返回一个time类型, "2018-01-14 14:53:21.053635209 +0800 CST m=+0.000259481"
	//time.Now().Unix() 会返回一个int64类型到 日历时间
	block := Block{}
	block.Timestamp = time.Now().Unix()
	block.Data = []byte(data)
	block.PrevBlockHash = prevBlockHash
	block.Hash = []byte{} //暂时当前到hash还没有计算

	block.SetHash()

	return &block
}

/*
	区块链结构
*/
type Blockchain struct {
	Blocks []*Block //有序到区块链
}

//添加一个区块到区块链中
func (this *Blockchain) AddBlock(data string) {
	//得到当前区块链中到最后一个区块
	prevBlock := this.Blocks[len(this.Blocks)-1]

	//创建一个区块 他的前一个区块,就是区块链中的最后一个
	newBlock := NewBlock(data, prevBlock.Hash)

	//将这个区块添加到区块链中
	this.Blocks = append(this.Blocks, newBlock)
}

//创建一个 创世块
func NewGenesisBlock() *Block {
	//创世块 是区块链第一个区块,当然他是没有前驱区块到,这里用[]byte{}表示
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	//创建一个区块链bc
	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to Aceld")
	bc.AddBlock("Send 2 ETH to 刘丹冰")
	bc.AddBlock("Send 3 XRP to ITCAST")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev.Hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Println("=======================")
	}
}
