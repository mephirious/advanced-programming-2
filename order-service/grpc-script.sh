#!/bin/bash

# Create an Order
echo "Creating Order..."
grpcurl -plaintext -d '{
  "userId": "67f2c352315ca4f05670da8a",
  "items": [
    {
      "productId": "67f27ebc364656484b26e7f7",
      "quantity": 1
    },
    {
      "productId": "67f27f85016b85b83acad8b8",
      "quantity": 2
    }
  ]
}' localhost:8002 order.OrderService/CreateOrder

echo "Order Created!"

# Fetch an Order by ID
echo "Fetching Order by ID..."
grpcurl -plaintext -d '{ "id": "67f2c352315ca4f05670da8b" }' localhost:8002 order.OrderService/GetOrderByID

echo "Fetched Order!"

# Update Order Status
echo "Updating Order Status..."
grpcurl -plaintext -d '{
  "id": "67f2c352315ca4f05670da8b",
  "status": "completed"
}' localhost:8002 order.OrderService/UpdateOrderStatus

echo "Order Status Updated!"

# Delete Order
echo "Find orders by user_id..."
grpcurl -plaintext -d '{ "user_id": "67f274e7807b59ec2efe8ab9" }' localhost:8002 order.OrderService/ListUserOrders

echo "Fetched users orders"

