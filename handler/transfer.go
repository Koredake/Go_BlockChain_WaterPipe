package handler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
)

var port string
var key string

func InitNodes(gethPort, privateKeys string) (err error) {
	port = gethPort
	key = privateKeys
	connRPC, err := rpc.Dial(gethPort)
	conn, err := ethclient.Dial(gethPort)
	var address []string
	errRPC := connRPC.Call(&address, "eth_accounts")
	fmt.Println(address)
	if errRPC != nil {
		return errRPC
		log.Fatalf("Failed to connect to the ethereum client : %v", errRPC)
	}
	if err != nil {
		return err
		log.Fatalf("Failed to connect to the ethereum client : %v", err)
	}
	defer connRPC.Close()
	defer conn.Close()
	return nil
}
func SendEth(to string) {
	client, err := ethclient.Dial(port)
	if err != nil {
		log.Fatal(err)
	}
	// 出账的私钥。
	privateKey, err := crypto.HexToECDSA(key) // from的私钥
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // 你要转多少wei
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}
