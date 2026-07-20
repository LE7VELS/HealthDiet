// Package store 集中实现 MongoDB 数据访问、索引初始化以及 Document 与业务 Model 的转换。
// MongoDB Driver、BSON 和 Collection 类型应限制在本包内，Handler 和 Service 只能接触业务 Model 与稳定错误。
package store
