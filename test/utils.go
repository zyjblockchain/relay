/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package test

import (
	"github.com/Loopring/relay/config"
	"github.com/Loopring/relay/crypto"
	"github.com/Loopring/relay/dao"
	"github.com/Loopring/relay/ethaccessor"
	"github.com/Loopring/relay/extractor"
	"github.com/Loopring/relay/log"
	"github.com/Loopring/relay/ordermanager"
	"github.com/Loopring/relay/types"
	"github.com/Loopring/relay/usermanager"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"strings"
	"time"
)

type TestParams struct {
	Accessor             *ethaccessor.EthNodeAccessor
	ImplAddress          common.Address
	MinerPrivateKey      []byte
	DelegateAddress      common.Address
	Owner                common.Address
	TokenRegistryAddress common.Address
	Accounts             map[string]string
	TokenAddrs           []string
	Config               *config.GlobalConfig
}

const (
	TokenAddressA = "0x937ff659c8a9d85aac39dfa84c4b49bb7c9b226e"
	TokenAddressB = "0x8711ac984e6ce2169a2a6bd83ec15332c366ee4f"
)

//const (
//	TokenAddressA = "0x359bbea6ade5155bce1e95918879903d3e93365f"
//	TokenAddressB = "0xc85819398e4043f3d951367d6d97bb3257b862e0"
//)

var (
	testAccounts = map[string]string{
		"0x48ff2269e58a373120FFdBBdEE3FBceA854AC30A": "07ae9ee56203d29171ce3de536d7742e0af4df5b7f62d298a0445d11e466bf9e",
		"0xb5fab0b11776aad5ce60588c16bd59dcfd61a1c2": "11293da8fdfe3898eae7637e429e7e93d17d0d8293a4d1b58819ac0ca102b446",
	}

	testTokens = []string{TokenAddressA, TokenAddressB}

	Ks *keystore.KeyStore
)

func Initialize() {
	c := loadConfig()
	Ks := keystore.NewKeyStore(c.Keystore.Keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	cyp := crypto.NewCrypto(true, Ks)
	crypto.Initialize(cyp)
}

func CreateOrder(tokenS, tokenB, protocol, owner common.Address, amountS, amountB *big.Int) *types.Order {
	order := &types.Order{}
	order.Protocol = protocol
	order.TokenS = tokenS
	order.TokenB = tokenB
	order.AmountS = amountS
	order.AmountB = amountB
	order.Timestamp = big.NewInt(time.Now().Unix())
	order.Ttl = big.NewInt(10000)
	order.Salt = big.NewInt(1000)
	order.LrcFee = big.NewInt(1000)
	order.BuyNoMoreThanAmountB = false
	order.MarginSplitPercentage = 0
	order.Owner = owner
	println(owner.Hex())
	if err := order.GenerateAndSetSignature(owner); nil != err {
		panic(err.Error())
	}
	return order
}

func LoadConfig() *config.GlobalConfig {
	return loadConfig()
}

func GenerateAccessor(c *config.GlobalConfig) (*ethaccessor.EthNodeAccessor, error) {
	ks := keystore.NewKeyStore(c.Keystore.Keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	accessor, err := ethaccessor.NewAccessor(c.Accessor, c.Common, ks)
	if nil != err {
		return nil, err
	}
	return accessor, nil
}

//
//func (testParams *TestParams) PrepareTestData() {
//	var err error
//	var hash string
//	accounts := []string{}
//	for k, _ := range testParams.Accounts {
//		accounts = append(accounts, k)
//	}
//
//	//delegate registry
//	delegateContract := &chainclient.TransferDelegate{}
//	testParams.Client.NewContract(delegateContract, testParams.DelegateAddress.Hex(), chainclient.TransferDelegateAbiStr)
//
//	hash, err = delegateContract.AddVersion.SendTransaction(testParams.Owner, common.HexToAddress(testParams.ImplAddress.Hex()))
//	if nil != err {
//		log.Errorf("delegate add version error:%s", err.Error())
//	} else {
//		log.Infof("delegate add version hash:%s", hash)
//	}
//	//
//	//tokenregistry
//	tokenRegistry := &chainclient.TokenRegistry{}
//	testParams.Client.NewContract(tokenRegistry, testParams.TokenRegistryAddress.Hex(), chainclient.TokenRegistryAbiStr)
//	for idx, tokenAddr := range testParams.TokenAddrs {
//		hash, err = tokenRegistry.RegisterToken.SendTransaction(testParams.Owner, common.HexToAddress(tokenAddr), "token"+strconv.Itoa(idx))
//		if nil != err {
//			log.Errorf("register token error:%s", err.Error())
//		} else {
//			log.Infof("register token hash:%s", hash)
//		}
//	}
//	testParams.approveToLoopring(accounts, testParams.TokenAddrs, big.NewInt(30000000))
//}
//
//func (testParams *TestParams) IsTestDataReady() {
//
//	accounts := []string{}
//	for k, _ := range testParams.Accounts {
//		accounts = append(accounts, k)
//	}
//
//	testParams.allowanceToLoopring(accounts, testParams.TokenAddrs)
//}
//
//func (testParams *TestParams) allowanceToLoopring(accounts []string, tokenAddrs []string) {
//	token := &chainclient.Erc20Token{}
//	for _, tokenAddr := range tokenAddrs {
//		testParams.Client.NewContract(token, tokenAddr, chainclient.Erc20TokenAbiStr)
//		for _, account := range accounts {
//			balance := &types.Big{}
//			if err := token.BalanceOf.Call(balance, "latest", common.HexToAddress(account)); nil != err {
//				log.Error(err.Error())
//			} else {
//				log.Infof("token: %s, balance %s : %s", tokenAddr, account, balance.BigInt().String())
//			}
//			if err := token.Allowance.Call(balance, "latest", common.HexToAddress(account), testParams.DelegateAddress); nil != err {
//				log.Error(err.Error())
//			} else {
//				log.Infof("token:%s, allowance: %s -> %s %s", tokenAddr, account, testParams.DelegateAddress.Hex(), balance.BigInt().String())
//			}
//		}
//	}
//}
//
//func (testParams *TestParams) approveToLoopring(accounts []string, tokenAddrs []string, amount *big.Int) {
//	token := &chainclient.Erc20Token{}
//	for _, tokenAddr := range tokenAddrs {
//		testParams.Client.NewContract(token, tokenAddr, chainclient.Erc20TokenAbiStr)
//		for _, account := range accounts {
//			if txHash, err := token.Approve.SendTransaction(types.HexToAddress(account), testParams.DelegateAddress, amount); nil != err {
//				log.Error(err.Error())
//			} else {
//				log.Info(txHash)
//			}
//		}
//
//	}
//}
//
//func (testParams *TestParams) CheckAllowance(tokenAddress, account string) {
//	var result types.Big
//	token := &chainclient.Erc20Token{}
//	testParams.Client.NewContract(token, tokenAddress, chainclient.Erc20TokenAbiStr)
//	if err := token.Allowance.Call(&result, "pending", types.HexToAddress(account), testParams.DelegateAddress); err != nil {
//		panic(err)
//	} else {
//		println(result.BigInt().String())
//	}
//}

func LoadConfigAndGenerateExtractor() *extractor.ExtractorServiceImpl {
	c := loadConfig()
	rds := LoadConfigAndGenerateDaoService()
	accessor, err := GenerateAccessor(c)
	if err != nil {
		panic(err)
	}
	l := extractor.NewExtractorService(c.Accessor, c.Common, accessor, rds)
	return l
}

func LoadConfigAndGenerateOrderManager() *ordermanager.OrderManagerImpl {
	c := loadConfig()
	rds := LoadConfigAndGenerateDaoService()
	um := usermanager.NewUserManager(rds)
	accessor, err := GenerateAccessor(c)
	if err != nil {
		panic(err)
	}
	ob := ordermanager.NewOrderManager(c.OrderManager, &c.Common, rds, um, accessor)
	return ob
}

func LoadConfigAndGenerateDaoService() *dao.RdsServiceImpl {
	c := loadConfig()
	return dao.NewRdsService(c.Mysql)
}

func loadConfig() *config.GlobalConfig {
	path := strings.TrimSuffix(os.Getenv("GOPATH"), "/") + "/src/github.com/Loopring/relay/config/relay.toml"
	c := config.LoadConfig(path)
	log.Initialize(c.Log)

	return c
}
