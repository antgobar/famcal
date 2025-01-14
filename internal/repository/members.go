package repository

var MembersStore Members = []Member{}

type Member struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Members []Member

func (m Members) CreateMember(name string) Member {
	member := Member{Members(MembersStore).nextId(), name}
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
