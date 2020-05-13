package middleware

import (
	"log"
	"net/http"
	"strconv"
	"surface_attack/consts"
	"surface_attack/providers/crudProvider"
	"time"
)

type TimerWrapper struct {
	CrudProvider crudProvider.Interface
}

func (t *TimerWrapper) Timer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)

		avgReqTime, e := t.CrudProvider.GetInt(consts.AVG_REQ_TIME)
		if e != nil {
			log.Println(e)
			return
		}

		oldReqCount, e := t.CrudProvider.GetInt(consts.REQ_COUNT_KEY_NAME)
		if e != nil {
			log.Println(e)
			return
		}

		newReqCount := oldReqCount + 1
		addOldAvgToVal := avgReqTime + int(duration)
		newAvg := addOldAvgToVal / newReqCount

		newAvgStr := strconv.Itoa(newAvg)
		e = t.CrudProvider.Set(consts.AVG_REQ_TIME, newAvgStr)
		if e != nil {
			log.Println(e)
			return
		}

		e = t.CrudProvider.Increment(consts.REQ_COUNT_KEY_NAME)
		if e != nil {
			log.Println(e)
			return
		}
	})
}
