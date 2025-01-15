package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antgobar/famcal/internal/models"
	"github.com/antgobar/famcal/internal/repository"
)

type CreateMemberRequest struct {
	Name string `json:"name"`
}

type MemberResponse struct {
	CalMember models.Member `json:"member"`
	Message   string        `json:"message"`
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repository.MembersStore)
}

func addMember(w http.ResponseWriter, r *http.Request) {
	req := &CreateMemberRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	member, err := repository.MembersStore.CreateMember(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	response := MemberResponse{
		CalMember: *member,
		Message:   "CalMember created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
