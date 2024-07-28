package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/weldonkipchirchir/simple_bank/util"
)

// makes test are independent
func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	user, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, newFullName, user.FullName)

	user2, err := testStore.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, oldUser.Username, user2.Username)
	require.Equal(t, oldUser.HashedPassword, user2.HashedPassword)
	require.NotEqual(t, oldUser.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	email := util.RandomEmail()
	user, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, email, user.Email)

	user2, err := testStore.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, oldUser.Username, user2.Username)
	require.Equal(t, oldUser.HashedPassword, user2.HashedPassword)
	require.Equal(t, oldUser.FullName, user2.FullName)
	require.NotEqual(t, oldUser.Email, user2.Email)
}
func TestUpdateUserHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	user, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, hashedPassword, user.HashedPassword)

	user2, err := testStore.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, oldUser.Username, user2.Username)
	require.NotEqual(t, oldUser.HashedPassword, user2.HashedPassword)
	require.Equal(t, oldUser.FullName, user2.FullName)
	require.Equal(t, oldUser.Email, user2.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	email := util.RandomEmail()
	fullName := util.RandomOwner()

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	user, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
		FullName: pgtype.Text{
			String: fullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, hashedPassword, user.HashedPassword)
	require.Equal(t, email, user.Email)
	require.Equal(t, fullName, user.FullName)

	user2, err := testStore.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)

	require.Equal(t, oldUser.Username, user2.Username)
	require.NotEqual(t, oldUser.HashedPassword, user2.HashedPassword)
	require.NotEqual(t, oldUser.FullName, user2.FullName)
	require.NotEqual(t, oldUser.Email, user2.Email)
}
