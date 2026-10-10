import { test, expect } from "@playwright/test";
import fs from "node:fs";
import path from "node:path";
const bindings = fs.readFileSync(
  path.resolve("bindings/github.com/SUDOKU-ASCII/sudoku-desktop/app.ts"),
  "utf8",
);
const ids = Object.fromEntries(
  [...bindings.matchAll(/export function (\w+)\([^]*?\$Call.ByID\((\d+)/g)].map(
    (m) => [Number(m[2]), m[1]],
  ),
);

test("import, connection state, settings drafts, routing and themes", async ({
  page,
}) => {
  let cfg: any = {
    version: 6,
    activeNodeId: "",
    nodes: [],
    routing: {
      proxyMode: "global",
      ruleUrls: [],
      customRulesEnabled: false,
      customRules: "",
    },
    tun: { enabled: false },
    core: { autoStart: false },
    ui: { theme: "auto", language: "ru", launchAtLogin: false },
    reverseClient: { routes: [] },
    reverseForward: {},
    portForwards: [],
  };
  let state: any = {
    running: false,
    coreRunning: false,
    tunRunning: false,
    startedAtUnix: 0,
    traffic: { recentBandwidth: [] },
    kernel: {},
    latencies: [],
    connections: [],
    recentLogs: [],
  };
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));
  await page.route("**/wails/runtime", async (route) => {
    const req = route.request().postDataJSON();
    let result: any = null;
    if (req.object === 0) {
      const method = ids[req.args.methodID],
        args = req.args.args || [];
      if (method === "GetConfig") result = cfg;
      else if (method === "GetState") result = state;
      else if (method === "SaveConfig") cfg = args[0];
      else if (method === "ImportShortLink") {
        const n = {
          id: "node1",
          name: args[0] || "Мой сервер",
          serverAddress: "vpn.example.com:443",
          key: "test",
          sourceUrl: "https://example.com/sub/test.yaml",
          httpMask: {},
          customTables: [],
        };
        cfg.nodes = [n];
        cfg.activeNodeId = n.id;
        result = n;
      } else if (method === "StartProxy")
        state = {
          ...state,
          running: true,
          coreRunning: true,
          tunRunning: args[0].withTun,
          startedAtUnix: Math.floor(Date.now() / 1000),
        };
      else if (method === "StopProxy")
        state = {
          ...state,
          running: false,
          coreRunning: false,
          tunRunning: false,
        };
      else if (method === "DeleteNode") {
        cfg.nodes = [];
        cfg.activeNodeId = "";
      }
    }
    await route.fulfill({
      contentType: "application/json",
      body: JSON.stringify(result),
    });
  });
  await page.route("**/wails/custom.js", (r) =>
    r.fulfill({ status: 404, body: "" }),
  );
  await page.goto("/");
  await expect(page.getByText("Добавьте сервер, чтобы начать")).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Подключить", exact: true }),
  ).toBeDisabled();
  await page.getByRole("button", { name: /Добавить соединение/ }).click();
  await page
    .getByLabel("Ссылка соединения")
    .fill("https://example.com/sub/test.yaml");
  await page.getByLabel("Название").fill("Нидерланды");
  await page
    .getByRole("button", { name: "Добавить соединение", exact: true })
    .click();
  await expect(page.getByText("Нидерланды", { exact: true })).toBeVisible();
  await page.getByRole("button", { name: "Подключить", exact: true }).click();
  await expect(page.getByText("Соединение активно")).toBeVisible();
  await page.getByRole("button", { name: "Настройки", exact: true }).click();
  await expect(
    page.getByRole("checkbox", { name: "VPN для всего устройства" }),
  ).toHaveCount(0);
  await expect(
    page.getByRole("switch", { name: "VPN для всего устройства" }),
  ).toBeDisabled();
  await page.getByRole("button", { name: "Назад", exact: true }).click();
  await page.getByRole("button", { name: "Отключить", exact: true }).click();
  await page.getByRole("button", { name: "Настройки", exact: true }).click();
  await page.getByRole("switch", { name: "Подключаться при запуске" }).check();
  await page.getByRole("button", { name: "Тёмная", exact: true }).click();
  expect(cfg.core.autoStart).toBe(false);
  await page
    .getByRole("combobox", { name: "Режим" })
    .selectOption({ label: "Все, кроме списка" });
  await page.getByLabel("Домены исключений").fill("https://bad.example");
  await page.getByRole("button", { name: "Сохранить настройки" }).click();
  await expect(page.getByRole("alert")).toContainText("Укажите домены");
  await page.getByLabel("Домены исключений").fill("hh.ru\navito.ru\nпример.рф");
  await page.getByRole("button", { name: "Сохранить настройки" }).click();
  expect(cfg.core.autoStart).toBe(true);
  expect(cfg.routing.customRules).toContain("DOMAIN-SUFFIX,hh.ru,DIRECT");
  expect(cfg.routing.ruleUrls).toEqual([]);
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await page.evaluate(() => document.fonts.ready);
  expect(
    await page.locator("body").evaluate((e) => getComputedStyle(e).fontFamily),
  ).toContain("Montserrat");
  if (process.env.SUDOKU_PREVIEW_DIR) {
    await page.evaluate(() => window.scrollTo(0, 0));
    await page.screenshot({
      fullPage: true,
      path: path.join(process.env.SUDOKU_PREVIEW_DIR, "sudoku-dark.png"),
    });
    await page.getByRole("button", { name: "Настройки", exact: true }).click();
    await page.evaluate(() => window.scrollTo(0, 0));
    await page.screenshot({
      fullPage: true,
      path: path.join(process.env.SUDOKU_PREVIEW_DIR, "settings-dark.png"),
    });
    await page.getByRole("button", { name: "Светлая", exact: true }).click();
    await page.evaluate(() => window.scrollTo(0, 0));
    await page.screenshot({
      fullPage: true,
      path: path.join(process.env.SUDOKU_PREVIEW_DIR, "settings-light.png"),
    });
  }
  expect(errors).toEqual([]);
});
