// Package handler 负责解析 HTTP 输入、调用 Service 或 Store，并按 API 合同输出响应。
// 本包只定义 HTTP DTO 和状态码映射，不直接操作 MongoDB，也不在响应中暴露密码哈希、数据库错误等内部信息。
package handler
