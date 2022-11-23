# 技術分析

- [技術分析](#技術分析)
  - [後端](#後端)
    - [API](#api)
      - [SignIn](#signin)
      - [SignOut](#signout)
      - [CreateUser](#createuser)
      - [UpdateUserConfiguration](#updateuserconfiguration)
      - [CreateAccount](#createaccount)
      - [UpdateAccount](#updateaccount)
      - [DeleteAccount](#deleteaccount)
      - [CreateCategory](#createcategory)
      - [UpdateCategory](#updatecategory)
      - [DeleteCategory](#deletecategory)
      - [CreateShop](#createshop)
      - [UpdateShop](#updateshop)
      - [DeleteShop](#deleteshop)
      - [CreateFee](#createfee)
      - [UpdateFee](#updatefee)
      - [DeleteFee](#deletefee)
      - [CreateRepeatingItem](#createrepeatingitem)
      - [UpdateRepeatingItem](#updaterepeatingitem)
      - [DeleteRepeatingItem](#deleterepeatingitem)
      - [CreateItem](#createitem)
      - [UpdateItem](#updateitem)
      - [DeleteItem](#deleteitem)
      - [ListDailyItems](#listdailyitems)
      - [ExportItems](#exportitems)
      - [ImportItems](#importitems)
    - [資料儲存](#資料儲存)
      - [System Icon Object](#system-icon-object)
      - [Account Object](#account-object)
      - [Category Object](#category-object)
      - [Shop Object](#shop-object)
      - [Shop Item Object](#shop-item-object)
      - [Fee Object](#fee-object)
      - [Repeating Item Object](#repeating-item-object)
      - [Item Object](#item-object)
      - [Daily Item Object](#daily-item-object)
      - [User Advance Item Object](#user-advance-item-object)
    - [技術問題](#技術問題)

## 後端

### API

#### SignIn

path: /sign-in

method: POST

request: TODO

#### SignOut

path: /sign-out

method: POST

request: TODO

#### CreateUser

path: /sign-up

method: POST

request:

```json
{
   "id": "",
   "password": ""
}
```

#### UpdateUserConfiguration

path: /${user-id}/config

method: PUT

request:

```json
{
   "compare_items_in_different_shop": false,
   "compare_items_in_same_shop": false,
}
```

#### CreateAccount

path: `/${user-id}/accounts`

method: POST

request:

```json
{
   "name": "account-name",
   "icon": "icon-id",
   "initial_balance": 0,
}
```

#### UpdateAccount

path: `/${user-id}/accounts/${account-id}`

method: PUT

request:

```json
{
   "name": "account-name",
   "icon": "icon-id",
   "initial_balance": 0,
}
```

#### DeleteAccount

path: `/${user-id}/accounts/${account-id}`

method: DELETE

#### CreateCategory

path: `/${user-id}/categories`

method: POST

request:

```json
{
   "type": 0, // enum: expense, income
   "name": "expense-or-income-name",
   "icon": "icon-id",
}
```

#### UpdateCategory

path: `/${user-id}/categories/${category-id}`

method: PUT

request:

```json
{
   "name": "expense-or-income-name",
   "icon": "icon-id",
}
```

#### DeleteCategory

path: `/${user-id}/categories/${category-id}`

method: DELETE

#### CreateShop

path: `/${user-id}/shops`

method: POST

request:

```json
{
   "name": "shop-name",
   "address": "shop-location",
}
```

#### UpdateShop

path: `/${user-id}/shops/${shop-id}`

method: PUT

request:

```json
{
   "name": "shop-name",
   "address": "shop-location",
}
```

#### DeleteShop

path: `/${user-id}/shops/${shop-id}`

method: DELETE

#### CreateFee

path: `/${user-id}/fee`

method: POST

request:

```json
{
   "name": "fee-name",
   "value": { // one of those
      "rate": 0.05,
      "fixed": 30
   }
}
```

#### UpdateFee

path: `/${user-id}/fee/${fee-id}`

method: PUT

request:

```json
{
   "name": "fee-name",
   "value": { // one of those
      "rate": 0.05,
      "fixed": 30
   }
}
```

#### DeleteFee

path: `/${user-id}/fee/${fee-id}`

method: DELETE

#### CreateRepeatingItem

path: `/${user-id}/repeating-items`

method: POST

request:

```json
{
   "item": {
      "name": "chocolate",
      "category": [1], // users.categories.[expense | income]
      "shop": 1, // users.shops, nullable
      "quantity": {
         "value": 10
      },
      "fee": 1, //  nullable
      "price": 50,
      "memo" : ""
   },
   "valid": {
      "start": "2022-01-01",
      "end": null
   },
   "frequency": { // one of those
      "duration": 0,
      "every_work_day": false
   }
}
```

#### UpdateRepeatingItem

path: `/${user-id}/repeating-items/${repeating-item-id}`

method:

request: PUT

```json
{
   "item": {
      "name": "chocolate",
      "category": [1], // users.categories.[expense | income]
      "shop": 1, // users.shops, nullable
      "quantity": { // nullable
         "value": 10
      },
      "fee": 1, //  nullable
      "price": 50,
      "memo" : ""
   },
   "valid": {
      "start": "2022-01-01",
      "end": null
   },
   "frequency": { // one of those
      "duration": 0,
      "every_work_day": false
   }
}
```

#### DeleteRepeatingItem

path: `/${user-id}/repeating-items/${repeating-item-id}`

method: DELETE

#### CreateItem

path: `/${user-id}/items`

method: POST

request:

```json
{
   "date": "2022-11-15",
   "name": "chocolate",
   "category": [1], // users.categories.[expense | income]
   "shop": 1, // users.shops, nullable
   "quantity": { // nullable
      "value": 10
   },
   "fee": 1, // nullable
   "price": 50,
   "memo" : ""
}
```

#### UpdateItem

path: `/${user-id}/items/${item-id}`

method: PUT

request:

```json
{
   "date": "2022-11-15",
   "name": "chocolate",
   "category": [1], // users.categories.[expense | income]
   "shop": 1, // users.shops, nullable
   "quantity": { // nullable
      "value": 10
   },
   "fee": 1, // nullable
   "price": 50,
   "memo" : ""
}
```

#### DeleteItem

path: `/${user-id}/items/${item-id}`

method: DELETE

#### ListDailyItems

path: `/${user-id}/items?date=2022-01-01`

method: GET

query parameters:

| param | default |
| --- | --- |
| date | today |

#### ExportItems

path: `/${user-id}/items/export?format=csv`

method: GET

query parameters:

| param | default | description |
| --- | --- | --- |
| format | csv | support `csv`, `json`, `xml` |

#### ImportItems

path: `/${user-id}/items/import`

method: POST

description: the upload file depends on header, it supports `csv`, `json`, `xml`.

### 資料儲存

boltDB, 分成以下 bucket:

1. system
   1. icons
      1. account
         1. [System Icon](#system-icon-object)
      2. category
         1. expense
            1. [System Icon](#system-icon-object)
         2. income
            1. [System Icon](#system-icon-object)
2. users
   1. accounts
      1. [Account](#account-object)
   2. categories
      1. expense
         1. [Category](#category-object)
      2. income
         1. [Category](#category-object)
   3. shops
      1. ${shop-id}
         1. [Shop](#shop-object)
         2. items
            1. [Shop Item](#shop-item-object)
   4. fee
      1. [Fee](#fee-object)
   5. repeating items
      1. [Repeating Item](#repeating-item-object)
   6. items
       1. [Item](#item-object)
   7. calendar
      1. ${year}
         1. ${month}
            1. [Daily Item](#daily-item-object)
   8. advance
      1. items
         1. [User Advance Item](#user-advance-item-object)

#### System Icon Object

key: ${sequence id}

value: binary-file

#### Account Object

key: ${sequence id}

value:

```json
{
   "name": "account-name",
   "icon": "icon-id", // system.icons.account
   "initial_balance": 0,
   "balance": 0,
}
```

#### Category Object

key: ${sequence id}

value:

```json
{
   "name": "expense-or-income-name",
   "icon": "icon-id", // system.icons.category.[expense | income]
   "items": { // map[item-id]empty
      "$item-id": {}
   }
}
```

#### Shop Object

key: information

value:

```json
{
   "name": "shop-name",
   "address": "shop-location",
}
```

#### Shop Item Object

key: ${item-name}

value:

```json
[1] // elements are `item-id`
```

#### Fee Object

key: ${sequence id}

value:

```json
{
   "name": "fee-name",
   "value": { // one of those
      "rate": 0.05,
      "fixed": 30
   }
}
```

#### Repeating Item Object

key: ${sequence id}

value:

```json
{
   "item": {
      "name": "chocolate",
      "category": [1], // users.categories.[expense | income]
      "shop": 1, // users.shops
      "quantity": {
         "value": 10
      },
      "fee": 1,
      "price": 50,
      "memo" : ""
   },
   "valid": {
      "start": "2022-01-01",
      "end": null
   },
   "frequency": { // one of those
      "duration": 0,
      "every_work_day": false
   }
}
```

#### Item Object

key: ${sequence id}

value:

```json
{
   "date": "2022-11-15",
   "name": "chocolate",
   "category": [1], // users.categories.[expense | income]
   "shop": 1, // users.shops
   "quantity": {
      "value": 10
   },
   "fee": 1,
   "price": 50,
   "memo" : ""
}
```

#### Daily Item Object

key: ${day}

value:

```json
[1] // elements are `item-id`
```

#### User Advance Item Object

key: ${item-name}

value:

```json
[1] // elements are `item-id`
```

### 技術問題

- [x] 如何計算帳戶餘額?
   1. 於 `users.account` 紀錄餘額
- [ ] 如何處理 repeating items?
   1. 檢查每筆 repeating item, 當日期等於今天, 將 repeating items 放到 `users.items.${year}.${month}.${day}`
- [x] 如何依據 category 篩選 item?
   1. 考慮把 item-id 放到 category, 因此需要重構 bucket 設計,
      將 Item Object 定義於另外一個 bucket, 其他 bucket 引用他的 id.
- [ ] 如何依據關鍵詞篩選 item (從 `item.name` 和 `item.memo` 篩選)?
   1. 粗估一位使用者一年的 item 總數約1,000, 因此目前想到逐筆 item 篩選,
      每次篩選出10筆, 類似於分頁效果, 如果 UI 需要再載入下一批資料.
      但是又需要考慮到可能發生搜尋了超過百筆資料仍不足10筆的對應方式.
- [x] 如何比較不同 shop 價格差異?
   1. 使用者設定要開啟這項功能, 避免過於占用空間(假如未來要設計於 client 儲存資料).
      將資料儲存於 `users.advance.items` bucket.
      以 `item.name` 為 key, `{"shop-id": 0, "price": 0, "date": "2000-01-01"}` 為 value.
      還要再設計多久以前的資料捨棄.
- [x] 如何計算同一 shop 同 `item.name` 價格差異(或漲幅)?
   1. 重構 `users.shops` 設計, 長得像:

      ```text
      1. users
         1. shops
            1. ${shop-id}
               1. Shop Object
               2. items
                  1. Shop Item Object

      Shop Object 的 key 固定為 `_key`

      Shop Item Object:
         key 為 item-name.
         value:
         {
            "price": 0
         }
      ```
