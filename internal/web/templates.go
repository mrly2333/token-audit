package web

const baseTemplate = `
{{define "base"}}
<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="csrf-token" content="{{.CSRFToken}}">
  <title>{{.Title}}</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f4f1ea;
      --panel: #fffdf9;
      --ink: #1f2933;
      --muted: #66707a;
      --line: #ddd4c7;
      --accent: #1f6f5f;
      --danger: #b42318;
      --ok: #166534;
      --shadow: 0 18px 40px rgba(31, 41, 51, 0.08);
      --shadow-strong: 0 26px 60px rgba(31, 41, 51, 0.14);
      --glass: rgba(255, 253, 249, 0.78);
      --glass-strong: rgba(255, 253, 249, 0.9);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      position: relative;
      font-family: "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
      background:
        radial-gradient(circle at top left, rgba(217, 144, 88, 0.18), transparent 34%),
        radial-gradient(circle at top right, rgba(31, 111, 95, 0.16), transparent 28%),
        var(--bg);
      color: var(--ink);
    }
    body::before {
      content: "";
      position: fixed;
      inset: -12vh -8vw auto;
      height: 46vh;
      background:
        radial-gradient(circle at 20% 20%, rgba(217, 144, 88, 0.18), transparent 34%),
        radial-gradient(circle at 78% 26%, rgba(31, 111, 95, 0.16), transparent 30%);
      filter: blur(20px);
      pointer-events: none;
      z-index: 0;
    }
    a { color: var(--accent); text-decoration: none; }
    a:hover { text-decoration: underline; }
    code {
      font-family: Consolas, "Courier New", monospace;
      font-size: 12px;
      word-break: break-all;
    }
    .shell {
      max-width: 1440px;
      margin: 0 auto;
      padding: 24px;
      position: relative;
      z-index: 1;
      animation: fade-up 0.45s ease both;
    }
    .nav {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;
      margin-bottom: 24px;
      padding: 18px 20px;
      background: rgba(255, 253, 249, 0.92);
      border: 1px solid var(--line);
      border-radius: 20px;
      box-shadow: var(--shadow);
      backdrop-filter: blur(18px) saturate(1.14);
      position: sticky;
      top: 10px;
      z-index: 12;
    }
    .brand {
      font-size: 22px;
      font-weight: 700;
      letter-spacing: 0.02em;
    }
    .nav-links {
      display: flex;
      gap: 16px;
      align-items: center;
      flex-wrap: wrap;
    }
    .toolbar-actions {
      display: flex;
      gap: 12px;
      align-items: center;
      flex-wrap: wrap;
    }
    .button, button {
      border: 0;
      border-radius: 12px;
      padding: 10px 16px;
      background: var(--accent);
      color: #fff;
      font: inherit;
      cursor: pointer;
      transition: transform 0.18s ease, box-shadow 0.18s ease, opacity 0.18s ease, background 0.18s ease;
      box-shadow: 0 10px 20px rgba(31, 111, 95, 0.18);
    }
    .button:hover, button:hover {
      transform: translateY(-1px);
      box-shadow: 0 16px 28px rgba(31, 111, 95, 0.22);
    }
    .button:disabled, button:disabled {
      opacity: 0.45;
      cursor: not-allowed;
      transform: none;
      box-shadow: none;
    }
    .button.secondary {
      background: rgba(255, 255, 255, 0.82);
      color: var(--accent);
      border: 1px solid var(--accent);
      box-shadow: 0 10px 18px rgba(31, 41, 51, 0.08);
    }
    .button.danger, button.danger {
      background: var(--danger);
      color: #fff;
    }
    .panel {
      background: var(--glass);
      border: 1px solid var(--line);
      border-radius: 22px;
      box-shadow: var(--shadow);
      padding: 20px;
      margin-bottom: 20px;
      backdrop-filter: blur(18px);
      animation: fade-up 0.42s ease both;
      position: relative;
      overflow: hidden;
      isolation: isolate;
    }
    .panel::before {
      content: "";
      position: absolute;
      inset: 0;
      background:
        linear-gradient(135deg, rgba(255, 255, 255, 0.18), transparent 44%),
        radial-gradient(circle at top right, rgba(255, 255, 255, 0.22), transparent 30%);
      opacity: 0.82;
      pointer-events: none;
      z-index: 0;
    }
    .panel > * {
      position: relative;
      z-index: 1;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 16px;
    }
    .stat {
      padding: 18px;
      border-radius: 18px;
      background: linear-gradient(160deg, rgba(31, 111, 95, 0.1), rgba(217, 144, 88, 0.12));
      border: 1px solid rgba(31, 111, 95, 0.16);
      backdrop-filter: blur(10px);
    }
    .stat h3 {
      margin: 0 0 8px;
      color: var(--muted);
      font-size: 14px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }
    .stat .value {
      font-size: 34px;
      font-weight: 800;
      line-height: 1.1;
      word-break: break-word;
    }
    .muted { color: var(--muted); }
    .hero {
      display: flex;
      justify-content: space-between;
      align-items: end;
      gap: 16px;
      flex-wrap: wrap;
      margin-bottom: 18px;
    }
    .hero h1 {
      margin: 0 0 8px;
      font-size: 30px;
    }
    .hero p {
      margin: 0;
      color: var(--muted);
    }
    .pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 6px 10px;
      border-radius: 999px;
      background: rgba(31, 111, 95, 0.08);
      color: var(--accent);
      font-size: 12px;
      font-weight: 700;
      transition: transform 0.24s ease, box-shadow 0.24s ease, background 0.24s ease, color 0.24s ease;
    }
    .pill.is-checking {
      background: rgba(217, 144, 88, 0.16);
      color: #8a4d16;
      box-shadow: 0 0 0 1px rgba(217, 144, 88, 0.14), 0 10px 24px rgba(217, 144, 88, 0.18);
      animation: pill-breathe 1.35s ease-in-out infinite;
    }
    .pill.has-updates {
      background: rgba(31, 111, 95, 0.14);
      box-shadow: 0 0 0 1px rgba(31, 111, 95, 0.16), 0 14px 28px rgba(31, 111, 95, 0.2);
      animation: pill-celebrate 0.9s cubic-bezier(.2,.8,.2,1);
    }
    .pill.is-error {
      background: rgba(180, 35, 24, 0.12);
      color: var(--danger);
      box-shadow: 0 0 0 1px rgba(180, 35, 24, 0.16), 0 12px 26px rgba(180, 35, 24, 0.12);
    }
    .form-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 12px;
      align-items: end;
    }
    label {
      display: block;
      margin-bottom: 6px;
      font-size: 13px;
      font-weight: 600;
      color: var(--muted);
    }
    input, select, textarea {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 12px;
      padding: 10px 12px;
      font: inherit;
      background: #fff;
      color: var(--ink);
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 14px;
    }
    th, td {
      padding: 12px 10px;
      border-bottom: 1px solid var(--line);
      vertical-align: top;
      text-align: left;
    }
    th {
      color: var(--muted);
      font-weight: 700;
      font-size: 12px;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .status-ok { color: var(--ok); font-weight: 700; }
    .status-err { color: var(--danger); font-weight: 700; }
    pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #f8f5ef;
      border: 1px solid var(--line);
      border-radius: 16px;
      padding: 16px;
      margin: 0;
      overflow-x: auto;
      font-family: Consolas, "Courier New", monospace;
      font-size: 13px;
      line-height: 1.55;
    }
    .detail-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
      gap: 20px;
    }
    .detail-grid.content-first {
      grid-template-columns: 1fr;
    }
    .token-form {
      display: flex;
      gap: 8px;
      align-items: center;
    }
    .token-form input {
      min-width: 140px;
    }
    .login-shell {
      min-height: 100vh;
      display: grid;
      place-items: center;
      padding: 24px;
    }
    .login-card {
      width: min(460px, 100%);
      background: var(--glass-strong);
      border: 1px solid var(--line);
      border-radius: 28px;
      box-shadow: var(--shadow-strong);
      padding: 28px;
      backdrop-filter: blur(22px);
      animation: modal-in 0.3s ease both;
    }
    .error-box {
      margin-bottom: 16px;
      padding: 12px 14px;
      border-radius: 14px;
      color: var(--danger);
      background: rgba(180, 35, 24, 0.08);
      border: 1px solid rgba(180, 35, 24, 0.16);
    }
    .success-box {
      margin-bottom: 16px;
      padding: 12px 14px;
      border-radius: 14px;
      color: var(--ok);
      background: rgba(22, 101, 52, 0.08);
      border: 1px solid rgba(22, 101, 52, 0.16);
    }
    .pagination {
      display: flex;
      gap: 12px;
      justify-content: flex-end;
      align-items: center;
      margin-top: 16px;
    }
    .shell-wide {
      max-width: 100%;
      padding: 14px 16px 20px;
    }
    .monitor-layout {
      display: grid;
      grid-template-columns: minmax(260px, 18%) minmax(0, 1fr);
      gap: 16px;
      align-items: start;
    }
    .log-sidebar, .log-main {
      min-height: calc(100vh - 170px);
    }
    .sidebar-header {
      display: flex;
      justify-content: space-between;
      align-items: end;
      gap: 12px;
      margin-bottom: 12px;
    }
    .log-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
      max-height: 72vh;
      overflow: auto;
      padding-right: 4px;
      position: relative;
      scroll-behavior: smooth;
    }
    .log-list::before {
      content: "";
      position: sticky;
      top: 0;
      display: block;
      height: 0;
      box-shadow: 0 0 90px 18px rgba(31, 111, 95, 0);
      pointer-events: none;
      z-index: 2;
      transition: box-shadow 0.45s ease;
    }
    .log-list.list-surge::before {
      box-shadow: 0 10px 120px 24px rgba(31, 111, 95, 0.28);
    }
    .log-sidebar.is-refreshing .log-list::before {
      box-shadow: 0 8px 80px 16px rgba(217, 144, 88, 0.18);
    }
    .log-item {
      width: 100%;
      text-align: left;
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.82);
      color: var(--ink);
      border-radius: 18px;
      padding: 14px;
      font: inherit;
      cursor: pointer;
      transition: transform 0.16s ease, border-color 0.16s ease, box-shadow 0.16s ease, background 0.16s ease;
      backdrop-filter: blur(16px) saturate(1.08);
      box-shadow: 0 8px 18px rgba(31, 41, 51, 0.05);
      position: relative;
      overflow: hidden;
      transform-origin: top center;
    }
    .log-item::after {
      content: "";
      position: absolute;
      inset: 0;
      background:
        linear-gradient(135deg, rgba(255, 255, 255, 0.22), transparent 38%),
        radial-gradient(circle at top right, rgba(31, 111, 95, 0.12), transparent 34%);
      opacity: 0.92;
      pointer-events: none;
      transition: opacity 0.24s ease;
    }
    .log-item > * {
      position: relative;
      z-index: 1;
    }
    .log-item:hover {
      transform: translateY(-2px);
      box-shadow: 0 16px 28px rgba(31, 41, 51, 0.1);
    }
    .log-item.active {
      border-color: var(--accent);
      box-shadow: 0 0 0 2px rgba(31, 111, 95, 0.12);
      background: rgba(31, 111, 95, 0.06);
    }
    .log-item.active::after {
      opacity: 1;
      background:
        linear-gradient(135deg, rgba(255, 255, 255, 0.24), transparent 34%),
        radial-gradient(circle at top right, rgba(31, 111, 95, 0.18), transparent 34%);
    }
    .log-item.log-item-enter {
      animation: log-card-enter 0.7s cubic-bezier(.18,.86,.28,1) both, log-card-glow 1.8s ease-out both;
    }
    .log-item.log-item-update {
      animation: log-card-update 0.55s ease both;
    }
    .log-item.log-item-selected-pulse {
      animation: selected-glow 0.52s ease both;
    }
    .log-item.is-new::before {
      content: "NEW";
      position: absolute;
      top: 12px;
      right: 12px;
      padding: 4px 8px;
      border-radius: 999px;
      background: linear-gradient(135deg, rgba(31, 111, 95, 0.96), rgba(217, 144, 88, 0.92));
      color: #fff;
      font-size: 10px;
      font-weight: 800;
      letter-spacing: 0.08em;
      box-shadow: 0 12px 24px rgba(31, 111, 95, 0.26);
      animation: new-badge-pop 0.78s cubic-bezier(.2,.9,.28,1) both;
      z-index: 3;
    }
    .count-burst {
      animation: pill-celebrate 0.76s cubic-bezier(.2,.8,.2,1);
    }
    .log-item-top {
      display: flex;
      justify-content: space-between;
      align-items: start;
      gap: 12px;
      margin-bottom: 8px;
    }
    .log-item-alias {
      font-weight: 800;
      font-size: 15px;
      line-height: 1.3;
      word-break: break-word;
    }
    .log-item-model {
      margin-top: 4px;
      color: var(--muted);
      font-size: 12px;
      word-break: break-word;
    }
    .log-item-time {
      color: var(--muted);
      font-size: 12px;
      white-space: nowrap;
    }
    .log-item-body {
      font-size: 13px;
      line-height: 1.55;
      color: var(--ink);
      word-break: break-word;
    }
    .log-item-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 10px;
    }
    .meta-badge {
      display: inline-flex;
      align-items: center;
      padding: 4px 8px;
      border-radius: 999px;
      background: #f4efe5;
      color: var(--muted);
      font-size: 12px;
      font-weight: 700;
    }
    .detail-empty {
      min-height: 64vh;
      display: grid;
      place-items: center;
      text-align: center;
      color: var(--muted);
      border: 1px dashed var(--line);
      border-radius: 22px;
      background: rgba(255, 255, 255, 0.55);
      padding: 40px 24px;
    }
    .detail-hidden {
      display: none;
    }
    .detail-animate {
      animation: detail-fade 0.42s cubic-bezier(.22,.8,.32,1);
    }
    .detail-hero {
      display: flex;
      justify-content: space-between;
      align-items: start;
      gap: 16px;
      flex-wrap: wrap;
      margin-bottom: 18px;
    }
    .pill-row {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      margin-bottom: 10px;
    }
    .conversation-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 18px;
      margin-bottom: 18px;
    }
    .conversation-card {
      padding: 18px;
      border-radius: 20px;
      border: 1px solid var(--line);
      background: rgba(255, 253, 251, 0.8);
      backdrop-filter: blur(14px);
      box-shadow: 0 12px 28px rgba(31, 41, 51, 0.06);
      position: relative;
      overflow: hidden;
    }
    .conversation-card::before {
      content: "";
      position: absolute;
      inset: 0;
      background: linear-gradient(135deg, rgba(255, 255, 255, 0.22), transparent 42%);
      pointer-events: none;
      opacity: 0.95;
    }
    .conversation-card.user {
      background: linear-gradient(180deg, rgba(217, 144, 88, 0.12), rgba(255, 253, 249, 0.95));
    }
    .conversation-card.assistant {
      background: linear-gradient(180deg, rgba(31, 111, 95, 0.10), rgba(255, 253, 249, 0.95));
    }
    .conversation-card pre {
      background: transparent;
      border: 0;
      padding: 0;
      min-height: 34vh;
      font-size: 15px;
      line-height: 1.72;
      position: relative;
      z-index: 1;
    }
    .section-kicker {
      margin-bottom: 10px;
      color: var(--muted);
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .meta-strip {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
      gap: 12px;
      margin-bottom: 16px;
    }
    .meta-chip {
      border: 1px solid var(--line);
      border-radius: 16px;
      padding: 12px 14px;
      background: rgba(255, 255, 255, 0.74);
      backdrop-filter: blur(12px);
    }
    .meta-chip span {
      display: block;
      color: var(--muted);
      font-size: 12px;
      margin-bottom: 6px;
    }
    .meta-chip strong {
      font-size: 18px;
      line-height: 1.25;
      word-break: break-word;
    }
    .details-stack {
      display: grid;
      gap: 16px;
    }
    .log-main {
      padding: 24px 26px;
    }
    .log-main .detail-hero {
      margin-bottom: 22px;
    }
    .log-main .meta-strip {
      grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
      margin-bottom: 14px;
    }
    .log-main .details-stack details {
      background: rgba(255, 255, 255, 0.82);
    }
    .mono-scroll {
      max-height: 280px;
      overflow: auto;
    }
    .empty-list {
      padding: 22px;
      text-align: center;
      border: 1px dashed var(--line);
      border-radius: 18px;
      color: var(--muted);
      background: rgba(255, 255, 255, 0.6);
    }
    .log-pagination {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      margin-top: 14px;
      padding-top: 14px;
      border-top: 1px solid var(--line);
    }
    .log-pagination span {
      text-align: center;
      flex: 1;
    }
    .text-pagination {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      margin-top: 12px;
    }
    .text-pagination span {
      flex: 1;
      text-align: center;
    }
    .paged-pre {
      margin-bottom: 0;
    }
    .token-manual-grid {
      display: grid;
      grid-template-columns: minmax(0, 2fr) minmax(180px, 1fr) auto;
      gap: 12px;
      align-items: end;
    }
    .gallery {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 14px;
      margin-bottom: 16px;
    }
    .gallery-card {
      border: 1px solid var(--line);
      border-radius: 18px;
      padding: 12px;
      background: rgba(255, 255, 255, 0.76);
      backdrop-filter: blur(12px);
    }
    .gallery-card img {
      display: block;
      width: 100%;
      height: auto;
      border-radius: 14px;
      background: #f5efe3;
    }
    .gallery-card p {
      margin: 10px 0 0;
      color: var(--muted);
      font-size: 12px;
      word-break: break-word;
    }
    .token-inline-note {
      margin: 8px 0 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.5;
    }
    details {
      border: 1px solid var(--line);
      border-radius: 16px;
      background: rgba(255, 255, 255, 0.72);
      padding: 12px 14px;
    }
    details summary {
      cursor: pointer;
      font-weight: 700;
      color: var(--ink);
    }
    details > *:not(summary) {
      margin-top: 12px;
    }
    .mini-table {
      width: 100%;
      border-collapse: collapse;
      font-size: 13px;
    }
    .mini-table th, .mini-table td {
      padding: 10px 8px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      vertical-align: top;
    }
    .modal-backdrop {
      position: fixed;
      inset: 0;
      background:
        radial-gradient(circle at top, rgba(255, 255, 255, 0.12), transparent 30%),
        rgba(31, 41, 51, 0.42);
      backdrop-filter: blur(12px) saturate(1.08);
      display: none;
      align-items: center;
      justify-content: center;
      z-index: 50;
      padding: 20px;
    }
    .modal-backdrop.show {
      display: flex;
    }
    .modal-card {
      width: min(1120px, calc(100vw - 24px));
      max-height: calc(100vh - 24px);
      overflow: auto;
      background: rgba(255, 253, 249, 0.88);
      border: 1px solid var(--line);
      border-radius: 26px;
      box-shadow: var(--shadow-strong);
      padding: 24px 26px;
      backdrop-filter: blur(22px);
      transform: translateY(18px) scale(0.985);
      opacity: 0;
    }
    #db-modal .modal-card {
      width: min(1460px, calc(100vw - 24px));
    }
    .modal-backdrop.show .modal-card {
      animation: modal-in 0.24s ease forwards;
    }
    .modal-header {
      display: flex;
      justify-content: space-between;
      align-items: start;
      gap: 16px;
      margin-bottom: 16px;
    }
    .modal-header h2 {
      margin: 0 0 8px;
      font-size: 24px;
    }
    .modal-header p {
      margin: 0;
      color: var(--muted);
    }
    .modal-close {
      background: #fff;
      color: var(--ink);
      border: 1px solid var(--line);
      min-width: 44px;
      padding: 10px 12px;
    }
    .modal-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 14px;
    }
    .modal-section {
      margin-top: 18px;
    }
    .modal-section h3 {
      margin: 0 0 10px;
      font-size: 18px;
    }
    .inline-message {
      margin: 12px 0 0;
      font-size: 13px;
      color: var(--muted);
    }
    .danger-box {
      margin-top: 12px;
      padding: 14px;
      border-radius: 16px;
      background: rgba(180, 35, 24, 0.08);
      border: 1px solid rgba(180, 35, 24, 0.16);
    }
    @keyframes fade-up {
      from {
        opacity: 0;
        transform: translateY(10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
    @keyframes modal-in {
      from {
        opacity: 0;
        transform: translateY(18px) scale(0.985);
      }
      to {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
    }
    @keyframes pill-breathe {
      0%, 100% {
        transform: translateY(0);
      }
      50% {
        transform: translateY(-1px) scale(1.02);
      }
    }
    @keyframes pill-celebrate {
      0% {
        transform: scale(0.94);
      }
      35% {
        transform: scale(1.06);
      }
      100% {
        transform: scale(1);
      }
    }
    @keyframes log-card-enter {
      0% {
        opacity: 0;
        transform: translateY(-18px) scale(0.97);
      }
      55% {
        opacity: 1;
        transform: translateY(2px) scale(1.01);
      }
      100% {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
    }
    @keyframes log-card-glow {
      0% {
        box-shadow: 0 0 0 rgba(31, 111, 95, 0);
      }
      35% {
        box-shadow: 0 24px 50px rgba(31, 111, 95, 0.18);
      }
      100% {
        box-shadow: 0 8px 18px rgba(31, 41, 51, 0.05);
      }
    }
    @keyframes log-card-update {
      0% {
        transform: scale(0.985);
        box-shadow: 0 0 0 rgba(217, 144, 88, 0);
      }
      55% {
        transform: scale(1.008);
        box-shadow: 0 18px 36px rgba(217, 144, 88, 0.16);
      }
      100% {
        transform: scale(1);
        box-shadow: 0 8px 18px rgba(31, 41, 51, 0.05);
      }
    }
    @keyframes selected-glow {
      0% {
        box-shadow: 0 0 0 0 rgba(31, 111, 95, 0.22);
      }
      100% {
        box-shadow: 0 0 0 2px rgba(31, 111, 95, 0.12);
      }
    }
    @keyframes new-badge-pop {
      0% {
        opacity: 0;
        transform: translateY(-10px) scale(0.82);
      }
      65% {
        opacity: 1;
        transform: translateY(0) scale(1.06);
      }
      100% {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
    }
    @keyframes detail-fade {
      0% {
        opacity: 0;
        transform: translateY(10px);
        filter: blur(8px);
      }
      100% {
        opacity: 1;
        transform: translateY(0);
        filter: blur(0);
      }
    }
    @media (max-width: 1180px) {
      .monitor-layout {
        grid-template-columns: 1fr;
      }
      .log-list {
        max-height: none;
      }
      .conversation-grid {
        grid-template-columns: 1fr;
      }
    }
    @media (max-width: 900px) {
      .token-form {
        flex-direction: column;
        align-items: stretch;
      }
    }
    @media (max-width: 720px) {
      .shell { padding: 16px; }
      .shell-wide { padding: 14px; }
      .nav { padding: 14px 16px; }
      .panel { padding: 16px; }
      .hero h1 { font-size: 24px; }
      th, td { padding: 10px 8px; }
      .token-manual-grid {
        grid-template-columns: 1fr;
      }
      .log-pagination, .text-pagination {
        flex-direction: column;
        align-items: stretch;
      }
    }
    @media (prefers-reduced-motion: reduce) {
      *,
      *::before,
      *::after {
        animation: none !important;
        transition: none !important;
        scroll-behavior: auto !important;
      }
    }
  </style>
</head>
<body data-audit-base="{{basePath}}">
  {{template "body" .}}
</body>
</html>
{{end}}
`

const loginTemplate = `
{{define "login"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<div class="login-shell">
  <div class="login-card">
    <div class="hero">
      <div>
        <h1>newapi-audit-proxy</h1>
        <p>登录后查看审计统计、请求日志、令牌用量和流式响应详情。</p>
      </div>
    </div>
    {{if .Error}}
    <div class="error-box">{{.Error}}</div>
    {{end}}
    <form method="post" action="{{path "/login"}}">
      <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
      <div style="margin-bottom: 16px;">
        <label for="username">用户名</label>
        <input id="username" name="username" autocomplete="username" required>
      </div>
      <div style="margin-bottom: 20px;">
        <label for="password">密码</label>
        <input id="password" name="password" type="password" autocomplete="current-password" required>
      </div>
      <button type="submit">登录</button>
    </form>
  </div>
</div>
{{end}}
`

const dashboardTemplate = `
{{define "dashboard"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<div class="shell shell-wide">
  <div class="nav">
    <div class="brand">newapi-audit-proxy</div>
    <div class="nav-links">
      <a href="{{path "/"}}">首页</a>
      <a href="{{path "/logs"}}">日志</a>
      <form method="post" action="{{path "/logout"}}" style="margin: 0;">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button class="button secondary" type="submit">退出登录</button>
      </form>
    </div>
  </div>

  <div class="hero">
    <div>
      <h1>审计概览</h1>
      <p>首页展示请求数量、token 用量以及各令牌的聚合统计。</p>
    </div>
    <div class="pill">刷新时间：{{formatTime .Now}}</div>
  </div>

  <div class="grid">
    <div class="stat">
      <h3>总请求数</h3>
      <div class="value">{{.Stats.TotalRequests}}</div>
    </div>
    <div class="stat">
      <h3>今日请求数</h3>
      <div class="value">{{.Stats.TodayRequests}}</div>
    </div>
    <div class="stat">
      <h3>错误数</h3>
      <div class="value">{{.Stats.ErrorCount}}</div>
    </div>
    <div class="stat">
      <h3>总 Token 数量</h3>
      <div class="value">{{.Stats.DistinctTokens}}</div>
    </div>
    <div class="stat">
      <h3>今日活跃 Token</h3>
      <div class="value">{{.Stats.TodayDistinctTokens}}</div>
    </div>
    <div class="stat">
      <h3>总 Token 数</h3>
      <div class="value">{{.Stats.TotalTokens}}</div>
    </div>
    <div class="stat">
      <h3>今日 Token 数</h3>
      <div class="value">{{.Stats.TodayTotalTokens}}</div>
    </div>
    <div class="stat">
      <h3>输入 / 输出 Tokens</h3>
      <div class="value" style="font-size: 22px;">{{.Stats.TotalPromptTokens}} / {{.Stats.TotalCompletionTokens}}</div>
    </div>
  </div>

  <div class="panel">
    <div class="hero">
      <div>
        <h1 style="font-size: 22px;">按令牌统计</h1>
        <p>按 token 指纹聚合请求次数、错误次数和用量统计。</p>
      </div>
    </div>
    <table>
      <thead>
        <tr>
          <th>令牌代号</th>
          <th>Token 指纹</th>
          <th>Token 预览</th>
          <th>请求数</th>
          <th>输入 Tokens</th>
          <th>输出 Tokens</th>
          <th>总 Tokens</th>
          <th>错误数</th>
          <th>最后时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        {{range .Stats.TokenGroups}}
        <tr>
          <td>{{if .TokenAlias}}{{.TokenAlias}}{{else}}-{{end}}</td>
          <td><code>{{if .TokenFingerprint}}{{.TokenFingerprint}}{{else}}（无）{{end}}</code></td>
          <td>{{if .TokenPreview}}{{.TokenPreview}}{{else}}-{{end}}</td>
          <td>{{.RequestCount}}</td>
          <td>{{.PromptTokens}}</td>
          <td>{{.CompletionTokens}}</td>
          <td>{{.TotalTokens}}</td>
          <td>{{.ErrorCount}}</td>
          <td>{{formatTime .LastSeen}}</td>
          <td><a href="{{path "/tokens"}}?token={{.TokenFingerprint}}">管理别名</a></td>
        </tr>
        {{else}}
        <tr><td colspan="10" class="muted">暂无数据</td></tr>
        {{end}}
      </tbody>
    </table>
  </div>
</div>
{{end}}
`

const logsTemplate = `
{{define "logs"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<div class="shell shell-wide">
  <div class="nav">
    <div class="brand">newapi-audit-proxy</div>
    <div class="toolbar-actions">
      <a class="button secondary" href="{{path "/logs"}}">日志</a>
      <button class="button secondary" type="button" id="open-stats-modal">统计</button>
      <button class="button secondary" type="button" id="open-settings-modal">设置</button>
      <button class="button secondary" type="button" id="open-db-modal">数据库管理</button>
      <form method="post" action="{{path "/logout"}}" style="margin: 0;">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button class="button secondary" type="submit">退出登录</button>
      </form>
    </div>
  </div>

  <div class="panel">
    <div class="hero">
      <div>
        <h1 style="font-size: 24px;">实时日志监控</h1>
        <p>每 5 秒检查一次是否有新日志，只在数据发生变化时刷新左侧列表。文本类筛选默认都是模糊匹配。</p>
      </div>
      <div class="pill" id="logs-refresh-state">监控启动中</div>
    </div>

    <form method="get" action="{{path "/logs"}}">
      <div class="form-grid">
        <div>
          <label for="from">开始时间</label>
          <input id="from" name="from" type="datetime-local" value="{{.Filters.From}}">
        </div>
        <div>
          <label for="to">结束时间</label>
          <input id="to" name="to" type="datetime-local" value="{{.Filters.To}}">
        </div>
        <div>
          <label for="token">Token 指纹（支持部分匹配）</label>
          <input id="token" name="token" value="{{.Filters.TokenFingerprint}}">
        </div>
        <div>
          <label for="alias">令牌代号</label>
          <input id="alias" name="alias" value="{{.Filters.TokenAlias}}">
        </div>
        <div>
          <label for="model">模型</label>
          <input id="model" name="model" value="{{.Filters.Model}}">
        </div>
        <div>
          <label for="status">状态码</label>
          <input id="status" name="status" value="{{.Filters.StatusCode}}">
        </div>
        <div>
          <label for="q">关键词</label>
          <input id="q" name="q" value="{{.Filters.Keyword}}">
        </div>
        <div>
          <button type="submit">搜索</button>
        </div>
      </div>
    </form>
  </div>

  <div class="monitor-layout">
    <div class="panel log-sidebar">
      <div class="sidebar-header">
        <div>
          <h2 style="margin: 0 0 6px;">日志列表</h2>
          <p class="muted" id="logs-count-text" style="margin: 0;">正在加载日志...</p>
        </div>
      </div>
      <div id="logs-list" class="log-list">
        <div class="empty-list">正在加载日志...</div>
      </div>
      <div class="log-pagination">
        <button class="button secondary" type="button" id="logs-prev-page" disabled>上一页</button>
        <span class="muted" id="logs-page-info">每页 100 条</span>
        <button class="button secondary" type="button" id="logs-next-page" disabled>下一页</button>
      </div>
    </div>

    <div class="panel log-main">
      <div id="log-detail-empty" class="detail-empty">
        <div>
          <div class="section-kicker">详情视图</div>
          <h2 style="margin: 0 0 10px;">左侧点选一条日志</h2>
          <p style="margin: 0;">右侧会优先展示用户发送内容和模型回复，其他元数据收在下方。</p>
        </div>
      </div>

      <div id="log-detail" class="detail-hidden">
        <div class="detail-hero">
          <div>
            <div class="pill-row">
              <span class="pill" id="detail-token-alias">未命名</span>
              <span class="pill" id="detail-model">-</span>
              <span class="pill" id="detail-status">-</span>
              <span class="pill" id="detail-response-type">-</span>
            </div>
            <h1 id="detail-time" style="margin: 0 0 8px; font-size: 28px;">-</h1>
            <p id="detail-path" class="muted" style="margin: 0;">-</p>
          </div>
          <div class="muted">
            <a id="detail-standalone-link" href="{{path "/logs/0"}}">单独打开详情</a>
          </div>
        </div>

        <div id="detail-images-wrap" class="detail-hidden">
          <div class="section-kicker">绘图结果</div>
          <div id="detail-images" class="gallery"></div>
        </div>

        <div class="conversation-grid">
          <div class="conversation-card user">
            <div class="section-kicker">用户发送</div>
            <pre id="detail-user-text" class="paged-pre">（空）</pre>
            <div id="detail-user-pager" class="text-pagination detail-hidden">
              <button class="button secondary" type="button" id="detail-user-prev">上一页</button>
              <span class="muted" id="detail-user-page-info">第 1 / 1 页</span>
              <button class="button secondary" type="button" id="detail-user-next">下一页</button>
            </div>
          </div>
          <div class="conversation-card assistant">
            <div class="section-kicker">模型回复</div>
            <pre id="detail-assistant-text" class="paged-pre">（空）</pre>
            <div id="detail-assistant-pager" class="text-pagination detail-hidden">
              <button class="button secondary" type="button" id="detail-assistant-prev">上一页</button>
              <span class="muted" id="detail-assistant-page-info">第 1 / 1 页</span>
              <button class="button secondary" type="button" id="detail-assistant-next">下一页</button>
            </div>
          </div>
        </div>

        <div class="meta-strip">
          <div class="meta-chip">
            <span>总 Tokens</span>
            <strong id="detail-total-tokens">0</strong>
          </div>
          <div class="meta-chip">
            <span>输入 Tokens</span>
            <strong id="detail-prompt-tokens">0</strong>
          </div>
          <div class="meta-chip">
            <span>输出 Tokens</span>
            <strong id="detail-completion-tokens">0</strong>
          </div>
          <div class="meta-chip">
            <span>请求耗时</span>
            <strong id="detail-duration">0 毫秒</strong>
          </div>
          <div class="meta-chip">
            <span>Token 预览</span>
            <strong id="detail-token-preview">-</strong>
          </div>
          <div class="meta-chip">
            <span>Token 指纹</span>
            <strong id="detail-token-fingerprint">-</strong>
          </div>
        </div>

        <div class="details-stack">
          <details>
            <summary>隐藏信息</summary>
            <p><strong>错误：</strong> <span id="detail-error-text">-</span></p>
            <p><strong>流式：</strong> <span id="detail-stream">-</span></p>
            <p><strong>请求字节：</strong> <span id="detail-request-bytes">0</span></p>
            <p><strong>响应字节：</strong> <span id="detail-response-bytes">0</span></p>
          </details>

          <details id="detail-raw-section">
            <summary>原始请求与回复</summary>
            <div class="details-stack">
              <div>
                <h3 style="margin-top: 0;">原始请求体</h3>
                <pre id="detail-request-body" class="mono-scroll">（空）</pre>
              </div>
              <div>
                <h3 style="margin-top: 0;">原始响应体</h3>
                <pre id="detail-response-body" class="mono-scroll">（空）</pre>
              </div>
            </div>
          </details>

          <details id="detail-json-section">
            <summary>JSON / 请求头 / 用量</summary>
            <div class="details-stack">
              <div>
                <h3 style="margin-top: 0;">请求 JSON</h3>
                <pre id="detail-request-json" class="mono-scroll">（空）</pre>
              </div>
              <div>
                <h3 style="margin-top: 0;">响应 JSON</h3>
                <pre id="detail-response-json" class="mono-scroll">（空）</pre>
              </div>
              <div>
                <h3 style="margin-top: 0;">用量 JSON</h3>
                <pre id="detail-usage-json" class="mono-scroll">（空）</pre>
              </div>
              <div>
                <h3 style="margin-top: 0;">请求头</h3>
                <pre id="detail-request-headers" class="mono-scroll">（空）</pre>
              </div>
              <div>
                <h3 style="margin-top: 0;">响应头</h3>
                <pre id="detail-response-headers" class="mono-scroll">（空）</pre>
              </div>
            </div>
          </details>
        </div>
      </div>
    </div>
  </div>
</div>

<div id="stats-modal" class="modal-backdrop">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>统计</h2>
        <p>展示当前审计主视图的聚合统计和 token 概览。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="stats-modal">关闭</button>
    </div>

    <div class="modal-grid">
      <div class="stat"><h3>总请求数</h3><div class="value" id="stats-total-requests">0</div></div>
      <div class="stat"><h3>今日请求数</h3><div class="value" id="stats-today-requests">0</div></div>
      <div class="stat"><h3>错误数</h3><div class="value" id="stats-error-count">0</div></div>
      <div class="stat"><h3>总 Token 数</h3><div class="value" id="stats-total-tokens">0</div></div>
      <div class="stat"><h3>输入 Tokens</h3><div class="value" id="stats-prompt-tokens">0</div></div>
      <div class="stat"><h3>输出 Tokens</h3><div class="value" id="stats-completion-tokens">0</div></div>
    </div>

    <div class="modal-section">
      <h3>高用量令牌</h3>
      <table class="mini-table">
        <thead>
          <tr>
            <th>令牌代号</th>
            <th>Token 预览</th>
            <th>请求数</th>
            <th>总 Tokens</th>
            <th>错误数</th>
            <th>最后时间</th>
          </tr>
        </thead>
        <tbody id="stats-token-table">
          <tr><td colspan="6" class="muted">加载中...</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</div>

<div id="settings-modal" class="modal-backdrop">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>设置</h2>
        <p>在此管理 token 代号映射，手动绑定或直接修改已出现的 token 。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="settings-modal">关闭</button>
    </div>

    <div class="modal-section">
      <h3>手动绑定 Token 代号</h3>
      <form id="settings-alias-form">
        <div class="token-manual-grid">
          <div>
            <label for="settings-token-value">Token 原文</label>
            <input id="settings-token-value" name="token_value" placeholder="sk-... / Bearer sk-...">
          </div>
          <div>
            <label for="settings-token-alias">令牌代号</label>
            <input id="settings-token-alias" name="token_alias" placeholder="例如：绘图A">
          </div>
          <div>
            <button type="submit">保存映射</button>
          </div>
        </div>
      </form>
      <p id="settings-message" class="inline-message">系统只保存 token 指纹和代号对应，不保存 token 原文。</p>
    </div>

    <div class="modal-section">
      <h3>令牌列表</h3>
      <form id="settings-token-filter-form">
        <div class="form-grid">
          <div>
            <label for="settings-filter-from">开始时间</label>
            <input id="settings-filter-from" name="from" type="datetime-local">
          </div>
          <div>
            <label for="settings-filter-to">结束时间</label>
            <input id="settings-filter-to" name="to" type="datetime-local">
          </div>
          <div>
            <label for="settings-filter-token">Token 指纹</label>
            <input id="settings-filter-token" name="token">
          </div>
          <div>
            <label for="settings-filter-alias">令牌代号</label>
            <input id="settings-filter-alias" name="alias">
          </div>
          <div>
            <label for="settings-filter-model">模型</label>
            <input id="settings-filter-model" name="model">
          </div>
          <div>
            <button type="submit">刷新 Token 列表</button>
          </div>
        </div>
      </form>

      <table class="mini-table" style="margin-top: 14px;">
        <thead>
          <tr>
            <th>令牌代号</th>
            <th>Token 预览</th>
            <th>Token 指纹</th>
            <th>请求数</th>
            <th>总 Tokens</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody id="settings-token-table">
          <tr><td colspan="6" class="muted">加载中...</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</div>

<div id="db-modal" class="modal-backdrop">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>数据库管理</h2>
        <p>查看当前入库量、今日入库和占用空间，并按条件清理审计日志。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="db-modal">关闭</button>
    </div>

    <div class="modal-grid">
      <div class="stat"><h3>审计日志条数</h3><div class="value" id="db-total-rows">0</div></div>
      <div class="stat"><h3>今日入库</h3><div class="value" id="db-today-rows">0</div></div>
      <div class="stat"><h3>audit_logs 总占用</h3><div class="value" id="db-audit-total-pretty">0</div></div>
      <div class="stat"><h3>数据主体</h3><div class="value" id="db-audit-table-pretty">0</div></div>
      <div class="stat"><h3>索引占用</h3><div class="value" id="db-audit-index-pretty">0</div></div>
      <div class="stat"><h3>TOAST 大字段</h3><div class="value" id="db-audit-toast-pretty">0</div></div>
      <div class="stat"><h3>死元组</h3><div class="value" id="db-dead-tuples">0</div></div>
      <div class="stat"><h3>存活元组</h3><div class="value" id="db-live-tuples">0</div></div>
      <div class="stat"><h3>整个数据库</h3><div class="value" id="db-size-pretty">0</div></div>
    </div>

    <p class="inline-message">
      说明：删除日志后，PostgreSQL 通常只会把空间标记为可复用，不会立刻把 audit_logs 文件缩回给操作系统。
      所以“条数少了但硬盘没少”通常是正常现象，需要再做空间维护。
    </p>

    <div class="modal-section">
      <h3>空间维护</h3>
      <div class="danger-box">
        <div class="toolbar-actions" style="margin-bottom: 12px;">
          <button class="button secondary" type="button" id="db-refresh-button">刷新统计</button>
          <button class="button secondary" type="button" data-db-maintenance="compact_payloads">历史记录瘦身</button>
          <button class="button secondary" type="button" data-db-maintenance="vacuum_analyze">整理空间</button>
          <button class="danger" type="button" data-db-maintenance="vacuum_full">强制缩盘</button>
        </div>
        <div class="modal-grid" style="margin-bottom: 12px;">
          <div class="meta-chip"><span>上次 VACUUM</span><strong id="db-last-vacuum">-</strong></div>
          <div class="meta-chip"><span>上次 AUTOVACUUM</span><strong id="db-last-autovacuum">-</strong></div>
          <div class="meta-chip"><span>上次 ANALYZE</span><strong id="db-last-analyze">-</strong></div>
          <div class="meta-chip"><span>上次 AUTOANALYZE</span><strong id="db-last-autoanalyze">-</strong></div>
        </div>
        <p class="inline-message">
          “历史记录瘦身”会按当前脱敏规则重写旧日志，把历史遗留的图片 base64 从已入库内容里清掉。
          做完这一步之后，再执行一次 <code>VACUUM FULL</code> 才能真正把这部分空间缩回去。
        </p>
        <p class="inline-message">
          “整理空间”执行 <code>VACUUM ANALYZE</code>，会整理可复用空间并刷新统计，但磁盘文件未必立刻变小。
          “强制缩盘”执行 <code>VACUUM FULL</code>，会锁住 audit_logs 表，但能尽量把空间还给操作系统。
        </p>
        <p id="db-maintenance-message" class="inline-message">建议先删除不需要的日志，再执行“整理空间”；如果明确要立刻缩小硬盘占用，再执行“强制缩盘”。</p>
      </div>
    </div>

    <div class="modal-section">
      <h3>按条件清理日志</h3>
      <div class="danger-box">
        <form id="db-cleanup-form">
          <div class="form-grid">
            <div>
              <label for="db-filter-from">开始时间</label>
              <input id="db-filter-from" name="from" type="datetime-local">
            </div>
            <div>
              <label for="db-filter-to">结束时间</label>
              <input id="db-filter-to" name="to" type="datetime-local">
            </div>
            <div>
              <label for="db-filter-token">Token 指纹</label>
              <input id="db-filter-token" name="token">
            </div>
            <div>
              <label for="db-filter-alias">令牌代号</label>
              <input id="db-filter-alias" name="alias">
            </div>
            <div>
              <label for="db-filter-model">模型</label>
              <input id="db-filter-model" name="model">
            </div>
            <div>
              <button class="danger" type="submit">清理符合条件的记录</button>
            </div>
          </div>
        </form>
        <p id="db-cleanup-message" class="inline-message">安全提示：至少需要输入一个条件才会执行清理。删除记录后如果想立刻缩小硬盘占用，请继续执行上方的“强制缩盘”。</p>
      </div>
    </div>
  </div>
</div>

<script>
(function () {
  var state = {
    version: '',
    selectedId: 0,
    detailId: 0,
    page: 1,
    listItems: Object.create(null),
    hasLoadedList: false
  };
  var textPagerState = {
    user: { kind: 'user', text: '（空）', page: 1, page_size: 4000, total_pages: 1, total_chars: 0 },
    assistant: { kind: 'assistant', text: '（空）', page: 1, page_size: 4000, total_pages: 1, total_chars: 0 }
  };
  var detailCache = Object.create(null);
  var rawDetailCache = Object.create(null);
  var textPageCache = Object.create(null);
  var pendingDetailController = null;
  var pendingRawController = null;
  var pendingTextControllers = {
    user: null,
    assistant: null
  };

  var refreshLabel = document.getElementById('logs-refresh-state');
  var countLabel = document.getElementById('logs-count-text');
  var listEl = document.getElementById('logs-list');
  var sidebarEl = document.querySelector('.log-sidebar');
  var detailEmptyEl = document.getElementById('log-detail-empty');
  var detailEl = document.getElementById('log-detail');
  var detailImagesWrapEl = document.getElementById('detail-images-wrap');
  var detailImagesEl = document.getElementById('detail-images');
  var detailRawSectionEl = document.getElementById('detail-raw-section');
  var detailJSONSectionEl = document.getElementById('detail-json-section');
  var auditBase = document.body.getAttribute('data-audit-base') || '';
  var searchParams = new URLSearchParams(window.location.search);
  state.selectedId = parseInt(searchParams.get('selected') || '0', 10) || 0;
  state.page = parseInt(searchParams.get('page') || '1', 10) || 1;
  if (state.page < 1) {
    state.page = 1;
  }

  function playTransientClass(el, className, duration) {
    if (!el || !className) {
      return;
    }
    el.classList.remove(className);
    void el.offsetWidth;
    el.classList.add(className);
    window.setTimeout(function () {
      el.classList.remove(className);
    }, duration || 900);
  }

  function setSidebarRefreshing(active) {
    if (sidebarEl) {
      sidebarEl.classList.toggle('is-refreshing', !!active);
    }
  }

  function setRefreshLabel(text, tone) {
    if (refreshLabel) {
      refreshLabel.textContent = text;
      refreshLabel.classList.remove('is-checking', 'has-updates', 'is-error');
      if (tone === 'checking') {
        refreshLabel.classList.add('is-checking');
      } else if (tone === 'updates') {
        refreshLabel.classList.add('has-updates');
      } else if (tone === 'error') {
        refreshLabel.classList.add('is-error');
      }
    }
  }

  function escapeHTML(value) {
    return String(value || '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;');
  }

  function paramsWithoutSelected() {
    var values = new URLSearchParams(window.location.search);
    values.delete('selected');
    return values.toString();
  }

  function withBase(path) {
    return auditBase ? auditBase + path : path;
  }

  function apiURL(path) {
    var query = paramsWithoutSelected();
    var fullPath = withBase(path);
    return query ? fullPath + '?' + query : fullPath;
  }

  function updateSelectedInURL(id) {
    var url = new URL(window.location.href);
    if (id > 0) {
      url.searchParams.set('selected', String(id));
    } else {
      url.searchParams.delete('selected');
    }
    var query = url.searchParams.toString();
    history.replaceState(null, '', url.pathname + (query ? '?' + query : ''));
  }

  function updatePageInURL(page) {
    var url = new URL(window.location.href);
    if (page > 1) {
      url.searchParams.set('page', String(page));
    } else {
      url.searchParams.delete('page');
    }
    url.searchParams.delete('selected');
    var query = url.searchParams.toString();
    history.replaceState(null, '', url.pathname + (query ? '?' + query : ''));
  }

  function showEmpty(message) {
    detailEmptyEl.classList.remove('detail-hidden');
    detailEl.classList.add('detail-hidden');
    detailEmptyEl.innerHTML = '<div><div class="section-kicker">详情视图</div><h2 style="margin: 0 0 10px;">' + escapeHTML(message) + '</h2><p style="margin: 0;">请在左侧选择日志，右侧会自动展示内容。</p></div>';
  }

  function showDetail() {
    detailEmptyEl.classList.add('detail-hidden');
    detailEl.classList.remove('detail-hidden');
    playTransientClass(detailEl, 'detail-animate', 460);
  }

  function csrfToken() {
    var meta = document.querySelector('meta[name="csrf-token"]');
    return meta ? meta.content : '';
  }

  function fetchJSON(url, options) {
    options = options || {};
    return fetch(url, {
      headers: {
        'Accept': 'application/json'
      },
      credentials: 'same-origin',
      signal: options.signal
    }).then(function (resp) {
      if (!resp.ok) {
        throw new Error('HTTP ' + resp.status);
      }
      return resp.json();
    });
  }

  function postForm(url, formData) {
    var body = new URLSearchParams();
    formData.forEach(function (value, key) {
      body.append(key, value);
    });
    return fetch(url, {
      method: 'POST',
      body: body.toString(),
      credentials: 'same-origin',
      headers: {
        'Accept': 'application/json',
        'X-CSRF-Token': csrfToken(),
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      }
    }).then(function (resp) {
      return resp.text().then(function (text) {
        var data = {};
        if (text) {
          try {
            data = JSON.parse(text);
          } catch (err) {
            data = { message: text };
          }
        }
        if (!resp.ok) {
          throw new Error(data.message || text || ('HTTP ' + resp.status));
        }
        return data;
      });
    });
  }

  function isAbortError(error) {
    if (!error) {
      return false;
    }
    if (error.name === 'AbortError') {
      return true;
    }
    return /abort/i.test(String(error.message || error));
  }

  function statusText(code) {
    if (!code) {
      return '-';
    }
    return String(code);
  }

  function setText(id, value) {
    var el = document.getElementById(id);
    if (el) {
      el.textContent = value;
    }
  }

  function normalizeTextPage(name, payload) {
    payload = payload || {};
    return {
      kind: payload.kind || name,
      text: payload.text || '（空）',
      page: parseInt(payload.page || 1, 10) || 1,
      page_size: parseInt(payload.page_size || 4000, 10) || 4000,
      total_pages: parseInt(payload.total_pages || 1, 10) || 1,
      total_chars: parseInt(payload.total_chars || 0, 10) || 0
    };
  }

  function textPageCacheKey(id, name, page) {
    return String(id) + ':' + name + ':' + String(page);
  }

  function cacheTextPage(id, payload) {
    if (!id || !payload) {
      return;
    }
    var normalized = normalizeTextPage(payload.kind || 'user', payload);
    textPageCache[textPageCacheKey(id, normalized.kind, normalized.page)] = normalized;
  }

  function renderTextPager(name) {
    var config = {
      user: {
        textId: 'detail-user-text',
        pagerId: 'detail-user-pager',
        prevId: 'detail-user-prev',
        nextId: 'detail-user-next',
        infoId: 'detail-user-page-info'
      },
      assistant: {
        textId: 'detail-assistant-text',
        pagerId: 'detail-assistant-pager',
        prevId: 'detail-assistant-prev',
        nextId: 'detail-assistant-next',
        infoId: 'detail-assistant-page-info'
      }
    }[name];
    if (!config) {
      return;
    }

    var pager = textPagerState[name];
    if (!pager) {
      return;
    }

    var page = parseInt(pager.page || 1, 10) || 1;
    if (page < 1) {
      page = 1;
    }
    if (page > pager.total_pages) {
      page = pager.total_pages;
    }
    pager.page = page;

    setText(config.textId, pager.text || '（空）');

    var pageInfo = '第 ' + page + ' / ' + pager.total_pages + ' 页';
    if (pager.total_chars > 0) {
      pageInfo += ' · ' + pager.total_chars + ' 字';
    }
    setText(config.infoId, pageInfo);

    var pagerEl = document.getElementById(config.pagerId);
    var prevEl = document.getElementById(config.prevId);
    var nextEl = document.getElementById(config.nextId);
    if (pagerEl) {
      pagerEl.classList.toggle('detail-hidden', pager.total_pages <= 1);
    }
    if (prevEl) {
      prevEl.disabled = page <= 1 || !!pendingTextControllers[name];
    }
    if (nextEl) {
      nextEl.disabled = page >= pager.total_pages || !!pendingTextControllers[name];
    }
  }

  function applyTextPage(name, payload) {
    textPagerState[name] = normalizeTextPage(name, payload);
    renderTextPager(name);
  }

  function setPagedText(name, text) {
    applyTextPage(name, {
      kind: name,
      text: text || '（空）',
      page: 1,
      total_pages: 1,
      total_chars: Array.from(String(text || '')).length
    });
  }

  function loadTextPage(name, page) {
    if (!state.selectedId) {
      return Promise.resolve(null);
    }

    var current = textPagerState[name];
    if (!current) {
      return Promise.resolve(null);
    }

    var nextPage = parseInt(page || 1, 10) || 1;
    if (nextPage < 1) {
      nextPage = 1;
    }
    if (nextPage > current.total_pages) {
      nextPage = current.total_pages;
    }

    var id = state.selectedId;
    var cacheKey = textPageCacheKey(id, name, nextPage);
    if (textPageCache[cacheKey]) {
      applyTextPage(name, textPageCache[cacheKey]);
      return Promise.resolve(textPageCache[cacheKey]);
    }

    if (pendingTextControllers[name]) {
      pendingTextControllers[name].abort();
    }

    pendingTextControllers[name] = new AbortController();
    applyTextPage(name, {
      kind: name,
      text: '正在加载...',
      page: nextPage,
      page_size: current.page_size,
      total_pages: current.total_pages,
      total_chars: current.total_chars
    });

    return fetchJSON(withBase('/api/logs/' + id + '/text?kind=' + encodeURIComponent(name) + '&page=' + nextPage), {
      signal: pendingTextControllers[name].signal
    }).then(function (detail) {
      pendingTextControllers[name] = null;
      if (state.selectedId !== id) {
        return null;
      }
      cacheTextPage(id, detail);
      applyTextPage(name, detail);
      return detail;
    }).catch(function (error) {
      pendingTextControllers[name] = null;
      if (isAbortError(error)) {
        return null;
      }
      applyTextPage(name, {
        kind: name,
        text: '文本分页加载失败',
        page: nextPage,
        page_size: current.page_size,
        total_pages: current.total_pages,
        total_chars: current.total_chars
      });
      return null;
    });
  }

  function changeTextPage(name, delta) {
    var pager = textPagerState[name];
    if (!pager) {
      return;
    }
    var nextPage = pager.page + delta;
    if (nextPage < 1 || nextPage > pager.total_pages) {
      return;
    }
    loadTextPage(name, nextPage);
  }

  function setRawPlaceholders(message) {
    var placeholder = message || '正在后台加载原始详情...';
    setText('detail-request-body', placeholder);
    setText('detail-response-body', placeholder);
    setText('detail-request-json', placeholder);
    setText('detail-response-json', placeholder);
    setText('detail-usage-json', placeholder);
    setText('detail-request-headers', placeholder);
    setText('detail-response-headers', placeholder);
    renderImages([]);
  }

  function renderQuickDetail(item) {
    showDetail();

    setText('detail-token-alias', item && item.token_alias ? item.token_alias : '加载中');
    setText('detail-model', item && item.model ? item.model : '-');
    setText('detail-status', 'HTTP ' + statusText(item ? item.status_code : 0));
    setText('detail-response-type', '加载中');
    setText('detail-time', item && item.started_at ? item.started_at : '-');
    setText('detail-path', '正在加载详细路径...');
    setPagedText('user', item && item.user_preview ? item.user_preview : '加载中...');
    setPagedText('assistant', item && item.assistant_preview ? item.assistant_preview : '加载中...');
    setText('detail-total-tokens', item ? String(item.total_tokens || 0) : '0');
    setText('detail-prompt-tokens', '加载中');
    setText('detail-completion-tokens', '加载中');
    setText('detail-duration', '加载中');
    setText('detail-token-preview', item && item.token_preview ? item.token_preview : '-');
    setText('detail-token-fingerprint', item && item.token_fingerprint ? item.token_fingerprint : '-');
    setText('detail-error-text', '-');
    setText('detail-stream', '-');
    setText('detail-request-bytes', '加载中');
    setText('detail-response-bytes', '加载中');
    setRawPlaceholders('正在后台载入原始请求/响应...');

    var standalone = document.getElementById('detail-standalone-link');
    if (standalone) {
      standalone.href = withBase('/logs/' + String(item && item.id ? item.id : 0));
    }
  }

  function applyDetailCore(detail) {
    showDetail();
    setText('detail-token-alias', detail.token_alias || '未命名');
    setText('detail-model', detail.model || '-');
    setText('detail-status', 'HTTP ' + statusText(detail.status_code));
    setText('detail-response-type', detail.response_type || '-');
    setText('detail-time', detail.started_at || '-');
    setText('detail-path', (detail.method || '-') + ' ' + (detail.path_with_query || '-'));
    cacheTextPage(detail.id, detail.user_text);
    cacheTextPage(detail.id, detail.assistant_text);
    applyTextPage('user', detail.user_text);
    applyTextPage('assistant', detail.assistant_text);
    setText('detail-total-tokens', String(detail.total_tokens || 0));
    setText('detail-prompt-tokens', String(detail.prompt_tokens || 0));
    setText('detail-completion-tokens', String(detail.completion_tokens || 0));
    setText('detail-duration', String(detail.duration_ms || 0) + ' 毫秒');
    setText('detail-token-preview', detail.token_preview || '-');
    setText('detail-token-fingerprint', detail.token_fingerprint || '-');
    setText('detail-error-text', detail.error_text || '-');
    setText('detail-stream', detail.stream || '-');

    var requestBytes = String(detail.request_bytes || 0);
    if (detail.request_truncated) {
      requestBytes += ' (已截断)';
    }
    var responseBytes = String(detail.response_bytes || 0);
    if (detail.response_truncated) {
      responseBytes += ' (已截断)';
    }

    setText('detail-request-bytes', requestBytes);
    setText('detail-response-bytes', responseBytes);

    var standalone = document.getElementById('detail-standalone-link');
    if (standalone) {
      standalone.href = withBase('/logs/' + String(detail.id || 0));
    }
  }

  function shouldLoadRawNow() {
    return !!((detailRawSectionEl && detailRawSectionEl.open) || (detailJSONSectionEl && detailJSONSectionEl.open));
  }

  function shouldAutoLoadRaw(detail) {
    if (shouldLoadRawNow()) {
      return true;
    }
    var model = String((detail && detail.model) || '').toLowerCase();
    return model.indexOf('image') >= 0;
  }

  function ensureRawDetailLoaded() {
    if (!state.selectedId) {
      return Promise.resolve(null);
    }
    if (rawDetailCache[state.selectedId]) {
      applyDetailRaw(rawDetailCache[state.selectedId]);
      return Promise.resolve(rawDetailCache[state.selectedId]);
    }
    if (pendingRawController) {
      return Promise.resolve(null);
    }
    setRefreshLabel('正在加载原始详情');
    return fetchRawDetail(state.selectedId);
  }

  function applyDetailRaw(detail) {
    setText('detail-request-body', detail.request_body || '（空）');
    setText('detail-response-body', detail.response_body || '（空）');
    setText('detail-request-json', detail.request_json || '（空）');
    setText('detail-response-json', detail.response_json || '（空）');
    setText('detail-usage-json', detail.usage_json || '（空）');
    setText('detail-request-headers', detail.request_headers || '（空）');
    setText('detail-response-headers', detail.response_headers || '（空）');
    renderImages(detail.images || []);
  }

  function fetchRawDetail(id) {
    if (!id) {
      return Promise.resolve();
    }
    if (rawDetailCache[id]) {
      if (state.selectedId === id) {
        applyDetailRaw(rawDetailCache[id]);
      }
      return Promise.resolve(rawDetailCache[id]);
    }

    if (pendingRawController) {
      pendingRawController.abort();
    }
    pendingRawController = new AbortController();

    return fetchJSON(withBase('/api/logs/' + id + '/raw'), { signal: pendingRawController.signal }).then(function (detail) {
      rawDetailCache[id] = detail;
      if (state.selectedId === id) {
        applyDetailRaw(detail);
        setRefreshLabel('已同步至最新');
      }
      pendingRawController = null;
      return detail;
    }).catch(function (error) {
      pendingRawController = null;
      if (isAbortError(error)) {
        return null;
      }
      if (state.selectedId === id) {
        setRawPlaceholders('原始详情加载失败');
        setRefreshLabel('原始详情加载失败');
      }
      return null;
    });
  }

  function renderImages(images) {
    if (!detailImagesWrapEl || !detailImagesEl) {
      return;
    }
    if (!images || !images.length) {
      detailImagesWrapEl.classList.add('detail-hidden');
      detailImagesEl.innerHTML = '';
      return;
    }

    detailImagesWrapEl.classList.remove('detail-hidden');
    detailImagesEl.innerHTML = images.map(function (image) {
      return '' +
        '<div class="gallery-card">' +
          '<img src="' + escapeHTML(image.data_url) + '" alt="' + escapeHTML(image.label || '解析图片') + '">' +
          '<p>' + escapeHTML(image.label || '解析图片') + '</p>' +
        '</div>';
    }).join('');
  }

  function renderListPagination(data) {
    state.page = parseInt(data.page || state.page || 1, 10) || 1;

    var prevButton = document.getElementById('logs-prev-page');
    var nextButton = document.getElementById('logs-next-page');
    var pageInfo = document.getElementById('logs-page-info');
    var totalPages = parseInt(data.total_pages || 1, 10) || 1;

    if (pageInfo) {
      pageInfo.textContent = '第 ' + state.page + ' / ' + totalPages + ' 页，每页 100 条';
    }
    if (prevButton) {
      prevButton.disabled = !data.has_prev;
    }
    if (nextButton) {
      nextButton.disabled = !data.has_next;
    }
  }

  function renderList(items, totalCount) {
    var summary = {
      newCount: 0,
      updatedCount: 0
    };

    if (!items.length) {
      listEl.innerHTML = '<div class="empty-list">没有匹配的日志</div>';
      countLabel.textContent = '匹配结果 0 条';
      state.listItems = Object.create(null);
      state.hasLoadedList = true;
      listEl.classList.remove('list-surge');
      return summary;
    }

    var shouldAnimateChanges = state.hasLoadedList;
    state.listItems = Object.create(null);
    countLabel.textContent = '匹配结果 ' + totalCount + ' 条';
    Array.prototype.forEach.call(listEl.querySelectorAll('.empty-list'), function (node) {
      node.remove();
    });
    var existingNodes = Object.create(null);
    Array.prototype.forEach.call(listEl.querySelectorAll('[data-log-id]'), function (node) {
      existingNodes[node.getAttribute('data-log-id')] = node;
    });

    var expectedIDs = Object.create(null);
    items.forEach(function (item) {
      var idKey = String(item.id);
      state.listItems[idKey] = item;
      expectedIDs[idKey] = true;

      var button = existingNodes[idKey];
      var isNewNode = !button;
      if (!button) {
        button = document.createElement('button');
        button.type = 'button';
        button.setAttribute('data-log-id', idKey);
      }
      button.onclick = function () {
        loadDetail(item.id, true);
      };

      var alias = item.token_alias || '未命名';
      var bodyPreview = item.user_preview || item.assistant_preview || '（空）';
      var renderKey = [
        item.id,
        item.started_at || '',
        alias,
        item.model || '',
        bodyPreview,
        item.status_code || 0,
        item.total_tokens || 0
      ].join('||');
      var previousRenderKey = button.getAttribute('data-render-key') || '';

      button.className = 'log-item' + (item.id === state.selectedId ? ' active' : '');
      if (previousRenderKey !== renderKey) {
        button.setAttribute('data-render-key', renderKey);
        button.innerHTML =
          '<div class="log-item-top">' +
            '<div>' +
              '<div class="log-item-alias">' + escapeHTML(alias) + '</div>' +
              '<div class="log-item-model">' + escapeHTML(item.model || '-') + '</div>' +
            '</div>' +
            '<div class="log-item-time">' + escapeHTML(item.started_at) + '</div>' +
          '</div>' +
          '<div class="log-item-body">' + escapeHTML(bodyPreview) + '</div>' +
          '<div class="log-item-meta">' +
            '<span class="meta-badge">HTTP ' + escapeHTML(statusText(item.status_code)) + '</span>' +
            '<span class="meta-badge">总量 ' + escapeHTML(String(item.total_tokens || 0)) + '</span>' +
          '</div>';
      }

      listEl.appendChild(button);

      if (isNewNode) {
        if (shouldAnimateChanges) {
          summary.newCount += 1;
          button.classList.add('is-new');
          playTransientClass(button, 'log-item-enter', 1600);
          window.setTimeout(function () {
            button.classList.remove('is-new');
          }, 4200);
        } else {
          button.classList.remove('is-new');
        }
      } else if (shouldAnimateChanges && previousRenderKey && previousRenderKey !== renderKey) {
        summary.updatedCount += 1;
        playTransientClass(button, 'log-item-update', 900);
      } else {
        button.classList.remove('is-new');
      }
    });

    Array.prototype.forEach.call(listEl.querySelectorAll('[data-log-id]'), function (node) {
      if (!expectedIDs[node.getAttribute('data-log-id')]) {
        node.remove();
      }
    });

    if (summary.newCount > 0) {
      playTransientClass(listEl, 'list-surge', 1200);
      playTransientClass(countLabel, 'count-burst', 900);
    } else if (summary.updatedCount > 0) {
      playTransientClass(countLabel, 'count-burst', 760);
    }

    state.hasLoadedList = true;
    return summary;
  }

  function highlightSelected() {
    var buttons = listEl.querySelectorAll('[data-log-id]');
    buttons.forEach(function (button) {
      var id = parseInt(button.getAttribute('data-log-id') || '0', 10);
      var shouldBeActive = id === state.selectedId;
      var wasActive = button.classList.contains('active');
      button.classList.toggle('active', shouldBeActive);
      if (shouldBeActive && !wasActive) {
        playTransientClass(button, 'log-item-selected-pulse', 620);
      }
    });
  }

  function loadDetail(id, updateURL) {
    if (!id) {
      if (pendingDetailController) {
        pendingDetailController.abort();
        pendingDetailController = null;
      }
      if (pendingRawController) {
        pendingRawController.abort();
        pendingRawController = null;
      }
      ['user', 'assistant'].forEach(function (name) {
        if (pendingTextControllers[name]) {
          pendingTextControllers[name].abort();
          pendingTextControllers[name] = null;
        }
      });
      state.selectedId = 0;
      state.detailId = 0;
      if (updateURL) {
        updateSelectedInURL(0);
      }
      showEmpty('没有可显示的日志');
      return Promise.resolve();
    }

    state.selectedId = id;
    if (updateURL) {
      updateSelectedInURL(id);
    }
    if (pendingDetailController) {
      pendingDetailController.abort();
      pendingDetailController = null;
    }
    if (pendingRawController) {
      pendingRawController.abort();
      pendingRawController = null;
    }
    ['user', 'assistant'].forEach(function (name) {
      if (pendingTextControllers[name]) {
        pendingTextControllers[name].abort();
        pendingTextControllers[name] = null;
      }
    });
    highlightSelected();
    renderQuickDetail(state.listItems[String(id)] || null);

    if (detailCache[id]) {
      state.detailId = id;
      applyDetailCore(detailCache[id]);
      if (rawDetailCache[id]) {
        if (shouldAutoLoadRaw(detailCache[id])) {
          applyDetailRaw(rawDetailCache[id]);
        }
        setRefreshLabel('已切换到缓存详情');
        return Promise.resolve(detailCache[id]);
      }

      if (shouldAutoLoadRaw(detailCache[id])) {
        setRefreshLabel('主要内容已加载，正在补全原始详情');
        return ensureRawDetailLoaded();
      }

      setRefreshLabel('主要内容已加载');
      return Promise.resolve(detailCache[id]);
    }

    pendingDetailController = new AbortController();

    setRefreshLabel('正在加载主要详情');
    return fetchJSON(withBase('/api/logs/' + id), { signal: pendingDetailController.signal }).then(function (detail) {
      pendingDetailController = null;
      if (state.selectedId !== id) {
        return null;
      }

      detailCache[id] = detail;
      state.detailId = id;
      applyDetailCore(detail);
      if (shouldAutoLoadRaw(detail)) {
        setRefreshLabel('主要内容已加载，正在补全原始详情');
        return ensureRawDetailLoaded();
      }
      setRefreshLabel('主要内容已加载');
      return detail;
    }).catch(function (error) {
      pendingDetailController = null;
      if (isAbortError(error)) {
        return null;
      }
      setRefreshLabel('详情加载失败');
      showEmpty('详情加载失败');
      return null;
    });
  }

  function openModal(id) {
    var modal = document.getElementById(id);
    if (!modal) {
      return;
    }
    modal.classList.add('show');
    if (id === 'stats-modal') {
      loadDashboardStats();
    } else if (id === 'settings-modal') {
      syncModalFilters();
      loadSettingsTokens();
    } else if (id === 'db-modal') {
      syncModalFilters();
      loadDBStats();
    }
  }

  function closeModal(id) {
    var modal = document.getElementById(id);
    if (modal) {
      modal.classList.remove('show');
    }
  }

  function wireModalButtons() {
    var statsButton = document.getElementById('open-stats-modal');
    var settingsButton = document.getElementById('open-settings-modal');
    var dbButton = document.getElementById('open-db-modal');

    if (statsButton) {
      statsButton.addEventListener('click', function () { openModal('stats-modal'); });
    }
    if (settingsButton) {
      settingsButton.addEventListener('click', function () { openModal('settings-modal'); });
    }
    if (dbButton) {
      dbButton.addEventListener('click', function () { openModal('db-modal'); });
    }

    document.querySelectorAll('[data-close-modal]').forEach(function (button) {
      button.addEventListener('click', function () {
        closeModal(button.getAttribute('data-close-modal'));
      });
    });

    document.querySelectorAll('.modal-backdrop').forEach(function (modal) {
      modal.addEventListener('click', function (event) {
        if (event.target === modal) {
          modal.classList.remove('show');
        }
      });
    });

    window.addEventListener('keydown', function (event) {
      if (event.key === 'Escape') {
        document.querySelectorAll('.modal-backdrop.show').forEach(function (modal) {
          modal.classList.remove('show');
        });
      }
    });
  }

  function wireListPagination() {
    var prevButton = document.getElementById('logs-prev-page');
    var nextButton = document.getElementById('logs-next-page');

    if (prevButton) {
      prevButton.addEventListener('click', function () {
        if (state.page <= 1) {
          return;
        }
        state.page -= 1;
        state.selectedId = 0;
        state.detailId = 0;
        updatePageInURL(state.page);
        refreshList(true);
      });
    }

    if (nextButton) {
      nextButton.addEventListener('click', function () {
        state.page += 1;
        state.selectedId = 0;
        state.detailId = 0;
        updatePageInURL(state.page);
        refreshList(true);
      });
    }
  }

  function wireTextPagers() {
    [
      ['detail-user-prev', 'user', -1],
      ['detail-user-next', 'user', 1],
      ['detail-assistant-prev', 'assistant', -1],
      ['detail-assistant-next', 'assistant', 1]
    ].forEach(function (item) {
      var button = document.getElementById(item[0]);
      if (!button) {
        return;
      }
      button.addEventListener('click', function () {
        changeTextPage(item[1], item[2]);
      });
    });
  }

  function wireRawDetails() {
    [detailRawSectionEl, detailJSONSectionEl].forEach(function (section) {
      if (!section) {
        return;
      }
      section.addEventListener('toggle', function () {
        if (section.open) {
          ensureRawDetailLoaded();
        }
      });
    });
  }

  function syncModalFilters() {
    var map = [
      ['from', 'settings-filter-from'],
      ['to', 'settings-filter-to'],
      ['token', 'settings-filter-token'],
      ['alias', 'settings-filter-alias'],
      ['model', 'settings-filter-model'],
      ['from', 'db-filter-from'],
      ['to', 'db-filter-to'],
      ['token', 'db-filter-token'],
      ['alias', 'db-filter-alias'],
      ['model', 'db-filter-model']
    ];

    map.forEach(function (entry) {
      var source = document.querySelector('[name="' + entry[0] + '"]');
      var target = document.getElementById(entry[1]);
      if (source && target) {
        target.value = source.value;
      }
    });
  }

  function loadDashboardStats() {
    return fetchJSON(withBase('/api/dashboard')).then(function (data) {
      setText('stats-total-requests', String(data.total_requests || 0));
      setText('stats-today-requests', String(data.today_requests || 0));
      setText('stats-error-count', String(data.error_count || 0));
      setText('stats-total-tokens', String(data.total_tokens || 0));
      setText('stats-prompt-tokens', String(data.total_prompt_tokens || 0));
      setText('stats-completion-tokens', String(data.total_completion_tokens || 0));

      var table = document.getElementById('stats-token-table');
      if (!table) {
        return;
      }
      if (!data.token_groups || !data.token_groups.length) {
        table.innerHTML = '<tr><td colspan="6" class="muted">暂无数据</td></tr>';
        return;
      }

      table.innerHTML = data.token_groups.map(function (item) {
        return '' +
          '<tr>' +
            '<td>' + escapeHTML(item.token_alias || '-') + '</td>' +
            '<td>' + escapeHTML(item.token_preview || '-') + '</td>' +
            '<td>' + escapeHTML(String(item.request_count || 0)) + '</td>' +
            '<td>' + escapeHTML(String(item.total_tokens || 0)) + '</td>' +
            '<td>' + escapeHTML(String(item.error_count || 0)) + '</td>' +
            '<td>' + escapeHTML(item.last_seen || '-') + '</td>' +
          '</tr>';
      }).join('');
    }).catch(function () {
      var table = document.getElementById('stats-token-table');
      if (table) {
        table.innerHTML = '<tr><td colspan="6" class="muted">加载统计失败</td></tr>';
      }
    });
  }

  function loadSettingsTokens() {
    var form = document.getElementById('settings-token-filter-form');
    var message = document.getElementById('settings-message');
    if (!form) {
      return Promise.resolve();
    }

    var params = new URLSearchParams(new FormData(form));
    var url = withBase('/api/tokens');
    var query = params.toString();
    return fetchJSON(query ? url + '?' + query : url).then(function (data) {
      var table = document.getElementById('settings-token-table');
      if (!table) {
        return;
      }
      if (!data.items || !data.items.length) {
        table.innerHTML = '<tr><td colspan="6" class="muted">暂无匹配的 token</td></tr>';
        if (message) {
          message.textContent = '可以使用上方表单手动绑定 token 代号。';
        }
        return;
      }

      table.innerHTML = data.items.map(function (item, index) {
        var inputId = 'settings-inline-alias-' + index;
        return '' +
          '<tr>' +
            '<td><input id="' + inputId + '" data-token-fingerprint="' + escapeHTML(item.token_fingerprint || '') + '" value="' + escapeHTML(item.token_alias || '') + '" placeholder="输入令牌代号"></td>' +
            '<td>' + escapeHTML(item.token_preview || '-') + '</td>' +
            '<td><code>' + escapeHTML(item.token_fingerprint || '-') + '</code></td>' +
            '<td>' + escapeHTML(String(item.request_count || 0)) + '</td>' +
            '<td>' + escapeHTML(String(item.total_tokens || 0)) + '</td>' +
            '<td><button class="button secondary" type="button" data-save-inline-alias="' + inputId + '">保存</button></td>' +
          '</tr>';
      }).join('');

      if (message) {
        message.textContent = '点击“保存”可直接修改 token 代号。';
      }
    }).catch(function () {
      var table = document.getElementById('settings-token-table');
      if (table) {
        table.innerHTML = '<tr><td colspan="6" class="muted">加载 token 列表失败</td></tr>';
      }
    });
  }

  function saveAlias(formData, messageEl) {
    return postForm(withBase('/api/tokens/alias'), formData).then(function (data) {
      if (messageEl) {
        messageEl.textContent = data.message || '令牌代号已保存';
      }
      return Promise.all([
        refreshList(true),
        loadSettingsTokens(),
        loadDashboardStats(),
        loadDBStats()
      ]).then(function () {
        return true;
      });
    }).catch(function (error) {
      if (messageEl) {
        messageEl.textContent = error.message || '保存令牌代号失败';
      }
      return false;
    });
  }

  function wireSettingsForms() {
    var aliasForm = document.getElementById('settings-alias-form');
    var settingsMessage = document.getElementById('settings-message');
    var tokenFilterForm = document.getElementById('settings-token-filter-form');
    var tokenTable = document.getElementById('settings-token-table');

    if (aliasForm) {
      aliasForm.addEventListener('submit', function (event) {
        event.preventDefault();
        saveAlias(new FormData(aliasForm), settingsMessage).then(function (ok) {
          if (ok) {
            aliasForm.reset();
          }
        });
      });
    }

    if (tokenFilterForm) {
      tokenFilterForm.addEventListener('submit', function (event) {
        event.preventDefault();
        loadSettingsTokens();
      });
    }

    if (tokenTable) {
      tokenTable.addEventListener('click', function (event) {
        var button = event.target.closest('[data-save-inline-alias]');
        if (!button) {
          return;
        }
        var input = document.getElementById(button.getAttribute('data-save-inline-alias'));
        if (!input) {
          return;
        }
        var formData = new FormData();
        formData.set('token_fingerprint', input.getAttribute('data-token-fingerprint') || '');
        formData.set('token_alias', input.value || '');
        saveAlias(formData, settingsMessage);
      });
    }
  }

  function loadDBStats() {
    return fetchJSON(withBase('/api/db/stats')).then(function (data) {
      setText('db-total-rows', String(data.total_rows || 0));
      setText('db-today-rows', String(data.today_rows || 0));
      setText('db-size-pretty', data.database_pretty || '0');
      setText('db-audit-total-pretty', data.audit_total_pretty || '0');
      setText('db-audit-table-pretty', data.audit_table_pretty || '0');
      setText('db-audit-index-pretty', data.audit_index_pretty || '0');
      setText('db-audit-toast-pretty', data.audit_toast_pretty || '0');
      setText('db-live-tuples', String(data.live_tuples || 0));
      setText('db-dead-tuples', String(data.dead_tuples || 0));
      setText('db-last-vacuum', data.last_vacuum || '-');
      setText('db-last-autovacuum', data.last_autovacuum || '-');
      setText('db-last-analyze', data.last_analyze || '-');
      setText('db-last-autoanalyze', data.last_autoanalyze || '-');
    });
  }

  function refreshDBPanels() {
    return Promise.all([
      refreshList(true),
      loadDBStats(),
      loadDashboardStats(),
      loadSettingsTokens()
    ]);
  }

  function setDBMaintenanceBusy(busy) {
    var controls = document.querySelectorAll('#db-refresh-button, [data-db-maintenance]');
    controls.forEach(function (element) {
      element.disabled = !!busy;
    });
  }

  function wireDBMaintenance() {
    var refreshButton = document.getElementById('db-refresh-button');
    var message = document.getElementById('db-maintenance-message');
    var actionButtons = document.querySelectorAll('[data-db-maintenance]');

    if (refreshButton) {
      refreshButton.addEventListener('click', function () {
        if (message) {
          message.textContent = '正在刷新数据库统计…';
        }
        setDBMaintenanceBusy(true);
        refreshDBPanels().then(function () {
          if (message) {
            message.textContent = '数据库统计已刷新。';
          }
        }).catch(function () {
          if (message) {
            message.textContent = '刷新数据库统计失败。';
          }
        }).finally(function () {
          setDBMaintenanceBusy(false);
        });
      });
    }

    actionButtons.forEach(function (button) {
      button.addEventListener('click', function () {
        var mode = button.getAttribute('data-db-maintenance') || '';
        var confirmMessage = '确认执行数据库维护吗？';
        if (mode === 'vacuum_full') {
          confirmMessage = '确认执行强制缩盘吗？这会运行 VACUUM FULL，并在执行期间锁住 audit_logs 表。';
        } else if (mode === 'compact_payloads') {
          confirmMessage = '确认执行历史记录瘦身吗？这会重写已入库旧日志里的大图片 base64，以便后续缩盘。';
        } else {
          confirmMessage = '确认执行整理空间吗？这会运行 VACUUM ANALYZE。';
        }
        if (!window.confirm(confirmMessage)) {
          return;
        }

        if (message) {
          if (mode === 'vacuum_full') {
            message.textContent = '正在执行 VACUUM FULL，这一步可能会持续较久…';
          } else if (mode === 'compact_payloads') {
            message.textContent = '正在重写历史记录中的大字段，这一步完成后建议再执行一次强制缩盘…';
          } else {
            message.textContent = '正在执行 VACUUM ANALYZE…';
          }
        }

        setDBMaintenanceBusy(true);
        var formData = new FormData();
        formData.set('mode', mode);
        postForm(withBase('/api/db/maintenance'), formData).then(function (data) {
          if (message) {
            message.textContent = data.message || '数据库维护已完成。';
          }
          return refreshDBPanels();
        }).catch(function (error) {
          if (message) {
            message.textContent = error.message || '数据库维护失败。';
          }
        }).finally(function () {
          setDBMaintenanceBusy(false);
        });
      });
    });
  }

  function wireDBCleanup() {
    var form = document.getElementById('db-cleanup-form');
    var message = document.getElementById('db-cleanup-message');
    if (!form) {
      return;
    }

    form.addEventListener('submit', function (event) {
      event.preventDefault();
      if (!window.confirm('确定要清理符合条件的审计记录吗？')) {
        return;
      }

      postForm(withBase('/api/db/cleanup'), new FormData(form)).then(function (data) {
        if (message) {
          message.textContent = data.message || '清理完成';
        }
        return refreshDBPanels();
      }).catch(function (error) {
        if (message) {
          message.textContent = error.message || '清理失败';
        }
      });
    });
  }

  function refreshList(force) {
    if (!force && !document.hidden) {
      setRefreshLabel('检查更新中');
    }

    return fetchJSON(apiURL('/api/logs')).then(function (data) {
      renderListPagination(data);

      if (!force && data.version === state.version) {
        setRefreshLabel('无新变化');
        return;
      }

      if (data.version !== state.version) {
        detailCache = Object.create(null);
        rawDetailCache = Object.create(null);
        textPageCache = Object.create(null);
      }
      state.version = data.version || '';
      renderList(data.items || [], data.total_count || 0);

      var exists = false;
      (data.items || []).forEach(function (item) {
        if (item.id === state.selectedId) {
          exists = true;
        }
      });

      if (!exists) {
        if (data.items && data.items.length > 0) {
          return loadDetail(data.items[0].id, true);
        }
        return loadDetail(0, true);
      }

      return loadDetail(state.selectedId, false).then(function () {
        setRefreshLabel('已同步至最新');
      });
    }).catch(function () {
      setRefreshLabel('自动刷新失败');
      if (!state.version) {
        listEl.innerHTML = '<div class="empty-list">加载日志失败，请稍后刷新页面。</div>';
      }
    });
  }

  function refreshVersionOnly() {
    return fetchJSON(apiURL('/api/logs/version')).then(function (data) {
      if (!state.version || data.version !== state.version) {
        return refreshList(true);
      }
      setRefreshLabel('无新变化');
    }).catch(function () {
      setRefreshLabel('自动刷新失败');
    });
  }

  function refreshList(force) {
    if (!force && !document.hidden) {
      setSidebarRefreshing(true);
      setRefreshLabel('检查更新中', 'checking');
    }

    return fetchJSON(apiURL('/api/logs')).then(function (data) {
      renderListPagination(data);

      if (!force && data.version === state.version) {
        setSidebarRefreshing(false);
        setRefreshLabel('无新变化');
        return;
      }

      if (data.version !== state.version) {
        detailCache = Object.create(null);
        rawDetailCache = Object.create(null);
        textPageCache = Object.create(null);
      }

      var previousVersion = state.version || '';
      state.version = data.version || '';
      var renderSummary = renderList(data.items || [], data.total_count || 0);
      setSidebarRefreshing(false);

      var exists = false;
      (data.items || []).forEach(function (item) {
        if (item.id === state.selectedId) {
          exists = true;
        }
      });

      function finalizeRefreshLabel() {
        if (renderSummary.newCount > 0) {
          setRefreshLabel('已接收 ' + renderSummary.newCount + ' 条新日志', 'updates');
        } else if (renderSummary.updatedCount > 0) {
          setRefreshLabel('有 ' + renderSummary.updatedCount + ' 条日志状态更新', 'updates');
        } else if (previousVersion) {
          setRefreshLabel('已同步至最新');
        } else {
          setRefreshLabel('日志已加载');
        }
      }

      if (!exists) {
        if (data.items && data.items.length > 0) {
          return loadDetail(data.items[0].id, true).then(function () {
            finalizeRefreshLabel();
          });
        }
        return loadDetail(0, true).then(function () {
          finalizeRefreshLabel();
        });
      }

      return loadDetail(state.selectedId, false).then(function () {
        finalizeRefreshLabel();
      });
    }).catch(function () {
      setSidebarRefreshing(false);
      setRefreshLabel('自动刷新失败', 'error');
      if (!state.version) {
        listEl.innerHTML = '<div class="empty-list">加载日志失败，请稍后刷新页面。</div>';
      }
    });
  }

  function refreshVersionOnly() {
    if (!document.hidden) {
      setSidebarRefreshing(true);
    }
    return fetchJSON(apiURL('/api/logs/version')).then(function (data) {
      if (!state.version || data.version !== state.version) {
        return refreshList(true);
      }
      setSidebarRefreshing(false);
      setRefreshLabel('无新变化');
    }).catch(function () {
      setSidebarRefreshing(false);
      setRefreshLabel('自动刷新失败', 'error');
    });
  }

  wireModalButtons();
  wireListPagination();
  wireTextPagers();
  wireRawDetails();
  wireSettingsForms();
  wireDBCleanup();
  wireDBMaintenance();
  syncModalFilters();

  refreshList(true).then(function () {
    if (state.selectedId > 0 && state.detailId !== state.selectedId) {
      return loadDetail(state.selectedId, false);
    }
    return null;
  });

  setInterval(function () {
    refreshVersionOnly();
  }, 5000);
})();
</script>
{{end}}
`

const tokensTemplate = `
{{define "tokens"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<div class="shell shell-wide">
  <div class="nav">
    <div class="brand">newapi-audit-proxy</div>
    <div class="nav-links">
      <a href="{{path "/"}}">首页</a>
      <a href="{{path "/logs"}}">日志</a>
      <a href="{{path "/tokens"}}">令牌管理</a>
      <form method="post" action="{{path "/logout"}}" style="margin: 0;">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button class="button secondary" type="submit">退出登录</button>
      </form>
    </div>
  </div>

  <div class="panel">
    <div class="hero">
      <div>
        <h1 style="font-size: 24px;">Token 映射与统计</h1>
        <p>可按时间、模型、token 指纹和令牌代号查看汇总用量，也可直接手动绑定 token 与代号的对应关系。</p>
      </div>
      <div class="pill">总令牌组数：{{.Result.TotalCount}}</div>
    </div>

    {{if .Saved}}
    <div class="success-box">令牌代号已保存</div>
    {{end}}
    {{if .Error}}
    <div class="error-box">{{.Error}}</div>
    {{end}}
    <div style="margin-bottom: 18px;">
      <h2 style="margin: 0 0 10px; font-size: 18px;">手动绑定 Token 代号</h2>
      <form method="post" action="{{path "/tokens/alias"}}">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <input type="hidden" name="redirect_to" value="{{$.CurrentURL}}">
        <div class="token-manual-grid">
          <div>
            <label for="token_value">Token 原文（只用来现场计算指纹）</label>
            <input id="token_value" name="token_value" placeholder="sk-... / Bearer sk-...">
          </div>
          <div>
            <label for="manual_alias">令牌代号</label>
            <input id="manual_alias" name="token_alias" placeholder="例如：主账号A">
          </div>
          <div>
            <button type="submit">保存映射</button>
          </div>
        </div>
      </form>
      <p class="token-inline-note">系统不会保存 token 原文，只会用 <code>hmac_secret</code> 现场计算指纹后保存映射关系。清空代号再保存，可删除已绑定的映射。</p>
    </div>

    <p class="muted" style="margin-top: 0;">提示：下方列表支持直接修改令牌代号，Token 指纹支持部分匹配。</p>

    <form method="get" action="{{path "/tokens"}}">
      <div class="form-grid">
        <div>
          <label for="from">开始时间</label>
          <input id="from" name="from" type="datetime-local" value="{{.Filters.From}}">
        </div>
        <div>
          <label for="to">结束时间</label>
          <input id="to" name="to" type="datetime-local" value="{{.Filters.To}}">
        </div>
        <div>
          <label for="token">Token 指纹（支持部分匹配）</label>
          <input id="token" name="token" value="{{.Filters.TokenFingerprint}}">
        </div>
        <div>
          <label for="alias">令牌代号</label>
          <input id="alias" name="alias" value="{{.Filters.TokenAlias}}">
        </div>
        <div>
          <label for="model">模型</label>
          <input id="model" name="model" value="{{.Filters.Model}}">
        </div>
        <div>
          <button type="submit">搜索</button>
        </div>
      </div>
    </form>
  </div>

  <div class="panel">
    <table>
      <thead>
        <tr>
          <th>令牌代号</th>
          <th>Token 指纹</th>
          <th>Token 预览</th>
          <th>请求数</th>
          <th>错误数</th>
          <th>输入 Tokens</th>
          <th>输出 Tokens</th>
          <th>总 Tokens</th>
          <th>首次时间</th>
          <th>最后时间</th>
        </tr>
      </thead>
      <tbody>
        {{range .Result.Items}}
        <tr>
          <td>
            <form class="token-form" method="post" action="{{path "/tokens/alias"}}">
              <input type="hidden" name="csrf_token" value="{{$.CSRFToken}}">
              <input type="hidden" name="token_fingerprint" value="{{.TokenFingerprint}}">
              <input type="hidden" name="redirect_to" value="{{$.CurrentURL}}">
              <input name="token_alias" value="{{.TokenAlias}}" placeholder="输入令牌代号">
              <button type="submit">保存</button>
            </form>
          </td>
          <td><code>{{if .TokenFingerprint}}{{.TokenFingerprint}}{{else}}（无）{{end}}</code></td>
          <td>{{if .TokenPreview}}{{.TokenPreview}}{{else}}-{{end}}</td>
          <td>{{.RequestCount}}</td>
          <td>{{.ErrorCount}}</td>
          <td>{{.PromptTokens}}</td>
          <td>{{.CompletionTokens}}</td>
          <td>{{.TotalTokens}}</td>
          <td>{{formatTime .FirstSeen}}</td>
          <td>{{formatTime .LastSeen}}</td>
        </tr>
        {{else}}
        <tr><td colspan="10" class="muted">暂无令牌数据</td></tr>
        {{end}}
      </tbody>
    </table>

    <div class="pagination">
      {{if gt .PrevPage 0}}
      <a class="button secondary" href="{{tokenPageURL .Filters .PrevPage}}">上一页</a>
      {{end}}
      <span class="muted">第 {{.Result.Page}} 页</span>
      {{if .HasNext}}
      <a class="button secondary" href="{{tokenPageURL .Filters .NextPage}}">下一页</a>
      {{end}}
    </div>
  </div>
</div>
{{end}}
`

const detailTemplate = `
{{define "detail"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<div class="shell shell-wide">
  <div class="nav">
    <div class="brand">newapi-audit-proxy</div>
    <div class="nav-links">
      <a href="{{path "/"}}">首页</a>
      <a href="{{path "/logs"}}">日志</a>
      <form method="post" action="{{path "/logout"}}" style="margin: 0;">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button class="button secondary" type="submit">退出登录</button>
      </form>
    </div>
  </div>

  <div class="panel">
    <div class="hero">
      <div>
        <h1 style="font-size: 24px;">日志详情 #{{.Log.ID}}</h1>
        <p>{{formatTime .Log.StartedAt}} - {{.Log.Method}} {{pathWithQuery .Log.Path .Log.QueryString}}</p>
      </div>
      <div class="pill">{{if .Log.ResponseIsSSE}}SSE{{else}}JSON/HTTP{{end}}</div>
    </div>

    <div class="grid">
      <div class="stat">
        <h3>状态码</h3>
        <div class="value">{{.Log.StatusCode}}</div>
      </div>
      <div class="stat">
        <h3>模型</h3>
        <div class="value" style="font-size: 22px;">{{if .Log.Model}}{{.Log.Model}}{{else}}-{{end}}</div>
      </div>
      <div class="stat">
        <h3>令牌代号</h3>
        <div class="value" style="font-size: 22px;">{{if .Log.TokenAlias}}{{.Log.TokenAlias}}{{else}}-{{end}}</div>
      </div>
      <div class="stat">
        <h3>总 Tokens</h3>
        <div class="value" style="font-size: 22px;">{{.Log.TotalTokens}}</div>
      </div>
    </div>
  </div>

  <div class="detail-grid content-first">
    <div class="panel">
      <h3>解析字段</h3>
      <p><strong>Token 指纹：</strong> <code>{{if .Log.TokenFingerprint}}{{.Log.TokenFingerprint}}{{else}}（无）{{end}}</code></p>
      <p><strong>Token 预览：</strong> {{if .Log.TokenPreview}}{{.Log.TokenPreview}}{{else}}-{{end}}</p>
      <p><strong>输入 Tokens：</strong> {{.Log.PromptTokens}}</p>
      <p><strong>输出 Tokens：</strong> {{.Log.CompletionTokens}}</p>
      <p><strong>总 Tokens：</strong> {{.Log.TotalTokens}}</p>
      <p><strong>流式：</strong> {{streamValue .Log.Stream}}</p>
      <p><strong>请求字节数：</strong> {{.Log.RequestBytes}}{{if .Log.RequestTruncated}}（已截断）{{end}}</p>
      <p><strong>响应字节数：</strong> {{.Log.ResponseBytes}}{{if .Log.ResponseTruncated}}（已截断）{{end}}</p>
      <p><strong>错误：</strong> {{if .Log.ErrorText}}{{.Log.ErrorText}}{{else}}-{{end}}</p>
      {{if .Log.TokenFingerprint}}
      <p><a href="{{path "/logs"}}?selected={{.Log.ID}}&token={{.Log.TokenFingerprint}}">回到日志主页，在“设置”中管理令牌代号</a></p>
      {{end}}
      {{if .Images}}
      <h3>绘图结果</h3>
      <div class="gallery">
        {{range .Images}}
        <div class="gallery-card">
          <img src="{{.DataURL}}" alt="{{.Label}}">
          <p>{{.Label}}</p>
        </div>
        {{end}}
      </div>
      {{end}}
      <h3>用户文本</h3>
      <pre>{{if .Log.UserText}}{{.Log.UserText}}{{else}}（空）{{end}}</pre>
      <h3>助手文本</h3>
      <pre>{{if .Log.AssistantText}}{{.Log.AssistantText}}{{else}}（空）{{end}}</pre>
      <h3>用量 JSON</h3>
      <pre>{{prettyJSON .Log.UsageJSON}}</pre>
    </div>

    <div class="panel">
      <h3>请求头</h3>
      <pre>{{prettyAny .Log.RequestHeaders}}</pre>
      <h3>响应头</h3>
      <pre>{{prettyAny .Log.ResponseHeaders}}</pre>
    </div>
  </div>

  <div class="detail-grid content-first">
    <div class="panel">
      <h3>原始请求体</h3>
      <pre>{{renderBody .Log.RequestBody}}</pre>
      <h3>请求 JSON</h3>
      <pre>{{prettyJSON .Log.RequestJSON}}</pre>
    </div>

    <div class="panel">
      <h3>原始响应体</h3>
      <pre>{{renderBody .Log.ResponseBody}}</pre>
      <h3>响应 JSON</h3>
      <pre>{{prettyJSON .Log.ResponseJSON}}</pre>
    </div>
  </div>
</div>
{{end}}
`
