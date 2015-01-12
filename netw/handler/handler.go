//Package handler provider multi base handler for netw connection.
//
//ChanH: the chan handler provide the feature of distributing command.
//
//RCH_*: the remove command handler provide the feature of remote command call.
package handler

func NewRC_S_V_S(h RC_V_H, v2b V2Byte, b2v Byte2V) *RC_S {
	return NewRC_S(NewRC_V_S(h, v2b, b2v))
}

func NewRC_S_Json_S(h RC_V_H) *RC_S {
	return NewRC_S(NewRC_Json_S(h))
}

func NewRC_S_V_M_S(h RC_V_M_H, v2b V2Byte, b2v Byte2V) *RC_S {
	return NewRC_S_V_S(NewRC_V_M_S(h), v2b, b2v)
}

func NewRC_S_Json_M_S(h RC_V_M_H) *RC_S {
	return NewRC_S_Json_S(NewRC_V_M_S(h))
}

func NewChan_Json_S(h RC_V_H) *ChanH {
	return NewChanH(NewRC_S_Json_S(h))
}
