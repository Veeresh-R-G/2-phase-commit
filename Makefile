.PHONY: delivery store order all

delivery:
	@echo "Starting delivery service"
	go run ./delivery/main.go

store:
	@echo "Starting store service"
	go run ./store/main.go

order:
	@echo "Starting order service"
	go run ./orders/main.go

all:
	@echo "Starting all services"
	@make delivery & make store & make order
	
