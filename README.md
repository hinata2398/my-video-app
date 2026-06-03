# My Video App

Go + Next.js で構築した動画投稿サイトのポートフォリオです。  
クリーンアーキテクチャを採用し、案件で想定される技術スタックを使って実装しています。

## 技術スタック

| レイヤー | 技術 |
|---|---|
| フロントエンド | Next.js 14 (App Router) / TypeScript |
| バックエンド | Go 1.22 / chi |
| アーキテクチャ | クリーンアーキテクチャ |
| インフラ | Docker / Docker Compose |

## ディレクトリ構成

```
my-video-app/
├── backend/
│   ├── domain/          # エンティティ・リポジトリインターフェース
│   ├── usecase/         # ビジネスロジック
│   ├── infrastructure/  # DB・外部API・HTTPハンドラ
│   └── main.go
├── frontend/
│   └── app/             # Next.js App Router
└── docker-compose.yml
```

## 起動方法

### 前提条件

- Docker
- Docker Compose

### 起動

```bash
docker compose up --build
```

- フロントエンド: http://localhost:3000
- バックエンド API: http://localhost:8080/api/health

## クリーンアーキテクチャについて

依存の方向を内側（domain）に向けることで、ビジネスロジックをフレームワークや DB から独立させています。

```
infrastructure → usecase → domain
```

- **domain**: エンティティとリポジトリのインターフェース定義
- **usecase**: ビジネスロジック（domain のみに依存）
- **infrastructure**: DB・HTTPなど外部との接続（usecase を呼び出す）
