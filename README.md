# Lua Actuator

Lua Actuator is a project developed using Golang that provides a Lua runtime, allowing for flexible addition of scheduled scripts, sending HTTP requests, and managing daily small scripts. The project includes a web interface for submitting, executing, scheduling, and managing scripts.

## Features

- **Lua Runtime**: Execute Lua scripts within a Golang environment.
- **Script Management**: Save, update, and delete scripts.
- **Scheduling**: Schedule scripts to run at specified times.
- **HTTP Requests**: Send HTTP requests from within Lua scripts.
- **Web Interface**: User-friendly web interface for managing scripts.

## Getting Started

### Prerequisites

- Golang 1.16 or higher
- Node.js and npm

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/lua-actuator.git
    cd lua-actuator
    ```

2. **Build the frontend:**

    ```bash
    cd frontend
    npm install
    npm run build
    cd ..
    ```

3. **Run the backend:**

    ```bash
    go run main.go
    ```

### Project Structure

- **main.go**: The main Golang file that sets up the server and routes.
- **frontend**: Contains the React frontend code.
- **scripts.db**: SQLite database file for storing scripts and key-value pairs.

### API Endpoints

- `POST /api/submit`: Submit a new script.
- `GET /api/scripts/:name/execute`: Execute a script by name.
- `GET /api/scripts`: Get all scripts.
- `POST /api/scripts/:name/schedule`: Schedule a script.
- `PUT /api/update/:name`: Update a script.

### Frontend Components

- **Home**: The home page of the application.
- **SubmitScript**: Page for submitting new scripts.
- **ExecuteScript**: Page for executing scripts.
- **ScheduleScript**: Page for scheduling scripts.
- **EditScript**: Modal for editing scripts.
- **Modal**: Generic modal component for displaying messages.

### Example Lua Script
```lua
-- Example Lua Script: Hello World

-- This script prints "Hello, World!" to the console
print("Hello, World!")

-- Example Lua Script: Sum of Two Numbers

-- This script calculates the sum of two numbers and prints the result
local function sum(a, b)
    return a + b
end

local num1 = 5
local num2 = 10
local result = sum(num1, num2)
print("The sum of " .. num1 .. " and " .. num2 .. " is " .. result)

-- Example Lua Script: Factorial

-- This script calculates the factorial of a number
local function factorial(n)
    if n == 0 then
        return 1
    else
        return n * factorial(n - 1)
    end
end

local number = 5
local fact = factorial(number)
print("The factorial of " .. number .. " is " .. fact)
```