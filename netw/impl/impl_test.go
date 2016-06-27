package impl

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"testing"
)

func TestNew(t *testing.T) {
	NewChanExecListener_m_j(nil, "", nil)
	NewChanExecListenerN_m_r(nil, "sdfs", nil, nil, nil, nil)
}

func TestB2V(t *testing.T) {
	var bys = []byte{123, 34, 97, 114, 103, 115, 34, 58, 110, 117, 108, 108, 44, 34, 110, 97, 109, 101, 34, 58, 34, 108, 105, 115, 116, 34, 125}
	fmt.Println("xxx")
	fmt.Println(string(bys), "--->")
	var mv = util.Map{}
	fmt.Println(Json_B2V(bys, mv))
}
