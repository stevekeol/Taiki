/*
Package account implements ...

Account Overview

1. masterchain和basic workchain中的account_id是标准的256位

2. 其它workchain中的account_id或长或短（但至少为64位）

3. 因此，account_id只有前64位可用于“消息路由”和“动态拆分”

*/

package account
