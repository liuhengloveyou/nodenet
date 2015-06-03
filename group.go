package nodenet

import (
	"sync"
)

type Group struct {
	Name     string       // 组名
	dispense string       // 分发策略
	members  []*Component // 包含的节点组件

	ch    chan *Component
	ready sync.Once
}

func (p *Group) getNode() *Component {
	p.ready.Do(func() {
		go func() {
			for i := 0; i <= len(p.members); i++ {
				if i >= len(p.members) {
					i = 0
				}

				p.ch <- p.members[i]
			}
		}()
	})

	return <-p.ch
}

func NewGroup(groupName, dispense string, members []*Component) (group *Group) {
	group = &Group{Name: groupName, dispense: dispense, members: members, ch: make(chan *Component)}
	groups[groupName] = group

	return
}

func GroupGetNext(groupName string) *Component {
	if g, ok := groups[groupName]; ok {
		return g.getNode()
	}

	return nil
}
