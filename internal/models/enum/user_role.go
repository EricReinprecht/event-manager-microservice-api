package enum

type UserRole string

const (
	UserRoleDefault        UserRole = "DEFAULT"
	UserRolePremium        UserRole = "PREMIUM"
	UserRolePlatinum       UserRole = "PLATINUM"
	UserRoleArtist         UserRole = "ARTIST"
	UserRoleVerifiedArtist UserRole = "VERIFIED_ARTIST"
	UserRolePartner        UserRole = "PARTNER"
	UserRoleModerator      UserRole = "MODERATOR"
)
