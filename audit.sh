#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_NAME="${SERVICE_NAME:-newapi-audit-proxy}"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
BIN_DIR="${SCRIPT_DIR}/bin"
BIN_PATH="${BIN_DIR}/newapi-audit-proxy"
CONFIG_PATH="${CONFIG_PATH:-${SCRIPT_DIR}/config.yaml}"
GO_CMD="${GO_CMD:-go}"
SERVICE_USER="${SERVICE_USER:-$(id -un)}"
SERVICE_GROUP="${SERVICE_GROUP:-$(id -gn)}"
TAIL_LINES="${TAIL_LINES:-200}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing command: $1" >&2
    exit 1
  fi
}

run_root() {
  if [[ "${EUID}" -eq 0 ]]; then
    "$@"
  else
    sudo "$@"
  fi
}

build_binary() {
  need_cmd "${GO_CMD}"
  mkdir -p "${BIN_DIR}"
  (
    cd "${SCRIPT_DIR}"
    "${GO_CMD}" build -o "${BIN_PATH}" ./cmd/newapi-audit-proxy
  )
  chmod 0755 "${BIN_PATH}"
}

write_service_file() {
  local tmp_file
  tmp_file="$(mktemp "${SCRIPT_DIR}/.${SERVICE_NAME}.XXXXXX")"

  cat > "${tmp_file}" <<EOF
[Unit]
Description=NewAPI Audit Proxy
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${SERVICE_USER}
Group=${SERVICE_GROUP}
WorkingDirectory=${SCRIPT_DIR}
ExecStart=${BIN_PATH} -config ${CONFIG_PATH}
Restart=always
RestartSec=3
TimeoutStopSec=50
KillSignal=SIGTERM
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

  run_root install -m 0644 "${tmp_file}" "${SERVICE_FILE}"
  rm -f "${tmp_file}"
}

install_service() {
  build_binary
  if [[ ! -f "${CONFIG_PATH}" ]]; then
    echo "config not found: ${CONFIG_PATH}" >&2
    exit 1
  fi
  write_service_file
  run_root systemctl daemon-reload
  run_root systemctl enable --now "${SERVICE_NAME}"
  run_root systemctl status "${SERVICE_NAME}" --no-pager
}

start_service() {
  build_binary
  run_root systemctl start "${SERVICE_NAME}"
}

stop_service() {
  run_root systemctl stop "${SERVICE_NAME}"
}

restart_service() {
  build_binary
  run_root systemctl restart "${SERVICE_NAME}"
}

status_service() {
  run_root systemctl status "${SERVICE_NAME}" --no-pager
}

log_service() {
  run_root journalctl -u "${SERVICE_NAME}" -n "${TAIL_LINES}" -f
}

uninstall_service() {
  run_root systemctl disable --now "${SERVICE_NAME}" || true
  run_root rm -f "${SERVICE_FILE}"
  run_root systemctl daemon-reload
}

usage() {
  cat <<'EOF'
Usage: bash audit.sh <command>

Commands:
  build      Build ./bin/newapi-audit-proxy
  install    Build binary, install systemd service, enable and start it
  start      Start the systemd service
  stop       Stop the systemd service
  restart    Restart the systemd service
  status     Show service status
  log        Follow recent service logs
  uninstall  Disable service and remove the systemd unit
EOF
}

case "${1:-help}" in
  build)
    build_binary
    ;;
  install)
    install_service
    ;;
  start)
    start_service
    ;;
  stop)
    stop_service
    ;;
  restart)
    restart_service
    ;;
  status)
    status_service
    ;;
  log|logs)
    log_service
    ;;
  uninstall)
    uninstall_service
    ;;
  help|-h|--help)
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
