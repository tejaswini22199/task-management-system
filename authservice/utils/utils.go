package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tejaswini22199/task-management-system/authservice/repository"
)

// ValidateUserID extracts user_id from Gin context
func ValidateUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("unauthorized")
	}

	authUserID, ok := userID.(int)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return authUserID, nil
}

func ValidateInputUserIDs(inputUserIds []int, c *gin.Context) ([]int, bool) {

	// Ensure input user IDs are unique using a set
	uniqueUserIDs := make(map[int]struct{})
	for _, id := range inputUserIds {
		uniqueUserIDs[id] = struct{}{}
	}

	// Convert set to slice for database validation
	validUserIDs := make([]int, 0, len(uniqueUserIDs))
	for id := range uniqueUserIDs {
		validUserIDs = append(validUserIDs, id)
	}

	// Validate if user IDs exist
	if len(validUserIDs) > 0 {
		existingUserIDs, err := repository.GetExistingUserIDs(validUserIDs)
		fmt.Println("line 44")
		fmt.Println(existingUserIDs, err)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Send the actual error message
			return nil, false
		}

		// Convert slice to a map for quick lookup
		existingUsers := make(map[int]struct{})
		for _, id := range existingUserIDs {
			existingUsers[id] = struct{}{}
		}

		// Find invalid users
		invalidUserIDs := []int{}
		for _, id := range inputUserIds {
			if _, exists := existingUsers[id]; !exists {
				invalidUserIDs = append(invalidUserIDs, id)
			}
		}

		// If any invalid user IDs exist, return an error
		if len(invalidUserIDs) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "Some user IDs do not exist",
				"invalid_users": invalidUserIDs,
			})
			return nil, false
		}
	}

	return validUserIDs, true
}
