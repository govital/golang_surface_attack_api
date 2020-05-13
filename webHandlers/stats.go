package webHandlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"surface_attack/consts"
	"surface_attack/providers/crudProvider"
	"time"
)

type StatsHandler struct {
	CrudProvider crudProvider.Interface
}

func (s *StatsHandler) Handle(w http.ResponseWriter, r *http.Request) {

	avgReqTime, e := s.CrudProvider.GetInt(consts.AVG_REQ_TIME)
	if webErrorExist(w, e) {
		return
	}

	sec := nanoSecToFloatSeconds(avgReqTime)

	oldReqCountStr, e := s.CrudProvider.Get(consts.REQ_COUNT_KEY_NAME)
	if webErrorExist(w, e) {
		return
	}

	vmCountStr, e := s.CrudProvider.Get(consts.MV_COUNT_KEY_NAME)
	if webErrorExist(w, e) {
		return
	}

	m := makeOutputMap(vmCountStr, oldReqCountStr, sec)
	b, e := json.Marshal(m)
	if webErrorExist(w, e) {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(b))
}

func makeOutputMap(vmCountStr, oldReqCountStr string, sec float64) map[string]interface{} {
	return map[string]interface{}{
		consts.MV_COUNT_KEY_NAME:  vmCountStr,
		consts.REQ_COUNT_KEY_NAME: oldReqCountStr,
		consts.AVG_REQ_TIME:       sec,
	}
}

func nanoSecToFloatSeconds(intgr int) float64 {
	d := time.Duration(intgr) * time.Nanosecond
	return float64(d) / float64(time.Second)
}

func webErrorExist(w http.ResponseWriter, e error) bool {
	if e != nil {
		log.Println(e)
		io.WriteString(w, e.Error())
		return true
	}
	return false
}
