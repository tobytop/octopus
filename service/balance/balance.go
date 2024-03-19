package balance

import (
	"fmt"
	"octopus/config"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type Balance interface {
	Add(addr string, weight int)
	Next() string
	Remove(addr string)
	SetWegiht(num int, addr string)
	GetAllAddress() []string
}

type roundRobinBalance struct {
	curIndex int
	addrList []string
	logger   *zerolog.Logger
}

func NewBalance(balancetype config.BalanceType, logger *zerolog.Logger) Balance {
	var balance Balance
	switch balancetype {
	case config.RoundRobin:
		balance = &roundRobinBalance{
			addrList: make([]string, 0),
			logger:   logger,
		}
	case config.WeightRobin:
		balance = &weightRoundRobinBalance{
			addrList: make(map[string]*node),
			logger:   logger,
		}
	default:
		balance = &roundRobinBalance{
			addrList: make([]string, 0),
			logger:   logger,
		}
	}
	return balance
}

func (b *roundRobinBalance) Add(addr string, weight int) {
	if !slices.Contains(b.addrList, addr) {
		b.addrList = append(b.addrList, addr)
	}
}
func (b *roundRobinBalance) SetWegiht(num int, addr string) {}

func (b *roundRobinBalance) Remove(addr string) {
	index := -1
	for key, value := range b.addrList {
		if value == addr {
			index = key
			break
		}
	}

	if index > -1 {
		b.logger.Info().Msg(fmt.Sprintf("begin delete the host %v", addr))
		if (index + 1) == len(b.addrList) {
			b.addrList = b.addrList[:index]
		} else {
			b.addrList = append(b.addrList[:index], b.addrList[index+1:]...)
		}
		b.logger.Info().Msg(fmt.Sprintf("the addrlist %v", b.addrList))
	}
}

func (b *roundRobinBalance) Next() string {
	len := len(b.addrList)
	if len == 0 {
		return ""
	}
	if b.curIndex >= len {
		b.curIndex = 0
	}
	addr := b.addrList[b.curIndex]
	b.curIndex = (b.curIndex + 1) % len
	return addr
}

func (b *roundRobinBalance) GetAllAddress() []string {
	return b.addrList
}

type weightRoundRobinBalance struct {
	curAddr  string
	addrList map[string]*node
	logger   *zerolog.Logger
}

type node struct {
	weght         int
	currentWeight int
	stepWeight    int
	addr          string
}

func (b *weightRoundRobinBalance) Add(addr string, weight int) {
	if _, ok := b.addrList[addr]; !ok {
		node := &node{
			weght:         weight,
			currentWeight: weight,
			stepWeight:    weight,
			addr:          addr,
		}
		b.addrList[addr] = node
	}
}

func (b *weightRoundRobinBalance) Next() string {
	if len(b.addrList) == 0 {
		return ""
	}
	totalWight := 0
	var maxWeghtNode *node
	for key, value := range b.addrList {
		totalWight += value.stepWeight
		value.currentWeight += value.stepWeight
		if maxWeghtNode == nil || maxWeghtNode.currentWeight < value.currentWeight {
			maxWeghtNode = value
			b.curAddr = key
		}
	}
	maxWeghtNode.currentWeight -= totalWight
	return maxWeghtNode.addr
}

func (b *weightRoundRobinBalance) Remove(addr string) {
	b.logger.Info().Msg(fmt.Sprintf("begin delete the host %v", addr))
	delete(b.addrList, addr)
	b.logger.Info().Msg(fmt.Sprintf("the addrlist %v", b.addrList))
}

func (b *weightRoundRobinBalance) SetWegiht(num int, addr string) {
	if node, ok := b.addrList[addr]; ok {
		if num > 0 && node.weght > node.stepWeight {
			if (node.stepWeight + num) > node.weght {
				node.stepWeight = node.weght
			} else {
				node.stepWeight += num
			}
		}
		if num == -1 {
			node.stepWeight = -1
		}
	}
}

func (b *weightRoundRobinBalance) GetAllAddress() []string {
	addrList := make([]string, 0)
	for addr := range b.addrList {
		addrList = append(addrList, addr)
	}
	return addrList
}
