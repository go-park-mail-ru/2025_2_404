package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := New(db)

	creds := `{"email":"test@example.com","password":"password123"}`

	passwordHash := "cbfdac6008f9cab4083784cbd1874f76618d2a97" 
	expectedUserID := "1"

	//проверка юзера по email 
	rowsUser := sqlmock.NewRows([]string{"id", "password"}).
		AddRow(expectedUserID, passwordHash)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email = $1`)).WithArgs("test@example.com").WillReturnRows(rowsUser)

	// проверка что сессия существует 
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT session_id FROM session WHERE user_id = $1`)).WithArgs(expectedUserID).WillReturnError(sql.ErrNoRows)

	// проверка для создании сессии
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO session (user_id, session_id) VALUES ($1, $2)`)).WithArgs(expectedUserID, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1)) 

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(creds))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	
	cookie := rr.Result().Cookies()[0]
	assert.Equal(t, "session_id", cookie.Name)
	assert.NotEmpty(t, cookie.Value) 
	
	assert.JSONEq(t, `{"message": "Successful authorization"}`, rr.Body.String())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestLoginHandler_InvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := New(db)

	creds := `{"email":"test@example.com","password":"wrongpassword"}`
	correctPasswordHash := "cbfdac6008f9cab4083784cbd1874f76618d2a97" // password123
	
	rowsUser := sqlmock.NewRows([]string{"id", "password"}).
		AddRow("1", correctPasswordHash)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email = $1`)).
		WithArgs("test@example.com").
		WillReturnRows(rowsUser)
	
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(creds))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid email or password\n", rr.Body.String())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := New(db)

	creds := `{"email":"notfound@example.com","password":"password123"}`

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password FROM users WHERE email = $1`)).
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(creds))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)
	
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}