package webHandlers

import (
	"io"
	"log"
	"net/http"
	"surface_attack/consts"
	"surface_attack/providers/crudProvider"
)

type AttackHandler struct {
	CrudProvider crudProvider.Interface
}

func (s *AttackHandler) Handle(w http.ResponseWriter, r *http.Request) {

	attackDestination, ok := r.URL.Query()[consts.ATTACK_URL_PARAM]

	if !ok || len(attackDestination[0]) < 1 {
		log.Println("Url Param ", consts.ATTACK_URL_PARAM, " is missing")
		io.WriteString(w, "Url Param "+consts.ATTACK_URL_PARAM+" is missing")
		return
	}

	attackSources, e := s.CrudProvider.Get(attackDestination[0])
	if e != nil {
		log.Println(e)
		io.WriteString(w, e.Error())
		return
	}

	attackSources = "[ " + attackSources + " ]"

	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(attackSources))
}
