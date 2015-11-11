package routing

import (
	"encoding/json"
	"github.com/Centny/gwf/util"
	"io/ioutil"
)

func (h *HTTPSession) UnmarshalJ(v interface{}) error {
	bys, err := ioutil.ReadAll(h.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bys, v)
	if err != nil {
		return util.Err("Unmarshal json error(%v) for data:%v", err.Error(), string(bys))
	}
	return nil
}
