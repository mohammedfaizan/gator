# Gator 🐊

Gator is a command-line RSS feed aggregator built in Go. It allows users to manage and view RSS feeds directly from the terminal. Keep up with your favorite blogs, news sites, podcasts, and more—all in one place.

## Features

- **Add RSS Feeds**: Easily add RSS feeds to be collected and aggregated.
- **Store Posts**: Collected posts are stored in a PostgreSQL database for easy access and management.
- **Follow/Unfollow Feeds**: Manage your feed subscriptions through the CLI.
- **View Summaries**: View summaries of aggregated posts directly in the terminal, with links to full content.

## Prerequisites

- **Go** (version 1.21.3 or later)
- **PostgreSQL** (version 15.x or later)

Ensure both Go and PostgreSQL are installed on your system before proceeding.

## Installation

### 1. Install Gator CLI

Use the `go install` command to install Gator:

```bash
go install github.com/mohammedfaizan/gator/cmd/gator@latest
```

This command downloads the Gator CLI and installs it in your Go `bin` directory. Make sure this directory is in your `PATH` environment variable.

### 2. Set Up the Configuration File

Create a configuration file named `.gator.toml` in your home directory:

```toml
db_url = "postgres://your_username:your_password@localhost:5432/gator?sslmode=disable"
current_user = ""
```

Replace `your_username`, `your_password`, and the database name (`gator`) with your actual PostgreSQL credentials and preferred database name.

### 3. Set Up the Database

- **Create a New Database**: Use PostgreSQL to create a new database named `gator` or your preferred name.
- **Update Credentials**: Ensure your database credentials in `.gator.toml` match those of your PostgreSQL setup.

### 4. Run Gator

You're ready to use Gator! Open your terminal and type:

```bash
gator
```

This will display the available commands.

## Usage

Gator offers several commands to manage your RSS feeds:

- **Register a New User**:

  ```bash
  gator register <username>
  ```

- **Log In as a User**:

  ```bash
  gator login <username>
  ```

- **Add a New Feed**:

  ```bash
  gator addfeed <name> <url>
  ```

- **List All Feeds You're Following**:

  ```bash
  gator feeds
  ```

- **Follow a Feed by URL**:

  ```bash
  gator follow <url>
  ```

- **Unfollow a Feed by URL**:

  ```bash
  gator unfollow <url>
  ```

- **Browse Recent Posts from Your Followed Feeds**:

  ```bash
  gator browse [limit]
  ```

  *Optional `limit` parameter specifies how many posts to display.*

- **Scrape Feeds and Update the Database**:

  ```bash
  gator scrape
  ```

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## Contact

For questions or feedback, please reach out to [Mohammed Faizan](mailto:faizan.mohammed7044@gmail.com).
