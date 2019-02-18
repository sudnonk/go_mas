1. エージェント同士がつながりを持つ。（`Barabashi-Albert Scale Free Network`）
2. 繋がってるエージェント同士でイデオロギーの交流が行われる
3. エージェントは、交流のたびに`Agent.HP -= 受容性 * イデオロギーの差`
4. 交流のたびに `Agent.Mix()` の中身 のアルゴリズムで自身のイデオロギーが変化
    4.1. イデオロギーの変化量は確率で変化
5. ターンのたびにエージェントの体力は `Agent.RecoveryRate` で回復
6. 

受容性は0-1のfloat64で、正規分布

## main.exe
メインのプログラム。この`exe`と同じ階層に`config.json`を置くこと。

`main.exe ourDir initPath`

- `outDir` : ログの出力先ディレクトリ。末尾に`/`を付けること。
- `initPath` : 初期値ファイルへのパス。指定しない場合はランダムな初期値で実行される。

## parser.exe
ログをパースするプログラム

`parser.exe -f filename -o outDir -t targets --type type`

- `filename` : パースするファイルへのパス。複数指定可
- `outDir` : パース結果を出力するディレクトリ。末尾に`/`を付けること。
- `targets`: パース対象のエージェントのID。複数指定するときは`-t agentId`を複数回指定する。
- `type` : 以下のどれかを指定する
-- `fanatic` : 各思想にどれだけ信者がいるか。`-t`不要。
-- `hp` : 特定のエージェントの体力の推移。`-t`必須。
-- `ideology`: 特定のエージェントの思想の推移。`-t`必須。
-- `range` : 特定のエージェントがフォローしてるエージェントの思想の範囲の推移。`-t`必須。
-- `list`: 全エージェントの受容度、思想、フォロー数、体力、回復率の一覧。`-t`不要。
-- `diversity` : 信者が0人でない思想の数。`-t`不要。
-- `all` : 上記すべてを実行。`-t`必須。

## init.exe
初期値を生成するためのスクリプト。

`init.exe outFile isNorm`

- `outFile` : 出力するファイル
- `isNorm` : `true`のとき、思想と受容度が正規分布。`false`のとき一様分布