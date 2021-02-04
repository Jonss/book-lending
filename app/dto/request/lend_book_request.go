package request

type LendBookRequest struct {
	Title           string `json:"book_external_id"`
	UserToLendEmail string `json:"user_to_lend_email"`
}
