# GitHub Persona

![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=kou7306)

GitHub Persona は、あなたの GitHub 活動を分析して、RPG 風のキャラクター画像を生成するサービスです。

## 🎮 判定システム

### 📊 レベル計算

レベルは以下の統計データから計算されます：

```
レベル = (スター数 + コントリビューション数 + Issue数 + PR数 + コミット数) ÷ 15
```

**最大レベル**: 100

### 🏆 ランク判定

| レベル範囲 | ランク | 称号                                   |
| ---------- | ------ | -------------------------------------- |
| 0-2        | C-     | 少年                                   |
| 3-9        | C      | 少年                                   |
| 10-14      | C+     | 冒険者見習い                           |
| 15-24      | B-     | 魔術師の見習い / 不良 / 駆け出し冒険者 |
| 25-34      | B      | 初級 [職業名]                          |
| 35-45      | B+     | 中級 [職業名]                          |
| 46-59      | A-     | 上級 [職業名]                          |
| 60-79      | A      | 特級 [職業名]                          |
| 80-99      | A+     | 特殊称号                               |
| 100        | S      | 神                                     |

### ⚔️ 職業判定

#### 🧙‍♂️ 魔法ルート

- **TypeScript**: 攻撃魔術師
- **R**: ネクロマンサー
- **Dart**: 防御魔術師
- **Go**: 召喚士
- **Scala**: 精霊魔法
- **Rust**: 回復術師

#### 🦹‍♂️ アウトルート

- **Assembly**: 賞金稼ぎ
- **C**: 犯罪者
- **C++**: 犯罪者
- **Objective-C**: 盗賊
- **Matlab**: 盗賊

#### ⚔️ 戦士ルート

- **C#**: 武闘家
- **Swift**: 弓使い
- **Kotlin**: 弓使い
- **Ruby**: 槍使い
- **PHP**: 槍使い
- **HTML**: 剣士
- **CSS**: 剣士
- **JavaScript**: 剣士
- **Java**: 騎士
- **Python**: 士官

### 🔥 ハイブリッド職業 (A-以上)

上位 2 つの言語が異なるルートで、それぞれ 15%以上の使用率の場合：

| 組み合わせ        | 職業       |
| ----------------- | ---------- |
| アウトロー + 戦士 | バーカーサ |
| 戦士 + アウトロー | 闇騎士     |
| 魔法 + アウトロー | 黒魔術師   |
| アウトロー + 魔法 | ライダー   |
| 戦士 + 魔法       | 魔法戦士   |
| 魔法 + 戦士       | 魔法騎士   |

### 👑 A+ランク特殊称号

| 基本職業                                       | 特殊称号     |
| ---------------------------------------------- | ------------ |
| 賞金稼ぎ/犯罪者/盗賊                           | 裏社会のボス |
| 攻撃魔術師/防御魔術師/召喚士/精霊魔法/回復術師 | 魔法帝       |
| 武闘家/弓使い/槍使い/剣士                      | 勇者         |
| 騎士/士官                                      | 騎士団長     |
| 魔法戦士/魔法騎士                              | 賢者         |
| バーカーサ/闇騎士                              | サイコパス   |
| その他                                         | 魔王         |

## 📈 統計データ

以下の GitHub GraphQL API から取得されるデータを使用：

- **Total Stars**: スターしたリポジトリ数
- **Total Commits**: 総コミット数
- **Total PRs**: 総プルリクエスト数
- **Total Issues**: 総 Issue 数
- **Contributed To**: コントリビューションしたリポジトリ数

## 🎨 画像生成

### 言語判定ルール

1. **HTML/CSS/JavaScript/TypeScript**は基本的に除外
2. 上位 2 つの言語を選択
3. 15%以上の使用率の言語を優先
4. 言語が見つからない場合は上位 2 つを使用

### キャラクター画像

- 職業に応じてキャラクター画像が選択
- レベルに応じて称号が表示
- 統計データの棒グラフ表示
- コミット履歴のカレンダー表示

## 🚀 使用方法

### GitHub README に表示

```markdown
![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=your-username)
```

### 例

![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=kou7306)

### 🔄 画像の自動更新

#### 方法 1: 強制更新パラメータ

```markdown
![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=your-username&update=1)
```

#### 方法 2: キャッシュ無効化

```markdown
![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=your-username&nocache=1)
```

#### 方法 3: タイムスタンプ付き（推奨）

```markdown
![GitHub Persona](https://github-persona-backend-vrkndjdhdq-uc.a.run.app/github?username=your-username&t=20240803)
```

**注意**: GitHub README の画像はキャッシュされるため、更新を確認したい場合は上記のパラメータを使用してください。

## 🌐 デモ

**フロントエンド**: https://github-persona-r6bzj8e4k-kous-projects-7736bb57.vercel.app

## 🔧 技術仕様

- **バックエンド**: Go (Cloud Run)
- **フロントエンド**: Next.js (Vercel)
- **画像生成**: GG (Go Graphics)
- **API**: GitHub GraphQL API

## 📝 開発

開発環境のセットアップについては [CONTRIBUTING.md](./CONTRIBUTING.md) を参照してください。

## 🎯 判定例

### 例 1: Python + TypeScript ユーザー

- **レベル**: 59 (A-)
- **職業**: 上級 士官
- **理由**: Python が士官ルート、TypeScript が魔法ルート、両方 15%以上

### 例 2: JavaScript + HTML ユーザー

- **レベル**: 25 (B)
- **職業**: 初級 剣士
- **理由**: JavaScript と HTML が戦士ルート、HTML は除外されるため JavaScript のみ

### 例 3: Go + Rust ユーザー

- **レベル**: 85 (A+)
- **職業**: 魔法帝
- **理由**: Go と Rust が魔法ルート、A+ランクのため魔法帝に昇格

---

**GitHub Persona** - あなたの GitHub 活動を RPG キャラクターに変換！
