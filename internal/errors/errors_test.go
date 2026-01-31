package errors

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	t.Run("with code and message", func(t *testing.T) {
		e := &Error{Code: CodeInvalid, Message: "bad input"}
		s := e.Error()
		assert.Contains(t, s, "<invalid>")
		assert.Contains(t, s, "bad input")
	})

	t.Run("with underlying error", func(t *testing.T) {
		e := &Error{Operation: "op", UnderlyingError: sql.ErrNoRows}
		s := e.Error()
		assert.Contains(t, s, "op")
		assert.Contains(t, s, sql.ErrNoRows.Error())
	})

	t.Run("underlying only", func(t *testing.T) {
		e := &Error{UnderlyingError: sql.ErrNoRows}
		assert.Equal(t, sql.ErrNoRows.Error(), e.Error())
	})
}

func TestError_Unwrap(t *testing.T) {
	inner := sql.ErrNoRows
	e := &Error{UnderlyingError: inner}
	assert.Equal(t, inner, e.Unwrap())
}

func TestErrorCode(t *testing.T) {
	assert.Empty(t, ErrorCode(nil))
	assert.Equal(t, CodeInvalid, ErrorCode(&Error{Code: CodeInvalid}))
	assert.Equal(t, CodeInternal, ErrorCode(&Error{Code: ""}))
	assert.Equal(t, CodeNotFound, ErrorCode(&Error{Code: CodeNotFound, UnderlyingError: sql.ErrNoRows}))
}

func TestErrorMessage(t *testing.T) {
	assert.Empty(t, ErrorMessage(nil))
	assert.Equal(t, "bad input", ErrorMessage(&Error{Message: "bad input"}))
	assert.Equal(t, "Invalid", ErrorMessage(&Error{Code: CodeInvalid}))
}

func TestMapSQLError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, MapSQLError(nil))
	})

	t.Run("ErrNoRows becomes NotFound", func(t *testing.T) {
		err := MapSQLError(sql.ErrNoRows)
		require.Error(t, err)
		assert.Equal(t, CodeNotFound, ErrorCode(err))
	})

	t.Run("other error becomes Internal", func(t *testing.T) {
		other := sql.ErrConnDone
		err := MapSQLError(other)
		require.Error(t, err)
		assert.Equal(t, CodeInternal, ErrorCode(err))
	})
}
