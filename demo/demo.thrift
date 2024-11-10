include "base.thrift"

enum SexType{
Male=0(tag.test="male")
    Female=1(tag.test="female")
}

struct CreateReq {
1:string name
    2:optional SexType sex
        3:   required string address

  255: optional base.Base base
}

struct CreateResp {
        1: required i32 code=0
}

const i64 T_I64 = 2024
const string T_STR   ="test"
const list<double> T_LIST = [1.0, 2.0, 3.0]
const map<string, string> T_MAP = {"k1":"v1","k2":"v2"}

struct ProductReq {
    1: required string product_id // 商品ID
    2: optional string sku_id // sku_id

    10: optional ProductApp app // app
    110: list<Property> option_list // 选项
}

struct ProductApp {
    1: string app_id
    255: optional base.Base base
}

struct Property {
    1: string key // key
    2: string name // name
}