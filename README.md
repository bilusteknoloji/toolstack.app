![Version](https://img.shields.io/badge/version-0.1.18-orange.svg)
[![Deploy Toolstack Apps](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/build-and-deploy.yml/badge.svg)](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/build-and-deploy.yml)
[![Dependabot Updates](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/dependabot/dependabot-updates)

# toolstack.app

Official website of toolstack.app

---

## Sites

- https://toolstack.app
- https://ibankeeper.toolstack.app
- https://ip.toolstack.app
- https://whatismyip.toolstack.app/

Find you remote IP via;

```bash
curl ip.toolstack.app
curl whatismyip.toolstack.app
http ip.toolstack.app
http whatismyip.toolstack.app

wget -qO- https://ip.toolstack.app
wget -qO- http://ip.toolstack.app
wget -qO- https://whatismyip.toolstack.app
wget -qO- http://whatismyip.toolstack.app
```

---

## Rake Tasks

```bash
$ rake -T

rake docker:build       # build docker image locally
rake docker:run         # run docker image locally
rake release[revision]  # release new version major,minor,patch, default: patch
rake run:infra          # run orbstack infra
rake run:server         # run server (default: :8000)
```

---

## Development

You need [orbstack](https://orbstack.dev/) (macOS) to accomplish this. Run the
`rake` in a tab, open other tab and run `rake run:infra` then open any of
these urls:

- https://proxy.local.orb.local
- https://ibankeeper.local.orb.local/
- https://ip.local.orb.local/
- https://reminder.local.orb.local/

---

## License

This project is licensed under MIT (MIT)

---

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].

[coc]: https://github.com/bilusteknoloji/toolstack.app/blob/main/CODE_OF_CONDUCT.md
