# lego-mydns

lego を使い MyDNS で Let's Encrypt ワイルドカード証明書を作る為のコマンド

## インストール方法

lego コマンドは以下でインストールします。

```
$ go get github.com/xenolf/lego
```

本コマンドは以下でインストールします。

```
$ go get github.com/mattn/go-mydns/cmd/lego-mydns
```

どちらも `$GOPATH/bin` (デフォルトは `$HOME/go/bin`) にインストールされるのでパスを通しておきます。

## 設定方法

まず MyDNS でアカウントを登録します。登録するとメールで MyDNS のマスターIDとパスワードが送られてきます。

このマスターIDとパスワードを環境変数に設定します。

### Windows の場合

```
set MYDNS_MASTERID=xxxxxxxxx
set MYDNS_PASSWORD=yyyyyyyyy
```

### UNIX の場合

```
export MYDNS_MASTERID=xxxxxxxxx
export MYDNS_PASSWORD=yyyyyyyyy
```

Windows で毎回設定するのが面倒な場合はシステム環境変数に登録して下さい。

## 使い方

初回だけ lego を run コマンドで実行する必要があります。例えば MyDNS で取得したドメイン名が `sugoi-domain.mydns.jp` であった場合は以下の様に実行します。

```
lego --accept-tos --dns exec --email メールアドレス --domains sugoi-domain.mydns.jp --domains *.sugoi-domain.mydns.jp run
```

以降は renew コマンドで3か月に1回くらい更新する必要があります。

```
lego --accept-tos --dns exec --email メールアドレス --domains sugoi-domain.mydns.jp --domains *.sugoi-domain.mydns.jp renew
```

Windows であればタスクスケジューラに登録しておくとよいでしょう。

生成された証明書は .lego/certificates に生成されます。必要であれば `--path` オプションを追加して下さい。

## ウェブサーバで使う

nginx の場合は以下の様に設定して下さい。

```
server {
    listen 443 ssl http2;
    server_name sugoi-domain.mydns.jp *.sugoi-domain.mydns.jp;

    ssl on;
    ssl_certificate /etc/letsencrypt/certificates/sugoi-domain.mydns.jp.crt;
    ssl_certificate_key /etc/letsencrypt/certificates/sugoi-domain.mydns.jp.key;

    ... 略 ...
}
```

## License

MIT

## Author

Yasuhiro Matsumoto
