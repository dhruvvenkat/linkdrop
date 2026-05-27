#!/usr/bin/env bash

set -euo pipefail

APP_NAME="linkdrop"
DAEMON_NAME="linkdropd"

INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.config/linkdrop"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"

CONFIG_FILE="$CONFIG_DIR/config.env"
SERVICE_FILE="$SYSTEMD_USER_DIR/linkdropd.service"

HTTP_PORT="${LINKDROP_PORT:-4545}"
DISCOVERY_PORT="${LINKDROP_DISCOVERY_PORT:-4546}"

echo "Installing LinkDrop..."

if ! command -v go >/dev/null 2>&1; then
    echo "error: Go is not installed or not in PATH"
    exit 1
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
mkdir -p "$SYSTEMD_USER_DIR"

echo "Building binaries..."

go build -o "$INSTALL_DIR/$DAEMON_NAME" ./cmd/linkdropd
go build -o "$INSTALL_DIR/$APP_NAME" ./cmd/linkdrop

chmod +x "$INSTALL_DIR/$DAEMON_NAME"
chmod +x "$INSTALL_DIR/$APP_NAME"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Creating config at $CONFIG_FILE"

    TOKEN="$(python3 -c 'import secrets; print(secrets.token_urlsafe(32))')"

    cat > "$CONFIG_FILE" <<EOF
BEARER_TOKEN=$TOKEN
LINKDROP_HOST=0.0.0.0
LINKDROP_PORT=$HTTP_PORT
LINKDROP_SERVER_ADDR=http://localhost:$HTTP_PORT
EOF

    chmod 600 "$CONFIG_FILE"
else
    echo "Config already exists at $CONFIG_FILE; leaving it unchanged"
fi

echo "Creating systemd user service..."

cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=LinkDrop daemon
After=network-online.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/$DAEMON_NAME
Restart=on-failure
RestartSec=2
EnvironmentFile=$CONFIG_FILE

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload
systemctl --user enable linkdropd
systemctl --user restart linkdropd

echo
echo "LinkDrop installed."
echo
echo "Binaries:"
echo "  $INSTALL_DIR/$APP_NAME"
echo "  $INSTALL_DIR/$DAEMON_NAME"
echo
echo "Config:"
echo "  $CONFIG_FILE"
echo
echo "Service:"
echo "  systemctl --user status linkdropd"
echo "  journalctl --user -u linkdropd -f"
echo
echo "Make sure $INSTALL_DIR is in your PATH."
echo "If not, add this to your shell config:"
echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
echo
echo "Current token:"
grep '^BEARER_TOKEN=' "$CONFIG_FILE" | sed 's/^/  /'