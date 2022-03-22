# core-api-template
Basic template to run core-api written in Go w/ go-flow


Environment variables
Variable	Required	Default Value	Description
ENV	NO	production	indicates in which environment app is running
LOG_LEVEL	NO	error	logging level
DB_DEV_CONNECTION	YES	-	Database connection string for DEV environment
DB_TEST_CONNECTION	YES	-	Database connection for TEST environment
DB_PRODUCTION_CONNECTION	YES	-	Database connection for PRODUCTION environment
