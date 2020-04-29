# gomigemo

ローマ字のまま日本語をインクリメンタル検索するためのツールであるMigemoを、Go言語で実装したものです。

## C/Migemoとの比較

| 項目 | C/Migemo | gomigemo |
| ---- | ---- | ---- |
| 実行ファイルサイズ(KB) | **72** | 2100 |
| 辞書ファイルサイズ(MB) | 4.78 MB | **2.03** |
| メモリ使用量(MB) | 26.1 | **11.1** |
| 起動時間(ms) | 198 | **78** |
| クエリ時間(ms) | **5053** | 9671 |

詳細は[migemo-benchmark](https://github.com/oguna/migemo-benchmark)を参照してください。

## ビルド方法

```
> go build -ldflags="-s -w" -trimpath
```

## 使い方

### CLI

gomigemoの利用には、辞書ファイルが必要です。
[migemo-compact-dict-latest](https://github.com/oguna/migemo-compact-dict-latest)
から `migemo-compact-dict` をダウンロードし、
作業フォルダ（シェルのカレントディレクトリ）に配置してください。
なお、`migemo-compact-dict` のライセンスはGPLv3のため、再配布する際はご注意ください。

```
> ./gomigemo.exe -h
Usage of C:\...\gomigemo.exe:
  -d string
        Use a file <dict> for dictionary. (default "migemo-compact-dict")
  -e    Use emacs style regexp.
  -n    Don't use newline match.
  -p int
        <port> number for HTTP server.
  -q    Show no message except results.
  -v    Use vim style regexp.
  -w string
        Expand a <word> and soon exit.
> ./gomigemo.exe -w kensaku
(kensaku|けんさく|ケンサク|建策|憲[作冊]|検索|献策|研削|羂索|ｋｅｎｓａｋｕ|ｹﾝｻｸ)
```

### HTTP

gomigemoはHTTPサーバとして動かすことができます。
コマンドライン引数 `-p <port>` で、サーバのポート番号を指定すると、HTTPサーバとして起動します。

```
> ./gomigemo.exe -p 8080
```

指定したポート番号に `/migemo?q=<word>` にGETメソッドでアクセスすると、
`<word>` の文字列がMigemoの正規表現として返ってきます。
クエリ文字列`q`を複数つなげると、複数Migemoの処理が同時に行われ、順に改行されて返されます。

```
> ./curl.exe  "localhost:8080/migemo?q=migemo&q=kensaku"
(migemo|みげも|ミゲモ|ｍｉｇｅｍｏ|ﾐｹﾞﾓ)
(kensaku|けんさく|ケンサク|建策|憲[作冊]|検索|献策|研削|羂索|ｋｅｎｓａｋｕ|ｹﾝｻｸ)
```

## ライセンス

**MIT License**の下で配布しています。