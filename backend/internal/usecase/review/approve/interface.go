package approve

import "context"

type ReviewApprover interface {
	ApproveReview(ctx context.Context, reviewID int) error
}
