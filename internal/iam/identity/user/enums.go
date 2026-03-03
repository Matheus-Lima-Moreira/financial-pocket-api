package user

type RegisterFrom string

const (
	RegisterFromInvite RegisterFrom = "INVITE"
	RegisterFromForm   RegisterFrom = "FORM"
)

type UserState string

const (
	UserStateActive        UserState = "ACTIVE"
	UserStateInactive      UserState = "INACTIVE"
	UserStateInvitePending UserState = "INVITE_PENDING"
)
