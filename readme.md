# Clean Architecture-Based Microservices

This assignment aims to build a basic e-commerce platform using three distinct microservices. The platform follows the principles of **clean architecture** and best practices for software design, and is implemented using **Gin** for the API Gateway and service layers. The three core services are:

- **API Gateway** - Handles routing, logging, telemetry, and authentication.
- **Inventory Service** - Manages product inventory, categories, stock levels, and prices.
- **Order Service** - Handles order creation, updates, payments, and order tracking.

## Project Structure
```bash
├── gateway-service
│   └── cmd
├── inventory-service
│   ├── cmd
│   ├── config
│   ├── internal
│   │   ├── adapter
│   │   │   └── http
│   │   │       └── service
│   │   │           └── handler
│   │   ├── app
│   │   ├── domain
│   │   │   └── dto
│   │   ├── repository
│   │   └── usecase
│   └── pkg
│       └── mongo
└── order-service
    ├── cmd
    ├── config
    ├── internal
    │   ├── adapter
    │   │   └── http
    │   │       └── service
    │   │           ├── gateway
    │   │           └── handler
    │   ├── app
    │   ├── domain
    │   │   └── dto
    │   ├── repository
    │   └── usecase
    └── pkg
        └── mongo
```

## Prerequisites

- Go 1.20+ installed
- Linux system (tested on Ubuntu)
- Terminal supporting `gnome-terminal` or similar

## How to run

- Create .env file using .env.examples in each service
- You can run all 3 services using Makefile

```bash
make run
```

## API Endpoints

### Inventory Service

**Products:**
| Method | Endpoint              | Description              |
|--------|-----------------------|--------------------------|
| POST   | `/products`           | Create a new product     |
| GET    | `/products`           | List all products        |
| GET    | `/products/:id`       | Get product by ID        |
| PATCH  | `/products/:id`       | Update product by ID     |
| DELETE | `/products/:id`       | Delete product by ID     |

**Categories:**
| Method | Endpoint              | Description               |
|--------|-----------------------|---------------------------|
| POST   | `/categories`         | Create a new category     |
| GET    | `/categories`         | List all categories       |
| GET    | `/categories/:id`     | Get category by ID        |
| PATCH  | `/categories/:id`     | Update category by ID     |
| DELETE | `/categories/:id`     | Delete category by ID     |

### Order Service

**Orders:**
| Method | Endpoint              | Description               |
|--------|-----------------------|---------------------------|
| POST   | `/orders`             | Create a new order        |
| GET    | `/orders`             | List all orders           |
| GET    | `/orders/:id`         | Get order by ID           |
| PATCH  | `/orders/:id`         | Update order status by ID |

## Usage Example

### Get all products

**Endpoint:** `GET http://0.0.0.0:8003/api/v1/products`

**Response Body:**

```json
[
    {
        "id": "67f2725f807b59ec2efe8ab5",
        "name": "Laptop",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T12:23:59.423Z",
        "updated_at": "2025-04-06T12:23:59.423Z"
    },
    {
        "id": "67f27264807b59ec2efe8ab6",
        "name": "Smartphone",
        "description": "Latest model smartphone with all the top features.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 999.99,
        "stock": 150,
        "created_at": "2025-04-06T12:24:04.929Z",
        "updated_at": "2025-04-06T12:24:04.929Z"
    },
    {
        "id": "67f27ebc364656484b26e7f7",
        "name": "Laptop_1",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T13:16:44.1Z",
        "updated_at": "2025-04-06T13:16:44.1Z"
    },
    {
        "id": "67f27f85016b85b83acad8b8",
        "name": "Laptop_2",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T13:20:05.346Z",
        "updated_at": "2025-04-06T13:20:05.346Z"
    },
    {
        "id": "67f2efa1634fc1779d42c964",
        "name": "Laptop",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T21:18:25.721Z",
        "updated_at": "2025-04-06T21:18:25.721Z"
    },
    {
        "id": "67f2efa7634fc1779d42c965",
        "name": "Laptop-1",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T21:18:31.712Z",
        "updated_at": "2025-04-06T21:18:31.712Z"
    },
    {
        "id": "67f2efad634fc1779d42c966",
        "name": "Laptop-2",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T21:18:37.05Z",
        "updated_at": "2025-04-06T21:18:37.05Z"
    },
    {
        "id": "67f2efb1634fc1779d42c967",
        "name": "Laptop-2",
        "description": "A powerful laptop for gaming and development.",
        "category_id": "67f271f8807b59ec2efe8ab1",
        "price": 1200,
        "stock": 50,
        "created_at": "2025-04-06T21:18:41.134Z",
        "updated_at": "2025-04-06T21:18:41.134Z"
    }
]
```