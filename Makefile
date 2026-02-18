new_schema:
	go run -mod=mod entgo.io/ent/cmd/ent new $(name)

gen_ent:
	go generate ./ent


drop-db:
	docker exec -it postgres16 dropdb passenger_service

create-db:
	docker exec -it postgres16 createdb --username=root --owner=root passenger_service

server:
	go run cmd/api/main.go


.PHONY: new_schema gen_ent create-db drob-db server
