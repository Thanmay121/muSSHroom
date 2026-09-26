# muSSHroom: An SSH-based terminal chat server with rooms and markdown

## What's muSSHroom?

muSSHroom is a chat server directly accessible right in your terminal, with a connection that’s set up using SSH protocol. There's no client to install, no web app to open — just enter the required `ssh` command and you’ll straight up enter a server with many other users along with you to chat with.

By building on the [Charm](https://charm.land/) stack (`wish` for the SSH server, `bubbletea` for the TUI, `lipgloss` for styling), the entire chat experience — private rooms, slash commands, colored system messages, emoji shortcodes — runs inside a standard terminal, rendered per-session as a proper interactive UI rather than a scroll of raw text. The result is a chat server that's trivial to connect.

## Running muSSHroom (skip to 4 if the server is already being hosted)
**1. Clone the repo**

```
git clone https://github.com/homebrew-ec-foss/muSSHroom.git
cd muSSHroom
```

**2. Build it**

You require Go 1.26.4 or newer.

```
go build -o musshroom .
```

**3. Run the server**

```bash
./musshroom
```

The default port is `3000` this way.

**Want a different port?** It's currently a hardcoded constant — open `internal/server/server.go` and change:

```go
const (
    Host = "localhost"
    Port = "3000"
)
```

to whatever port you want, then rebuild with `go build -o musshroom .`.

**Hosting it somewhere other than your laptop**

- **Home box / Raspberry Pi:** run the binary, forward the port on your router, connect using your public IP or a dynamic DNS hostname.
- **Free-tier VPS:** something like Oracle Cloud's free tier, or any other free server provider.

Whichever route you pick, the only real requirement is that the chosen port is reachable from outside (open in the firewall / security group), since that's the only thing a client needs to reach.

**4. Connect from the client side**

Assuming it's already hosted on `<hostname>`:

```
ssh <hostname> -p <port>
```

## Working of muSSHroom

`muSSHroom` is built as a per-connection TUI server — every SSH session gets its own isolated Bubbletea program, coordinated through a small set of shared, mutex-protected state. Here's the breakdown.

### 1. Connection handling via Wish

Every incoming SSH connection is intercepted by `wish` middleware, which hands it off to a fresh `bubbletea` program running in "SSH mode" instead of a local terminal. This is what lets each client get their own independently-rendered UI.

### 2. Session management

Each connected user is tracked as a session, keyed by username in a `map[string]*userSession`. Access to this map is mutex-protected to keep concurrent connect/disconnect/room-switch events from racing each other.

### 3. Rooms and message routing

Rooms are created via `/room USER1 USER2 ...`, and every session carries a reference to its currently active room so system messages (joins, leaves, renames) get routed to the room the event actually happened in.

### 4. Slash commands

A small command layer parses `/`-prefixed input (room creation, room deletion, renaming) separately from normal chat messages, dispatched through the Bubbletea `Update` loop (MVU pattern) rather than special-cased string matching scattered through the render path.

### 5. Emoji shortcodes

Message text is scanned for `:shortcode:` patterns and expanded before rendering, giving slack/discord-style emoji support without needing a graphical client.

### 6. Concurrency model

Because every session is its own Bubbletea program, most UI state is naturally isolated per-connection — the only shared, contended state is the session map and room membership, both behind mutexes. This keeps the concurrency surface small and auditable instead of threading locks through the whole message pipeline.

## Resources

- Project Repository: <https://github.com/homebrew-ec-foss/muSSHroom>
- [Charm](https://charm.land/) — `bubbletea`, `lipgloss`, `wish`, `bubbles`
