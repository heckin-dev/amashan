package bnet

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"github.com/heckin-dev/amashan/test"
	"net/http"
	"strconv"
)

type BattleNetMock struct{}

func (b *BattleNetMock) CharacterSummary(w http.ResponseWriter, r *http.Request) {
	res := &CharacterSummaryResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterSummary)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterSummary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterStatus(w http.ResponseWriter, r *http.Request) {
	res := &CharacterStatusResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterStatus)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterStatus", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterEquipmentSummary(w http.ResponseWriter, r *http.Request) {
	res := &CharacterEquipmentResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterEquipment)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterEquipment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterMedia(w http.ResponseWriter, r *http.Request) {
	res := &CharacterMediaResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterMedia)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterMedia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterStatistics(w http.ResponseWriter, r *http.Request) {
	res := &CharacterStatisticsResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterStatistics)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterStatistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterDungeonEncounters(w http.ResponseWriter, r *http.Request) {
	res := &CharacterDungeonEncountersResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterDungeonEncounters)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterDungeonEncounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) CharacterRaidEncounters(w http.ResponseWriter, r *http.Request) {
	res := &CharacterRaidEncountersResponse{}
	err := json.NewDecoder(bytes.NewReader(test.CharacterRaidEncounters)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.CharacterRaidEncounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) MythicKeystoneIndex(w http.ResponseWriter, r *http.Request) {
	res := &MythicKeystoneIndexResponse{}
	err := json.NewDecoder(bytes.NewReader(test.MythicKeystoneIndex)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.MythicKeystoneIndex", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) MythicKeystoneSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonStr, ok := vars["seasonID"]
	if !ok {
		http.Error(w, "seasonID not provided in route parameter", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(seasonStr); err != nil {
		http.Error(w, "failed to parse seasonID to integer", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := &MythicKeystoneSeasonResponse{}
	err := json.NewDecoder(bytes.NewReader(test.MythicKeystoneSeason)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.MythicKeystoneSeason", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (b *BattleNetMock) Route(r *mux.Router) {
	publicProfile := r.PathPrefix("/profile/wow").Subrouter()
	publicProfile.Use(middleware.UseRealm().Middleware)
	publicProfile.Use(middleware.UseCharacter().Middleware)

	publicProfile.HandleFunc("/character/{realm}/{character}", b.CharacterSummary)
	publicProfile.HandleFunc("/character/{realm}/{character}/status", b.CharacterStatus)
	publicProfile.HandleFunc("/character/{realm}/{character}/equipment", b.CharacterEquipmentSummary)
	publicProfile.HandleFunc("/character/{realm}/{character}/character-media", b.CharacterMedia)
	publicProfile.HandleFunc("/character/{realm}/{character}/statistics", b.CharacterStatistics)
	publicProfile.HandleFunc("/character/{realm}/{character}/encounters/dungeons", b.CharacterDungeonEncounters)
	publicProfile.HandleFunc("/character/{realm}/{character}/encounters/raids", b.CharacterRaidEncounters)
	publicProfile.HandleFunc("/character/{realm}/{character}/mythic-keystone-profile", b.MythicKeystoneIndex)
	publicProfile.HandleFunc("/character/{realm}/{character}/mythic-keystone-profile/season/{seasonID}", b.MythicKeystoneSeason)
}

func NewBattleNetMock() *BattleNetMock {
	return &BattleNetMock{}
}
