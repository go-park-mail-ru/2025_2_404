package handlers

import (
	"2025_2_404/internal/models"
	"bytes"
	"database/sql"
	"errors"
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
	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForSelectUsers)).WithArgs("test@example.com").WillReturnRows(rowsUser)

	// проверка что сессия существует 
	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForCheckSession)).WithArgs(expectedUserID).WillReturnError(sql.ErrNoRows)

	// проверка для создании сессии
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO session (user_id, session_id) VALUES ($1, $2)`)).WithArgs(expectedUserID, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1)) 

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(creds))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	
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
	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForSelectUsers)).
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

	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForSelectUsers)).
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(creds))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)
	
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRegisterHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := New(db)

	userData := `{"email":"newuser@example.com", "password":"password123", "user_name":"Newbie"}`
	passwordHash := "cbfdac6008f9cab4083784cbd1874f76618d2a97" 
	expectedUserID := 1

	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForInsertUsers)).
		WithArgs("newuser@example.com", passwordHash, "Newbie").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedUserID))

	mock.ExpectExec(regexp.QuoteMeta(sqlTextForInsertSession)).
		WithArgs(expectedUserID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(userData))
	rr := httptest.NewRecorder()
	h.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	cookie := rr.Result().Cookies()[0]
	assert.Equal(t, "session_id", cookie.Name)
	assert.NotEmpty(t, cookie.Value)
	assert.JSONEq(t, `{"message": "User created"}`, rr.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := New(db)

	userData := `{"email":"existing@example.com", "password":"password123", "user_name":"ExistingUser"}`
	passwordHash := "cbfdac6008f9cab4083784cbd1874f76618d2a97"

	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForInsertUsers)).
		WithArgs("existing@example.com", passwordHash, "ExistingUser").
		WillReturnError(errors.New("UNIQUE constraint failed: users.email")) // Симулируем ошибку БД

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(userData))
	rr := httptest.NewRecorder()
	h.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, "User already registered\n", rr.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandle_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	h := New(db)

	sessionID := "valid-session-id"
	userID := "123"
	expectedAd := models.Ads{
		ID: 1,
		CreatorID:  userID, 
		FilePath:  "/files/ad.jpg",
		Title:     "Тестовая реклама",
		Text:      "Лучшее предложение!",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForFoundUser)).WithArgs(sessionID).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	rows := sqlmock.NewRows([]string{"id", "file_path", "title", "text"}).AddRow(expectedAd.ID, expectedAd.FilePath, expectedAd.Title, expectedAd.Text)
	mock.ExpectQuery(regexp.QuoteMeta(sqlTextForSelectAds)).WithArgs(userID).WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr := httptest.NewRecorder()
	h.Handle(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedBody := `{
		"message": "Successful authorization",
		"ads": [
			{
				"add_id": "ad-001",
				"creater_id": "123",
				"file_path": "/files/ad.jpg",
				"title": "Тестовая реклама",
				"text": "Лучшее предложение!"
			}
		]
	}`
	assert.JSONEq(t, expectedBody, rr.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}