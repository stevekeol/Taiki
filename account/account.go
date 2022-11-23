package account

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AccountID is decided by Workchain and Address together
type AccountId struct {
	Workchain int32
	Address   [32]byte
}

// New returns a new instance reference of AccountId
func New(id int32, addr [32]byte) *AccountId {
	return &AccountId{Workchain: id, Address: addr}
}

// IsZero judges whether address equals to zero
func (id AccountId) IsZero() bool {
	for i := range id.Address {
		if id.Address[i] != 0 {
			return false
		}
	}
	return true
}

// String returns combanation Workchain and Address as string
func (id AccountId) String() string {
	return fmt.Sprintf("%v:%x", id.Workchain, id.Address)
}

// MarshalJSON converts
func (id AccountId) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// TODO
func (id *AccountId) UnmarshalJSON(data []byte) error {
	a, err := AccountIdFromRaw(strings.Trim(string(data), "\"\n "))
	if err != nil {
		return err
	}
	id.Workchain = a.Workchain
	id.Address = a.Address
	return nil
}

// AccountIdFromRaw retrieves the AccountId type from string
// Such string looks like "-1:5461696b6920697320666f722065726572796f6e6520696e20776562332e2e2e"
func AccountIdFromRaw(s string) (*AccountId, error) {
	if len(s) == 0 {
		return nil, nil
	}
	var (
		workchain int32
		address   []byte
		aa        AccountId
	)
	_, err := fmt.Sscanf(s, "%d:%x", &workchain, &address) // fmt.Scanf用于扫描解析赋值字符串，很值得使用
	if err != nil {
		return nil, err
	}
	if len(address) != 32 {
		return nil, fmt.Errorf("address len must be 32 bytes")
	}
	aa.Workchain = workchain
	copy(aa.Address[:], address) // []byte是引用类型，必须拷贝过来，而不是直接赋值；Address是数组，因此需要转换成slice
	return &aa, nil
}
