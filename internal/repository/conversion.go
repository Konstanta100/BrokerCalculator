package repository

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertProtoTimestampToPgType(protoTs *timestamppb.Timestamp) pgtype.Timestamp {
	if protoTs == nil || !protoTs.IsValid() {
		return pgtype.Timestamp{Valid: false}
	}

	goTime := protoTs.AsTime()
	return pgtype.Timestamp{
		Time:  goTime,
		Valid: true,
	}
}

func NumericToFloat64(num pgtype.Numeric) (float64, error) {
	if !num.Valid {
		return 0, errors.New("NULL numeric value")
	}

	f, err := num.Float64Value()
	if num.NaN {
		return 0, fmt.Errorf("conversion error: %w", err)
	}

	if !f.Valid {
		return 0, errors.New("invalid float64 value")
	}
	return f.Float64, nil
}
