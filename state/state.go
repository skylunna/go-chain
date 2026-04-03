package state

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

/*
	记录每个地址有多少钱。真实区块链用 Merkle Trie，我们用简单的 Map + LevelDB。
*/

type StateManager struct {
	db *leveldb.DB
}

func NewStateManager(dbPath string) (*StateManager, error) {
	db, err := leveldb.OpenFile(dbPath+"/state", nil)
	if err != nil {
		return nil, err
	}

	// 初始化创世地址余额
	sm := &StateManager{db: db}
	if !sm.HasState() {
		sm.SetBalance("genesis", 1000000.0)
	}

	return sm, nil
}

// 获取余额
func (sm *StateManager) GetBalance(addr string) (float64, error) {
	data, err := sm.db.Get([]byte("balance:"+addr), nil)
	if err != nil {
		return 0, nil // 不存在视为0
	}
	var bal float64
	json.Unmarshal(data, &bal)
	return bal, nil
}

// 设置余额
func (sm *StateManager) SetBalance(addr string, amount float64) error {
	data, _ := json.Marshal(amount)
	return sm.db.Put([]byte("balance:"+addr), data, nil)
}

// 交易
func (sm *StateManager) Transfer(from, to string, amount float64) error {
	balFrom, _ := sm.GetBalance(from)
	if balFrom < amount {
		return fmt.Errorf("余额不足")
	}

	sm.SetBalance(from, balFrom-amount)
	balTo, _ := sm.GetBalance(to)
	sm.SetBalance(to, balTo+amount)
	return nil
}

// 检查数据库里有没有初始化过区块链状态
func (sm *StateManager) HasState() bool {
	has, _ := sm.db.Has([]byte("state_init"), nil)
	return has
}
