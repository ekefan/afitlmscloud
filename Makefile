# create migration file for users table
migrate-users:
	migrate create -ext sql -dir internal/db/migrations/user -seq init_users_table

migrate-students:
	migrate create -ext sql -dir internal/db/migrations/student -seq init_students_table

migrate-lecturers:
	migrate create -ext sql -dir internal/db/migrations/lecturer -seq init_lecturers_table

migrate-courses:
	migrate create -ext sql -dir internal/db/migrations/course -seq init_courses_table

migrate-attendances:
	migrate create -ext sql -dir internal/db/migrations/attendance -seq init_attendances_table

postgres:
	docker run --name postgres-test -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=afitlms -p 5432:5432 -d postgres:latest
