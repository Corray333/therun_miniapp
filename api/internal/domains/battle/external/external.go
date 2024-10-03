package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
)

type BattleExternal struct {
}

func New() *BattleExternal {
	return &BattleExternal{}
}

func (e *BattleExternal) GetNewBattles() ([]types.Battle, error) {

	req, err := http.NewRequest("GET", os.Getenv("MAIN_THERUN_SERVER")+"/explorer/battle/special/new", nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("auth", os.Getenv("MAIN_SERVER_AUTH"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	response := struct {
		Battles []types.Battle `json:"battles"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	e.AckBattles(response.Battles)

	return response.Battles, nil
}

func (e *BattleExternal) AckBattles(battles []types.Battle) error {
	battleIDs := []int{}
	for _, battle := range battles {
		battleIDs = append(battleIDs, battle.ID)
	}

	data, err := json.Marshal(struct {
		ID []int `json:"id"`
	}{
		ID: battleIDs,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("MAIN_THERUN_SERVER")+"/explorer/battle/special/new/ack", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("auth", os.Getenv("MAIN_SERVER_AUTH"))

	// Send the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (e *BattleExternal) GetBattlesByID(ids []int) ([]types.Battle, error) {
	url := os.Getenv("MAIN_THERUN_SERVER") + "/explorer/battle/special?"
	for _, id := range ids {
		url += fmt.Sprintf("id[]=%d&", id)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("auth", os.Getenv("MAIN_SERVER_AUTH"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	response := struct {
		Battles []types.Battle `json:"battles"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Battles, nil
}
