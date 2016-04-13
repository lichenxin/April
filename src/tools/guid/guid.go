package guid

import (
	"sync"
	"time"
	. "tools"
)

type Guid struct {
	sequence      int32
	mx            sync.Mutex
	lastTimestamp int64
}

func (this *Guid) NewID(area_id uint16) uint64 {
	this.mx.Lock()
	defer this.mx.Unlock()

	if area_id > 4095 {
		ERR("area_id超出最大值")
		return 0
	}

	timestamp := time.Now().Unix()
	if timestamp < this.lastTimestamp {
		ERR("请调整服务器时间!")
		return 0
	}

	if timestamp == this.lastTimestamp {
		//当前毫秒内, 则+1
		this.sequence += 1
		if this.sequence > 4095 {
			//当前毫秒内计数满了, 则等待下一秒
			this.sequence = 0
			for {
				timestamp = time.Now().Unix()
				if timestamp > this.lastTimestamp {
					break
				}
			}
		}
	} else {
		this.sequence = 0
	}
	this.lastTimestamp = timestamp

	//40(毫秒) + 12(serverID) + 12(重复累加)
	return ((uint64(timestamp) << 40) | (uint64(area_id) << 12) | uint64(this.sequence))
}

func NewGuid() *Guid {
	return &Guid{
		sequence:      0,
		lastTimestamp: -1,
	}
}
