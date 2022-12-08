# Vkonity

Vkonity is a service for searching for new posts from VKontakte groups and sending them a short preview in private
messages.

# Installation

1. Grab the latest version from the releases
2. Create a config file
3. Run with `vkonity -config config.toml`

# Configuration

Vkonity uses `config.toml` file in the working directory as the config. You can change it by the flag `-config foo.toml`
.

```toml
# Notification chat IDs
Admins = [
    1,
]
# Listen group IDs
Groups = [
    22822305, #@vk
    1, #@apiclub
]

# User or service VK token
ServiceToken = ""
# VK Bot token
MessageToken = ""

# Frequency of requests to API VK. Valid time units are "ms", "s", "m", "h".
Period = "10s"
# Message template sent to PM
Message = "🆕 Новый пост в @public%v\n🌎 Быстрый переход: https://vk.com/wall%v_%v\n🖌 Содержание:\n%s"
```

# Contributors

Vkonity uses some resources from [Acamar](https://github.com/xtrafrancyz/acamar)