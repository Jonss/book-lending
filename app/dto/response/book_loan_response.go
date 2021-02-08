package response

import "github.com/Jonss/book-lending/domain/models"

type BookLoanResponse struct {
	BookResponse *BookResponse `json:"book"`
	FromUserID   string        `json:"fromUser"`
	ToUserID     string        `json:"toUser"`
	LentAt       string        `json:"lentAt"`
	ReturnedAt   string        `json:"returnedAt"`
}

func ToBookLoanResponse(bookStatus models.BookStatus,
	fromUser string,
	toUser string,
	status string,
	returnedAt string,
	lentAt string,
) BookLoanResponse {
	return BookLoanResponse{
		BookResponse: &BookResponse{
			Title:      bookStatus.Book.Title,
			ExternalID: bookStatus.Book.Slug,
			Pages:      bookStatus.Book.Pages,
			CreatedAt:  bookStatus.Book.CreatedAt,
			Status:     status,
		},
		FromUserID: fromUser,
		ToUserID:   toUser,
		ReturnedAt: returnedAt,
		LentAt:     lentAt,
	}
}
