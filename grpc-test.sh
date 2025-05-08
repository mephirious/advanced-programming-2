#!/bin/bash
# protoc   --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative   proto/inventory.proto

GRPC_SERVER=localhost:8001
SERVICE=inventory.InventoryService


echo "Creating Product..."
grpcurl -plaintext -d '{
  "name": "Smartphone",
  "description": "Latest model smartphone with all the top features.",
  "category_id": "67f271f8807b59ec2efe8ab1",
  "price": 999.99,
  "stock": 150
}' $GRPC_SERVER $SERVICE/CreateProduct

echo -e "\nGetting Product By ID..."
grpcurl -plaintext -d '{
  "id": "67f27264807b59ec2efe8ab6"
}' $GRPC_SERVER $SERVICE/GetProductByID

echo -e "\nUpdating Product..."
grpcurl -plaintext -d '{
  "id": "681c4cc31fb89e229e153ae2",
  "price": 899.99,
  "stock": 100
}' $GRPC_SERVER $SERVICE/UpdateProduct

echo -e "\nDeleting Product..."
grpcurl -plaintext -d '{
  "id": "681c4cc31fb89e229e153ae2"
}' $GRPC_SERVER $SERVICE/DeleteProduct

echo -e "\nListing Products..."
grpcurl -plaintext -d '{
  "limit": 10,
  "page": 1,
  "sort_by": "price",
  "sort_order": "asc"
}' $GRPC_SERVER $SERVICE/ListProduct

echo -e "\nCreating Category..."
grpcurl -plaintext -d '{
  "name": "Electronics",
  "description": "All kinds of electronic gadgets and accessories."
}' $GRPC_SERVER $SERVICE/CreateCategory

echo -e "\nGetting Category By ID..."
grpcurl -plaintext -d '{
  "id": "67fb844575ff1d6d8cd19193"
}' $GRPC_SERVER $SERVICE/GetCategoryByID

echo -e "\nUpdating Category..."
grpcurl -plaintext -d '{
  "id": "67fb844575ff1d6d8cd19193",
  "description": "Updated description for electronics."
}' $GRPC_SERVER $SERVICE/UpdateCategory

echo -e "\nDeleting Category..."
grpcurl -plaintext -d '{
  "id": "67fb844575ff1d6d8cd19193"
}' $GRPC_SERVER $SERVICE/DeleteCategory

echo -e "\nListing Categories..."
grpcurl -plaintext -d '{}' $GRPC_SERVER $SERVICE/ListCategories

echo -e "\nDone."


#!/bin/bash

GRPC_SERVER=localhost:8002
SERVICE=order.OrderService

  echo "Creating Order..."
  grpcurl -plaintext -d '{
    "user_id": "67f271f8807b59ec2efe8ab1",
    "items": [
      {
        "product_id": "67f2efa1634fc1779d42c964",
        "quantity": 2
      },
      {
        "product_id": "67f27ebc364656484b26e7f7",
        "quantity": 1
      }
    ]
  }' $GRPC_SERVER $SERVICE/CreateOrder

echo -e "\nGetting Order By ID..."
grpcurl -plaintext -d '{
  "id": "681c5a19cb5964c723ac8006"
}' $GRPC_SERVER $SERVICE/GetOrderByID

echo -e "\nUpdating Order Status..."
grpcurl -plaintext -d '{
  "id": "681c5a19cb5964c723ac8006",
  "status": "COMPLETED"
}' $GRPC_SERVER $SERVICE/UpdateOrderStatus

echo -e "\nListing Orders for User..."
grpcurl -plaintext -d '{
  "user_id": "67f271f8807b59ec2efe8ab1",
  "page": 1,
  "limit": 10
}' $GRPC_SERVER $SERVICE/ListUserOrders

echo -e "\nDone."
