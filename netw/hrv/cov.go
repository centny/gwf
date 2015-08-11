package hrv

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
	"github.com/golang/protobuf/proto"
)

func HRV_V2B(v interface{}) ([]byte, error) {
	switch v.(type) {
	case *Res:
		bys, err := proto.Marshal(v.(*Res))
		if err == nil {
			return bys, nil
		} else {
			log.D("HRV_V2B(proto) by v(%v) err:%v", v, err.Error())
			return bys, err
		}
	default:
		bys, err := impl.Json_V2B(v)
		if err == nil {
			return bys, nil
		} else {
			log.D("HRV_V2B(json) by v(%v) err:%v", v, err.Error())
			return bys, err
		}
	}
}

func HRV_B2V(bys []byte, v interface{}) (interface{}, error) {
	switch v.(type) {
	case *Res:
		err := proto.Unmarshal(bys, v.(*Res))
		if err == nil {
			return v, nil
		} else {
			log.D("HRV_B2V(proto) by []byte(%v) err:%v", bys, err.Error())
			return v, err
		}
	default:
		_, err := impl.Json_B2V(bys, v)
		if err == nil {
			return v, nil
		} else {
			log.D("HRV_B2V(json) by []byte(%v) err:%v", bys, err.Error())
			return v, err
		}
	}
}
