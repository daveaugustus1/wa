# Check if Chocolatey is installed
if (!(Get-Command "choco.exe" -ErrorAction SilentlyContinue)) {
    # Install Chocolatey
    Set-ExecutionPolicy Bypass -Scope Process -Force
    iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
}

# Install Nmap using Chocolatey
choco install nmap -y --force

# Refresh PowerShell session
powershell.exe -Command "exit"

# Set the current working directory to the location of the script
Set-Location $PSScriptRoot

# Move config.toml to C:\Program Files\GoAgent\config\config.toml
$destinationDir = "C:\Program Files\GoAgent\config"
$destinationFile = Join-Path $destinationDir "config.toml"

# Create the destination directory if it doesn't exist
if (-not (Test-Path -Path $destinationDir)) {
    New-Item -ItemType Directory -Path $destinationDir | Out-Null
}

$etcDir = "C:\Program Files\GoAgent\etc"
# Create the etc directory if it doesn't exist
if (-not (Test-Path -Path $etcDir)) {
    New-Item -ItemType Directory -Path $etcDir | Out-Null
}

# Copy and overwrite the config.toml file
Copy-Item -Path ".\config.toml" -Destination $destinationFile -Force

Write-Host "config.toml has been moved to '$destinationFile'."


# Create a 'bin' folder in C:\Program Files\GoAgent directory
$binDir = "C:\Program Files\GoAgent\bin"
New-Item -ItemType Directory -Path $binDir -Force | Out-Null

# Build main.go in the bin folder and rename it to agent.go
$agentSourceFile = ".\goagent\cmd\main.go"
$agentDestinationFile = "C:\Program Files\GoAgent\bin\agent.go"
# go build -o $agentDestinationFile $agentSourceFile


# Check if main.exe exists in the current directory
if (Test-Path -Path $agentDestinationFile) {
    # If it does, ask the user if they want to rebuild
    $buildResponse = Read-Host "The exe file already exists. Do you want to rebuild? (y/n)"
    if ($buildResponse -eq "y" -or $buildResponse -eq "Y") {
        # If the user wants to rebuild, run the go build command
        go build -o $agentDestinationFile $agentSourceFile
    } else {
        # If the user doesn't want to rebuild, skip the go build command
        Write-Host "Skipping build."
    }
} else {
    # If main.exe doesn't exist, run the go build command
    go build -o $agentDestinationFile $agentSourceFile
}

# Define variables
$serviceName = "GoAgent"
$displayName = "GoAgent"
$description = "My custom service"

# Check if the service exists
$existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue

if ($existingService) {
    # Stop the service if it's running
    if ($existingService.Status -eq "Running") {
        Stop-Service -Name $serviceName
    }

    # Delete the service
    & sc.exe delete $serviceName
    Write-Host "The '$serviceName' service has been deleted."
}

# The service does not exist, so create it
Write-Host "Creating the '$displayName' service, pointing to executable '$agentDestinationFile'"
New-Service -Name $serviceName -BinaryPathName $agentDestinationFile -DisplayName $displayName -Description $description -StartupType Automatic

# Configure service recovery options
$failureActions = @"
    <sc.exe path> failure "$serviceName" reset= 86400 actions= restart/0
"@

$failureActionsPath = Join-Path $PSScriptRoot "failureactions.txt"
$failureActions | Out-File -FilePath $failureActionsPath -Encoding ASCII

& sc.exe failure $serviceName reset= 86400 actions= restart/0
& sc.exe failureflag $serviceName 1

# Set the service to restart on system restart
Set-Service -Name $serviceName -StartupType Automatic

# Start the service
Write-Host "Starting the '$displayName' service..."
Start-Service -Name $serviceName
