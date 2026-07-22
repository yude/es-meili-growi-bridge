# es-meili-growi-bridge

GROWI の Elasticsearch 互換 API サーバ。Meilisearch を検索バックエンドとして使用。

```
GROWI → Elasticsearch API → es-meili-growi-bridge → Meilisearch
```

GROWI 側は Elasticsearch に接続していると認識したまま Meilisearch を利用できます。

## アーキテクチャ

```
cmd/es-meili-growi-bridge/
├── main.go              # エントリポイント
internal/
├── config/              # 設定 (環境変数)
├── server/              # HTTP サーバ (graceful shutdown)
├── handler/             # HTTP ハンドラ + ルーティング
├── elasticsearch/       # ES 互換リクエスト/レスポンス構造体
├── meilisearch/         # Meilisearch REST API クライアント
└── translator/          # ES DSL ↔ Meilisearch 変換 (最重要)
```

### リクエストフロー

```
Request (JSON)
  → elasticsearch.Request (パース)
  → translator (ES DSL → Meilisearch パラメータ変換)
  → meilisearch.Client (Meilisearch API 呼び出し)
  → meilisearch.Response
  → translator (Meilisearch → ES 互換レスポンス変換)
  → elasticsearch.Response (JSON シリアライズ)
  → Response (JSON)
```

## API 対応状況

| Elasticsearch API | 対応 | 備考 |
|---|---|---|
| `GET /` (Root/Version) | ✅ | 固定値を返却 |
| `GET /_cluster/health` | ✅ | 固定値を返却 (status: green) |
| `GET /_nodes/info` | ✅ | 固定値を返却 |
| `HEAD /{index}` | ✅ | Meilisearch の index 存在確認に変換 |
| `PUT /{index}` | ✅ | Meilisearch の Index 作成に変換 |
| `DELETE /{index}` | ✅ | Meilisearch の Index 削除に変換 |
| `GET /{index}/_stats` | ✅ | Meilisearch の Stats API に変換 |
| `HEAD /{index}/_alias/{name}` | ✅ | メモリ上の仮想エイリアスマップで管理 |
| `PUT /{index}/_alias/{name}` | ✅ | 仮想エイリアス作成 |
| `GET /{index}/_alias` | ✅ | 仮想エイリアス一覧取得 |
| `POST /_aliases` | ✅ | 仮想エイリアス更新 (add/remove) |
| `POST /{index}/_search` | ✅ | DSL → Meilisearch 検索パラメータ変換 |
| `POST /_bulk` | ✅ | NDJSON → Meilisearch 文書操作変換 |
| `POST /_reindex` | ✅ | 即座に成功を返す (wait_for_completion=false) |
| `GET /{index}/_validate/query` | ✅ | 常に valid=true を返却 |
| `GET /_cat/indices` | ✅ | Meilisearch の index 一覧を cat 形式で返却 |
| `GET /_cat/aliases` | ✅ | 仮想エイリアス一覧を cat 形式で返却 |

### ES Query DSL 変換マッピング

| ES Query DSL | Meilisearch | 方針 |
|---|---|---|
| `multi_match` (most_fields) | `q` | 全文検索クエリとして抽出 |
| `multi_match` (phrase) | 引用符でラップ | Meilisearch フレーズ検索に対応 |
| `bool.must_not[].multi_match` | `-word` | 否定語として抽出 |
| `bool.filter[].prefix` | `field STARTS WITH` | filter 式に変換 |
| `bool.filter[].term` | `field = value` | filter 式に変換 |
| `bool.filter[].terms` | `field IN [values]` | filter 式に変換 |
| `bool.filter[].bool.should` | `OR` 結合 | filter 式に変換 |
| `bool.filter[].bool.must` | `AND` 結合 | filter 式に変換 |
| `function_score` | 無視 | デフォルトスコアを使用 |
| `sort._score` | デフォルト | Meilisearch デフォルトの関連順 |
| `sort.field` | `field:asc/desc` | そのまま変換 |
| `highlight` | `attributesToHighlight` | pre/post_tags を維持 |

## ビルド方法

```bash
# ビルド
make build

# テスト
make test

# リンター
make lint
```

または:

```bash
go build -o es-meili-growi-bridge ./cmd/es-meili-growi-bridge/
```

## 実行方法

```bash
# 環境変数を設定
export MEILISEARCH_URL=http://localhost:7700
export MEILISEARCH_API_KEY=your-api-key
export LISTEN_ADDR=:9200

# 実行
./es-meili-growi-bridge
```

### 環境変数

| 変数 | デフォルト | 説明 |
|---|---|---|
| `MEILISEARCH_URL` | `http://localhost:7700` | Meilisearch の URL |
| `MEILISEARCH_API_KEY` | `""` | Meilisearch API Key |
| `LISTEN_ADDR` | `:9200` | サーバの listen アドレス |
| `LOG_LEVEL` | `info` | ログレベル |

## GROWI との接続方法

1. `es-meili-growi-bridge` を起動します。

```bash
MEILISEARCH_URL=http://localhost:7700 LISTEN_ADDR=:9200 ./es-meili-growi-bridge
```

2. GROWI の環境変数 `ELASTICSEARCH_URI` を `http://localhost:9200` に設定します。

```bash
ELASTICSEARCH_URI=http://localhost:9200/growi
```

3. GROWI を起動し、管理画面の全文検索設定が機能することを確認します。

## 注意事項

- Meilisearch には Elasticsearch のエイリアス機能が存在しないため、エイリアスはサーバのメモリ上で仮想的に管理されます。再起動時に GROWI がエイリアスを再作成するため、データの損失は発生しません。
- `field_value_factor` によるブーストスコアリングは Meilisearch に相当機能がないため無視されます。
- Reindex API は即座に成功を返します。実際のデータ移行は GROWI の `addAllPages` で MongoDB から直接行われます。
