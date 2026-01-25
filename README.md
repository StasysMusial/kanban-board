# Minimal Kanban TUI
A simple somewhat configurable TUI for managing tasks built using the [Bubbletea][tea] and [Lip Gloss][lip] framework.

## Demo
![demo_gif](https://github.com/StasysMusial/kanban-board/blob/main/demo/demo.gif?raw=true)

## Features
</details>
<details><summary>Multiple Boards</summary>

kanban-board manages a unique board for every directory you execute it in. This allows separate task management spaces for each of your projects. When executed for the first time in a given directory kanban-board will prompt you if to create a new board.
</details>

</details>
<details><summary>Board Customization</summary>

Each board has it's own configuration file which is copied from a default config that is located in `~/.config/kanban`. Using this config file the user can adjust the following properties:

- Board (title and color)
- Tags (icon and color)
- Columns (title, icon and color)

The default config uses the name of the current directory as the board title and comes with four columns (`IDEAS`, `TO DO`, `IN PROGRESS` and `DONE`).
</details>

</details>
<details><summary>VIM-style Workflow</summary>

Boards are navigated and edited using a VIM-like input scheme:

```
h/j/k/l - select task
H/J/K/L - move task
g/G     - go to top/bottom
s/S     - sort by tags (descending/ascending)
a       - add task
x       - cut task
y       - yank task
p       - paste task
```

Similar to VIM, cutting a task will store it in the clipboard, allowing the user to paste it elsewhere.
</details>

</details>
<details><summary>Lip Gloss Color Support</summary>

All configurable colors support ANSI 16 (4-bit), ANSI 256 (8-bit) and True Color (24-bit) as defined by [Lip Gloss][lipcolors].
</details>

## Installation
You can install kanban-board by downloading a prebuilt binary from the [releases page][releases]. Prebuilt binaries are available for Windows and MacOS.

## Building From Source
Building from sources requires [the Go programming language][goinstall].

If you're on Linux or would like to build kanban-board yourself start by cloning the repo:

```bash
git clone https://github.com/StasysMusial/kanban-board
```

Then navigate into the directory:

```bash
cd kanban-board
```

Fetch the necessary Go packages:

```bash
go mod tidy
```

Build the application using the following command:

```bash
go build -o "build/" .
```

If the build succeeded you will find the executable in `kanban-board/builds`. Then add the executable to your `PATH`.

You should now be able to run the following command in any directory to start kanban-board and initialize a board:

```bash
kanban-board
```

## Configuration
Boards can be configured using a TOML file with the following keys:

</details>
<details><summary>Title</summary>

Format: `title: string`

The title of the board, which will be displayed in the bottom left corner of the screen. When set to `"DEFAULT"` the name of the working directory will be used instead.
</details>

</details>
<details><summary>Color</summary>

Format: `color: string`

The color used for rendering the project title. When set to `"DEFAULT"` the terminal foreground color will be used instead.
</details>

</details>
<details><summary>Tags</summary>

Format: `tags: [{ icon: string, color: string},...]`

Tags that can be added to tasks for sorting and organizational purposes. The order in which they are listed determines their sorting significance. Up to 8 tags are supported. Adding more will have no effect.
</details>

</details>
<details><summary>Columns</summary>

Format: `columns: [{ name: string, icon: string, color: string },...]`

The columns in your board from left to right. How many can fit on screen depends on your terminal window size.
</details>

All `color` fields support ANSI 16 (4-bit), ANSI 256 (8-bit) and True Color (24-bit) as defined by [Lip Gloss Colors][lipcolors].

### Default Config
The generated `default_config.toml` is located in `~/.config/kanban`.

```toml
title = "DEFAULT"
color = "DEFAULT"
tags = [
	{ icon="󰫢", color="#ff4cc4" },
	{ icon="󰅩", color="#89d789" },
	{ icon="", color="#5c84d6" },
	{ icon="󰃣", color="#f5d33d" },
]
columns = [
	{
		name="IDEAS",
		icon="",
		color="#5d4cff",
	},
	{
		name="TO DO",
		icon="",
		color="#ff4cc4",
	},
	{
		name="IN PROGRESS",
		icon="",
		color="#ffcb4c",
	},
	{
		name="FINISHED",
		icon="",
		color="#bcff4c",
	},
]
```

### Board Config
When you first initialize a board, `.kanban/config.toml` will be created in your working directory. `config.toml` is a copy of `default_config.toml` so modify the ladder if you want to change the blueprint for new boards and modify the former to configure your current board.

## Notes
This project is functional but lacks some features which might make or break viability for you:

- undo and redo
- input customization
- full color customization and theme support
- adding or removing columns directly inside the app
<!--external Links-->
[tea]: https://github.com/charmbracelet/bubbletea
[lip]: https://github.com/charmbracelet/lipgloss
[lipcolors]: https://github.com/charmbracelet/lipgloss#colors
[releases]: https://github.com/StasysMusial/kanban-board/releases
[goinstall]: https://go.dev/doc/install
