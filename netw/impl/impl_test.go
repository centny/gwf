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
	var bys = []byte{50, 50, 55, 18, 5, 85, 45, 49, 45, 52, 26, 5, 85, 45, 49, 45, 49, 26, 5, 85, 45, 49, 45, 50, 26, 5, 85, 45, 49, 45, 51, 32, 0, 42, 5, 85, 45, 49, 45, 51, 50, 6, 80, 117, 115, 104, 45, 62, 58, 5, 85, 45, 49, 45, 52, 64, 187, 253, 241, 244, 186, 42}
	fmt.Println("xxx")
	fmt.Println(string(bys), "--->")
	var mv = util.Map{}
	fmt.Println(Json_B2V(bys, mv))
}
