# Overview

This is a Go-based terminal user interface (TUI) application that creates a scanner interface similar to gobuster. The project uses the Bubble Tea framework for building interactive terminal applications and Lip Gloss for styling and visual formatting. The application simulates URL scanning and displays HTTP status codes in a colorful, interactive terminal interface.

# Recent Changes

**September 10, 2025**: Complete Bubble Tea Scanner setup and deployment
- Installed Go 1.24 language module
- Created main.go with full TUI scanner implementation
- Initialized Go module (bubbletea-scan) with required dependencies
- Successfully installed Bubble Tea and Lip Gloss libraries
- Configured and started "Bubble Tea Scanner" workflow
- Application running successfully with colorful URL scan simulation

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## Core Framework
The application is built on the Bubble Tea framework, which follows the Elm architecture pattern with a model-view-update cycle. This provides a reactive and predictable way to handle terminal user interface state management.

## Application Structure
- **Model**: Contains the application state including scan results and quit status
- **Messages**: Uses a tickMsg type to handle asynchronous result updates
- **Styling**: Leverages Lip Gloss for consistent visual formatting with predefined color schemes

## Data Flow
The application uses a result struct to encapsulate URL scanning outcomes, storing both the target URL and HTTP status code. Results are collected in a slice within the main model, allowing for real-time display updates as scanning progresses.

## Visual Design
The styling system uses color-coded output to differentiate between different HTTP response types:
- Green for successful responses (2xx)
- Red for not found errors (404)
- Purple for forbidden access (403)
- Gray for neutral/other status codes

## Terminal Interface Pattern
The application follows the standard Bubble Tea TUI pattern where the interface updates reactively based on incoming messages, providing a smooth user experience similar to modern CLI tools like gobuster.

# External Dependencies

## Go Modules
- **github.com/charmbracelet/bubbletea**: Core TUI framework for building interactive terminal applications
- **github.com/charmbracelet/lipgloss**: Styling library for terminal output formatting and color management

## Runtime Requirements
- Go programming language environment
- Terminal with color support for proper visual rendering
- No external databases or web services required for core functionality