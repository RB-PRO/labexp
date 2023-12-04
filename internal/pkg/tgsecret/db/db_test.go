package JsonBase

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	bs, ErrBS := NewBase("test.json")
	if ErrBS != nil {
		t.Error(ErrBS)
	}
	fmt.Println(bs.FileName)
	fmt.Println(bs.Data)
	bs.Data.Keys = []string{"1", "2"}
	bs.Data.Kkeys = make(map[string]Info)
	bs.Data.Kkeys["asd"] = Info{
		Access: true,
		VisitHistory: []Visit{
			{
				UserPC:   "TESTUSER",
				FileName: "TESTFILE",
			},
		},
	}

	errSave := bs.Save()
	if errSave != nil {
		t.Error(errSave)
	}

}
