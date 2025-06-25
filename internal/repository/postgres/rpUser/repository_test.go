package rpUser

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-storage/internal/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *RepositoryUser) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	repo := NewRepository(db)
	return db, mock, repo
}

func TestNewRepository(t *testing.T) {
	db, _, _ := setupMockDB(t)
	defer db.Close()

	repo := NewRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCreateUser_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	user := &domain.User{
		ID:         "test-id",
		FirstName:  "John",
		SecondName: "William",
		LastName:   "Doe",
		Username:   "johndoe",
		Email:      "john@example.com",
		Phone:      "+1234567890",
		Password:   "hashedpassword",
		CompanyId:  "company-id",
		RoleId:     "role-id",
		LastLogin:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
	}

	mock.ExpectExec(`insert into users \( id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active \) values \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13, \$14\)`).
		WithArgs(
			user.ID, user.FirstName, user.SecondName, user.LastName, user.Username,
			user.Email, user.Phone, user.Password,
			user.CompanyId, user.RoleId,
			user.LastLogin, user.CreatedAt, user.UpdatedAt, user.IsActive,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	user := &domain.User{
		ID:        "test-id",
		Username:  "johndoe",
		Email:     "john@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(`insert into users`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"
	expectedUser := &domain.User{
		ID:         userID,
		FirstName:  "John",
		SecondName: "William",
		LastName:   "Doe",
		Username:   "johndoe",
		Email:      "john@example.com",
		Phone:      "+1234567890",
		Password:   "hashedpassword",
		CompanyId:  "company-id",
		RoleId:     "role-id",
		IsActive:   true,
	}

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "second_name", "last_name", "username",
		"email", "phone", "password", "company_id", "role_id",
		"last_login", "created_at", "updated_at", "is_active",
	}).AddRow(
		expectedUser.ID, expectedUser.FirstName, expectedUser.SecondName, expectedUser.LastName, expectedUser.Username,
		expectedUser.Email, expectedUser.Phone, expectedUser.Password, expectedUser.CompanyId, expectedUser.RoleId,
		time.Now(), time.Now(), time.Now(), expectedUser.IsActive,
	)

	mock.ExpectQuery(`SELECT .+ FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	result, err := repo.GetUserByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "nonexistent-id"

	mock.ExpectQuery(`SELECT .+ FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUserByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	email := "john@example.com"
	expectedUser := &domain.User{
		ID:       "test-id",
		Username: "johndoe",
		Email:    email,
		IsActive: true,
	}

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "second_name", "last_name", "username",
		"email", "phone", "password", "company_id", "role_id",
		"last_login", "created_at", "updated_at", "is_active",
	}).AddRow(
		expectedUser.ID, "", "", "", expectedUser.Username,
		expectedUser.Email, "", "", "", "",
		time.Now(), time.Now(), time.Now(), expectedUser.IsActive,
	)

	mock.ExpectQuery(`SELECT .+ FROM users WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(rows)

	result, err := repo.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	username := "johndoe"
	expectedUser := &domain.User{
		ID:       "test-id",
		Username: username,
		Email:    "john@example.com",
		IsActive: true,
	}

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "second_name", "last_name", "username",
		"email", "phone", "password", "company_id", "role_id",
		"last_login", "created_at", "updated_at", "is_active",
	}).AddRow(
		expectedUser.ID, "", "", "", expectedUser.Username,
		expectedUser.Email, "", "", "", "",
		time.Now(), time.Now(), time.Now(), expectedUser.IsActive,
	)

	mock.ExpectQuery(`SELECT .+ FROM users WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(rows)

	result, err := repo.GetUserByUsername(context.Background(), username)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersByCompanyID_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	companyID := "company-id"
	users := []*domain.User{
		{ID: "user1", Username: "user1", CompanyId: companyID, IsActive: true},
		{ID: "user2", Username: "user2", CompanyId: companyID, IsActive: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "second_name", "last_name", "username",
		"email", "phone", "password", "company_id", "role_id",
		"last_login", "created_at", "updated_at", "is_active",
	})

	for _, user := range users {
		rows.AddRow(
			user.ID, "", "", "", user.Username,
			"", "", "", user.CompanyId, "",
			time.Now(), time.Now(), time.Now(), user.IsActive,
		)
	}

	mock.ExpectQuery(`SELECT .+ FROM users WHERE company_id = \$1`).
		WithArgs(companyID).
		WillReturnRows(rows)

	result, err := repo.GetUsersByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[1].ID, result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	user := &domain.User{
		ID:         "test-id",
		FirstName:  "John",
		SecondName: "William",
		LastName:   "Doe",
		Username:   "johndoe",
		Email:      "john@example.com",
		Phone:      "+1234567890",
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec(`UPDATE users.*SET.*WHERE.*id.*is_active`).
		WithArgs(
			user.ID, user.FirstName, user.SecondName, user.LastName, user.Username,
			user.Email, user.Phone, sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.UpdateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePassword_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"
	hashedPassword := "newhashed"

	mock.ExpectExec(`UPDATE users.*SET.*"password".*updated_at.*WHERE.*id`).
		WithArgs(userID, hashedPassword, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePassword(context.Background(), userID, hashedPassword)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateIsActive_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"
	isActive := false

	mock.ExpectExec(`UPDATE users SET is_active = \$2, updated_at = \$3 WHERE id = \$1`).
		WithArgs(userID, isActive, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateIsActive(context.Background(), userID, isActive)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLastLogin_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"

	mock.ExpectExec(`UPDATE users SET last_login = \$2, updated_at = \$3 WHERE id = \$1`).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateLastLogin(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserRole_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"
	roleID := "new-role-id"

	mock.ExpectExec(`UPDATE users SET role_id = \$2, updated_at = \$3 WHERE id = \$1`).
		WithArgs(userID, roleID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateUserRole(context.Background(), userID, roleID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	users := []*domain.User{
		{ID: "user1", Username: "user1", IsActive: true},
		{ID: "user2", Username: "user2", IsActive: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "second_name", "last_name", "username",
		"email", "phone", "password", "company_id", "role_id",
		"last_login", "created_at", "updated_at", "is_active",
	})

	for _, user := range users {
		rows.AddRow(
			user.ID, "", "", "", user.Username,
			"", "", "", "", "",
			time.Now(), time.Now(), time.Now(), user.IsActive,
		)
	}

	mock.ExpectQuery(`SELECT .+ FROM users`).
		WillReturnRows(rows)

	result, err := repo.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[1].ID, result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserCompany_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := "test-id"
	companyID := "new-company-id"

	mock.ExpectExec(`UPDATE users SET company_id = \$2, updated_at = \$3 WHERE id = \$1`).
		WithArgs(userID, companyID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateUserCompany(context.Background(), userID, companyID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
