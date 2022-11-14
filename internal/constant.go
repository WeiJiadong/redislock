package internal

// 分布式锁相关常量定义
const (
	OK               = "ok"
	ErrTypeNotMatch  = "redis return type not match"           // redis返回数据类型不匹配
	ErrWLockConflict = "write lock conflict"                   // 写锁冲突
	ErrHadWLock      = "had write lock, lock read lock failed" // 已经有写锁，加读锁失败
	ErrLockNotExist  = "lock not exist"                        // 锁不存在
	ErrLockFailed    = "lock failed"                           // 加锁失败
)
