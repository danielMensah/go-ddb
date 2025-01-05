![coverage](https://raw.githubusercontent.com/danielMensah/onetable-go/badges/.badges/main/coverage.svg)
# **OneTable: A Go Library to Simplify DynamoDB Operations**

**OneTable** is a lightweight library for **DynamoDB** in Go. It simplifies the use of the AWS DynamoDB SDK by abstracting common complexities and avoiding repetitive code, allowing developers to focus on building applications instead of dealing with DynamoDB's low-level details.

---

## **Why OneTable?**

Working with DynamoDB's SDK often involves writing repetitive boilerplate code and managing low-level operations such as query conditions, attribute value mappings, and table configurations. **OneTable** aims to:
- Abstract away complexity.
- Simplify querying and item creation.
- Provide reusable logic to reduce boilerplate.
- Make DynamoDB development faster and more enjoyable.

---

## **Features**

1. **Item Creation**:
    - Easily insert new items into your DynamoDB table without worrying about manual attribute mappings.

2. **Query Items**:
    - Retrieve a single item or a list of items with flexible query support.
    - Automatically maps DynamoDB results into Go structs.

3. **SQL-Like Query Filtering**:
    - Use familiar SQL-like conditions for filtering results.
    - Supports placeholders for safe and dynamic query construction (e.g., `"userId = $1, age = $2"`).

4. **Flexible API**:
    - Designed to work with both single-item and multi-item queries based on the provided output type (`out`).

---

## **Installation**

Add OneTable to your Go project using:
```bash
go get github.com/danielMensah/onetable-go
```

---

## **Usage**

### **Initialize OneTable Client**

```go
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/danielMensah/onetable-go"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Initialize DynamoDB client
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	// Create OneTable client
	client := onetable.New(dynamoDBClient, "YourTableName")

	// Use the client for CRUD operations
}
```

---

### **Creating an Item**

```go
type User struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

func createUser(ctx context.Context, client onetable.Client) {
	user := User{
		ID:   "123",
		Name: "John Doe",
	}

	_, err := client.CreateItem(ctx, user)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Println("User created successfully!")
}
```

---

### **Finding Items**

#### **Retrieve a Single Item**
```go
func getUser(ctx context.Context, client onetable.Client) {
	var user User
	err := client.Find(ctx, &user, "id = $1", "123")
	if err != nil {
		log.Fatalf("failed to retrieve user: %v", err)
	}

	log.Printf("Retrieved user: %+v\n", user)
}
```

#### **Retrieve Multiple Items**
```go
func listUsers(ctx context.Context, client onetable.Client) {
	var users []User
	err := client.Find(ctx, &users, "name = $1", "John Doe")
	if err != nil {
		log.Fatalf("failed to list users: %v", err)
	}

	log.Printf("Retrieved users: %+v\n", users)
}
```

---

## **How It Works**

- **`CreateItem`**:
    - Automatically marshals the provided struct into DynamoDB attributes and inserts it into the specified table.

- **`Find`**:
    - Dynamically determines whether to retrieve a single item or multiple items based on the type of `out`.
    - Simplifies query construction using SQL-like placeholders (e.g., `id = $1`) and automatically maps results to Go structs.

---

## **Roadmap**

The following features are planned for future releases:
- **UpdateItem**: Support for updating existing items.
- **DeleteItem**: Simplify item deletion.
- **Batch Operations**: Add support for `BatchWrite` and `BatchGet`.
- **Expression Builders**: Simplify conditional expressions for complex queries.
- **Pagination Support**: Enable seamless handling of paginated results.
- **Transactions**: Add support for DynamoDB transactions.

---

## **Contributing**

Contributions are welcome! If you encounter issues or have feature suggestions, feel free to open an issue or submit a pull request.

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**

- Built with ❤️ by [Daniel Mensah](https://github.com/danielMensah).
- Powered by [AWS DynamoDB SDK for Go](https://github.com/aws/aws-sdk-go-v2).

---

Let me know if there’s anything else to refine or improve!