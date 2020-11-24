gauth: replace Google Authenticator
===================================

Usage
-----

- In web interfaces, pretend you can't read QR codes, get a secret like `hret 3ij7 kaj4 2jzg` instead.
- Store one secret per line in `~/.gauth/gauth.csv`, in the format `name:secret`. For example:

        AWS:   ABCDEFGHIJKLMNOPQRSTUVWXYZ234567ABCDEFGHIJKLMNOPQRSTUVWXYZ234567
        Airbnb:abcd efgh ijkl mnop
        Google:a2b3c4d5e6f7g8h9
        Github:234567qrstuvwxyz

- Run `gauth`. The progress bar indicates how far the next change is.

        $ gauth
                   prev   curr   next
        AWS        315306 135387 483601
        Airbnb     563728 339206 904549
        Google     453564 477615 356846
        Github     911264 548790 784099
        [=======                      ]


- Remember to keep your system clock synchronized and to **lock your computer when brewing your tea!**

Encryption
----------

`gauth` supports password-based encryption of `gauth.csv`. To encrypt, use:

        $ openssl enc -aes-128-cbc -md sha256 -in gauth.csv -out ~/.gauth/gauth.csv
        enter aes-128-cbc encryption password:
        Verifying - enter aes-128-cbc encryption password:

`gauth` will then prompt you for that password on every run:

        $ gauth
        Encryption password: 
                   prev   curr   next
        LastPass   915200 479333 408710

Note that this encryption mechanism is far from ideal from a pure security standpoint.
Please read [OpenSSL's notes on the subject](http://www.openssl.org/docs/crypto/EVP_BytesToKey.html#NOTES).

Build
-----

You can build the binary by using the make commands (you must have make and go installed)

For Linux binary :

```make build```

For Windows binary :

```make build_windows```

For Mac binary :

```make build_mac```

Compatibility
-------------

Tested with:

- Airbnb
- Apple
- AWS
- DreamHost
- Dropbox
- Evernote
- Facebook
- Gandi
- Github
- Google
- LastPass
- Linode
- Microsoft
- Okta (reported by Bryan Baldwin)
- WP.com
- bittrex.com
- poloniex.com
- Slack

Please report further results to pierre@gcarrier.fr.
