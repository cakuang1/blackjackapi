package handlers

import (
	"blackjackapi/models"

	"fmt"
	"github.com/gorilla/mux"

	"net/http"
)

// documentation from upstash kafka.

func (h *Handler) KafkaHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	tableId := v["tableID"]

	if tableId == "" {
		http.Error(w, "tableID is a  required query parameters", http.StatusBadRequest)
		return
	}
	_, err := models.GetTable(h.Context, tableId, h.Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Table %s does not exist", tableId), http.StatusBadRequest)
		return
	}
	models.KafkaConsumer(tableId)
}


