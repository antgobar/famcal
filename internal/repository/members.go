package repository

import (
	"errors"

	"github.com/antgobar/famcal/internal/models"
)

type Members []models.Member

var MembersStore Members = []models.Member{}

func (m Members) CreateMember(name string) (*models.Member, error) {
	if len(MembersStore) > 10 {
		return nil, errors.New("reached max allowed members")
	}
	member := models.Member{ID: Members(MembersStore).nextId(), Name: name}
	MembersStore = append(MembersStore, member)
	return &member, nil
}

func (u Members) nextId() int {
	maxId := 0
	for _, member := range u {
		if member.ID >= maxId {
			maxId = member.ID
		}
	}
	maxId++
	return maxId
}
