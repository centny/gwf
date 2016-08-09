package routing

import (
	"encoding/json"
	"encoding/xml"
	"github.com/Centny/gwf/util"
	"io/ioutil"
)

func (h *HTTPSession) UnmarshalJ(v interface{}) error {
	_, err := h.UnmarshalJ_v(v)
	return err
}

func (h *HTTPSession) UnmarshalJ_v(v interface{}) ([]byte, error) {
	bys, err := ioutil.ReadAll(h.R.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bys, v)
	if err != nil {
		return bys, util.Err("Unmarshal json error(%v) for data:%v", err.Error(), string(bys))
	}
	return bys, nil
}

func (h *HTTPSession) UnmarshalX(v interface{}) error {
	_, err := h.UnmarshalX_v(v)
	return err
}
func (h *HTTPSession) UnmarshalX_v(v interface{}) ([]byte, error) {
	bys, err := ioutil.ReadAll(h.R.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(bys, v)
	if err != nil {
		return bys, util.Err("Unmarshal xml error(%v) for data:%v", err.Error(), string(bys))
	}
	return bys, nil
}
