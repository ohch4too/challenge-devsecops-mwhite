## PR#2: REST API Compliance and HTTP Standards

This pull request transforms the API from a non-compliant interface into a proper RESTful service following HTTP standards and best practices. The changes address fundamental issues with endpoint design, status code usage, and error handling that prevented the API from functioning as a standard REST service. By implementing correct HTTP verbs, RESTful resource naming, and appropriate status codes for different operations, the API becomes predictable and usable by standard HTTP clients. Improved error messages and logging provide better observability for debugging and monitoring. These corrections ensure the API behaves correctly according to REST principles and HTTP specifications, making integration with frontend applications and external services straightforward.

### RESTful Endpoint Design

Restructured endpoints to follow REST conventions with proper resource naming and HTTP verbs. Changed from non-standard paths like `GET /user/delete/:id` to `DELETE /users/:id`. This aligns with REST principles where resources are nouns and HTTP verbs indicate actions. Consistent `/users` resource naming improves API discoverability and usability.

```go
v1.GET("/users", ListUsers)
v1.POST("/users", AddUser)
v1.GET("/users/:id", GetUser)
v1.DELETE("/users/:id", DelUser)
```

### HTTP Status Code Corrections

Implemented appropriate status codes for different operations: 200 for successful GET, 201 for resource creation, 204 for successful deletion, 404 for not found, and 500 for server errors. Previously all responses returned 200 or incorrect codes. Proper status codes enable correct client-side error handling and follow HTTP specifications.

```go
// User creation returns 201 Created
RespondJSON(c, 201, user)

// Deletion returns 204 No Content
RespondJSON(c, 204, nil)

// Not found returns 404
RespondJSON(c, 404, gin.H{"error": errMsg})
```

### Error Response Standardization

Standardized error responses to include descriptive error messages in JSON format instead of returning resource IDs. Changed from returning the ID on error to returning structured error objects. This provides clients with actionable error information and maintains consistent response structure. Improved error messages aid debugging and user experience.

```go
errMsg := "User not found"
log.Errorf("%s: %v", errMsg, id)
RespondJSON(c, 404, gin.H{"error": errMsg})
```

### Enhanced Error Logging

Improved error messages and logging format for better observability. Replaced vague messages like "Impossible to fetch all users" with clear descriptive messages. Added context to log entries by including relevant data. Consistent error message formatting makes log analysis and debugging more efficient for operational monitoring.

```go
errMsg := "Failed to fetch users"
log.Errorf(errMsg)

errMsg := "Failed to add user"
log.Errorf("%s: %v\n", errMsg, user)
```

### Comprehensive Integration Tests

Added integration tests for all CRUD operations validating correct HTTP verbs and status codes. Tests cover GET, POST, DELETE operations with proper status code assertions. This ensures API endpoints behave according to REST standards and catch regressions in endpoint behavior during future changes.

```go
func TestAddUser(t *testing.T) {
    user := map[string]string{"firstname": "Test", "lastname": "User"}
    body, _ := json.Marshal(user)
    resp, err := http.Post(apiBaseURL+"/users", "application/json", bytes.NewBuffer(body))
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", resp.StatusCode)
    }
}
```