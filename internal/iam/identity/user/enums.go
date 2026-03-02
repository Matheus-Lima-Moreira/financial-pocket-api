package user

type RegisterFrom string

const (
	RegisterFromInvite RegisterFrom = "INVITE"
	RegisterFromForm   RegisterFrom = "FORM"
)
