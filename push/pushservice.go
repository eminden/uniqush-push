package push

import (
	"encoding/json"
	"fmt"
)

type PushService interface {
	Type() string
	MarshalDeliveryPoint(dp DeliveryPoint) (data []byte, err error)
	UnmarshalDeliveryPoint(data []byte, dp DeliveryPoint) error
	UnmarshalDeliveryPointFromMap(data map[string]string, dp DeliveryPoint) error
}

type UnmarshalFromMapToStructPushService struct {
}

func (self *UnmarshalFromMapToStructPushService) UnmarshalDeliveryPointFromMap(data map[string]string, dp DeliveryPoint) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, dp)
}

type mapToPushPeer interface {
	UnmarshalDeliveryPointFromMap(data map[string]string, dp DeliveryPoint) error
}

type BasicPushService struct {
	This mapToPushPeer
}

func (self *BasicPushService) MarshalDeliveryPoint(dp DeliveryPoint) (data []byte, err error) {
	data, err = json.Marshal(dp)
	return
}

func (self *BasicPushService) oldDataToMap(data []byte) (m map[string]string, err error) {
	var mapslice []map[string]string
	err = json.Unmarshal(data, &mapslice)
	if err != nil {
		err = fmt.Errorf("Unable to use old unmarshal technique. %v: %v", err, string(data))
		return
	}

	// merge these data into one big map
	if len(mapslice) > 1 {
		for _, m := range mapslice[1:] {
			for k, v := range m {
				mapslice[0][k] = v
			}
		}
	} else {
		err = fmt.Errorf("Unable to use old unmarshal technique. Has to be a 2-element slice: %v", string(data))
		return
	}
	m = mapslice[0]
	return
}

func (self *BasicPushService) UnmarshalDeliveryPoint(data []byte, dp DeliveryPoint) error {
	// Backward compatible!
	if self.This != nil && len(data) > 0 && data[0] != '{' {
		m, err := self.oldDataToMap(data)
		if err != nil {
			return err
		}
		return self.This.UnmarshalDeliveryPointFromMap(m, dp)
	}
	return json.Unmarshal(data, dp)
}
