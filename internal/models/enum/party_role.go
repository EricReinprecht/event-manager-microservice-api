package enum

type PartyMemberRole string

const (
	PartyRoleOrganizer PartyMemberRole = "ORGANIZER"
	PartyRoleAdmin     PartyMemberRole = "ADMIN"
	PartyRoleRefunder  PartyMemberRole = "REFUNDER"
	PartyRoleScanner   PartyMemberRole = "SCANNER"
)
