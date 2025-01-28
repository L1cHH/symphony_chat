# Load env file and read variables
function Read-EnvVariables {

    Write-Host "Reading environment variables from .env file..."

    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            $key = $matches[1]
            $value = $matches[2]

            [Environment]::SetEnvironmentVariable($key, $value)
        }
    }

    Write-Host "Environment variables read successfully" -ForegroundColor Green
}

function New-TestDatabase {
    $conn_to_postgres = "host=$($env:TEST_DB_HOST) port=$($env:TEST_DB_PORT) user=$($env:TEST_DB_USER) password=$($env:TEST_DB_PASSWORD) dbname=postgres sslmode=$($env:TEST_DB_SSLMODE)"

    Write-Host "Creating test database..."

    psql "$conn_to_postgres" -c "DROP DATABASE IF EXISTS $($env:TEST_DB_NAME) WITH (FORCE)"
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to drop existing test database"
    } 

    psql "$conn_to_postgres" -c "CREATE DATABASE $($env:TEST_DB_NAME)"
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to create test database"
    } 

    Write-Host "Test database created successfully" -ForegroundColor Green
}

function Start-Migrations {
    Write-Host "Applying migrations..."

    $conn_to_test_db = "postgres://$($env:TEST_DB_USER):$($env:TEST_DB_PASSWORD)@$($env:TEST_DB_HOST):$($env:TEST_DB_PORT)/$($env:TEST_DB_NAME)?sslmode=$($env:TEST_DB_SSLMODE)"

    migrate -database $conn_to_test_db -path migrations up
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to apply migrations"
    } 

    Write-Host "Migrations applied successfully" -ForegroundColor Green
}

function Remove_TestDatabase {
    Write-Host "Dropping test database..."

    $conn_to_postgres = "host=$($env:TEST_DB_HOST) port=$($env:TEST_DB_PORT) user=$($env:TEST_DB_USER) password=$($env:TEST_DB_PASSWORD) dbname=postgres sslmode=$($env:TEST_DB_SSLMODE)"

    psql "$conn_to_postgres" -c "DROP DATABASE IF EXISTS $($env:TEST_DB_NAME) WITH (FORCE)"
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to drop test database"
    } 

    Write-Host "Test database dropped successfully" -ForegroundColor Green
}

function Main {
    try {

        # Loading .env file
        Read-EnvVariables

        # Creating test database
        New-TestDatabase

        # Applying migrations
        Start-Migrations

        # Running integration tests
        Write-Host "Running integration tests..."
        go test ./tests/integration/... -v
        $testResult = $LASTEXITCODE

        exit $testResult
    } 
    catch {
        Write-Host "Failed to run script 'run_integration_tests': $_" -ForegroundColor Red
        return 1
    }
    finally {
        # Dropping test database
        Remove_TestDatabase
    }
}

# Running the main function
Main