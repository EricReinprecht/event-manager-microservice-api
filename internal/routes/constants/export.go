package constants

var Routes = map[string]string{

	// * ---------------------------------------------------------------------
	// * BASE
	// * ---------------------------------------------------------------------

	"Health": Health,

	// * ---------------------------------------------------------------------
	// * PREFIXES
	// * ---------------------------------------------------------------------

	"API": API,

	// * ---------------------------------------------------------------------
	// * AUTH
	// * ---------------------------------------------------------------------

	"AuthRegister": AuthRegister,

	"AuthVerifyEmail": AuthVerifyEmail,

	"AuthLogin": AuthLogin,

	"AuthRefresh": AuthRefresh,

	"AuthLogout": AuthLogout,

	"AuthForgotPassword": AuthForgotPassword,

	"AuthResetPassword": AuthResetPassword,

	"AuthResendVerification": AuthResendVerification,

	// * ---------------------------------------------------------------------
	// * USERS
	// * ---------------------------------------------------------------------

	"UserSessions": UserSessions,

	"UserSessionByFamilyID": UserSessionByFamilyID,

	"UserMe": UserMe,

	"UserCompleteProfile": UserCompleteProfile,

	"UserPassword": UserPassword,

	"UserParties": UserParties,

	// * ---------------------------------------------------------------------
	// * CATEGORIES
	// * ---------------------------------------------------------------------

	"CategoryUpdate": CategoryUpdate,

	"CategoryDelete": CategoryDelete,

	"CategoryCreate": CategoryCreate,

	"CategoryByID": CategoryByID,

	"CategoryList": CategoryList,

	// * ---------------------------------------------------------------------
	// * PARTIES
	// * ---------------------------------------------------------------------

	"PartyCreate": PartyCreate,

	"PartyUpdate": PartyUpdate,

	"PartyDelete": PartyDelete,

	"PartyByID": PartyByID,

	"PartyList": PartyList,

	"PublicParties": PublicParties,

	"PublicPartyByID": PublicPartyByID,

	"PartyTicketCategories": PartyTicketCategories,

	"PublicPartyTicketCategories": PublicPartyTicketCategories,

	// * ---------------------------------------------------------------------
	// * PARTY MEMBERS
	// * ---------------------------------------------------------------------

	"PartyMembers": PartyMembers,

	"PartyMemberByID": PartyMemberByID,

	"PartyMemberRoles": PartyMemberRoles,

	// * ---------------------------------------------------------------------
	// * MEDIA
	// * ---------------------------------------------------------------------

	"MediaUpload": MediaUpload,

	// * ---------------------------------------------------------------------
	// * TICKET CATEGORIES
	// * ---------------------------------------------------------------------

	"TicketCategoryByID": TicketCategoryByID,

	"TicketCategoryCreate": TicketCategoryCreate,

	"TicketCategoryUpdate": TicketCategoryUpdate,

	"TicketCategoryDelete": TicketCategoryDelete,

	// * ---------------------------------------------------------------------
	// * PURCHASES
	// * ---------------------------------------------------------------------

	"PurchaseCreate": PurchaseCreate,

	"PurchaseByID": PurchaseByID,

	// * ---------------------------------------------------------------------
	// * TICKETS
	// * ---------------------------------------------------------------------

	"MyTickets": MyTickets,

	"TicketScan": TicketScan,

	"TicketVerifyScan": TicketVerifyScan,

	// * ---------------------------------------------------------------------
	// * PAYMENTS
	// * ---------------------------------------------------------------------

	"PaymentCheckout": PaymentCheckout,

	"PaymentRefund": PaymentRefund,

	"PayPalWebhook": PayPalWebhook,

	"Checkout": Checkout,

	"Refund": Refund,
}
