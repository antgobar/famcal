package repository

import (
	"errors"
	"fmt"

	"github.com/antgobar/famcal/internal/models"
)

type Members []models.Member

var MembersStore Members = []models.Member{}

func (m Members) CreateMember(name string, colour int) (*models.Member, error) {
	if len(MembersStore) > 10 {
		return nil, errors.New("reached max allowed members")
	}
	member := models.Member{ID: Members(MembersStore).nextId(), Name: name, Colour: colour}
	MembersStore = append(MembersStore, member)
	return &member, nil
}

func (m Members) DeleteMember(id int) error {
	for index, member := range m {
		if member.ID == id {
			MembersStore = append(m[:index], m[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("member with id %v, does not exist", id)

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
