package fees

import (
	"context"
	"errors"
	"fmt"

	"github.com/n101661/maney/pkg/utils"
	"github.com/n101661/maney/pkg/utils/slugid"
	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/repository"
	"github.com/samber/lo"
)

type service struct {
	repository repository.FeeRepository

	opts *FeeServiceOptions
}

func NewService(
	repository repository.FeeRepository,
	opts ...utils.Option[FeeServiceOptions],
) (Service, error) {
	return &service{
		repository: repository,
		opts:       utils.ApplyOptions(defaultFeeServiceOptions(), opts),
	}, nil
}

func (s *service) Create(ctx context.Context, r *CreateRequest) (*CreateReply, error) {
	if r.UserID == "" {
		return nil, fmt.Errorf("%w: missing user id", ErrDataInsufficient)
	}
	if r.Fee == nil {
		return nil, fmt.Errorf("%w: missing fee", ErrDataInsufficient)
	}
	switch models.FeeType(r.Fee.Type) {
	case models.FeeTypeRate:
		if r.Fee.Rate == nil {
			return nil, fmt.Errorf("%w: missing fee.rate", ErrDataInsufficient)
		}
	case models.FeeTypeFixed:
		if r.Fee.Fixed == nil {
			return nil, fmt.Errorf("%w: missing fee.fixed", ErrDataInsufficient)
		}
	}

	rows, err := s.repository.Create(ctx, &repository.CreateFeesRequest{
		UserID: r.UserID,
		Fees: []*repository.BaseCreateFee{
			parseBaseCreateFee(r.Fee, s.opts.genPublicID),
		},
	})
	if err != nil {
		return nil, err
	}

	return &CreateReply{
		Fee: parseFee(rows[0]),
	}, nil
}

func (s *service) List(ctx context.Context, r *ListRequest) (*ListReply, error) {
	if r.UserID == "" {
		return nil, fmt.Errorf("%w: missing user id", ErrDataInsufficient)
	}

	reply, err := s.repository.List(ctx, &repository.ListFeesRequest{
		UserID: r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return &ListReply{
				Fees: []*Fee{},
			}, nil
		}
		return nil, err
	}

	return &ListReply{
		Fees: lo.Map(reply.Fees, func(item *repository.Fee, _ int) *Fee {
			return parseFee(item)
		}),
	}, nil
}

func (s *service) Update(ctx context.Context, r *UpdateRequest) (*UpdateReply, error) {
	if r.UserID == "" {
		return nil, fmt.Errorf("%w: missing user id", ErrDataInsufficient)
	}
	if r.FeePublicID == "" {
		return nil, fmt.Errorf("%w: missing public id", ErrDataInsufficient)
	}
	if r.Fee == nil {
		return nil, fmt.Errorf("%w: missing fee", ErrDataInsufficient)
	}
	switch models.FeeType(r.Fee.Type) {
	case models.FeeTypeRate:
		if r.Fee.Rate == nil {
			return nil, fmt.Errorf("%w: missing fee.rate", ErrDataInsufficient)
		}
	case models.FeeTypeFixed:
		if r.Fee.Fixed == nil {
			return nil, fmt.Errorf("%w: missing fee.fixed", ErrDataInsufficient)
		}
	}

	row, err := s.repository.Update(ctx, &repository.UpdateFeeRequest{
		UserID:      r.UserID,
		FeePublicID: r.FeePublicID,
		Fee:         parseBaseFee(r.Fee),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrFeeNotFound
		}
		return nil, err
	}

	return &UpdateReply{
		Fee: parseFee(row),
	}, nil
}

func (s *service) Delete(ctx context.Context, r *DeleteRequest) (*DeleteReply, error) {
	if r.UserID == "" {
		return nil, fmt.Errorf("%w: missing user id", ErrDataInsufficient)
	}
	if r.FeePublicID == "" {
		return nil, fmt.Errorf("%w: missing public id", ErrDataInsufficient)
	}

	_, err := s.repository.Delete(ctx, &repository.DeleteFeesRequest{
		FeePublicIDs: []string{r.FeePublicID},
		UserID:       r.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataNotFound) {
			return nil, ErrFeeNotFound
		}
		return nil, err
	}
	return &DeleteReply{}, nil
}

func parseFee(v *repository.Fee) *Fee {
	return &Fee{
		ID:       v.ID,
		PublicID: v.PublicID,
		BaseFee:  lo.ToPtr(BaseFee(*v.BaseFee)),
	}
}

func parseBaseFee(v *BaseFee) *repository.BaseFee {
	if v == nil {
		return nil
	}
	return lo.ToPtr(repository.BaseFee(*v))
}

func parseBaseCreateFee(v *BaseFee, genPublicID func() string) *repository.BaseCreateFee {
	if v == nil {
		return nil
	}
	return &repository.BaseCreateFee{
		PublicID: genPublicID(),
		BaseFee:  parseBaseFee(v),
	}
}

type FeeServiceOptions struct {
	genPublicID func() string
}

func defaultFeeServiceOptions() *FeeServiceOptions {
	return &FeeServiceOptions{
		genPublicID: func() string {
			return slugid.New("fee", 11)
		},
	}
}

func WithFeeServiceGenPublicID(f func() string) utils.Option[FeeServiceOptions] {
	return func(o *FeeServiceOptions) {
		o.genPublicID = f
	}
}
