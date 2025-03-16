![Version](https://img.shields.io/badge/version-0.0.9-orange.svg)
[![Deploy Toolstack Apps](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/build-and-deploy.yml/badge.svg)](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/build-and-deploy.yml)
[![Dependabot Updates](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/bilusteknoloji/toolstack.app/actions/workflows/dependabot/dependabot-updates)

# toolstack.app

Official website of toolstack.app

---

## Rake Tasks

```bash
$ rake -T

rake docker:build       # build docker image locally
rake docker:run         # run docker image locally
rake release[revision]  # release new version major,minor,patch, default: patch
rake run                # run server (default: :8000)
```

---

## License

This project is licensed under MIT (MIT)

---

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].

[coc]: https://github.com/bilusteknoloji/toolstack.app/blob/main/CODE_OF_CONDUCT.md
