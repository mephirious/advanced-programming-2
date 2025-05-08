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
}' $GRPC_SERVER $SERVICE/ListProducts

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

