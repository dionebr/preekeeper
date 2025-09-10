# Overview

This is a professional-grade web directory scanner (similar to gobuster/dirb) built with Go and featuring an advanced Bubble Tea terminal interface. The Preekeeper Scanner combines powerful concurrent HTTP scanning capabilities with a beautiful, interactive TUI that provides real-time progress tracking, color-coded results, and comprehensive filtering options.

# Recent Changes

**September 10, 2025**: Complete Preekeeper Scanner with Bubble Tea UI implementation
- Integrated advanced Preekeeper scanner tool with stunning Bubble Tea interface
- Implemented full-featured web directory scanner similar to gobuster
- Added FastHTTP-based concurrent scanning with configurable threads
- Created interactive TUI with real-time progress, results filtering, and controls
- Supports wordlist scanning, extensions, status code filtering, recursion
- Beautiful color-coded results display with scroll functionality
- Added comprehensive filtering options (size, lines, regex)
- Configured complete workflow for terminal-based scanning interface

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