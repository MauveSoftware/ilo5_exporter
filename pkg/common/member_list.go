package common

type MemberList struct {
	Members []Member
}

type Member struct {
	Path string `json:"@odata.id"`
}
