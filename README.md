# Celerix Git GUI

> [!CAUTION]
> **Pre-Alpha State**: This software is in a very early pre-alpha state and is NOT ready for day-to-day use yet! Use it only if you enjoy cutting-edge software and fully understand the risks involved.

<img src="celerix-logo.png" alt="Celerix Logo" width="200" />

Celerix is a high-performance, cross-platform Git desktop client built with **Wails**, **Go**, and **Vue 3**. It offers a clean, modern interface designed for developers who want a fast and reliable way to manage their repositories without the bloat of traditional GUIs.

## ✨ Features

### 🚀 Repository Management
- **Multiple Tabs**: Work on several repositories simultaneously with a tabbed interface.
- **Recent Repositories**: Quickly jump back into your most recently used projects.
- **Initialization**: Easily initialize new Git repositories directly from the UI.
- **File Manager Integration**: Open repository folders in your system's default file manager with one click.

### 📜 Commit History & Visualization
- **Interactive Git Graph**: Visualize your branch history, merges, and tags with a clean, color-coded graph.
- **Rich Commit Details**: View full commit information including subject, body, author, and date.
- **File Changes List**: See exactly which files were Added, Modified, or Deleted in each commit.
- **Side-by-Side Diff**: High-performance diff viewer to compare changes between commits.
- **Flexible Layout**: Toggle between list-focused and detail-focused views to see more of what matters.

### 🛠️ Working Copy & Committing
- **Staging Management**: Stage or unstage individual files or all changes at once.
- **Smart Commits**: Support for commit subjects and detailed descriptions.
- **Amend Support**: Easily amend your last commit, including pre-filling the previous commit message.
- **Local Diff**: Review your unstaged and staged changes before committing.

### 🌿 Branching & Tagging
- **Branch Management**: Create, checkout, and track local and remote branches.
- **Tagging**: Create lightweight or annotated tags to mark important milestones.
- **Contextual Actions**: Create branches and tags directly from any point in the history.

### ☁️ Remote Operations
- **Fetch, Pull & Push**: Seamlessly synchronize with remote repositories (origin).
- **Progress Tracking**: Real-time progress bar and status updates for long-running Git operations.
- **SSH Authentication**: Built-in SSH key generation and management (Ed25519) for secure remote access.

### 🛡️ Performance & Safety
- **Thread-Safe Backend**: Per-repository locking mechanism ensures that concurrent operations from multiple tabs never corrupt your Git state.
- **Optimized Stats**: Fast calculation of repository size and statistics.
- **Clean Architecture**: Refactored component-based frontend and a robust Go backend.

## 💻 Tech Stack

- **Backend**: [Go](https://go.dev/) with [go-git](https://github.com/go-git/go-git)
- **Frontend**: [Vue 3](https://vuejs.org/), [TypeScript](https://www.typescriptlang.org/), [Vite](https://vitejs.dev/)
- **Framework**: [Wails v2](https://wails.io/)
- **Styling**: [Bootstrap 5](https://getbootstrap.com/)
- **Icons**: [Tabler Icons](https://tabler-icons.io/)

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (1.24+ recommended)
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Development Mode

To run the application in development mode with hot-reloading:

```bash
wails dev
```

### Building for Production

To create a production-ready redistribution package:

```bash
wails build
```

The resulting binary will be in the `build/bin` directory.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details (if applicable).

---
*Built with ❤️ by the Celerix team.*