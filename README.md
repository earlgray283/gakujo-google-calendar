# gakujo-google-calendar

![学情のスクショ](https://i.imgur.com/TrjZ2GA.png)

## Features

- クロスプラットフォームです(Windows / Linux / macOS)
- 定期的に学務情報システムから情報を取得して、Google カレンダー(学情カレンダー)に追加します
- アイコンをクリックすると直近の課題と締め切りまでの残り時間を表示します

## 使い方

- Windows ... [使い方](./WINDOWS.md)
- Linux ... Windows を参考にしてください
- MacOS ... Windows を参考にしてください

## For Developer

### Build

各プラットフォームごとの build から zip は次のコマンドで可能です。

```console
$ make all
```

#### Windows

**Dependencies**

- mingw-w64(https://www.mingw-w64.org/)

#### Linux

**Dependencies**

- Docker(https://docs.docker.com/engine/install/)

#### macOS

**Dependencies**

- clang

```console
$ sudo apt update && sudo apt install -y clang
```