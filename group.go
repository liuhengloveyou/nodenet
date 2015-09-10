package nodenet

import (
	"fmt"
	"hash/crc32"
)

const (
	DISPENSE_POLLING    = "polling"    // 轮询
	DISPENSE_HASH       = "hash"       // CRC32哈希
	DISPENSE_CONSISTENT = "consistent" // 一致性哈希 //@@

)

type ComponentGroup struct {
	Name     string       // 组名
	dispense string       // 均衡分发策略
	members  []*Component // 包含的节点组件
	ch       chan *Component
}

func NewGroup(groupName, dispense string, members []*Component) (group *ComponentGroup) {
	if groupName == "" || dispense == "" {
		return nil
	}
	if len(members) < 1 {
		return nil
	}

	group = &ComponentGroup{
		Name:     groupName,
		dispense: dispense,
		members:  members,
		ch:       make(chan *Component)}

	if group.ready() != nil {
		group = nil
		return nil
	}

	return
}

func (p *ComponentGroup) ready() error {
	switch p.dispense {
	case DISPENSE_POLLING:
		go func() {
			for i := 0; i <= len(p.members); i++ {
				if i >= len(p.members) {
					i = 0
				}

				p.ch <- p.members[i]
			}
		}()
	case DISPENSE_HASH:
		// 没啥好准备的
	default:
		return fmt.Errorf("Dispense of ComponentGroup %s is unknown: [%s]", p.Name, p.dispense)
	}

	return nil
}

func (p *ComponentGroup) GetNode(key string) (com *Component) {
	switch p.dispense {
	case DISPENSE_POLLING:
		com = <-p.ch
	case DISPENSE_HASH:
		if key != "" {
			com = p.members[int(crc32.ChecksumIEEE([]byte(key)))%len(p.members)]
		}
	}

	return
}
