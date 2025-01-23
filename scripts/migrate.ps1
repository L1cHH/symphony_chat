param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("up", "down", "create")]
    [string]$Command,
    
    [Parameter(Mandatory=$false)]
    [string]$Name,
    
    [Parameter(Mandatory=$false)]
    [int]$Steps
)

# Загрузка переменных окружения из .env файла
$envContent = Get-Content .env
foreach ($line in $envContent) {
    if ($line -match '^([^=]+)=(.*)$') {
        $key = $matches[1]
        $value = $matches[2]
        [Environment]::SetEnvironmentVariable($key, $value)
    }
}

# Формируем строку подключения
$DATABASE_URL = "postgres://$($env:DB_USER):$($env:DB_PASSWORD)@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)?sslmode=$($env:DB_SSLMODE)"
Write-Host "Connection string: $DATABASE_URL"

switch ($Command) {
    "up" {
        if ($Steps) {
            Write-Host "Applying $Steps migrations..."
            migrate -database $DATABASE_URL -path migrations up $Steps
        }
        else {
            Write-Host "Applying all migrations..."
            migrate -database $DATABASE_URL -path migrations up
        }
    }
    "down" {
        if ($Steps) {
            Write-Host "Rolling back $Steps migrations..."
            migrate -database $DATABASE_URL -path migrations down $Steps
        }
        else {
            Write-Host "Rolling back all migrations..."
            migrate -database $DATABASE_URL -path migrations down
        }
    }
    "create" {
        if (-not $Name) {
            Write-Error "Name parameter is required for create command"
            exit 1
        }
        Write-Host "Creating new migration: $Name"
        migrate create -ext sql -dir migrations -seq $Name
    }
}
