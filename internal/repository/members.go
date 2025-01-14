package repository

import "github.com/antgobar/famcal/internal/models"

type Members []models.Member

var MembersStore Members = []models.Member{}

func (m Members) CreateMember(name string) models.Member {
	member := models.Member{ID: Members(MembersStore).nextId(), Name: name}
	MembersStore = append(MembersStore, member)
	return member
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
