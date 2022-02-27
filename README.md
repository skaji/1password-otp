# 1password otp

This is a simple tool that generates all One-Time Password (OTP) from 1Password Interchange Format (.1pif).

# Usage

First export 1pif file from 1Passwrod application. Then execute `one-password-otp`:

```console
❯ 1password-otp data.1pif
* amazon (you@example.com)
  155785
* aws (you@example.com)
  683458
* dropbox (you@exmaple.com)
  725998
```

You can also use `-json` option so that the output is in JSON:

```console
❯ 1password-otp -json data.1pif
[
  {
    "Title": "amazon",
    "Username": "you@example.com",
    "URI": "otpauth://totp/Amazon:you@example.com?secret=XXXXXXXXXXXXXXX\u0026issuer=Amazon",
    "Number": "448598",
    "Expiration": 1645985490
  },
  {
    "Title": "aws",
    "Username": "you@example.com",
    "URI": "otpauth://totp/Amazon%20Web%20Services:root-account-mfa-device@123456?secret=XXXXXXXXXXXX\u0026issuer=Amazon%20Web%20Services",
    "Number": "136824",
    "Expiration": 1645985490
  },
  {
    "Title": "dropbox",
    "Username": "you@example.com",
    "URI": "otpauth://totp/Dropbox:you@example.com?secret=XXXXXXXXXXXXXX4\u0026issuer=Dropbox",
    "Number": "001862",
    "Expiration": 1645985490
  }
]
```

# Install

Download appropriate binaries from [the release page](https://github.com/skaji/1password-otp/releases/latest).

# Author

Shoichi Kaji

# License

Apache 2.0
