# create migration file for users table
migrate-users:
	migrate create -ext sql -dir internal/db/migrations/users -seq init_users_table

migrate-students:
	migrate create -ext sql -dir internal/db/migrations/students -seq init_students_table

postgres:
	docker run --name postgres-test -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=afitlms -p 5432:5432 -d postgres:latest
