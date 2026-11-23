# churkabot-tg

A Telegram bot for a group of friends: it answers trigger words in group chats with stickers and posts a sticker every morning.

## Behaviour

**Groups.** A message whose text or caption contains a trigger word gets a sticker as a reply. Rules match whole words in a listed set of forms (`точно`, `точный`, `точняк`; not `достаточно`), lower-cased with `ё` folded to `е`; the first rule in [`internal/rules/rules.go`](internal/rules/rules.go) wins. A message counts as shouting when its upper-case words add up to six letters or more; short Latin acronyms (`API`, `HTTPS`) do not count. Words inside links, mentions, commands and code are ignored.

The bot stays silent on messages from other bots, on channel posts forwarded into a linked group, on messages older than ten minutes (Telegram re-delivers a day of updates after downtime), and for fifteen seconds after its last sticker in the same chat.

**Private chat.** Only admins get answers: a sticker or a photo is answered with its `file_id`, text with the rule it would trigger in a group. Anyone else gets one line on `/start` and nothing otherwise.

**Daily sticker.** `churkabot daily` sends one sticker and exits; the schedule belongs to whatever runs it (a Kubernetes CronJob). If the chat has been upgraded to a supergroup, it sends to the new ID and tells the admins to update the setting.

## Configuration

| Variable | Required | Meaning |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | yes | Bot token from BotFather |
| `CHURKABOT_ADMIN_IDS` | no | Comma-separated numeric user IDs of admins |
| `CHURKABOT_DAILY_CHAT_ID` | for `daily` | Chat that gets the daily sticker |

Privacy mode must be off in BotFather (`/setprivacy` → Disable), or in groups the bot sees only commands; it logs a warning at start when it is on.

Exit codes: `1` for a failure worth retrying, `2` for one that is not (bad configuration, revoked token, bot removed from the chat).

## Development

Tools are pinned in `mise.toml`.

```sh
mise run ci      # lint, test, build
```

Do not run the bot locally with the production token: Telegram serves `getUpdates` to one client at a time, so a second poller steals updates from the deployed one.

Adding a trigger: add the forms to the table in `internal/rules/rules.go` and examples to `internal/rules/rules_test.go`, including words that contain the root but must not fire.

## Releases

Conventional Commits; release-please opens the release PR, and merging it tags the version. CI publishes `ghcr.io/sanchpet/churkabot-tg:<version>`.
