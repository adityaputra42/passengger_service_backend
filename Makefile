new_schema:
	go run -mod=mod entgo.io/ent/cmd/ent new $(name)

gen_ent:
	go generate ./ent

.PHONY: new_schema gen_ent
