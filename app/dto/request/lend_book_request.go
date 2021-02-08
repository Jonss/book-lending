package request

type LendBookRequest struct {
	BookID       string `json:"book_external_id"`
	UserToLendID string `json:"user_to_lend_id"`
}
