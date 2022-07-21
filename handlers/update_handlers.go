package handlers

import (
	"encoding/json"
	"io/ioutil"
	"manager/models"
	"manager/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h Handler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	// shape of request payload
	// the JSON tags identify what part of the incoming payload
	// to assign to the field in the `json.Unmarshal` method
	type flagUpdate struct {
		Name   string `json:"name"`
		SdkKey string `json:"sdkKey"`
		Status bool   `json:"status"`
		// Audiences []string `json:"audiences"`
	}

	var flagReq flagUpdate
	// Read to request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.NoRecordResponse(w, r, err)
		return
	}

	err = json.Unmarshal(body, &flagReq)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	// get audience objects to insert join reference
	// (GORM model for flags expects Audience objects, not key strings)
	// var dbAuds []models.Audience
	// h.DB.Where("key in ?", flagReq.Audiences).Find(&dbAuds)

	// h.DB.Association("Audiences")
	var updatedFlag models.Flag

	// update the table `updatedFlag` belongs to with the
	// fields existent in flagReq object (these have to map accurately)
	result := h.DB.Model(&updatedFlag).Updates(&flagReq)

	if result.Error != nil {
		utils.HandleErr(result.Error, "should put a failed to update")
		return
	}

	// Send a 201 created response
	utils.UpdatedResponse(w, r, &updatedFlag)
}

func (h Handler) ToggleFlag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	togglef := struct {
		Status bool `json:"status"`
	}{}

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &togglef)

	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	var flag models.Flag
	update := map[string]interface{}{"status": togglef.Status}
	result := h.DB.Model(&flag).Where("id = ?", id).Updates(update)
	if result.Error != nil {
		utils.NoRecordResponse(w, r, result.Error)
		return
	}

	h.DB.First(&flag, id)
	response := models.FlagNoAudsResponse{Flag: &flag}

	utils.UpdatedResponse(w, r, &response)
}
