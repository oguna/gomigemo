# gomigemo

MigemoをGo言語で実装したものです。
辞書ファイルの構造を効率化し、メモリ使用量の削減を図りました。

## C/Migemoとの比較

| 項目 | C/Migemo | gomigemo |
| ---- | ---- | ---- |
| 実行ファイルサイズ | **72 KB** | 1.86 MB |
| 辞書ファイルサイズ | 4.78 MB | **2.03 MB** |
| メモリ使用量 | 26.1 MB | **10.9 MB** |
| 起動時間 | 141 ms | **60 ms** |
| 検索時間 | **1.738 s** | 4.734 s |

### 検索時間

夏目漱石「こころ」に含まれている4524個のルビをローマ字で入力し、
すべての正規表現の出力に要した時間を比較しています。
ベンチマークの設定等は公開予定です。

平均して約1msで1件の検索が完了するので、実用的な処理速度です。

## ビルド方法

```
> go build -ldflags="-s -w" -trimpath
```

## 使い方

```
> ./gomigemo.exe -h
Usage of C:\...\main.exe:
  -d string
        Use a file <dict> for dictionary. (default "migemo-compact-dict")
  -e    Use emacs style regexp.
  -n    Don't use newline match.
  -q    Show no message except results.
  -v    Use vim style regexp.
  -w string
        Expand a <word> and soon exit.
exit status 2
> ./gomigemo.exe -w kensaku
(kensaku|けんさく|ケンサク|建策|憲[作冊]|検索|献策|研削|羂索|ｋｅｎｓａｋｕ|ｹﾝｻｸ)
```

## ライセンス

`main.go`ファイル及び`migemo`ディレクトリは、**MIT License**の下で配布しています。

`migemo-compact-dict`ファイルは、SKK辞書のライセンスを継承し、**GPL**の下で配布しています。