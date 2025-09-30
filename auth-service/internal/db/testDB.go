package db

import (
	"auth-service/internal/auth/models"
	"auth-service/internal/crypto"
	"fmt"
)

// testInsertUsers creates dummy users and inserts them into the database.
// It returns the created user objects for use in subsequent tests.
func testInsertUsers(database *Database) (models.User, models.User, models.User, error) {
	fmt.Println("======== TEST INSERT USERS =========")

	// Define dummy users
	params := crypto.GetDefaultParams()
	hashedPassword, err := crypto.HashPassword("pass1", &params)
	if err != nil {
		return models.User{}, models.User{}, models.User{}, fmt.Errorf("cannot hash User1")
	}
	dummyUser1 := models.User{Username: "User1", Password: hashedPassword, Email: "Test1@email.com", Role: models.ADMIN}
	dummyUser2 := models.User{Username: "User2", Password: "pass2", Email: "Test2@email.com", Role: models.USER}
	dummyUser3 := models.User{Username: "User3", Password: "pass3", Email: "Test3@email.com", Role: models.USER}

	usersToInsert := []models.User{dummyUser1, dummyUser2, dummyUser3}

	for i, user := range usersToInsert {
		result, err := database.CreateUser(&user)
		if err != nil {
			return models.User{}, models.User{}, models.User{}, fmt.Errorf("cannot INSERT dummyUser%d: %w", i+1, err)
		}
		fmt.Printf("INSERT dummyUser%d Result: %s\n", i+1, result.String())
	}

	fmt.Println()

	return dummyUser1, dummyUser2, dummyUser3, nil
}

// testSelectUsers performs various select queries to fetch user data.
func testSelectUsers(database *Database, user models.User) error {
	fmt.Println("======== TEST SELECT QUERIES =========")

	// Test GetUserByID
	fmt.Println("--- Testing SELECT USER BY ID ---")
	userShell, err := database.SelectUserByID(1)
	if err != nil {
		return fmt.Errorf("failed to SELECT dummyUser1 by ID: %w", err)
	}
	fmt.Println("SELECT dummyUser1 by ID Result:", userShell)
	fmt.Println()

	// Test GetUserByUsernameAndPass
	fmt.Println("--- Testing SELECT USER BY Username & PASSWORD ---")
	userShell, err = database.SelectUserByUsernameAndPass(user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to SELECT dummyUser2 by username and password: %w", err)
	}
	fmt.Println("SELECT dummyUser2 by Username & Pass Result:", userShell)
	fmt.Println()

	return nil
}

// testUpdateUser tests updating a user's information.
func testUpdateUser(database *Database, userToUpdate models.User) error {
	fmt.Println("======== TEST UPDATE USER =========")
	// First, get the user to ensure we have the correct ID
	userShell, err := database.SelectUserByUsername(userToUpdate.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch user for update: %w", err)
	}
	fmt.Println("SELECT dummyUser2 for update Result:", userShell)

	// Now, update the user
	userShell.Username = "User2_UPDATED"
	result, err := database.UpdateUser(userShell)
	if err != nil {
		return fmt.Errorf("failed to UPDATE dummyUser2: %w", err)
	}
	fmt.Println("UPDATE dummyUser2 Result:", result.String())
	fmt.Println()

	return nil
}

// testDeleteUser tests deleting a user from the database.
func testDeleteUser(database *Database, userToDelete models.User) error {
	fmt.Println("======== TEST DELETE USER =========")
	// First, get the user to ensure we have the correct ID
	userShell, err := database.SelectUserByUsername(userToDelete.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch user for deletion: %w", err)
	}
	fmt.Println("SELECT dummyUser3 for deletion Result:", userShell)

	// Now, delete the user
	result, err := database.DeleteUser(userShell)
	if err != nil {
		return fmt.Errorf("failed to DELETE dummyUser3: %w", err)
	}
	fmt.Println("DELETE dummyUser3 Result:", result.String())
	fmt.Println()

	return nil
}

// testDropTable handles dropping the users table if the flag is set.
func testDropTable(database *Database, dropFlag bool) error {
	fmt.Println("======== TEST DROP users TABLE =========")
	const DropTable = `DROP TABLE users`
	if dropFlag {
		fmt.Println("DROP TABLE flag is ON, deleting users TABLE...")
		result, err := database.ExecuteQuery(DropTable)
		if err != nil {
			return fmt.Errorf("failed to DROP users TABLE: %w", err)
		}
		fmt.Println("DROP user TABLE Result:", result.String())
	} else {
		fmt.Println("DROP TABLE flag is OFF, users TABLE still exists.")
	}
	fmt.Println()

	return nil
}

// TestDB orchestrates the execution of all database tests in sequence.
func TestDB(database *Database, dropFlag bool) error {
	_, dummyUser2, dummyUser3, err := testInsertUsers(database)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = testSelectUsers(database, dummyUser2)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = testUpdateUser(database, dummyUser2)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = testDeleteUser(database, dummyUser3)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = testDropTable(database, dropFlag)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("All database tests completed successfully!")

	return nil
}
