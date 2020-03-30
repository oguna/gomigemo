# gomigemo

MigemoをGo言語で実装したものです。
辞書ファイルの構造を効率化し、メモリ使用量の削減を図りました。

| 項目 | C/Migemo | gomigemo |
| ---- | ---- | ---- |
| 実行ファイルサイズ | **72 KB** | 1.86 MB |
| 辞書ファイルサイズ | 4.78 MB | **2.03 MB** |
| メモリ使用量 | 26.1 MB | **10.9 MB** |

## 使い方

```
> go run main.go -h
Usage of C:\..\main.exe:
  -d string
        Use a file <dict> for dictionary. (default "migemo-compact-dict")
  -e    Use emacs style regexp.
  -n    Don't use newline match.
  -q    Show no message except results.
  -v    Use vim style regexp.
  -w string
        Expand a <word> and soon exit.
exit status 2
> go run main.go -w kensaku
(kensaku|けんさく|ケンサク|建策|憲[作冊]|検索|献策|研削|羂索|ｋｅｎｓａｋｕ|ｹﾝｻｸ)
```

## ライセンス
`main.go`ファイル及び`migemo`ディレクトリは、**MIT License**の下で配布しています。
`migemo-compact-dict`ファイルは、SKK辞書のライセンスを継承し、**GPL**の下で配布しています。