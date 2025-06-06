# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [role.proto](#role-proto)
    - [CreateRoleRequest](#role-CreateRoleRequest)
    - [DeleteRoleRequest](#role-DeleteRoleRequest)
    - [DeleteRoleResponse](#role-DeleteRoleResponse)
    - [GetRoleRequest](#role-GetRoleRequest)
    - [ListRolesRequest(未実装)](#role-ListRolesRequest)
    - [ListRolesResponse(未実装)](#role-ListRolesResponse)
    - [Role](#role-Role)
    - [RoleResponse](#role-RoleResponse)
    - [UpdateRoleRequest](#role-UpdateRoleRequest)
  
    - [RoleService](#role-RoleService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="role-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## role.proto



<a name="role-CreateRoleRequest"></a>

### CreateRoleRequest
================================
リクエスト / レスポンスメッセージ
================================


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="role-DeleteRoleRequest"></a>

### DeleteRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="role-DeleteRoleResponse"></a>

### DeleteRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="role-GetRoleRequest"></a>

### GetRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="role-ListRolesRequest"></a>

### ListRolesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [int32](#int32) |  |  |
| page_size | [int32](#int32) |  |  |






<a name="role-ListRolesResponse"></a>

### ListRolesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| roles | [Role](#role-Role) | repeated |  |
| total_count | [int32](#int32) |  |  |






<a name="role-Role"></a>

### Role



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="role-RoleResponse"></a>

### RoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [Role](#role-Role) |  |  |






<a name="role-UpdateRoleRequest"></a>

### UpdateRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |





 

 

 


<a name="role-RoleService"></a>

### RoleService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateRole | [CreateRoleRequest](#role-CreateRoleRequest) | [RoleResponse](#role-RoleResponse) | Roleの生成 |
| GetRole | [GetRoleRequest](#role-GetRoleRequest) | [RoleResponse](#role-RoleResponse) | Roleの取得 |
| UpdateRole | [UpdateRoleRequest](#role-UpdateRoleRequest) | [RoleResponse](#role-RoleResponse) | Roleの更新 |
| DeleteRole | [DeleteRoleRequest](#role-DeleteRoleRequest) | [DeleteRoleResponse](#role-DeleteRoleResponse) | Roleの削除 |
| ListRoles | [ListRolesRequest](#role-ListRolesRequest) | [ListRolesResponse](#role-ListRolesResponse) | Roleの一覧取得 |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

