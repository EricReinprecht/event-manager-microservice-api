package enum

type PartyRole string

const (
	RoleOrganizer PartyRole = "ORGANIZER"
	RoleAdmin     PartyRole = "ADMIN"
	RoleStaff     PartyRole = "STAFF"
	RoleAttendee  PartyRole = "ATTENDEE"
)
