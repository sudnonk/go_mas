1. エージェント同士がつながりを持つ。（`Barabashi-Albert Scale Free Network`）
2. 繋がってるエージェント同士でイデオロギーの交流が行われる
3. エージェントは、交流のたびに`Agent.HP -= 受容性 * イデオロギーの差`
4. 交流のたびに `Agent.Mix()` の中身 のアルゴリズムで自身のイデオロギーが変化
5. ターンのたびにエージェントの体力は `Agent.RecoveryRate` で回復

受容性は0-1のfloat64で、正規分布
