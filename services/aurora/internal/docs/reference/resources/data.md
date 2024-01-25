---
title: Data
replacement: https://developers.diamcircle.org/api/resources/accounts/data/
---

Each account in Diamcircle network can contain multiple key/value pairs associated with it. Aurora can be used to retrieve value of each data key.

When aurora returns information about a single account data key it uses the following format:

## Attributes

| Attribute | Type | | 
| --- | --- | --- |
| value | base64-encoded string | The base64-encoded value for the key |

## Example

```json
{
  "value": "MTAw"
}
```
