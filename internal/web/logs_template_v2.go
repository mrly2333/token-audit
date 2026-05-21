package web

const logsTemplateV2 = `
{{define "logs"}}
{{template "base" .}}
{{end}}

{{define "body"}}
<style>
  html,
  body {
    height: 100%;
    overflow: hidden;
  }
  body {
    background:
      radial-gradient(circle at 10% 8%, rgba(255, 171, 206, 0.42), transparent 30%),
      radial-gradient(circle at 88% 6%, rgba(255, 219, 232, 0.48), transparent 24%),
      radial-gradient(circle at 22% 88%, rgba(255, 192, 218, 0.32), transparent 28%),
      radial-gradient(circle at 74% 72%, rgba(255, 233, 240, 0.42), transparent 22%),
      linear-gradient(180deg, #fff9fc 0%, #fff2f7 50%, #fff7fb 100%);
    background-size: 120% 120%, 110% 110%, 130% 130%, 115% 115%, 100% 100%;
  }
  body::before {
    inset: -10vh -8vw auto;
    height: 52vh;
    background:
      radial-gradient(circle at 15% 24%, rgba(255, 167, 196, 0.42), transparent 30%),
      radial-gradient(circle at 82% 18%, rgba(255, 226, 237, 0.42), transparent 28%);
    filter: blur(34px);
  }
  body::after {
    content: "";
    position: fixed;
    inset: auto -8vw -14vh;
    height: 42vh;
    background:
      radial-gradient(circle at 18% 42%, rgba(255, 196, 220, 0.24), transparent 26%),
      radial-gradient(circle at 76% 28%, rgba(255, 223, 236, 0.32), transparent 24%);
    filter: blur(42px);
    pointer-events: none;
    z-index: 0;
  }
  .logs-screen {
    height: 100vh;
    padding: 6px 10px 10px;
    overflow: hidden;
  }
  .logs-screen,
  .monitor-modal {
    --rose-900: #6f213f;
    --rose-800: #843051;
    --rose-700: #a94368;
    --rose-600: #d75f90;
    --rose-500: #ef79aa;
    --rose-400: #ff9fc3;
    --rose-300: #ffc0d6;
    --rose-200: #ffdde9;
    --rose-100: #fff1f7;
    --rose-surface: rgba(255, 255, 255, 0.72);
    --rose-surface-strong: rgba(255, 249, 252, 0.92);
    --rose-line: rgba(239, 177, 203, 0.74);
    --rose-shadow: 0 18px 42px rgba(215, 95, 144, 0.14);
    --rose-scroll: rgba(239, 121, 170, 0.82);
    --rose-scroll-track: rgba(255, 232, 241, 0.72);
  }
  .monitor-shell {
    height: calc(100vh - 16px);
    min-height: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
    overflow: hidden;
  }
  .monitor-topbar {
    flex: 0 0 auto;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
    align-items: center;
    gap: 8px;
    margin-bottom: 0;
    padding: 7px 9px;
    border-radius: 16px;
    top: auto;
    position: relative;
    border-color: rgba(239, 177, 203, 0.78);
    background: linear-gradient(180deg, rgba(255, 249, 252, 0.9), rgba(255, 238, 246, 0.84));
    box-shadow:
      0 10px 28px rgba(214, 121, 160, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.82);
    backdrop-filter: blur(28px) saturate(1.08);
    overflow: hidden;
  }
  .monitor-topbar::before {
    content: "";
    position: absolute;
    inset: 0;
    background:
      radial-gradient(circle at 15% 18%, rgba(255, 202, 222, 0.62), transparent 24%),
      radial-gradient(circle at 88% 16%, rgba(255, 255, 255, 0.72), transparent 22%);
    opacity: 0.86;
    pointer-events: none;
  }
  .monitor-topbar > * {
    position: relative;
    z-index: 1;
  }
  .topbar-overview {
    display: flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
    flex-wrap: wrap;
    justify-self: start;
  }
  .topbar-status {
    display: flex;
    align-items: center;
    justify-content: center;
    justify-self: center;
    min-width: 0;
  }
  .topbar-actions {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
    justify-content: flex-end;
    justify-self: end;
  }
  .sync-pill {
    padding: 4px 8px;
    min-width: 132px;
    text-align: center;
    font-size: 11px;
    line-height: 1.15;
    background: rgba(255, 224, 236, 0.96);
    color: var(--rose-700);
    border: 1px solid rgba(228, 126, 170, 0.22);
    box-shadow: 0 6px 14px rgba(225, 138, 178, 0.14);
  }
  .sync-pill.is-checking {
    background: rgba(255, 214, 228, 0.92);
    color: #b03f6b;
    box-shadow: 0 0 0 1px rgba(231, 131, 169, 0.16), 0 14px 28px rgba(236, 145, 179, 0.22);
  }
  .sync-pill.has-updates {
    background: rgba(255, 199, 221, 0.96);
    color: #a13f64;
    box-shadow: 0 0 0 1px rgba(224, 119, 160, 0.18), 0 16px 30px rgba(224, 119, 160, 0.24);
  }
  .sync-pill.is-error {
    background: rgba(255, 226, 230, 0.94);
    color: #bc4055;
    box-shadow: 0 0 0 1px rgba(188, 64, 85, 0.16), 0 14px 28px rgba(188, 64, 85, 0.12);
  }
  .topbar-stat {
    min-width: 82px;
    padding: 5px 8px;
    border-radius: 13px;
    background: rgba(255, 255, 255, 0.62);
    border: 1px solid rgba(235, 193, 209, 0.82);
    backdrop-filter: blur(18px);
    box-shadow: 0 6px 14px rgba(216, 131, 168, 0.08);
    transition: box-shadow 0.18s ease, background 0.18s ease, border-color 0.18s ease;
  }
  .topbar-stat:hover {
    background: rgba(255, 255, 255, 0.78);
    border-color: rgba(235, 171, 197, 0.9);
    box-shadow: 0 10px 24px rgba(216, 131, 168, 0.14);
  }
  .topbar-stat span {
    display: block;
    font-size: 9px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #9d7885;
    margin-bottom: 2px;
    font-weight: 700;
  }
  .topbar-stat strong {
    display: block;
    font-size: 14px;
    color: var(--rose-900);
    line-height: 1.2;
    white-space: nowrap;
  }
  .topbar-stat.pulse strong {
    animation: rose-count-pop 0.38s ease-out;
  }
  .logs-screen .button,
  .logs-screen button,
  .monitor-modal .button,
  .monitor-modal button {
    appearance: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    background: linear-gradient(180deg, #ffa9ca 0%, #ef7fac 100%);
    color: #fff;
    border: 1px solid rgba(218, 101, 149, 0.36);
    min-height: 32px;
    padding: 7px 13px;
    border-radius: 999px;
    box-shadow:
      0 8px 16px rgba(221, 105, 154, 0.14),
      inset 0 1px 0 rgba(255, 255, 255, 0.46);
    font-weight: 800;
    font-size: 12px;
    line-height: 1;
    text-decoration: none;
    white-space: nowrap;
    cursor: pointer;
    position: relative;
    overflow: visible;
    transform: none;
    transition: transform 0.14s ease, box-shadow 0.14s ease, background 0.14s ease, border-color 0.14s ease, color 0.14s ease;
  }
  .logs-screen .button:hover,
  .logs-screen button:hover,
  .monitor-modal .button:hover,
  .monitor-modal button:hover {
    transform: translateY(-1px);
    border-color: rgba(203, 78, 128, 0.42);
    background: linear-gradient(180deg, #ffb7d3 0%, #f08ab6 100%);
    box-shadow:
      0 10px 20px rgba(221, 105, 154, 0.18),
      inset 0 1px 0 rgba(255, 255, 255, 0.56);
  }
  .logs-screen .button:active,
  .logs-screen button:active,
  .monitor-modal .button:active,
  .monitor-modal button:active {
    transform: translateY(0);
    box-shadow:
      0 4px 10px rgba(221, 105, 154, 0.12),
      inset 0 1px 3px rgba(132, 48, 81, 0.12);
  }
  .logs-screen .button:focus-visible,
  .logs-screen button:focus-visible,
  .monitor-modal .button:focus-visible,
  .monitor-modal button:focus-visible {
    outline: none;
    box-shadow:
      0 0 0 3px rgba(255, 196, 220, 0.7),
      0 8px 18px rgba(221, 105, 154, 0.16);
  }
  .logs-screen .button:disabled,
  .logs-screen button:disabled,
  .monitor-modal .button:disabled,
  .monitor-modal button:disabled {
    cursor: not-allowed;
    opacity: 0.54;
    transform: none;
    box-shadow: none;
  }
  .logs-screen .button.secondary,
  .logs-screen button.secondary,
  .monitor-modal .button.secondary,
  .monitor-modal button.secondary {
    background: rgba(255, 255, 255, 0.76);
    color: #9b3f63;
    border-color: rgba(231, 145, 181, 0.48);
    box-shadow:
      0 7px 14px rgba(221, 105, 154, 0.1),
      inset 0 1px 0 rgba(255, 255, 255, 0.82);
  }
  .logs-screen .button.secondary:hover,
  .logs-screen button.secondary:hover,
  .monitor-modal .button.secondary:hover,
  .monitor-modal button.secondary:hover {
    background: rgba(255, 247, 251, 0.94);
    color: #873354;
    border-color: rgba(222, 116, 159, 0.58);
  }
  .logs-screen .button.danger,
  .logs-screen button.danger,
  .monitor-modal .button.danger,
  .monitor-modal button.danger {
    background: linear-gradient(180deg, #ff8fb0 0%, #e55e86 100%);
    color: #fff;
    border-color: rgba(196, 65, 105, 0.34);
  }
  .monitor-grid {
    flex: 1 1 0;
    display: grid;
    grid-template-columns: minmax(300px, 20%) minmax(0, 1fr);
    gap: 10px;
    height: auto;
    min-height: 0;
    max-height: 100%;
    overflow: hidden;
    align-items: stretch;
  }
  .monitor-grid > * {
    min-height: 0;
    min-width: 0;
  }
  .log-lane,
  .detail-lane {
    margin-bottom: 0;
    min-height: 0;
    height: auto;
    max-height: 100%;
  }
  .log-lane {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background:
      radial-gradient(circle at 14% 0%, rgba(255, 219, 232, 0.72), transparent 34%),
      linear-gradient(180deg, rgba(255, 252, 254, 0.9), rgba(255, 239, 246, 0.86));
    border-color: var(--rose-line);
    backdrop-filter: blur(30px) saturate(1.1);
    overflow: hidden;
  }
  .detail-lane {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    padding: 16px 18px;
    background:
      radial-gradient(circle at 84% 0%, rgba(255, 224, 236, 0.68), transparent 30%),
      linear-gradient(180deg, rgba(255, 253, 254, 0.92), rgba(255, 244, 249, 0.88));
    border-color: var(--rose-line);
    backdrop-filter: blur(32px) saturate(1.12);
    overflow: hidden;
  }
  .lane-head {
    flex: 0 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: start;
    gap: 14px;
    margin-bottom: 0;
  }
  .lane-head h1,
  .lane-head h2,
  .detail-toolbar h2 {
    margin: 0;
    font-size: 20px;
    color: #762f4b;
  }
  .lane-head p,
  .detail-toolbar p {
    margin: 6px 0 0;
    color: #936d7c;
    font-size: 13px;
    line-height: 1.5;
  }
  .lane-head .pill {
    background: rgba(255, 228, 237, 0.9);
    color: #b04b6f;
  }
  .list-summary {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    margin-top: 10px;
  }
  .summary-chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.68);
    border: 1px solid rgba(236, 192, 208, 0.82);
    color: #8f6073;
    font-size: 12px;
    font-weight: 700;
  }
  .scroll-fade {
    position: relative;
    display: flex;
    flex-direction: column;
    flex: 1 1 0;
    height: auto;
    min-height: 0;
    max-height: 100%;
    overflow: hidden;
  }
  .scroll-fade::before,
  .scroll-fade::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    height: 32px;
    z-index: 5;
    pointer-events: none;
    opacity: 0.95;
  }
  .scroll-fade::before {
    top: 0;
    background: linear-gradient(180deg, rgba(255, 246, 250, 0.98), rgba(255, 246, 250, 0));
  }
  .scroll-fade::after {
    bottom: 0;
    background: linear-gradient(0deg, rgba(255, 244, 249, 0.98), rgba(255, 244, 249, 0));
  }
  .lane-scroll,
  .detail-scroll {
    flex: 1 1 0;
    min-height: 0;
    overflow: auto;
    padding-right: 4px;
    scroll-behavior: smooth;
  }
  .lane-scroll {
    display: flex;
    flex-direction: column;
    gap: 10px;
    height: auto;
    max-height: 100%;
    padding-bottom: 8px;
    overscroll-behavior: contain;
    overflow-x: hidden;
    overflow-y: scroll;
    -webkit-overflow-scrolling: touch;
    touch-action: pan-y;
    scrollbar-gutter: stable;
  }
  .log-list-frame {
    position: relative;
    z-index: 1;
    flex: 1 1 auto;
    min-height: 120px;
    height: auto;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-radius: 24px;
    border: 1px solid rgba(239, 177, 203, 0.66);
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.54), rgba(255, 243, 248, 0.42));
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72), 0 16px 32px rgba(215, 95, 144, 0.08);
    backdrop-filter: blur(18px);
  }
  .log-list-frame::before,
  .log-list-frame::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    height: 30px;
    z-index: 5;
    pointer-events: none;
  }
  .log-list-frame::before {
    top: 0;
    background: linear-gradient(180deg, rgba(255, 246, 250, 0.98), rgba(255, 246, 250, 0));
  }
  .log-list-frame::after {
    bottom: 0;
    background: linear-gradient(0deg, rgba(255, 244, 249, 0.98), rgba(255, 244, 249, 0));
  }
  .log-list {
    position: relative;
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    gap: 8px;
    height: auto;
    max-height: none;
    min-height: 0;
    overflow-x: hidden;
    overflow-y: auto;
    padding: 10px 8px 12px;
    overscroll-behavior: contain;
    -webkit-overflow-scrolling: touch;
    scrollbar-gutter: stable;
    scroll-behavior: auto;
    scroll-padding: 18px;
    overflow-anchor: none;
    contain: layout paint;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .log-list::-webkit-scrollbar,
  .detail-scroll::-webkit-scrollbar {
    width: 9px;
  }
  .log-list::-webkit-scrollbar-track,
  .detail-scroll::-webkit-scrollbar-track {
    border-radius: 999px;
    background: linear-gradient(180deg, rgba(255, 241, 247, 0.86), rgba(255, 223, 235, 0.68));
  }
  .log-list::-webkit-scrollbar-thumb,
  .detail-scroll::-webkit-scrollbar-thumb {
    border-radius: 999px;
    background: linear-gradient(180deg, rgba(255, 174, 205, 0.96), rgba(239, 121, 170, 0.96));
    border: 2px solid rgba(255, 244, 249, 0.92);
    box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.32);
  }
  .log-list::-webkit-scrollbar-thumb:hover,
  .detail-scroll::-webkit-scrollbar-thumb:hover,
  .speech-card pre::-webkit-scrollbar-thumb:hover,
  .modal-table-wrap::-webkit-scrollbar-thumb:hover,
  .suggest-menu::-webkit-scrollbar-thumb:hover,
  .gallery-grid::-webkit-scrollbar-thumb:hover,
  .data-fold pre::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(180deg, #ff8fbd, #e65395);
  }
  .lane-scroll::-webkit-scrollbar,
  .detail-scroll::-webkit-scrollbar {
    width: 10px;
  }
  .lane-scroll::-webkit-scrollbar-track,
  .detail-scroll::-webkit-scrollbar-track {
    border-radius: 999px;
    background: rgba(255, 228, 238, 0.62);
  }
  .lane-scroll::-webkit-scrollbar-thumb,
  .detail-scroll::-webkit-scrollbar-thumb {
    border-radius: 999px;
    background: linear-gradient(180deg, rgba(255, 159, 195, 0.9), rgba(239, 121, 170, 0.9));
    border: 2px solid rgba(255, 244, 249, 0.92);
  }
  .detail-scroll-wrap {
    min-height: 0;
    position: relative;
    z-index: 1;
  }
  .detail-scroll {
    height: 100%;
    padding: 6px 6px 6px 0;
    overflow-x: hidden;
    overflow-y: auto;
    overscroll-behavior: contain;
    scroll-behavior: smooth;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .detail-scroll.is-swapping {
    animation: detail-pane-swap 0.22s ease-out;
  }
  .log-cluster {
    border-radius: 22px;
    border: 1px solid rgba(235, 191, 206, 0.92);
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.76), rgba(255, 245, 249, 0.7));
    backdrop-filter: blur(20px) saturate(1.05);
    box-shadow: 0 14px 32px rgba(214, 136, 169, 0.12);
    overflow: hidden;
    transition: box-shadow 0.18s ease, border-color 0.18s ease, background 0.18s ease;
  }
  .log-cluster:hover {
    box-shadow: 0 18px 38px rgba(214, 136, 169, 0.16);
  }
  .log-cluster.is-selected {
    border-color: rgba(221, 113, 153, 0.94);
    background: linear-gradient(180deg, rgba(255, 241, 247, 0.94), rgba(255, 248, 251, 0.84));
    box-shadow: 0 0 0 1px rgba(221, 113, 153, 0.12), 0 18px 34px rgba(221, 113, 153, 0.12);
  }
  .log-cluster.is-fresh {
    animation: cluster-appear 0.34s ease-out both, fresh-halo 1.2s ease-out both;
  }
  .logs-screen .cluster-main {
    display: block;
    width: 100%;
    border: 0;
    background: transparent;
    text-align: left;
    padding: 14px 14px 12px;
    cursor: pointer;
    color: inherit;
    box-shadow: none;
    transform: none;
    white-space: normal;
  }
  .logs-screen .cluster-main:hover,
  .logs-screen .cluster-main.active,
  .logs-screen .cluster-child:hover,
  .logs-screen .cluster-child.active,
  .logs-screen .cluster-group-toggle:hover {
    transform: none;
  }
  .logs-screen .cluster-main.active {
    background: linear-gradient(180deg, rgba(255, 233, 242, 0.56), rgba(255, 255, 255, 0.14));
  }
  .cluster-topline {
    display: flex;
    justify-content: space-between;
    align-items: start;
    gap: 12px;
    margin-bottom: 10px;
  }
  .cluster-user {
    min-width: 0;
  }
  .cluster-user strong {
    display: block;
    font-size: 16px;
    line-height: 1.25;
    color: #6f2f49;
    word-break: break-word;
  }
  .cluster-user code {
    display: inline-block;
    margin-top: 5px;
    color: #a0687e;
    font-size: 11px;
  }
  .cluster-time {
    white-space: nowrap;
    font-size: 12px;
    color: #9c7585;
  }
  .cluster-preview {
    font-size: 13px;
    line-height: 1.58;
    color: #704357;
    min-height: 40px;
  }
  .cluster-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 12px;
  }
  .cluster-badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 9px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.66);
    border: 1px solid rgba(236, 192, 208, 0.72);
    color: #906373;
    font-size: 11px;
    font-weight: 700;
  }
  .cluster-badge.status-ok {
    color: #8f4662;
  }
  .cluster-badge.status-err {
    color: #c04468;
    background: rgba(255, 232, 238, 0.92);
  }
  .logs-screen .cluster-group-toggle {
    display: block;
    width: 100%;
    border: 0;
    border-top: 1px solid rgba(236, 192, 208, 0.56);
    background: rgba(255, 244, 248, 0.78);
    color: #a44970;
    font-weight: 700;
    font-size: 13px;
    padding: 11px 14px;
    box-shadow: none;
    transform: none;
    white-space: normal;
  }
  .logs-screen .cluster-group-toggle:hover {
    transform: none;
    box-shadow: none;
    background: rgba(255, 235, 243, 0.92);
  }
  .cluster-children {
    padding: 8px 8px 10px;
    display: grid;
    gap: 8px;
    border-top: 1px solid rgba(236, 192, 208, 0.42);
    background: linear-gradient(180deg, rgba(255, 246, 250, 0.58), rgba(255, 250, 252, 0.88));
  }
  .logs-screen .cluster-child {
    display: block;
    width: 100%;
    border: 1px solid rgba(236, 192, 208, 0.62);
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.7);
    text-align: left;
    padding: 10px 12px;
    color: inherit;
    box-shadow: none;
    transform: none;
    white-space: normal;
  }
  .logs-screen .cluster-child.active {
    border-color: rgba(221, 113, 153, 0.88);
    background: rgba(255, 233, 241, 0.82);
  }
  .logs-screen .cluster-child:hover {
    transform: none;
    background: rgba(255, 239, 245, 0.88);
  }
  .cluster-child-row {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    align-items: center;
    margin-bottom: 6px;
  }
  .cluster-child-row strong {
    color: #7a3551;
    font-size: 13px;
  }
  .cluster-child-preview {
    font-size: 12px;
    color: #906476;
    line-height: 1.5;
  }
  .log-row {
    flex: 0 0 auto;
    border-radius: 16px;
    border: 1px solid rgba(236, 192, 208, 0.78);
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.78), rgba(255, 246, 250, 0.72));
    overflow: hidden;
    box-shadow: 0 8px 18px rgba(214, 136, 169, 0.08);
    transform: none;
    transition: border-color 0.16s ease, background 0.16s ease, box-shadow 0.16s ease;
  }
  .log-row:hover {
    border-color: rgba(232, 135, 174, 0.86);
    background: rgba(255, 247, 250, 0.9);
    box-shadow: 0 10px 22px rgba(214, 136, 169, 0.11);
  }
  .log-row.is-selected {
    border-color: rgba(221, 113, 153, 0.96);
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.8), rgba(255, 246, 250, 0.74));
    box-shadow: 0 0 0 1px rgba(221, 113, 153, 0.1), 0 10px 20px rgba(221, 113, 153, 0.09);
  }
  .log-row.is-fresh {
    animation: cluster-appear 0.42s cubic-bezier(.18,.86,.26,1) both, fresh-halo 1.15s ease-out both;
  }
  .logs-screen .log-row-main,
  .logs-screen .log-row-child,
  .logs-screen .log-row-toggle {
    display: block;
    width: 100%;
    border: 0;
    border-radius: 0;
    background: transparent;
    color: inherit;
    box-shadow: none;
    text-align: left;
    position: relative;
    transform: none;
    white-space: normal;
    transition: background 0.18s ease, color 0.18s ease;
  }
  .logs-screen .log-row-main {
    padding: 10px 11px;
    border-radius: 14px;
  }
  .logs-screen .log-row-main:hover,
  .logs-screen .log-row-main.active,
  .logs-screen .log-row-child:hover,
  .logs-screen .log-row-child.active,
  .logs-screen .log-row-toggle:hover {
    transform: none;
    box-shadow: none;
    background: rgba(255, 238, 246, 0.72);
  }
  .logs-screen .log-row-main.active,
  .logs-screen .log-row-child.active {
    background: rgba(255, 226, 239, 0.92);
  }
  .logs-screen .log-row-main.active::after,
  .logs-screen .log-row-child.active::after {
    content: "";
    position: absolute;
    inset: 0;
    border: 2px solid rgba(120, 18, 46, 0.96);
    border-radius: inherit;
    box-sizing: border-box;
    pointer-events: none;
  }
  .log-row-line {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 10px;
    min-width: 0;
  }
  .log-row-title {
    min-width: 0;
  }
  .log-row-title strong {
    display: block;
    color: #6f2f49;
    font-size: 14px;
    line-height: 1.25;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .log-row-title code {
    display: block;
    margin-top: 3px;
    color: #a06b7f;
    font-size: 10px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .log-row-time {
    flex: 0 0 auto;
    color: #9c7585;
    font-size: 11px;
    white-space: nowrap;
  }
  .log-row-preview {
    margin-top: 7px;
    color: #704357;
    font-size: 12px;
    line-height: 1.45;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .log-row-preview.small {
    -webkit-line-clamp: 1;
    font-size: 11px;
  }
  .log-row-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 5px;
    margin-top: 8px;
  }
  .log-row-badge {
    display: inline-flex;
    align-items: center;
    max-width: 100%;
    padding: 3px 7px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.72);
    border: 1px solid rgba(236, 192, 208, 0.7);
    color: #8f6073;
    font-size: 10px;
    font-weight: 800;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .log-row-badge.status-ok {
    color: #8f4662;
  }
  .log-row-badge.status-err {
    color: #c04468;
    background: rgba(255, 232, 238, 0.92);
  }
  .logs-screen .log-row-toggle {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    border-top: 1px solid rgba(236, 192, 208, 0.52);
    padding: 7px 11px;
    color: #a44970;
    font-size: 12px;
    font-weight: 800;
  }
  .logs-screen .log-row-toggle::after {
    content: "";
    flex: 0 0 auto;
    width: 9px;
    height: 9px;
    margin-right: 2px;
    border-right: 2px solid currentColor;
    border-bottom: 2px solid currentColor;
    transform: rotate(45deg);
    transform-origin: 50% 50%;
    opacity: 0.78;
    transition: transform 0.24s ease, opacity 0.24s ease;
  }
  .logs-screen .log-row-toggle[aria-expanded="true"]::after {
    transform: rotate(225deg) translateY(-1px);
    opacity: 1;
  }
  .log-row-children {
    display: grid;
    gap: 6px;
    padding: 7px;
    border-top: 1px solid rgba(236, 192, 208, 0.42);
    background: rgba(255, 248, 251, 0.72);
    transform-origin: top center;
    animation: log-children-reveal 0.26s cubic-bezier(.2,.85,.22,1) both;
    overflow: clip;
  }
  .logs-screen .log-row-child {
    padding: 8px 9px;
    border-radius: 12px;
    border: 1px solid rgba(236, 192, 208, 0.58);
    background: rgba(255, 255, 255, 0.72);
    animation: log-child-rise 0.3s cubic-bezier(.22,.84,.28,1) both;
    animation-delay: calc(var(--child-index, 0) * 0.035s);
  }
  .log-row-child strong {
    color: #7a3551;
    font-size: 12px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .lane-empty,
  .detail-empty {
    min-height: 100%;
    display: grid;
    place-items: center;
    text-align: center;
    padding: 30px 20px;
    border-radius: 24px;
    border: 1px dashed rgba(228, 183, 202, 0.82);
    background: rgba(255, 252, 253, 0.62);
    color: #8f6878;
  }
  .lane-empty h3,
  .detail-empty h2 {
    margin: 0 0 8px;
    color: #7f3652;
  }
  .lane-foot {
    flex: 0 0 auto;
    position: relative;
    z-index: 20;
    isolation: isolate;
    pointer-events: auto;
    display: flex;
    align-items: center;
    gap: 10px;
    border-top: 1px solid rgba(236, 192, 208, 0.68);
    padding-top: 12px;
    margin-top: 0;
  }
  .lane-foot span {
    flex: 1;
    text-align: center;
    color: #936c7e;
    font-weight: 700;
  }
  .lane-foot .button,
  .lane-foot button {
    position: relative;
    z-index: 1;
    pointer-events: auto;
  }
  .lane-foot .button[aria-disabled="true"] {
    opacity: 0.45;
    cursor: not-allowed;
    pointer-events: none;
  }
  .detail-toolbar {
    position: relative;
    z-index: 100;
    overflow: visible;
    isolation: isolate;
    padding: 0 0 8px;
    border-bottom: 1px solid rgba(236, 192, 208, 0.72);
    margin-bottom: 2px;
  }
  .monitor-filters {
    position: relative;
    z-index: 25;
    overflow: visible;
    isolation: isolate;
    display: grid;
    grid-template-columns: minmax(150px, 1.05fr) minmax(150px, 1.05fr) minmax(140px, 1fr) minmax(140px, 1fr) minmax(118px, 0.85fr) minmax(180px, 1.35fr) auto;
    gap: 8px;
    align-items: end;
  }
  .filter-field {
    position: relative;
    z-index: 1;
    overflow: visible;
  }
  .filter-field.is-open,
  .filter-field:focus-within {
    z-index: 80;
  }
  .filter-field label {
    color: #8e6678;
    margin-bottom: 4px;
    font-size: 12px;
  }
  .filter-field input {
    background: rgba(255, 255, 255, 0.82);
    border-color: rgba(236, 192, 208, 0.88);
    color: #6a2942;
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72);
    min-height: 34px;
    padding: 7px 10px;
  }
  .filter-field input:focus {
    outline: none;
    border-color: rgba(223, 112, 153, 0.86);
    box-shadow: 0 0 0 3px rgba(223, 112, 153, 0.12);
  }
  .suggest-shell {
    position: relative;
    z-index: 1;
    overflow: visible;
  }
  .suggest-shell.is-open,
  .suggest-shell:focus-within {
    z-index: 90;
  }
  .suggest-shell input {
    padding-right: 42px;
  }
  .logs-screen .suggest-trigger {
    position: absolute;
    top: 50%;
    right: 6px;
    transform: translateY(-50%);
    width: 30px;
    height: 30px;
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.78);
    color: #9b3f63;
    border: 1px solid rgba(231, 145, 181, 0.44);
    box-shadow: 0 4px 10px rgba(221, 105, 154, 0.08);
    padding: 0;
    font-size: 16px;
    line-height: 1;
  }
  .logs-screen .suggest-trigger:hover {
    transform: translateY(-50%);
    box-shadow: 0 6px 14px rgba(221, 105, 154, 0.12);
    background: rgba(255, 247, 251, 0.94);
  }
  .logs-screen .suggest-trigger:disabled {
    opacity: 0.38;
    cursor: default;
    box-shadow: none;
    transform: translateY(-50%);
  }
  .suggest-menu {
    position: absolute;
    left: 0;
    right: 0;
    top: calc(100% + 6px);
    background: rgba(255, 248, 251, 0.96);
    border: 1px solid rgba(236, 192, 208, 0.92);
    border-radius: 16px;
    box-shadow: 0 22px 44px rgba(194, 113, 145, 0.18);
    backdrop-filter: blur(22px);
    max-height: 260px;
    overflow: auto;
    display: none;
    z-index: 300;
    padding: 6px;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .suggest-menu.show {
    display: block;
  }
  .logs-screen .suggest-option {
    display: block;
    width: 100%;
    border: 0;
    background: transparent;
    color: #7c4058;
    text-align: left;
    padding: 10px 12px;
    border-radius: 12px;
    box-shadow: none;
    transition: background 0.14s ease, color 0.14s ease;
  }
  .logs-screen .suggest-option:hover {
    transform: none;
    box-shadow: none;
    background: rgba(255, 234, 241, 0.94);
  }
  .filter-actions {
    position: relative;
    z-index: 2;
    display: flex;
    gap: 6px;
    align-items: center;
    align-self: end;
    justify-content: flex-end;
    white-space: nowrap;
  }
  .filter-actions button {
    min-height: 34px;
    padding: 8px 12px;
  }
  .detail-shell {
    display: flex;
    flex-direction: column;
    gap: 8px;
    height: 100%;
    min-height: 0;
  }
  .detail-compact {
    flex: 0 0 auto;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: center;
    gap: 10px;
    padding: 9px 10px;
    border-radius: 18px;
    border: 1px solid rgba(236, 192, 208, 0.78);
    background:
      radial-gradient(circle at 8% 0%, rgba(255, 224, 236, 0.62), transparent 34%),
      rgba(255, 255, 255, 0.66);
    backdrop-filter: blur(18px) saturate(1.05);
    box-shadow: 0 10px 22px rgba(206, 131, 161, 0.08);
  }
  .detail-compact-line {
    display: flex;
    align-items: center;
    gap: 7px;
    min-width: 0;
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;
    scrollbar-width: none;
  }
  .detail-compact-line::-webkit-scrollbar {
    display: none;
  }
  .detail-compact-path {
    flex: 0 1 240px;
    min-width: 96px;
    color: #6e2a45;
    font-size: 15px;
    font-weight: 900;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .detail-compact-chip {
    flex: 0 0 auto;
    display: inline-flex;
    align-items: center;
    max-width: 210px;
    padding: 4px 8px;
    border-radius: 999px;
    border: 1px solid rgba(236, 192, 208, 0.72);
    background: rgba(255, 243, 248, 0.86);
    color: #86566b;
    font-size: 11px;
    font-weight: 800;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: background 0.16s ease, border-color 0.16s ease;
  }
  .detail-compact-chip:hover {
    background: rgba(255, 235, 244, 0.94);
    border-color: rgba(232, 135, 174, 0.82);
  }
  .detail-compact-chip.status-err {
    color: #c24b67;
    background: rgba(255, 231, 238, 0.96);
  }
  .detail-compact .button.secondary {
    padding: 7px 10px;
    border-radius: 12px;
    font-size: 12px;
    line-height: 1;
  }
  .detail-head {
    display: flex;
    justify-content: space-between;
    align-items: start;
    gap: 16px;
    flex-wrap: wrap;
  }
  .detail-head h2 {
    margin: 0 0 8px;
    color: #6e2a45;
    font-size: 28px;
    line-height: 1.2;
  }
  .detail-head p {
    margin: 0;
    color: #916a7a;
    font-size: 14px;
    line-height: 1.55;
  }
  .detail-head .button.secondary {
    color: #9b3f63;
    border-color: rgba(231, 145, 181, 0.48);
    background: rgba(255, 255, 255, 0.76);
  }
  .meta-ribbon {
    flex: 0 0 auto;
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  .meta-pill-soft {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 7px 12px;
    border-radius: 999px;
    background: rgba(255, 242, 247, 0.92);
    border: 1px solid rgba(236, 192, 208, 0.82);
    color: #8e566d;
    font-size: 12px;
    font-weight: 700;
  }
  .meta-pill-soft.status-err {
    color: #c24b67;
    background: rgba(255, 231, 238, 0.96);
  }
  .detail-info-grid {
    flex: 0 0 auto;
    display: grid;
    grid-template-columns: repeat(6, minmax(0, 1fr));
    gap: 10px;
  }
  .detail-info-chip {
    padding: 14px 15px;
    border-radius: 18px;
    background: rgba(255, 255, 255, 0.72);
    border: 1px solid rgba(236, 192, 208, 0.82);
    backdrop-filter: blur(16px);
  }
  .detail-info-chip span {
    display: block;
    color: #987082;
    font-size: 12px;
    margin-bottom: 6px;
  }
  .detail-info-chip strong {
    display: block;
    color: #6e2d46;
    font-size: 18px;
    line-height: 1.28;
    word-break: break-word;
  }
  #detail-stage {
    height: 100%;
    min-height: 0;
  }
  .conversation-stack {
    flex: 1 1 auto;
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    gap: 14px;
    height: 100%;
    min-height: 0;
    align-items: stretch;
  }
  .conversation-column {
    min-width: 0;
    min-height: 0;
  }
  .speech-card {
    display: flex;
    flex-direction: column;
    padding: 12px 13px 11px;
    border-radius: 26px;
    border: 1px solid rgba(236, 192, 208, 0.86);
    backdrop-filter: blur(18px);
    box-shadow: 0 16px 34px rgba(206, 131, 161, 0.08);
    overflow: hidden;
    position: relative;
    height: 100%;
    min-height: 0;
    transition: border-color 0.18s ease, box-shadow 0.18s ease;
  }
  .speech-card:hover {
    border-color: rgba(232, 135, 174, 0.92);
    box-shadow: 0 18px 34px rgba(206, 131, 161, 0.11);
  }
  .speech-card::before {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.24), transparent 42%);
    opacity: 0.95;
  }
  .speech-card > * {
    position: relative;
    z-index: 1;
  }
  .speech-card.user {
    background: linear-gradient(180deg, rgba(255, 225, 235, 0.92), rgba(255, 248, 251, 0.88));
  }
  .speech-card.assistant {
    background: linear-gradient(180deg, rgba(255, 239, 245, 0.92), rgba(255, 252, 253, 0.92));
  }
  .speech-card h3 {
    margin: 0;
    font-size: 15px;
    color: #6d2a43;
  }
  .speech-card pre {
    background: transparent;
    border: 0;
    padding: 0;
    margin: 0;
    min-height: 0;
    flex: 1 1 auto;
    overflow: auto;
    padding-right: 6px;
    font-size: 15px;
    line-height: 1.75;
    color: #5d233b;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .speech-card pre::-webkit-scrollbar {
    width: 8px;
  }
  .speech-card pre::-webkit-scrollbar-thumb {
    border-radius: 999px;
    background: linear-gradient(180deg, rgba(255, 174, 205, 0.86), rgba(239, 121, 170, 0.86));
  }
  .text-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
    margin-bottom: 7px;
    flex-wrap: wrap;
  }
  .text-toolbar .muted {
    color: #977082;
    font-size: 12px;
  }
  .mini-pagination {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 14px;
  }
  .mini-pagination span {
    min-width: 96px;
    text-align: center;
    color: #8f6779;
    font-weight: 700;
  }
  .gallery-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 14px;
    margin-bottom: 14px;
    max-height: 220px;
    overflow: auto;
    padding-right: 6px;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .gallery-card {
    padding: 12px;
    border-radius: 22px;
    background: rgba(255, 255, 255, 0.76);
    border: 1px solid rgba(236, 192, 208, 0.84);
    box-shadow: 0 16px 30px rgba(211, 139, 168, 0.08);
  }
  .gallery-card img {
    width: 100%;
    display: block;
    border-radius: 18px;
    background: linear-gradient(180deg, rgba(255, 238, 244, 0.9), rgba(255, 250, 252, 0.94));
  }
  .gallery-card p {
    margin: 10px 0 0;
    font-size: 12px;
    color: #916b7d;
  }
  .secondary-stack {
    flex: 0 0 auto;
    display: grid;
    gap: 14px;
  }
  .data-fold {
    border-radius: 20px;
    border: 1px solid rgba(236, 192, 208, 0.84);
    background: rgba(255, 255, 255, 0.7);
    padding: 14px 16px;
  }
  .data-fold summary {
    color: #74334d;
    font-weight: 800;
    cursor: pointer;
  }
  .data-fold pre {
    background: rgba(255, 250, 252, 0.86);
    border-color: rgba(236, 192, 208, 0.84);
    color: #643148;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .pre-grid {
    display: grid;
    gap: 12px;
    margin-top: 12px;
  }
  .pre-block h4 {
    margin: 0 0 8px;
    font-size: 13px;
    color: #956d7f;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }
  .raw-placeholder {
    margin-top: 12px;
    padding: 18px;
    border-radius: 18px;
    border: 1px dashed rgba(228, 183, 202, 0.82);
    background: rgba(255, 250, 252, 0.62);
    color: #936d80;
  }
  .raw-placeholder p {
    margin: 0 0 10px;
    line-height: 1.55;
  }
  .message-strip {
    padding: 12px 14px;
    border-radius: 16px;
    font-size: 13px;
    line-height: 1.55;
    border: 1px solid rgba(236, 192, 208, 0.82);
    background: rgba(255, 250, 252, 0.76);
    color: #8a6074;
  }
  .message-strip.hidden {
    display: none;
  }
  .message-strip.success {
    background: rgba(255, 234, 241, 0.9);
    color: #a1476a;
  }
  .message-strip.error {
    background: rgba(255, 230, 236, 0.92);
    color: #be4966;
  }
  .modal-backdrop.monitor-modal {
    background:
      radial-gradient(circle at top, rgba(255, 228, 237, 0.26), transparent 28%),
      rgba(97, 56, 72, 0.34);
    backdrop-filter: blur(18px) saturate(1.08);
  }
  .monitor-modal .modal-card {
    width: min(1500px, calc(100vw - 28px));
    background: linear-gradient(180deg, rgba(255, 250, 252, 0.94), rgba(255, 241, 247, 0.92));
    border-color: rgba(236, 192, 208, 0.86);
    box-shadow: 0 24px 70px rgba(182, 107, 138, 0.18);
    backdrop-filter: blur(30px) saturate(1.08);
  }
  .monitor-modal .modal-header h2 {
    color: #6d2d46;
  }
  .monitor-modal .modal-header p {
    color: #926d7e;
  }
  .monitor-modal .modal-close {
    background: rgba(255, 255, 255, 0.78);
    color: #9b3f63;
    border-color: rgba(231, 145, 181, 0.48);
    box-shadow:
      0 7px 14px rgba(221, 105, 154, 0.1),
      inset 0 1px 0 rgba(255, 255, 255, 0.82);
  }
  .monitor-modal .modal-close:hover {
    background: rgba(255, 247, 251, 0.94);
    color: #873354;
    border-color: rgba(222, 116, 159, 0.58);
  }
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(6, minmax(0, 1fr));
    gap: 12px;
  }
  .stats-card {
    padding: 16px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.72);
    border: 1px solid rgba(236, 192, 208, 0.84);
    box-shadow: 0 12px 30px rgba(205, 118, 151, 0.08);
  }
  .stats-card span {
    display: block;
    color: #9b7485;
    font-size: 12px;
    margin-bottom: 8px;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    font-weight: 700;
  }
  .stats-card strong {
    display: block;
    font-size: 28px;
    color: #6f2e49;
    line-height: 1.1;
  }
  .modal-two-col {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    gap: 16px;
  }
  .modal-panel {
    padding: 16px;
    border-radius: 22px;
    border: 1px solid rgba(236, 192, 208, 0.82);
    background: rgba(255, 255, 255, 0.7);
  }
  .modal-panel h3 {
    margin: 0 0 12px;
    color: #71314b;
  }
  .modal-table-wrap {
    max-height: 430px;
    overflow: auto;
    border-radius: 18px;
    scrollbar-color: var(--rose-scroll) var(--rose-scroll-track);
    scrollbar-width: thin;
  }
  .modal-table-wrap table {
    background: rgba(255, 255, 255, 0.5);
  }
  .modal-table-wrap thead th {
    position: sticky;
    top: 0;
    background: rgba(255, 245, 249, 0.95);
    z-index: 1;
  }
  .token-manage-grid {
    display: grid;
    grid-template-columns: minmax(240px, 1.5fr) minmax(180px, 1fr) auto;
    gap: 12px;
    align-items: end;
  }
  .token-row-input {
    width: 100%;
    min-width: 0;
  }
  .row-actions {
    display: flex;
    gap: 8px;
    align-items: center;
  }
  .row-actions .button,
  .row-actions button {
    padding: 8px 12px;
    box-shadow: none;
  }
  .db-stats-grid {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 12px;
  }
  .db-actions {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    margin-top: 14px;
  }
  .db-actions button {
    box-shadow: 0 10px 22px rgba(205, 118, 151, 0.12);
  }
  .cleanup-grid {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 12px;
    align-items: end;
  }
  .cleanup-grid .filter-actions {
    grid-column: span 5;
    justify-content: flex-start;
  }
  .subtle {
    color: #977082;
    font-size: 12px;
    line-height: 1.55;
  }
  .cell-note {
    margin-top: 5px;
    color: #9a7485;
    font-size: 12px;
    line-height: 1.45;
  }
  .mono-box {
    font-family: Consolas, "Courier New", monospace;
    font-size: 12px;
    word-break: break-all;
  }
  .detail-error-banner {
    padding: 14px 16px;
    border-radius: 20px;
    border: 1px solid rgba(213, 109, 139, 0.24);
    background: rgba(255, 232, 238, 0.86);
    color: #bf4d69;
  }
  .detail-error-banner strong {
    display: block;
    margin-bottom: 6px;
  }
  .loading-glass {
    position: relative;
    overflow: hidden;
  }
  .loading-glass::after {
    content: "";
    position: absolute;
    inset: 0;
    background: linear-gradient(115deg, transparent 18%, rgba(255, 255, 255, 0.46) 48%, transparent 80%);
    transform: translateX(-120%);
    animation: shimmer 1.3s ease infinite;
    pointer-events: none;
  }
  @keyframes shimmer {
    to {
      transform: translateX(120%);
    }
  }
  @keyframes cluster-appear {
    0% {
      opacity: 0;
      transform: translateY(-12px) scale(0.985);
    }
    58% {
      opacity: 1;
      transform: translateY(2px) scale(1);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
    }
  }
  @keyframes fresh-halo {
    0% {
      box-shadow: 0 0 0 rgba(221, 113, 153, 0);
    }
    38% {
      box-shadow: 0 0 0 4px rgba(239, 121, 170, 0.18), 0 16px 30px rgba(221, 113, 153, 0.15);
    }
    100% {
      box-shadow: 0 12px 28px rgba(208, 136, 165, 0.09);
    }
  }
  @keyframes rose-count-pop {
    0% {
      transform: scale(0.98);
      color: #b84d74;
    }
    50% {
      transform: scale(1.035);
      color: #8d3658;
    }
    100% {
      transform: scale(1);
      color: #7d334f;
    }
  }
  @keyframes detail-pane-swap {
    0% {
      opacity: 0.64;
      transform: translateY(4px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
    }
  }
  @keyframes log-children-reveal {
    0% {
      opacity: 0;
      transform: translateY(-8px) scaleY(0.96);
    }
    100% {
      opacity: 1;
      transform: translateY(0) scaleY(1);
    }
  }
  @keyframes log-child-rise {
    0% {
      opacity: 0;
      transform: translateY(-7px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
    }
  }
  @media (prefers-reduced-motion: reduce) {
    body,
    body::before,
    body::after,
    .logs-screen .button,
    .logs-screen button,
    .monitor-modal .button,
    .monitor-modal button,
    .topbar-stat,
    .log-row,
    .log-row.is-fresh,
    .log-row-children,
    .log-row-child,
    .detail-scroll.is-swapping,
    .speech-card,
    .loading-glass::after {
      animation: none !important;
      transition-duration: 0.01ms !important;
    }
    .log-list,
    .detail-scroll,
    .speech-card pre {
      scroll-behavior: auto;
    }
  }
  @media (max-width: 1500px) {
    .detail-info-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .stats-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .db-stats-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .conversation-stack {
      min-height: 0;
    }
  }
  @media (max-width: 1180px) {
    .monitor-grid {
      grid-template-columns: 1fr;
      grid-template-rows: minmax(250px, 34vh) minmax(0, 1fr);
    }
    .monitor-topbar {
      grid-template-columns: minmax(0, 1fr) auto minmax(0, auto);
    }
    .topbar-overview,
    .topbar-status,
    .topbar-actions {
      justify-self: auto;
    }
    .topbar-overview {
      justify-content: flex-start;
    }
    .topbar-status {
      justify-content: center;
    }
    .monitor-filters {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .filter-actions {
      grid-column: span 3;
      justify-content: flex-start;
    }
    .modal-two-col {
      grid-template-columns: 1fr;
    }
    .cleanup-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
    .cleanup-grid .filter-actions {
      grid-column: span 2;
    }
    .conversation-stack {
      grid-template-columns: 1fr;
      min-height: 0;
      overflow: auto;
      padding-right: 4px;
    }
    .speech-card {
      height: 38vh;
      min-height: 240px;
    }
  }
  @media (max-width: 820px) {
    .logs-screen {
      padding: 6px 8px 8px;
    }
    .topbar-actions {
      width: auto;
      justify-content: flex-start;
    }
    .monitor-filters {
      grid-template-columns: 1fr 1fr;
    }
    .filter-actions {
      grid-column: span 2;
    }
    .detail-info-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
    .stats-grid,
    .db-stats-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
    .token-manage-grid {
      grid-template-columns: 1fr;
    }
    .speech-card {
      height: 36vh;
    }
  }
  @media (max-width: 640px) {
    .monitor-topbar {
      align-items: flex-start;
    }
    .monitor-filters,
    .cleanup-grid {
      grid-template-columns: 1fr;
    }
    .filter-actions,
    .cleanup-grid .filter-actions {
      grid-column: span 1;
    }
    .detail-info-grid,
    .stats-grid,
    .db-stats-grid {
      grid-template-columns: 1fr;
    }
    .lane-foot {
      flex-direction: column;
    }
  }
</style>

<div
  class="shell shell-wide logs-screen"
  data-log-page="{{.Page}}"
  data-log-page-size="{{.PageSize}}"
  data-log-current-count="{{.CurrentCount}}"
  data-log-total-count="{{.TotalCount}}"
  data-log-total-pages="{{.TotalPages}}"
  data-log-has-prev="{{if .HasPrev}}true{{else}}false{{end}}"
  data-log-has-next="{{if .HasNext}}true{{else}}false{{end}}">
  <div class="monitor-shell">
    <div class="nav monitor-topbar">
      <div class="topbar-overview">
        <div class="topbar-stat" id="quick-requests-card">
          <span>总请求</span>
          <strong id="quick-total-requests">-</strong>
        </div>
        <div class="topbar-stat" id="quick-today-card">
          <span>今日</span>
          <strong id="quick-today-requests">-</strong>
        </div>
        <div class="topbar-stat" id="quick-errors-card">
          <span>错误</span>
          <strong id="quick-error-count">-</strong>
        </div>
        <div class="topbar-stat" id="quick-tokens-card">
          <span>总 Tokens</span>
          <strong id="quick-total-tokens">-</strong>
        </div>
      </div>
      <div class="topbar-status">
        <div class="pill sync-pill" id="logs-refresh-state">实时审计在线</div>
      </div>
      <div class="topbar-actions">
        <button class="button secondary" type="button" id="open-stats-modal">统计</button>
        <button class="button secondary" type="button" id="open-settings-modal">设置</button>
        <button class="button secondary" type="button" id="open-db-modal">数据库</button>
        <form method="post" action="{{path "/logout"}}" style="margin: 0;">
          <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
          <button class="button secondary" type="submit">退出登录</button>
        </form>
      </div>
    </div>

    <div class="monitor-grid">
      <aside class="panel log-lane">
        <div class="lane-head">
          <div>
            <h1>日志流</h1>
            <p>左侧按连续同一令牌自动折叠，默认每页 100 条，只在检测到变化时增量刷新。</p>
            <div class="list-summary">
              <div class="summary-chip">当前页 <strong id="list-current-count">{{.CurrentCount}}</strong></div>
              <div class="summary-chip">总匹配 <strong id="list-total-count">{{.TotalCount}}</strong></div>
            </div>
          </div>
          <div class="pill">列表</div>
        </div>
        <div class="log-list-frame">
          <div class="log-list" id="log-groups" role="list">
            <div class="lane-empty">
              <div>
                <h3>正在加载日志</h3>
                <p>稍候会自动选中最新一条，并在右侧展示完整内容。</p>
              </div>
            </div>
          </div>
        </div>
        <div class="lane-foot">
          <a class="button secondary" id="prev-page" href="{{pageURL .Filters .PrevPage}}" aria-disabled="{{if .HasPrev}}false{{else}}true{{end}}">上一页</a>
          <span id="list-page-text">第 {{.Page}} / {{.TotalPages}} 页</span>
          <a class="button secondary" id="next-page" href="{{pageURL .Filters .NextPage}}" aria-disabled="{{if .HasNext}}false{{else}}true{{end}}">下一页</a>
        </div>
      </aside>

      <section class="panel detail-lane">
        <div class="detail-toolbar">
          <form id="filters-form" class="monitor-filters" method="get" action="{{path "/logs"}}">
            <div class="filter-field">
              <label for="filter-from">开始时间</label>
              <input id="filter-from" name="from" type="datetime-local" value="{{.Filters.From}}">
            </div>
            <div class="filter-field">
              <label for="filter-to">结束时间</label>
              <input id="filter-to" name="to" type="datetime-local" value="{{.Filters.To}}">
            </div>
            <div class="filter-field">
              <label for="filter-alias">令牌代号</label>
              <div class="suggest-shell">
                <input id="filter-alias" name="alias" type="text" autocomplete="off" value="{{.Filters.TokenAlias}}" data-suggest="token_aliases" placeholder="支持片段匹配">
                <button class="suggest-trigger" type="button" data-clear-filter="filter-alias" title="清空令牌代号" aria-label="清空令牌代号">×</button>
                <div class="suggest-menu" id="suggest-filter-alias"></div>
              </div>
            </div>
            <div class="filter-field">
              <label for="filter-model">模型</label>
              <div class="suggest-shell">
                <input id="filter-model" name="model" type="text" autocomplete="off" value="{{.Filters.Model}}" data-suggest="models" placeholder="支持片段匹配">
                <button class="suggest-trigger" type="button" data-clear-filter="filter-model" title="清空模型" aria-label="清空模型">×</button>
                <div class="suggest-menu" id="suggest-filter-model"></div>
              </div>
            </div>
            <div class="filter-field">
              <label for="filter-status">状态码</label>
              <div class="suggest-shell">
                <input id="filter-status" name="status" type="text" autocomplete="off" value="{{.Filters.StatusCode}}" data-suggest="status_codes" placeholder="如 200 / 4xx">
                <button class="suggest-trigger" type="button" data-clear-filter="filter-status" title="清空状态码" aria-label="清空状态码">×</button>
                <div class="suggest-menu" id="suggest-filter-status"></div>
              </div>
            </div>
            <div class="filter-field">
              <label for="filter-keyword">关键词</label>
              <input id="filter-keyword" name="q" type="text" autocomplete="off" value="{{.Filters.Keyword}}" placeholder="路径、文本、模型、代号均可匹配">
            </div>
            <div class="filter-actions">
              <button type="submit" id="apply-filters">应用筛选</button>
              <button type="button" class="button secondary" id="clear-filters">清空</button>
            </div>
          </form>
        </div>

        <div class="scroll-fade detail-scroll-wrap">
          <div class="detail-scroll" id="detail-scroll">
            <div id="detail-stage">
              <div class="detail-empty">
                <div>
                  <h2>等待日志详情</h2>
                  <p>登录后会默认打开最新一条日志。这里主要显示用户发送和模型回复，其余原始字段会折叠到下方。</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</div>

<div class="modal-backdrop monitor-modal" id="stats-modal">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>统计</h2>
        <p>查看总体请求、今日请求、错误量、令牌分组和模型调用排行。Token 统一按 M 单位展示。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="stats-modal">关闭</button>
    </div>
    <div class="stats-grid" id="stats-summary-grid"></div>
    <div class="modal-section">
      <div class="modal-two-col">
        <div class="modal-panel">
          <h3>模型调用排行</h3>
          <div class="modal-table-wrap">
            <table>
              <thead>
                <tr>
                  <th>模型</th>
                  <th>请求数</th>
                  <th>总 Tokens</th>
                  <th>错误</th>
                  <th>最后调用</th>
                </tr>
              </thead>
              <tbody id="stats-model-body"></tbody>
            </table>
          </div>
        </div>
        <div class="modal-panel">
          <h3>令牌用量排行</h3>
          <div class="modal-table-wrap">
            <table>
              <thead>
                <tr>
                  <th>令牌</th>
                  <th>指纹</th>
                  <th>请求数</th>
                  <th>总 Tokens</th>
                  <th>最后调用</th>
                </tr>
              </thead>
              <tbody id="stats-token-body"></tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="modal-backdrop monitor-modal" id="settings-modal">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>设置</h2>
        <p>这里集成令牌管理。可以手动为 token 指纹设置代号，也可以直接在表格里改别名并保存。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="settings-modal">关闭</button>
    </div>

    <div class="message-strip hidden" id="settings-message"></div>

    <div class="modal-panel">
      <h3>手动设置令牌代号</h3>
      <form id="token-alias-form">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <div class="token-manage-grid">
          <div>
            <label for="alias-token-fingerprint">Token 指纹</label>
            <input id="alias-token-fingerprint" name="token_fingerprint" type="text" placeholder="输入完整或已复制的 token 指纹">
          </div>
          <div>
            <label for="alias-token-name">令牌代号</label>
            <input id="alias-token-name" name="token_alias" type="text" placeholder="例如 主账号 / 绘图机 / 备用 key">
          </div>
          <div>
            <button type="submit" id="save-alias-button">保存代号</button>
          </div>
        </div>
        <p class="subtle">如果代号留空并保存，会清除该 token 指纹已绑定的代号。</p>
      </form>
    </div>

    <div class="modal-section">
      <div class="modal-panel">
        <div style="display:flex;justify-content:space-between;align-items:center;gap:12px;flex-wrap:wrap;">
          <h3>最近活跃令牌</h3>
          <button class="button secondary" type="button" id="refresh-token-directory">刷新列表</button>
        </div>
        <div class="modal-table-wrap" style="max-height:520px;">
          <table>
            <thead>
              <tr>
                <th>令牌</th>
                <th>预览</th>
                <th>请求数</th>
                <th>总 Tokens</th>
                <th>最后调用</th>
                <th>代号编辑</th>
              </tr>
            </thead>
            <tbody id="token-directory-body"></tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="modal-backdrop monitor-modal" id="db-modal">
  <div class="modal-card">
    <div class="modal-header">
      <div>
        <h2>数据库管理</h2>
        <p>查看当前入库量、今日入库、占用空间和大字段情况，并按时间 / 令牌 / 模型清理日志。对数据库执行操作后，统计会立刻刷新。</p>
      </div>
      <button class="modal-close" type="button" data-close-modal="db-modal">关闭</button>
    </div>

    <div class="message-strip hidden" id="db-message"></div>

    <div class="db-stats-grid" id="db-stats-grid"></div>

    <div class="modal-section">
      <div class="modal-panel">
        <h3>空间维护</h3>
        <p class="subtle">删除日志后，PostgreSQL 通常只会把空间标记为可复用，不会立刻把文件缩回给操作系统。需要时可以继续执行整理空间或强制缩盘。</p>
        <div class="db-actions">
          <button type="button" data-db-maintenance="vacuum_analyze">整理空间</button>
          <button type="button" class="button secondary" data-db-maintenance="analyze">仅刷新统计</button>
          <button type="button" class="button secondary" data-db-maintenance="compact_payloads">压缩历史大字段</button>
          <button type="button" class="button danger" data-db-maintenance="vacuum_full">强制缩盘</button>
          <button type="button" class="button secondary" id="refresh-db-stats">刷新统计</button>
        </div>
      </div>
    </div>

    <div class="modal-section">
      <div class="modal-panel">
        <h3>按条件清理日志</h3>
        <form id="db-cleanup-form">
          <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
          <div class="cleanup-grid">
            <div>
              <label for="cleanup-from">开始时间</label>
              <input id="cleanup-from" name="from" type="datetime-local">
            </div>
            <div>
              <label for="cleanup-to">结束时间</label>
              <input id="cleanup-to" name="to" type="datetime-local">
            </div>
            <div>
              <label for="cleanup-token">Token 指纹</label>
              <input id="cleanup-token" name="token" type="text" placeholder="支持片段">
            </div>
            <div>
              <label for="cleanup-alias">令牌代号</label>
              <input id="cleanup-alias" name="alias" type="text" placeholder="支持片段">
            </div>
            <div>
              <label for="cleanup-model">模型</label>
              <input id="cleanup-model" name="model" type="text" placeholder="支持片段">
            </div>
            <div class="filter-actions">
              <button type="submit" class="button danger" id="run-db-cleanup">清理符合条件的记录</button>
              <button type="button" class="button secondary" id="copy-main-filters-to-cleanup">使用当前筛选</button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>

<script>
(function () {
  var basePath = document.body.getAttribute('data-audit-base') || '';
  var logsRoot = document.querySelector('.logs-screen');
  var refs = {
    refreshState: document.getElementById('logs-refresh-state'),
    quickTotalRequests: document.getElementById('quick-total-requests'),
    quickTodayRequests: document.getElementById('quick-today-requests'),
    quickErrorCount: document.getElementById('quick-error-count'),
    quickTotalTokens: document.getElementById('quick-total-tokens'),
    quickRequestCard: document.getElementById('quick-requests-card'),
    quickTodayCard: document.getElementById('quick-today-card'),
    quickErrorCard: document.getElementById('quick-errors-card'),
    quickTokensCard: document.getElementById('quick-tokens-card'),
    listCurrentCount: document.getElementById('list-current-count'),
    listTotalCount: document.getElementById('list-total-count'),
    listPageText: document.getElementById('list-page-text'),
    logGroups: document.getElementById('log-groups'),
    prevPage: document.getElementById('prev-page'),
    nextPage: document.getElementById('next-page'),
    filtersForm: document.getElementById('filters-form'),
    filterFrom: document.getElementById('filter-from'),
    filterTo: document.getElementById('filter-to'),
    filterAlias: document.getElementById('filter-alias'),
    filterModel: document.getElementById('filter-model'),
    filterStatus: document.getElementById('filter-status'),
    filterKeyword: document.getElementById('filter-keyword'),
    clearFilters: document.getElementById('clear-filters'),
    detailScroll: document.getElementById('detail-scroll'),
    detailStage: document.getElementById('detail-stage'),
    openStatsModal: document.getElementById('open-stats-modal'),
    openSettingsModal: document.getElementById('open-settings-modal'),
    openDBModal: document.getElementById('open-db-modal'),
    statsModal: document.getElementById('stats-modal'),
    settingsModal: document.getElementById('settings-modal'),
    dbModal: document.getElementById('db-modal'),
    statsSummaryGrid: document.getElementById('stats-summary-grid'),
    statsModelBody: document.getElementById('stats-model-body'),
    statsTokenBody: document.getElementById('stats-token-body'),
    settingsMessage: document.getElementById('settings-message'),
    aliasForm: document.getElementById('token-alias-form'),
    aliasFingerprint: document.getElementById('alias-token-fingerprint'),
    aliasName: document.getElementById('alias-token-name'),
    saveAliasButton: document.getElementById('save-alias-button'),
    tokenDirectoryBody: document.getElementById('token-directory-body'),
    refreshTokenDirectory: document.getElementById('refresh-token-directory'),
    dbStatsGrid: document.getElementById('db-stats-grid'),
    dbMessage: document.getElementById('db-message'),
    refreshDBStats: document.getElementById('refresh-db-stats'),
    dbCleanupForm: document.getElementById('db-cleanup-form'),
    cleanupFrom: document.getElementById('cleanup-from'),
    cleanupTo: document.getElementById('cleanup-to'),
    cleanupToken: document.getElementById('cleanup-token'),
    cleanupAlias: document.getElementById('cleanup-alias'),
    cleanupModel: document.getElementById('cleanup-model'),
    copyMainFiltersToCleanup: document.getElementById('copy-main-filters-to-cleanup')
  };

  var state = {
    filters: {
      from: refs.filterFrom.value || '',
      to: refs.filterTo.value || '',
      alias: refs.filterAlias.value || '',
      model: refs.filterModel.value || '',
      status: refs.filterStatus.value || '',
      q: refs.filterKeyword.value || ''
    },
    page: parsePositiveInt((logsRoot && logsRoot.getAttribute('data-log-page')) || new URLSearchParams(window.location.search).get('page')) || 1,
    version: '',
    items: [],
    totalCount: parsePositiveInt(logsRoot && logsRoot.getAttribute('data-log-total-count')),
    pageSize: parsePositiveInt(logsRoot && logsRoot.getAttribute('data-log-page-size')) || 100,
    totalPages: parsePositiveInt(logsRoot && logsRoot.getAttribute('data-log-total-pages')) || 1,
    hasPrev: logsRoot && logsRoot.getAttribute('data-log-has-prev') === 'true',
    hasNext: logsRoot && logsRoot.getAttribute('data-log-has-next') === 'true',
    selectedId: 0,
    detail: null,
    raw: null,
    rawLoadedFor: 0,
    rawLoadingFor: 0,
    rawError: '',
    textPages: {
      user: null,
      assistant: null
    },
    dashboard: null,
    tokenDirectory: [],
    dbStats: null,
    filterOptions: {
      models: [],
      token_aliases: [],
      token_fingerprints: [],
      status_codes: []
    },
    groupCollapsed: {},
    groupMeta: {},
    groupCache: {},
    controllers: {},
    activeSuggestInput: '',
    logListHeight: 0,
    lastLogUpdateAt: 0,
    pollingTimer: 0
  };

  function parsePositiveInt(value) {
    var n = parseInt(value || '', 10);
    return isFinite(n) && n > 0 ? n : 0;
  }

  function apiPath(path) {
    return (basePath || '') + path;
  }

  function escapeHTML(value) {
    return String(value == null ? '' : value)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  function formatNumber(value) {
    var n = Number(value || 0);
    return n.toLocaleString('zh-CN');
  }

  function formatTokensM(value) {
    var n = Number(value || 0);
    return (n / 1000000).toFixed(n >= 10000000 ? 1 : 2) + ' M';
  }

  function formatBytes(value) {
    var n = Number(value || 0);
    if (n <= 0) {
      return '0 B';
    }
    var units = ['B', 'KB', 'MB', 'GB', 'TB'];
    var idx = Math.min(Math.floor(Math.log(n) / Math.log(1024)), units.length - 1);
    return (n / Math.pow(1024, idx)).toFixed(idx === 0 ? 0 : 2) + ' ' + units[idx];
  }

  function truncateText(value, maxLen) {
    var text = String(value || '').trim();
    if (text.length <= maxLen) {
      return text;
    }
    return text.slice(0, Math.max(0, maxLen - 1)) + '…';
  }

  function pickAlias(item) {
    if (item && item.token_alias) {
      return item.token_alias;
    }
    if (item && item.token_preview) {
      return item.token_preview;
    }
    if (item && item.token_fingerprint) {
      return truncateText(item.token_fingerprint, 18);
    }
    return '未命名令牌';
  }

  function renderFingerprint(item) {
    if (!item || !item.token_fingerprint) {
      return '未记录指纹';
    }
    return item.token_fingerprint;
  }

  function isAbortError(err) {
    return err && (err.name === 'AbortError' || String(err.message || '').indexOf('aborted') >= 0);
  }

  function getController(name) {
    if (state.controllers[name]) {
      state.controllers[name].abort();
    }
    var controller = new AbortController();
    state.controllers[name] = controller;
    return controller;
  }

  async function requestJSON(url, options) {
    options = options || {};
    var requestOptions = Object.assign({}, options);
    var method = String(requestOptions.method || 'GET').toUpperCase();
    var headers = new Headers(requestOptions.headers || {});
    headers.set('Accept', 'application/json');
    if (method !== 'GET' && method !== 'HEAD') {
      var csrfMeta = document.querySelector('meta[name="csrf-token"]');
      if (csrfMeta && csrfMeta.content) {
        headers.set('X-CSRF-Token', csrfMeta.content);
      }
    }
    requestOptions.headers = headers;
    requestOptions.credentials = requestOptions.credentials || 'same-origin';
    var resp = await fetch(url, requestOptions);
    if (resp.status === 401) {
      window.location.href = apiPath('/login');
      throw new Error('登录状态已失效，请重新登录。');
    }
    var text = await resp.text();
    var payload = null;
    if (text) {
      try {
        payload = JSON.parse(text);
      } catch (err) {
        if (!resp.ok) {
          throw new Error(text || '请求失败');
        }
        throw new Error('返回内容不是合法 JSON。');
      }
    }
    if (!resp.ok) {
      if (payload && payload.message) {
        throw new Error(payload.message);
      }
      throw new Error(text || '请求失败');
    }
    return payload || {};
  }

  function setRefreshState(text, mode) {
    refs.refreshState.textContent = text;
    refs.refreshState.classList.remove('is-checking', 'has-updates', 'is-error');
    if (mode) {
      refs.refreshState.classList.add(mode);
    }
  }

  function formatElapsedSince(timestamp) {
    if (!timestamp) {
      return '等待更新';
    }
    var seconds = Math.max(0, Math.floor((Date.now() - timestamp) / 1000));
    if (seconds < 5) {
      return '刚刚';
    }
    if (seconds < 60) {
      return seconds + ' 秒前';
    }
    var minutes = Math.floor(seconds / 60);
    if (minutes < 60) {
      return minutes + ' 分钟前';
    }
    var hours = Math.floor(minutes / 60);
    if (hours < 24) {
      return hours + ' 小时前';
    }
    return Math.floor(hours / 24) + ' 天前';
  }

  function renderLastUpdateState(mode) {
    setRefreshState('上次更新 · ' + formatElapsedSince(state.lastLogUpdateAt), mode || '');
  }

  function pulseStatCard(card) {
    if (!card) {
      return;
    }
    card.classList.remove('pulse');
    window.requestAnimationFrame(function () {
      card.classList.add('pulse');
      window.setTimeout(function () {
        card.classList.remove('pulse');
      }, 800);
    });
  }

  function getLogPageFromURL() {
    return parsePositiveInt(new URLSearchParams(window.location.search).get('page')) || 1;
  }

  function buildLogQuery(includePage, pageOverride) {
    var targetPage = parsePositiveInt(pageOverride) || state.page || 1;
    var params = new URLSearchParams();
    if (state.filters.from) {
      params.set('from', state.filters.from);
    }
    if (state.filters.to) {
      params.set('to', state.filters.to);
    }
    if (state.filters.alias) {
      params.set('alias', state.filters.alias);
    }
    if (state.filters.model) {
      params.set('model', state.filters.model);
    }
    if (state.filters.status) {
      params.set('status', state.filters.status);
    }
    if (state.filters.q) {
      params.set('q', state.filters.q);
    }
    if (includePage && targetPage > 1) {
      params.set('page', String(targetPage));
    }
    return params;
  }

  function buildLogPageURL(page) {
    var params = buildLogQuery(true, page);
    return window.location.pathname + (params.toString() ? '?' + params.toString() : '');
  }

  function syncPageStateFromURL() {
    state.page = getLogPageFromURL();
  }

  function syncFilterInputsFromState() {
    refs.filterFrom.value = state.filters.from;
    refs.filterTo.value = state.filters.to;
    refs.filterAlias.value = state.filters.alias;
    refs.filterModel.value = state.filters.model;
    refs.filterStatus.value = state.filters.status;
    refs.filterKeyword.value = state.filters.q;
    updateSuggestClearButtons();
  }

  function syncStateFromFilterInputs() {
    state.filters.from = refs.filterFrom.value.trim();
    state.filters.to = refs.filterTo.value.trim();
    state.filters.alias = refs.filterAlias.value.trim();
    state.filters.model = refs.filterModel.value.trim();
    state.filters.status = refs.filterStatus.value.trim();
    state.filters.q = refs.filterKeyword.value.trim();
  }

  function hasActiveFilters() {
    return !!(state.filters.from || state.filters.to || state.filters.alias || state.filters.model || state.filters.status || state.filters.q);
  }

  function updateSuggestClearButtons() {
    document.querySelectorAll('[data-clear-filter]').forEach(function (button) {
      var input = document.getElementById(button.getAttribute('data-clear-filter'));
      button.disabled = !input || !input.value.trim();
    });
  }

  function updateFilterSummary() {
    return;
  }

  function buildLogGroups(items) {
    var groups = [];
    var current = null;
    if (hasActiveFilters()) {
      groups = items.map(function (item) {
        return {
          identity: 'filtered-' + item.id,
          items: [item]
        };
      });
    } else {
      items.forEach(function (item) {
        var identity = [item.token_alias || '', item.token_fingerprint || '', item.token_preview || ''].join('|');
        if (!current || current.identity !== identity) {
          current = {
            identity: identity,
            items: [item]
          };
          groups.push(current);
        } else {
          current.items.push(item);
        }
      });
    }
    groups.forEach(function (group) {
      group.latest = group.items[0];
      group.oldest = group.items[group.items.length - 1];
      group.key = 'group-' + group.latest.id + '-' + group.oldest.id;
      group.totalTokens = group.items.reduce(function (sum, item) {
        return sum + Number(item.total_tokens || 0);
      }, 0);
      group.hasSelected = group.items.some(function (item) {
        return item.id === state.selectedId;
      });
    });
    return groups;
  }

  function groupMetaFor(group) {
    var ids = {};
    group.items.forEach(function (item) {
      ids[String(item.id)] = true;
    });
    return {
      identity: group.identity,
      ids: ids
    };
  }

  function groupOverlapsMeta(group, meta) {
    if (!meta || meta.identity !== group.identity || !meta.ids) {
      return false;
    }
    return group.items.some(function (item) {
      return !!meta.ids[String(item.id)];
    });
  }

  function inheritGroupCollapseStates(groups) {
    var previousMeta = state.groupMeta || {};
    var previousCollapsed = state.groupCollapsed || {};
    var nextMeta = {};
    var nextCollapsed = {};
    var previousKeys = Object.keys(previousMeta);

    groups.forEach(function (group) {
      if (Object.prototype.hasOwnProperty.call(previousCollapsed, group.key)) {
        nextCollapsed[group.key] = previousCollapsed[group.key];
      } else {
        previousKeys.some(function (previousKey) {
          if (!Object.prototype.hasOwnProperty.call(previousCollapsed, previousKey)) {
            return false;
          }
          if (!groupOverlapsMeta(group, previousMeta[previousKey])) {
            return false;
          }
          nextCollapsed[group.key] = previousCollapsed[previousKey];
          return true;
        });
      }
      nextMeta[group.key] = groupMetaFor(group);
    });

    state.groupCollapsed = nextCollapsed;
    state.groupMeta = nextMeta;
  }

  function isGroupCollapsed(group) {
    if (group.items.length <= 1) {
      return false;
    }
    if (hasActiveFilters()) {
      return false;
    }
    if (Object.prototype.hasOwnProperty.call(state.groupCollapsed, group.key)) {
      return state.groupCollapsed[group.key];
    }
    return true;
  }

  function buildGroupSignature(group) {
    var parts = group.items.map(function (item) {
      return [
        item.id,
        item.started_at,
        item.token_alias || '',
        item.model || '',
        item.status_code,
        item.total_tokens,
        item.user_preview || '',
        item.assistant_preview || ''
      ].join('|');
    });
    return [
      group.key,
      group.hasSelected ? '1' : '0',
      isGroupCollapsed(group) ? '1' : '0',
      parts.join('~')
    ].join('::');
  }

  function toggleGroupCollapse(key) {
    var current = Object.prototype.hasOwnProperty.call(state.groupCollapsed, key) ? state.groupCollapsed[key] : true;
    state.groupCollapsed[key] = !current;
  }

  function renderPreviewText(item) {
    if (item.user_preview) {
      return escapeHTML(item.user_preview);
    }
    if (item.assistant_preview) {
      return escapeHTML(item.assistant_preview);
    }
    return '暂无文本预览';
  }

  function renderGroupNode(group, isFresh) {
    var root = document.createElement('div');
    var latest = group.latest;
    var collapsed = isGroupCollapsed(group);
    var freshClass = isFresh ? ' is-fresh' : '';
    var toggleHTML = '';
    var childHTML = '';

    if (group.items.length > 1) {
      toggleHTML = '<button type="button" class="log-row-toggle" data-toggle-group="' + escapeHTML(group.key) + '" aria-expanded="' + (!collapsed ? 'true' : 'false') + '">' +
        (collapsed ? '展开连续 ' + group.items.length + ' 条请求' : '收起连续 ' + group.items.length + ' 条请求') +
        '</button>';
      if (!collapsed) {
        childHTML = '<div class="log-row-children">' + group.items.map(function (item, index) {
          var activeClass = item.id === state.selectedId ? ' active' : '';
          return '' +
            '<button type="button" class="log-row-child' + activeClass + '" data-select-log="' + item.id + '" style="--child-index:' + index + ';">' +
              '<div class="log-row-line">' +
                '<strong>' + escapeHTML(item.model || '未标记模型') + '</strong>' +
                '<span class="log-row-time">' + escapeHTML(item.started_at || '-') + '</span>' +
              '</div>' +
              '<div class="log-row-meta">' +
                '<span class="log-row-badge ' + (Number(item.status_code || 0) >= 400 ? 'status-err' : 'status-ok') + '">HTTP ' + escapeHTML(item.status_code || 0) + '</span>' +
                '<span class="log-row-badge">' + escapeHTML(formatTokensM(item.total_tokens)) + '</span>' +
              '</div>' +
              '<div class="log-row-preview small">' + renderPreviewText(item) + '</div>' +
            '</button>';
        }).join('') + '</div>';
      }
    }

    root.className = 'log-row' + freshClass;
    root.setAttribute('data-group-key', group.key);
    root.setAttribute('data-expanded', collapsed ? 'false' : 'true');
    root.setAttribute('role', 'listitem');
    root.innerHTML = '' +
      '<button type="button" class="log-row-main' + (latest.id === state.selectedId ? ' active' : '') + '" data-select-log="' + latest.id + '">' +
        '<div class="log-row-line">' +
          '<div class="log-row-title">' +
            '<strong>' + escapeHTML(pickAlias(latest)) + '</strong>' +
            '<code>' + escapeHTML(renderFingerprint(latest)) + '</code>' +
          '</div>' +
          '<div class="log-row-time">' + escapeHTML(latest.started_at || '-') + '</div>' +
        '</div>' +
        '<div class="log-row-preview">' + renderPreviewText(latest) + '</div>' +
        '<div class="log-row-meta">' +
          '<span class="log-row-badge">' + escapeHTML(latest.model || '未标记模型') + '</span>' +
          '<span class="log-row-badge ' + (Number(latest.status_code || 0) >= 400 ? 'status-err' : 'status-ok') + '">HTTP ' + escapeHTML(latest.status_code || 0) + '</span>' +
          '<span class="log-row-badge">' + escapeHTML(formatTokensM(group.totalTokens)) + '</span>' +
          '<span class="log-row-badge">' + (group.items.length > 1 ? group.items.length + ' 条连续' : '单条') + '</span>' +
        '</div>' +
      '</button>' +
      toggleHTML +
      childHTML;
    return root;
  }

  function markFreshNode(node) {
    if (!node) {
      return;
    }
    node.classList.remove('is-fresh');
    void node.offsetWidth;
    node.classList.add('is-fresh');
    window.setTimeout(function () {
      node.classList.remove('is-fresh');
    }, 1300);
  }

  function getLogScrollAnchor() {
    if (!refs.logGroups) {
      return null;
    }
    var listRect = refs.logGroups.getBoundingClientRect();
    var rows = Array.prototype.slice.call(refs.logGroups.querySelectorAll('.log-row'));
    for (var i = 0; i < rows.length; i += 1) {
      var rect = rows[i].getBoundingClientRect();
      if (rect.bottom >= listRect.top + 8) {
        return {
          key: rows[i].getAttribute('data-group-key'),
          offset: rect.top - listRect.top
        };
      }
    }
    return null;
  }

  function escapeCSSIdent(value) {
    if (window.CSS && typeof window.CSS.escape === 'function') {
      return window.CSS.escape(value);
    }
    return String(value || '').replace(/["\\]/g, '\\$&');
  }

  function restoreLogScrollAnchor(anchor) {
    if (!anchor || !anchor.key || !refs.logGroups) {
      return false;
    }
    var node = refs.logGroups.querySelector('[data-group-key="' + escapeCSSIdent(anchor.key) + '"]');
    if (!node) {
      return false;
    }
    var listRect = refs.logGroups.getBoundingClientRect();
    var rect = node.getBoundingClientRect();
    refs.logGroups.scrollTop += rect.top - listRect.top - anchor.offset;
    return true;
  }

  function renderLogList(freshIDs, preserveScroll) {
    var previousScrollTop = preserveScroll && refs.logGroups ? refs.logGroups.scrollTop : 0;
    var previousScrollHeight = preserveScroll && refs.logGroups ? refs.logGroups.scrollHeight : 0;
    var keepViewportAnchored = preserveScroll && previousScrollTop > 24 && freshIDs && freshIDs.length > 0;
    var scrollAnchor = preserveScroll && (!freshIDs || !freshIDs.length || previousScrollTop > 24) ? getLogScrollAnchor() : null;
    updatePaginationControls();

    if (!state.items.length) {
      refs.logGroups.innerHTML = '' +
        '<div class="lane-empty">' +
          '<div>' +
            '<h3>没有匹配日志</h3>' +
            '<p>可以调整时间、模型、令牌代号或关键词，再自动刷新最新结果。</p>' +
          '</div>' +
        '</div>';
      state.groupCache = {};
      state.groupMeta = {};
      resizeLogListViewport();
      return;
    }

    var groups = buildLogGroups(state.items);
    inheritGroupCollapseStates(groups);
    var nextCache = {};
    var cursor = refs.logGroups.firstElementChild;
    groups.forEach(function (group) {
      var signature = buildGroupSignature(group);
      var cached = state.groupCache[group.key];
      var node;
      var shouldMarkFresh = freshIDs.some(function (id) {
        return group.items.some(function (item) {
          return item.id === id;
        });
      });
      if (cached && cached.signature === signature) {
        node = cached.node;
      } else {
        node = renderGroupNode(group, shouldMarkFresh);
      }
      if (shouldMarkFresh) {
        markFreshNode(node);
      }
      nextCache[group.key] = {
        node: node,
        signature: signature
      };
      if (node === cursor) {
        cursor = cursor.nextElementSibling;
      } else {
        refs.logGroups.insertBefore(node, cursor);
      }
    });
    Array.prototype.slice.call(refs.logGroups.children).forEach(function (child) {
      var key = child.getAttribute('data-group-key');
      if (!key || !nextCache[key]) {
        child.remove();
      }
    });
    state.groupCache = nextCache;
    resizeLogListViewport();
    if (preserveScroll) {
      if (restoreLogScrollAnchor(scrollAnchor)) {
        return syncLogSelectionState();
      }
      if (keepViewportAnchored) {
        previousScrollTop += Math.max(0, refs.logGroups.scrollHeight - previousScrollHeight);
      }
      var maxScrollTop = Math.max(0, refs.logGroups.scrollHeight - refs.logGroups.clientHeight);
      refs.logGroups.scrollTop = Math.min(previousScrollTop, maxScrollTop);
    } else {
      refs.logGroups.scrollTop = 0;
    }
    syncLogSelectionState();
  }

  function syncLogSelectionState() {
    if (!refs.logGroups) {
      return;
    }
    refs.logGroups.querySelectorAll('[data-select-log]').forEach(function (button) {
      var isActive = parsePositiveInt(button.getAttribute('data-select-log')) === state.selectedId;
      button.classList.toggle('active', isActive);
    });
    refs.logGroups.querySelectorAll('.log-row.is-selected').forEach(function (cluster) {
      cluster.classList.remove('is-selected');
    });
  }

  function renderDetailSkeleton() {
    refs.detailScroll.scrollTop = 0;
    refs.detailScroll.classList.remove('is-swapping');
    window.requestAnimationFrame(function () {
      refs.detailScroll.classList.add('is-swapping');
    });
    refs.detailStage.innerHTML = '' +
      '<div class="detail-shell">' +
        '<div class="loading-glass" style="height:88px;border-radius:24px;background:rgba(255,255,255,0.58);border:1px solid rgba(236,192,208,0.76);"></div>' +
        '<div class="loading-glass" style="height:104px;border-radius:24px;background:rgba(255,255,255,0.58);border:1px solid rgba(236,192,208,0.76);"></div>' +
        '<div class="loading-glass" style="height:290px;border-radius:28px;background:rgba(255,255,255,0.58);border:1px solid rgba(236,192,208,0.76);"></div>' +
        '<div class="loading-glass" style="height:290px;border-radius:28px;background:rgba(255,255,255,0.58);border:1px solid rgba(236,192,208,0.76);"></div>' +
      '</div>';
  }

  function renderTextPanel(kind, title, pageData, cssClass, extraHTML) {
    if (!pageData) {
      return '' +
        '<div class="speech-card ' + cssClass + '">' +
          '<h3>' + title + '</h3>' +
          '<pre>暂无内容</pre>' +
        '</div>';
    }
    var totalPages = Number(pageData.total_pages || 1);
    var page = Number(pageData.page || 1);
    var pagination = '';
    if (totalPages > 1) {
      pagination = '' +
        '<div class="mini-pagination">' +
          '<button type="button" class="button secondary" data-text-kind="' + kind + '" data-text-page="' + (page - 1) + '"' + (page <= 1 ? ' disabled' : '') + '>上一页</button>' +
          '<span>第 ' + page + ' / ' + totalPages + ' 页</span>' +
          '<button type="button" class="button secondary" data-text-kind="' + kind + '" data-text-page="' + (page + 1) + '"' + (page >= totalPages ? ' disabled' : '') + '>下一页</button>' +
        '</div>';
    }
    return '' +
      '<div class="speech-card ' + cssClass + '" id="' + kind + '-text-panel">' +
        '<div class="text-toolbar">' +
          '<h3>' + title + '</h3>' +
          '<div class="muted">共 ' + formatNumber(pageData.total_chars || 0) + ' 字</div>' +
        '</div>' +
        (extraHTML || '') +
        '<pre>' + escapeHTML(pageData.text || '暂无内容') + '</pre>' +
        pagination +
      '</div>';
  }

  function renderImageGallery(raw) {
    if (!raw || !raw.images || !raw.images.length) {
      return '';
    }
    return '' +
      '<div class="gallery-grid">' +
        raw.images.map(function (image) {
          return '' +
            '<div class="gallery-card">' +
              '<img src="' + escapeHTML(image.data_url || '') + '" alt="' + escapeHTML(image.label || '生成图片') + '">' +
              '<p>' + escapeHTML(image.label || '模型图片') + ' · ' + escapeHTML(image.mime || '') + '</p>' +
            '</div>';
        }).join('') +
      '</div>';
  }

  function renderRawSections(detail, raw) {
    if (!detail) {
      return '';
    }
    if (state.rawLoadingFor === detail.id && !raw) {
      return '' +
        '<details class="data-fold" open>' +
          '<summary>原始与结构化数据</summary>' +
          '<div class="raw-placeholder loading-glass"><p>正在加载原始 request / response、JSON、usage 和 headers…</p></div>' +
        '</details>';
    }
    if (state.rawError && state.rawLoadedFor !== detail.id) {
      return '' +
        '<details class="data-fold" open>' +
          '<summary>原始与结构化数据</summary>' +
          '<div class="raw-placeholder">' +
            '<p>' + escapeHTML(state.rawError) + '</p>' +
            '<button type="button" class="button secondary" data-load-raw="' + detail.id + '">重新加载</button>' +
          '</div>' +
        '</details>';
    }
    if (!raw) {
      return '' +
        '<details class="data-fold" data-needs-raw="1" data-log-id="' + detail.id + '">' +
          '<summary>原始与结构化数据</summary>' +
          '<div class="raw-placeholder">' +
            '<p>默认先显示用户内容和模型回复，点开这里后再按需加载原始 body、JSON、headers 和 usage，切换日志时响应会更快。</p>' +
            '<button type="button" class="button secondary" data-load-raw="' + detail.id + '">立即加载</button>' +
          '</div>' +
        '</details>';
    }
    return '' +
      '<details class="data-fold">' +
        '<summary>原始与结构化数据</summary>' +
        '<div class="pre-grid">' +
          '<div class="pre-block"><h4>请求 JSON</h4><pre>' + escapeHTML(raw.request_json || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>响应 JSON</h4><pre>' + escapeHTML(raw.response_json || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>Usage</h4><pre>' + escapeHTML(raw.usage_json || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>请求 Headers</h4><pre>' + escapeHTML(raw.request_headers || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>响应 Headers</h4><pre>' + escapeHTML(raw.response_headers || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>原始请求 Body</h4><pre>' + escapeHTML(raw.request_body || '（空）') + '</pre></div>' +
          '<div class="pre-block"><h4>原始响应 Body</h4><pre>' + escapeHTML(raw.response_body || '（空）') + '</pre></div>' +
        '</div>' +
      '</details>';
  }

  function renderDetail() {
    if (!state.detail) {
      refs.detailStage.innerHTML = '' +
        '<div class="detail-empty">' +
          '<div>' +
            '<h2>暂无详情</h2>' +
            '<p>左侧选中任意一条日志后，这里会立即显示用户发送与模型回复。</p>' +
          '</div>' +
        '</div>';
      return;
    }
    var detail = state.detail;
    var raw = state.rawLoadedFor === detail.id ? state.raw : null;
    var errorBanner = '';
    if (detail.error_text) {
      errorBanner = '' +
        '<div class="detail-error-banner">' +
          '<strong>请求或响应过程中有错误</strong>' +
          '<div>' + escapeHTML(detail.error_text) + '</div>' +
        '</div>';
    }
    var compactStats = 'Tokens ' + formatTokensM(detail.total_tokens) + ' / 入 ' + formatTokensM(detail.prompt_tokens) + ' / 出 ' + formatTokensM(detail.completion_tokens);
    var compactSizes = formatBytes(detail.request_bytes) + ' → ' + formatBytes(detail.response_bytes);
    refs.detailStage.innerHTML = '' +
      '<div class="detail-shell">' +
        '<div class="detail-compact">' +
          '<div class="detail-compact-line" title="' + escapeHTML((detail.path_with_query || '') + ' · ' + (detail.token_fingerprint || '')) + '">' +
            '<strong class="detail-compact-path">' + escapeHTML(detail.path_with_query || '') + '</strong>' +
            '<span class="detail-compact-chip">' + escapeHTML(detail.started_at || '-') + '</span>' +
            '<span class="detail-compact-chip">' + escapeHTML(detail.method || 'POST') + '</span>' +
            '<span class="detail-compact-chip">令牌 ' + escapeHTML(detail.token_alias || detail.token_preview || '未命名') + '</span>' +
            '<span class="detail-compact-chip">模型 ' + escapeHTML(detail.model || '未标记') + '</span>' +
            '<span class="detail-compact-chip ' + (Number(detail.status_code || 0) >= 400 ? 'status-err' : '') + '">HTTP ' + escapeHTML(detail.status_code || 0) + '</span>' +
            '<span class="detail-compact-chip">' + escapeHTML(detail.response_type || '未知') + '</span>' +
            '<span class="detail-compact-chip">' + escapeHTML(compactStats) + '</span>' +
            '<span class="detail-compact-chip">' + escapeHTML(detail.duration_ms || 0) + ' ms</span>' +
            '<span class="detail-compact-chip">' + escapeHTML(compactSizes) + '</span>' +
          '</div>' +
          '<a class="button secondary" href="' + escapeHTML(apiPath('/logs/' + detail.id)) + '" target="_blank" rel="noreferrer">详情</a>' +
        '</div>' +
        errorBanner +
        '<div class="conversation-stack">' +
          '<div class="conversation-column">' +
            renderTextPanel('user', '用户发送', state.textPages.user || detail.user_text, 'user', '') +
          '</div>' +
          '<div class="conversation-column">' +
            renderTextPanel('assistant', '模型回复', state.textPages.assistant || detail.assistant_text, 'assistant', renderImageGallery(raw)) +
          '</div>' +
        '</div>' +
        '<div class="secondary-stack">' +
          renderRawSections(detail, raw) +
        '</div>' +
      '</div>';
  }

  function shouldPrefetchRaw(detail) {
    if (!detail) {
      return false;
    }
    return /image/i.test(detail.model || '');
  }

  async function ensureRawLoaded(id, silent) {
    if (!id) {
      return;
    }
    if (state.rawLoadedFor === id || state.rawLoadingFor === id) {
      return;
    }
    state.rawError = '';
    state.rawLoadingFor = id;
    if (!silent) {
      renderDetail();
    }
    var controller = getController('raw');
    try {
      var raw = await requestJSON(apiPath('/api/logs/' + id + '/raw'), { signal: controller.signal });
      if (state.selectedId !== id) {
        return;
      }
      state.raw = raw;
      state.rawLoadedFor = id;
      state.rawLoadingFor = 0;
      state.rawError = '';
      renderDetail();
    } catch (err) {
      if (isAbortError(err)) {
        return;
      }
      state.rawLoadingFor = 0;
      state.rawError = err.message || '加载原始数据失败。';
      renderDetail();
    }
  }

  async function loadDetail(id) {
    if (!id) {
      return;
    }
    state.selectedId = id;
    state.detail = null;
    state.raw = null;
    state.rawLoadedFor = 0;
    state.rawLoadingFor = 0;
    state.rawError = '';
    state.textPages.user = null;
    state.textPages.assistant = null;
    syncLogSelectionState();
    renderDetailSkeleton();
    var controller = getController('detail');
    if (state.controllers.raw) {
      state.controllers.raw.abort();
    }
    try {
      var detail = await requestJSON(apiPath('/api/logs/' + id), { signal: controller.signal });
      if (state.selectedId !== id) {
        return;
      }
      state.detail = detail;
      state.textPages.user = detail.user_text;
      state.textPages.assistant = detail.assistant_text;
      renderDetail();
      if (shouldPrefetchRaw(detail)) {
        window.setTimeout(function () {
          ensureRawLoaded(id, true);
        }, 60);
      }
    } catch (err) {
      if (isAbortError(err)) {
        return;
      }
      refs.detailStage.innerHTML = '' +
        '<div class="detail-empty">' +
          '<div>' +
            '<h2>详情加载失败</h2>' +
            '<p>' + escapeHTML(err.message || '未知错误') + '</p>' +
          '</div>' +
        '</div>';
    }
  }

  async function loadTextPage(kind, page) {
    if (!state.selectedId || !kind || !page) {
      return;
    }
    var controller = getController('text-' + kind);
    try {
      var payload = await requestJSON(apiPath('/api/logs/' + state.selectedId + '/text?kind=' + encodeURIComponent(kind) + '&page=' + encodeURIComponent(page)), {
        signal: controller.signal
      });
      if (state.selectedId !== Number(payload.id || state.selectedId)) {
        return;
      }
      state.textPages[kind] = payload;
      renderDetail();
    } catch (err) {
      if (isAbortError(err)) {
        return;
      }
      setRefreshState('文本分页加载失败', 'is-error');
    }
  }

  async function loadTextPageV2(kind, page) {
    if (!state.selectedId || !kind || !page) {
      return;
    }
    var currentPageData = state.textPages[kind] || (state.detail && state.detail[kind + '_text']) || null;
    var targetPage = parsePositiveInt(page) || 1;
    if (currentPageData) {
      var maxPage = Math.max(1, parsePositiveInt(currentPageData.total_pages) || 1);
      targetPage = Math.min(Math.max(1, targetPage), maxPage);
    }
    var controller = getController('text-v2-' + kind);
    try {
      var payload = await requestJSON(apiPath('/api/logs/' + state.selectedId + '/text?kind=' + encodeURIComponent(kind) + '&page=' + encodeURIComponent(targetPage)), {
        signal: controller.signal
      });
      if (state.selectedId !== Number(payload.id || state.selectedId)) {
        return;
      }
      state.textPages[kind] = payload;
      renderDetail();
    } catch (err) {
      if (isAbortError(err)) {
        return;
      }
      var message = String(err && err.message ? err.message : '');
      if (message.indexOf('404') >= 0 || message.indexOf('不存在') >= 0 || message.indexOf('分页已失效') >= 0) {
        setRefreshState('文本分页已失效，正在刷新当前日志', 'is-error');
        await loadDetail(state.selectedId);
        return;
      }
      setRefreshState('文本分页加载失败', 'is-error');
    }
  }

  async function loadFilterOptions() {
    try {
      state.filterOptions = await requestJSON(apiPath('/api/filter-options'));
    } catch (err) {
      if (!isAbortError(err)) {
        setRefreshState('筛选候选加载失败', 'is-error');
      }
    }
  }

  function renderQuickStats() {
    if (!state.dashboard) {
      refs.quickTotalRequests.textContent = '-';
      refs.quickTodayRequests.textContent = '-';
      refs.quickErrorCount.textContent = '-';
      refs.quickTotalTokens.textContent = '-';
      return;
    }
    refs.quickTotalRequests.textContent = formatNumber(state.dashboard.total_requests);
    refs.quickTodayRequests.textContent = formatNumber(state.dashboard.today_requests);
    refs.quickErrorCount.textContent = formatNumber(state.dashboard.error_count);
    refs.quickTotalTokens.textContent = formatTokensM(state.dashboard.total_tokens);
  }

  function renderStatsModal() {
    if (!state.dashboard) {
      refs.statsSummaryGrid.innerHTML = '';
      refs.statsModelBody.innerHTML = '';
      refs.statsTokenBody.innerHTML = '';
      return;
    }
    var cards = [
      ['总请求', formatNumber(state.dashboard.total_requests)],
      ['今日请求', formatNumber(state.dashboard.today_requests)],
      ['错误数', formatNumber(state.dashboard.error_count)],
      ['活跃令牌', formatNumber(state.dashboard.distinct_tokens)],
      ['今日活跃令牌', formatNumber(state.dashboard.today_distinct_tokens)],
      ['总 Tokens', formatTokensM(state.dashboard.total_tokens)],
      ['今日 Tokens', formatTokensM(state.dashboard.today_total_tokens)],
      ['输入 Tokens', formatTokensM(state.dashboard.total_prompt_tokens)],
      ['输出 Tokens', formatTokensM(state.dashboard.total_completion_tokens)]
    ];
    refs.statsSummaryGrid.innerHTML = cards.map(function (card) {
      return '<div class="stats-card"><span>' + escapeHTML(card[0]) + '</span><strong>' + escapeHTML(card[1]) + '</strong></div>';
    }).join('');

    var modelRows = state.dashboard.model_groups || [];
    refs.statsModelBody.innerHTML = modelRows.length ? modelRows.map(function (item) {
      return '' +
        '<tr>' +
          '<td>' + escapeHTML(item.model || '未标记模型') + '</td>' +
          '<td>' + escapeHTML(formatNumber(item.request_count)) + '</td>' +
          '<td>' + escapeHTML(formatTokensM(item.total_tokens)) + '</td>' +
          '<td>' + escapeHTML(formatNumber(item.error_count)) + '</td>' +
          '<td>' + escapeHTML(item.last_seen || '-') + '</td>' +
        '</tr>';
    }).join('') : '<tr><td colspan="5" class="muted">暂无模型统计</td></tr>';

    var tokenRows = state.dashboard.token_groups || [];
    refs.statsTokenBody.innerHTML = tokenRows.length ? tokenRows.map(function (item) {
      return '' +
        '<tr>' +
          '<td>' + escapeHTML(item.token_alias || item.token_preview || '未命名') + '</td>' +
          '<td><code>' + escapeHTML(item.token_fingerprint || '未记录') + '</code></td>' +
          '<td>' + escapeHTML(formatNumber(item.request_count)) + '</td>' +
          '<td>' + escapeHTML(formatTokensM(item.total_tokens)) + '</td>' +
          '<td>' + escapeHTML(item.last_seen || '-') + '</td>' +
        '</tr>';
    }).join('') : '<tr><td colspan="5" class="muted">暂无令牌统计</td></tr>';
  }

  async function loadDashboard(silent) {
    try {
      state.dashboard = await requestJSON(apiPath('/api/dashboard'));
      renderQuickStats();
      renderStatsModal();
      if (!silent) {
        pulseStatCard(refs.quickRequestCard);
        pulseStatCard(refs.quickTodayCard);
        pulseStatCard(refs.quickErrorCard);
        pulseStatCard(refs.quickTokensCard);
      }
    } catch (err) {
      if (!isAbortError(err)) {
        setRefreshState('统计加载失败', 'is-error');
      }
    }
  }

  function renderTokenDirectory() {
    var rows = state.tokenDirectory || [];
    refs.tokenDirectoryBody.innerHTML = rows.length ? rows.map(function (item) {
      return '' +
        '<tr data-token-fingerprint="' + escapeHTML(item.token_fingerprint || '') + '">' +
          '<td>' +
            '<strong>' + escapeHTML(item.token_alias || '未命名') + '</strong>' +
            '<div class="cell-note mono-box">' + escapeHTML(item.token_fingerprint || '未记录') + '</div>' +
          '</td>' +
          '<td>' + escapeHTML(item.token_preview || '-') + '</td>' +
          '<td>' + escapeHTML(formatNumber(item.request_count)) + '</td>' +
          '<td>' + escapeHTML(formatTokensM(item.total_tokens)) + '</td>' +
          '<td>' + escapeHTML(item.last_seen || '-') + '</td>' +
          '<td>' +
            '<div class="row-actions">' +
              '<input class="token-row-input" type="text" value="' + escapeHTML(item.token_alias || '') + '" placeholder="设置或清空代号">' +
              '<button type="button" class="button secondary" data-save-row-alias="' + escapeHTML(item.token_fingerprint || '') + '">保存</button>' +
            '</div>' +
          '</td>' +
        '</tr>';
    }).join('') : '<tr><td colspan="6" class="muted">暂无令牌数据</td></tr>';
  }

  async function loadTokenDirectory() {
    refs.refreshTokenDirectory.disabled = true;
    try {
      var payload = await requestJSON(apiPath('/api/tokens'));
      state.tokenDirectory = payload.items || [];
      renderTokenDirectory();
    } catch (err) {
      if (!isAbortError(err)) {
        showMessage(refs.settingsMessage, err.message || '加载令牌目录失败。', 'error');
      }
    } finally {
      refs.refreshTokenDirectory.disabled = false;
    }
  }

  function renderDBStats() {
    if (!state.dbStats) {
      refs.dbStatsGrid.innerHTML = '';
      return;
    }
    var cards = [
      ['审计日志条数', formatNumber(state.dbStats.total_rows)],
      ['今日入库', formatNumber(state.dbStats.today_rows)],
      ['audit_logs 总占用', state.dbStats.audit_total_pretty || formatBytes(state.dbStats.audit_total_size)],
      ['数据主体', state.dbStats.audit_table_pretty || formatBytes(state.dbStats.audit_table_size)],
      ['索引占用', state.dbStats.audit_index_pretty || formatBytes(state.dbStats.audit_index_size)],
      ['TOAST 大字段', state.dbStats.audit_toast_pretty || formatBytes(state.dbStats.audit_toast_size)],
      ['死元组', formatNumber(state.dbStats.dead_tuples)],
      ['存活元组', formatNumber(state.dbStats.live_tuples)],
      ['整个数据库', state.dbStats.database_pretty || formatBytes(state.dbStats.database_size)],
      ['上次 VACUUM', state.dbStats.last_vacuum || '-'],
      ['上次 AUTOVACUUM', state.dbStats.last_autovacuum || '-'],
      ['上次 ANALYZE', state.dbStats.last_analyze || '-'],
      ['上次 AUTOANALYZE', state.dbStats.last_autoanalyze || '-']
    ];
    refs.dbStatsGrid.innerHTML = cards.map(function (card) {
      return '<div class="stats-card"><span>' + escapeHTML(card[0]) + '</span><strong>' + escapeHTML(card[1]) + '</strong></div>';
    }).join('');
  }

  async function loadDBStats() {
    refs.refreshDBStats.disabled = true;
    try {
      state.dbStats = await requestJSON(apiPath('/api/db/stats'));
      renderDBStats();
    } catch (err) {
      if (!isAbortError(err)) {
        showMessage(refs.dbMessage, err.message || '加载数据库统计失败。', 'error');
      }
    } finally {
      refs.refreshDBStats.disabled = false;
    }
  }

  function showMessage(el, text, type) {
    if (!el) {
      return;
    }
    if (!text) {
      el.textContent = '';
      el.className = 'message-strip hidden';
      return;
    }
    el.textContent = text;
    el.className = 'message-strip ' + (type || '');
  }

  async function saveAlias(fingerprint, alias) {
    var value = (alias || '').trim();
    if (!fingerprint) {
      showMessage(refs.settingsMessage, '请先提供 token 指纹。', 'error');
      return;
    }
    refs.saveAliasButton.disabled = true;
    try {
      var body = new URLSearchParams();
      body.set('token_fingerprint', fingerprint);
      body.set('token_alias', value);
      await requestJSON(apiPath('/api/tokens/alias'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: body.toString()
      });
      showMessage(refs.settingsMessage, value ? '令牌代号已保存。' : '令牌代号已清除。', 'success');
      await loadFilterOptions();
      await loadTokenDirectory();
      await loadDashboard(true);
      await loadList(false, true);
      if (state.selectedId) {
        await loadDetail(state.selectedId);
      }
    } catch (err) {
      showMessage(refs.settingsMessage, err.message || '保存令牌代号失败。', 'error');
    } finally {
      refs.saveAliasButton.disabled = false;
    }
  }

  async function runDBMaintenance(mode) {
    var label = mode === 'vacuum_full' ? '强制缩盘' : mode === 'compact_payloads' ? '压缩历史大字段' : mode === 'analyze' ? '仅刷新统计' : '整理空间';
    if (mode === 'vacuum_full' && !window.confirm('执行 VACUUM FULL 会锁住 audit_logs 表，确认继续吗？')) {
      return;
    }
    var buttons = document.querySelectorAll('[data-db-maintenance]');
    buttons.forEach(function (button) {
      button.disabled = true;
    });
    try {
      var body = new URLSearchParams();
      body.set('mode', mode);
      var payload = await requestJSON(apiPath('/api/db/maintenance'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: body.toString()
      });
      showMessage(refs.dbMessage, payload.message || (label + '已完成。'), 'success');
      await loadDBStats();
    } catch (err) {
      showMessage(refs.dbMessage, err.message || (label + '失败。'), 'error');
    } finally {
      buttons.forEach(function (button) {
        button.disabled = false;
      });
    }
  }

  async function runDBCleanup() {
    var from = refs.cleanupFrom.value.trim();
    var to = refs.cleanupTo.value.trim();
    var token = refs.cleanupToken.value.trim();
    var alias = refs.cleanupAlias.value.trim();
    var model = refs.cleanupModel.value.trim();
    if (!(from || to || token || alias || model)) {
      showMessage(refs.dbMessage, '请至少填写一个清理条件。', 'error');
      return;
    }
    if (!window.confirm('确认按当前条件清理审计日志吗？这一步不可撤销。')) {
      return;
    }
    var button = document.getElementById('run-db-cleanup');
    button.disabled = true;
    try {
      var body = new URLSearchParams();
      if (from) {
        body.set('from', from);
      }
      if (to) {
        body.set('to', to);
      }
      if (token) {
        body.set('token', token);
      }
      if (alias) {
        body.set('alias', alias);
      }
      if (model) {
        body.set('model', model);
      }
      var payload = await requestJSON(apiPath('/api/db/cleanup'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: body.toString()
      });
      showMessage(refs.dbMessage, payload.message || '日志清理完成。', 'success');
      await Promise.all([
        loadDBStats(),
        loadDashboard(true),
        loadFilterOptions(),
        loadTokenDirectory(),
        loadList(false, true)
      ]);
    } catch (err) {
      showMessage(refs.dbMessage, err.message || '清理日志失败。', 'error');
    } finally {
      button.disabled = false;
    }
  }

  async function loadList(triggerPulse, keepSelection) {
    syncPageStateFromURL();
    var controller = getController('list');
    var params = buildLogQuery(true);
    if (!triggerPulse && !keepSelection) {
      setRefreshState('正在同步日志…', 'is-checking');
    }
    var previousIDs = state.items.map(function (item) { return item.id; });
    var selectedBefore = state.selectedId;
    try {
      var payload = await requestJSON(apiPath('/api/logs' + (params.toString() ? '?' + params.toString() : '')), {
        signal: controller.signal
      });
      state.version = payload.version || state.version;
      state.items = payload.items || [];
      state.totalCount = Number(payload.total_count || 0);
      state.pageSize = parsePositiveInt(payload.page_size) || state.pageSize || 100;
      state.totalPages = parsePositiveInt(payload.total_pages) || 1;
      state.page = Math.min(parsePositiveInt(payload.page) || parsePositiveInt(state.page) || 1, state.totalPages);
      state.hasPrev = !!payload.has_prev || state.page > 1;
      state.hasNext = !!payload.has_next || state.page < state.totalPages;
      var currentIDs = state.items.map(function (item) { return item.id; });
      var previousSet = {};
      previousIDs.forEach(function (id) {
        previousSet[id] = true;
      });
      var freshIDs = currentIDs.filter(function (id) {
        return !previousSet[id];
      });
      var visibleFreshIDs = triggerPulse ? freshIDs : [];
      if (keepSelection && selectedBefore && currentIDs.indexOf(selectedBefore) >= 0) {
        state.selectedId = selectedBefore;
      } else if (!state.selectedId || currentIDs.indexOf(state.selectedId) < 0) {
        state.selectedId = currentIDs[0] || 0;
      }
      renderLogList(visibleFreshIDs, keepSelection || triggerPulse);
      updateFilterSummary();
      if (visibleFreshIDs.length) {
        state.lastLogUpdateAt = Date.now();
        pulseStatCard(refs.quickRequestCard);
        refs.refreshState.classList.remove('has-updates');
        window.requestAnimationFrame(function () {
          refs.refreshState.classList.add('has-updates');
        });
      } else if (!state.lastLogUpdateAt) {
        state.lastLogUpdateAt = Date.now();
      }
      if (state.selectedId) {
        if (!state.detail || state.detail.id !== state.selectedId) {
          loadDetail(state.selectedId);
        } else {
          syncLogSelectionState();
        }
      } else {
        state.detail = null;
        renderDetail();
      }
      if (visibleFreshIDs.length) {
        setRefreshState('发现 ' + visibleFreshIDs.length + ' 条变化 · 刚刚', 'has-updates');
      } else {
        renderLastUpdateState('');
      }
      await loadDashboard(true);
    } catch (err) {
      if (isAbortError(err)) {
        return;
      }
      setRefreshState(err.message || '日志加载失败', 'is-error');
      refs.logGroups.innerHTML = '' +
        '<div class="lane-empty">' +
          '<div>' +
            '<h3>日志列表加载失败</h3>' +
            '<p>' + escapeHTML(err.message || '未知错误') + '</p>' +
          '</div>' +
        '</div>';
    }
  }

  async function checkVersion() {
    var controller = getController('version');
    try {
      var params = buildLogQuery(false);
      var payload = await requestJSON(apiPath('/api/logs/version' + (params.toString() ? '?' + params.toString() : '')), {
        signal: controller.signal
      });
      if (payload.version && payload.version !== state.version) {
        await loadList(true, true);
      } else {
        renderLastUpdateState('');
      }
    } catch (err) {
      if (!isAbortError(err)) {
        setRefreshState(err.message || '更新检查失败', 'is-error');
      }
    }
  }

  function openModal(id) {
    var modal = document.getElementById(id);
    if (!modal) {
      return;
    }
    modal.classList.add('show');
  }

  function closeModal(id) {
    var modal = document.getElementById(id);
    if (!modal) {
      return;
    }
    modal.classList.remove('show');
  }

  function closeAllSuggestMenus() {
    document.querySelectorAll('.suggest-menu.show').forEach(function (menu) {
      menu.classList.remove('show');
      menu.innerHTML = '';
    });
    document.querySelectorAll('.suggest-shell.is-open').forEach(function (shell) {
      shell.classList.remove('is-open');
    });
    document.querySelectorAll('.filter-field.is-open').forEach(function (field) {
      field.classList.remove('is-open');
    });
    state.activeSuggestInput = '';
  }

  function getSuggestOptions(kind) {
    var list = state.filterOptions[kind] || [];
    return list.filter(function (item) {
      return !!item;
    });
  }

  function openSuggestForInput(input) {
    if (!input) {
      return;
    }
    var kind = input.getAttribute('data-suggest');
    var menu = document.getElementById('suggest-' + input.id);
    if (!kind || !menu) {
      return;
    }
    var shell = input.closest('.suggest-shell');
    var field = input.closest('.filter-field');
    var keyword = input.value.trim().toLowerCase();
    var options = getSuggestOptions(kind).filter(function (item) {
      return !keyword || String(item).toLowerCase().indexOf(keyword) >= 0;
    }).slice(0, 30);
    closeAllSuggestMenus();
    if (!options.length) {
      return;
    }
    if (shell) {
      shell.classList.add('is-open');
    }
    if (field) {
      field.classList.add('is-open');
    }
    menu.innerHTML = options.map(function (item) {
      return '<button type="button" class="suggest-option" data-suggest-value="' + escapeHTML(item) + '" data-target-input="' + escapeHTML(input.id) + '">' + escapeHTML(item) + '</button>';
    }).join('');
    menu.classList.add('show');
    state.activeSuggestInput = input.id;
  }

  function bindSuggestInput(input) {
    if (!input) {
      return;
    }
    input.addEventListener('focus', function () {
      openSuggestForInput(input);
    });
    input.addEventListener('input', function () {
      updateSuggestClearButtons();
      openSuggestForInput(input);
    });
    input.addEventListener('keydown', function (event) {
      if (event.key === 'Escape') {
        closeAllSuggestMenus();
      }
    });
  }

  function resizeLogListViewport() {
    if (!refs.logGroups) {
      return;
    }
    var wrap = refs.logGroups.closest('.log-list-frame');
    if (!wrap) {
      return;
    }
    state.logListHeight = 0;
    wrap.style.height = '';
    wrap.style.minHeight = '';
    refs.logGroups.style.height = '';
    refs.logGroups.style.maxHeight = '';
    refs.logGroups.style.overflowY = 'auto';
  }

  function updatePaginationControls() {
    var pageSize = parsePositiveInt(state.pageSize) || 100;
    var computedTotalPages = state.totalCount > 0 ? Math.ceil(Math.max(0, Number(state.totalCount || 0)) / pageSize) : 1;
    state.totalPages = Math.max(parsePositiveInt(state.totalPages) || 1, computedTotalPages);
    state.page = Math.min(parsePositiveInt(state.page) || 1, state.totalPages);
    state.hasPrev = state.page > 1;
    state.hasNext = state.page < state.totalPages;

    refs.listCurrentCount.textContent = String(state.items.length);
    refs.listTotalCount.textContent = formatNumber(state.totalCount);
    refs.listPageText.textContent = '第 ' + state.page + ' / ' + state.totalPages + ' 页';
  }

  function attachEventHandlers() {
    resizeLogListViewport();
    window.addEventListener('resize', resizeLogListViewport);

    refs.filtersForm.addEventListener('submit', function (event) {
      event.preventDefault();
      syncStateFromFilterInputs();
      window.location.assign(buildLogPageURL(1));
    });

    refs.clearFilters.addEventListener('click', function () {
      state.filters = {
        from: '',
        to: '',
        alias: '',
        model: '',
        status: '',
        q: ''
      };
      syncFilterInputsFromState();
      closeAllSuggestMenus();
      window.location.assign(buildLogPageURL(1));
    });

    refs.logGroups.addEventListener('click', function (event) {
      var selectTarget = event.target.closest('[data-select-log]');
      if (selectTarget) {
        var id = parsePositiveInt(selectTarget.getAttribute('data-select-log'));
        if (id) {
          loadDetail(id);
        }
        return;
      }
      var toggleTarget = event.target.closest('[data-toggle-group]');
      if (toggleTarget) {
        var key = toggleTarget.getAttribute('data-toggle-group');
        toggleGroupCollapse(key);
        renderLogList([], true);
      }
    });

    refs.detailStage.addEventListener('click', function (event) {
      var pageButton = event.target.closest('[data-text-kind][data-text-page]');
      if (pageButton) {
        var kind = pageButton.getAttribute('data-text-kind');
        var page = parsePositiveInt(pageButton.getAttribute('data-text-page'));
        if (kind && page) {
          loadTextPageV2(kind, page);
        }
        return;
      }
      var rawButton = event.target.closest('[data-load-raw]');
      if (rawButton) {
        var id = parsePositiveInt(rawButton.getAttribute('data-load-raw'));
        if (id) {
          ensureRawLoaded(id, false);
        }
      }
    });

    refs.detailStage.addEventListener('toggle', function (event) {
      var fold = event.target;
      if (!fold || fold.tagName !== 'DETAILS') {
        return;
      }
      if (fold.hasAttribute('data-needs-raw') && fold.open && state.detail) {
        ensureRawLoaded(state.detail.id, false);
      }
    }, true);

    refs.openStatsModal.addEventListener('click', function () {
      openModal('stats-modal');
      loadDashboard(true);
    });

    refs.openSettingsModal.addEventListener('click', function () {
      if (state.detail && state.detail.token_fingerprint) {
        refs.aliasFingerprint.value = state.detail.token_fingerprint;
        if (state.detail.token_alias) {
          refs.aliasName.value = state.detail.token_alias;
        }
      }
      openModal('settings-modal');
      loadTokenDirectory();
    });

    refs.openDBModal.addEventListener('click', function () {
      refs.cleanupFrom.value = state.filters.from;
      refs.cleanupTo.value = state.filters.to;
      refs.cleanupToken.value = '';
      refs.cleanupAlias.value = state.filters.alias;
      refs.cleanupModel.value = state.filters.model;
      openModal('db-modal');
      loadDBStats();
    });

    document.querySelectorAll('[data-close-modal]').forEach(function (button) {
      button.addEventListener('click', function () {
        closeModal(button.getAttribute('data-close-modal'));
      });
    });

    [refs.statsModal, refs.settingsModal, refs.dbModal].forEach(function (modal) {
      modal.addEventListener('click', function (event) {
        if (event.target === modal) {
          modal.classList.remove('show');
        }
      });
    });

    document.addEventListener('keydown', function (event) {
      if (event.key === 'Escape') {
        closeAllSuggestMenus();
        closeModal('stats-modal');
        closeModal('settings-modal');
        closeModal('db-modal');
      }
    });

    refs.aliasForm.addEventListener('submit', function (event) {
      event.preventDefault();
      saveAlias(refs.aliasFingerprint.value.trim(), refs.aliasName.value.trim());
    });

    refs.tokenDirectoryBody.addEventListener('click', function (event) {
      var button = event.target.closest('[data-save-row-alias]');
      if (!button) {
        return;
      }
      var fingerprint = button.getAttribute('data-save-row-alias');
      var row = button.closest('tr');
      var input = row ? row.querySelector('input') : null;
      saveAlias(fingerprint, input ? input.value.trim() : '');
    });

    refs.refreshTokenDirectory.addEventListener('click', function () {
      loadTokenDirectory();
    });

    refs.refreshDBStats.addEventListener('click', function () {
      loadDBStats();
    });

    document.querySelectorAll('[data-db-maintenance]').forEach(function (button) {
      button.addEventListener('click', function () {
        runDBMaintenance(button.getAttribute('data-db-maintenance'));
      });
    });

    refs.dbCleanupForm.addEventListener('submit', function (event) {
      event.preventDefault();
      runDBCleanup();
    });

    refs.copyMainFiltersToCleanup.addEventListener('click', function () {
      refs.cleanupFrom.value = state.filters.from;
      refs.cleanupTo.value = state.filters.to;
      refs.cleanupToken.value = '';
      refs.cleanupAlias.value = state.filters.alias;
      refs.cleanupModel.value = state.filters.model;
      showMessage(refs.dbMessage, '已把右侧当前筛选复制到数据库清理条件。', 'success');
    });

    document.addEventListener('click', function (event) {
      var option = event.target.closest('[data-suggest-value][data-target-input]');
      if (option) {
        var input = document.getElementById(option.getAttribute('data-target-input'));
        if (input) {
          input.value = option.getAttribute('data-suggest-value');
          updateSuggestClearButtons();
          closeAllSuggestMenus();
          input.focus();
        }
        return;
      }
      var clearButton = event.target.closest('[data-clear-filter]');
      if (clearButton) {
        var targetInput = document.getElementById(clearButton.getAttribute('data-clear-filter'));
        if (targetInput) {
          targetInput.value = '';
          updateSuggestClearButtons();
          closeAllSuggestMenus();
          targetInput.focus();
        }
        return;
      }
      if (!event.target.closest('.suggest-shell')) {
        closeAllSuggestMenus();
      }
    });

    bindSuggestInput(refs.filterAlias);
    bindSuggestInput(refs.filterModel);
    bindSuggestInput(refs.filterStatus);
  }

  async function bootstrap() {
    attachEventHandlers();
    syncFilterInputsFromState();
    updateFilterSummary();
    await Promise.all([
      loadFilterOptions(),
      loadDashboard(true)
    ]);
    await loadList(false, true);
    state.pollingTimer = window.setInterval(function () {
      checkVersion();
    }, 5000);
  }

  bootstrap();
})();
</script>
{{end}}
`
