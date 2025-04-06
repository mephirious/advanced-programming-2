.PHONY: run

run:
	gnome-terminal -- bash -c "cd gateway-service && go run ./cmd; exec bash" &
	gnome-terminal -- bash -c "cd inventory-service && go run ./cmd; exec bash" &
	gnome-terminal -- bash -c "cd order-service && go run ./cmd; exec bash"
