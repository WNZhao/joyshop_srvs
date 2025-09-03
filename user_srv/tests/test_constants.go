/*
 * @Description: 用户服务测试数据常量
 * 基于SQL脚本中的真实测试数据定义
 */
package tests

// TestUserConstants 测试用户数据常量
type TestUserConstants struct {
	ID       int32
	Username string
	Password string
	Mobile   string
	Email    string
	Nickname string
	Role     int32
}

// 预定义的测试用户数据（与SQL脚本保持一致）
var (
	// 管理员用户
	AdminUser = TestUserConstants{
		ID:       1,
		Username: "admin",
		Password: "123456", // 统一密码
		Mobile:   "18888888888",
		Email:    "admin@joyshop.com",
		Nickname: "系统管理员",
		Role:     2, // 管理员角色
	}

	// 普通测试用户列表（ID: 2-12）
	RegularUsers = []TestUserConstants{
		{ID: 2, Username: "zhangsan", Password: "123456", Mobile: "18782222220", Email: "zhangsan@example.com", Nickname: "张三", Role: 1},
		{ID: 3, Username: "lisi", Password: "123456", Mobile: "18782222221", Email: "lisi@example.com", Nickname: "李四", Role: 1},
		{ID: 4, Username: "wangwu", Password: "123456", Mobile: "18782222222", Email: "wangwu@example.com", Nickname: "王五", Role: 1},
		{ID: 5, Username: "zhaoliu", Password: "123456", Mobile: "18782222223", Email: "zhaoliu@example.com", Nickname: "赵六", Role: 1},
		{ID: 6, Username: "sunqi", Password: "123456", Mobile: "18782222224", Email: "sunqi@example.com", Nickname: "孙七", Role: 1},
		{ID: 7, Username: "zhouba", Password: "123456", Mobile: "18782222225", Email: "zhouba@example.com", Nickname: "周八", Role: 1},
		{ID: 8, Username: "wujiu", Password: "123456", Mobile: "18782222226", Email: "wujiu@example.com", Nickname: "吴九", Role: 1},
		{ID: 9, Username: "zhengshi", Password: "123456", Mobile: "18782222227", Email: "zhengshi@example.com", Nickname: "郑十", Role: 1},
		{ID: 10, Username: "qianyi", Password: "123456", Mobile: "18782222228", Email: "qianyi@example.com", Nickname: "钱一", Role: 1},
		{ID: 11, Username: "chener", Password: "123456", Mobile: "18782222229", Email: "chener@example.com", Nickname: "陈二", Role: 1},
		{ID: 12, Username: "liusan", Password: "123456", Mobile: "18782222230", Email: "liusan@example.com", Nickname: "刘三", Role: 1},
	}

	// VIP用户列表（ID: 13-14）
	VIPUsers = []TestUserConstants{
		{ID: 13, Username: "vip_liu", Password: "123456", Mobile: "18888888880", Email: "vip.liu@joyshop.com", Nickname: "VIP刘总", Role: 1},
		{ID: 14, Username: "vip_wang", Password: "123456", Mobile: "18888888881", Email: "vip.wang@joyshop.com", Nickname: "VIP王总", Role: 1},
	}

	// 测试用新用户（用于注册测试）
	NewTestUser = TestUserConstants{
		ID:       15,
		Username: "testuser",
		Password: "123456",
		Mobile:   "18900000000",
		Email:    "testuser@example.com",
		Nickname: "测试用户",
		Role:     1,
	}
)

// GetTestUserByID 根据ID获取测试用户
func GetTestUserByID(id int32) *TestUserConstants {
	if id == AdminUser.ID {
		return &AdminUser
	}
	
	for i := range RegularUsers {
		if RegularUsers[i].ID == id {
			return &RegularUsers[i]
		}
	}
	
	for i := range VIPUsers {
		if VIPUsers[i].ID == id {
			return &VIPUsers[i]
		}
	}
	
	if id == NewTestUser.ID {
		return &NewTestUser
	}
	
	return nil
}

// GetTestUserByUsername 根据用户名获取测试用户
func GetTestUserByUsername(username string) *TestUserConstants {
	if username == AdminUser.Username {
		return &AdminUser
	}
	
	for i := range RegularUsers {
		if RegularUsers[i].Username == username {
			return &RegularUsers[i]
		}
	}
	
	for i := range VIPUsers {
		if VIPUsers[i].Username == username {
			return &VIPUsers[i]
		}
	}
	
	if username == NewTestUser.Username {
		return &NewTestUser
	}
	
	return nil
}

// GetRandomTestUser 获取随机测试用户（从普通用户中选择）
func GetRandomTestUser() TestUserConstants {
	if len(RegularUsers) == 0 {
		return NewTestUser
	}
	// 简单取第一个用户，避免使用随机数
	return RegularUsers[0]
}

// GetVIPTestUser 获取VIP测试用户
func GetVIPTestUser() TestUserConstants {
	if len(VIPUsers) == 0 {
		return NewTestUser
	}
	return VIPUsers[0]
}

// AllTestUsers 获取所有测试用户
func AllTestUsers() []TestUserConstants {
	all := make([]TestUserConstants, 0)
	all = append(all, AdminUser)
	all = append(all, RegularUsers...)
	all = append(all, VIPUsers...)
	return all
}

// TestScenarios 测试场景常量
var TestScenarios = struct {
	AdminLogin    string
	UserLogin     string
	VIPUserLogin  string
	UserRegistration string
	UserUpdate    string
	UserDelete    string
}{
	AdminLogin:       "管理员登录测试",
	UserLogin:        "普通用户登录测试", 
	VIPUserLogin:     "VIP用户登录测试",
	UserRegistration: "用户注册测试",
	UserUpdate:       "用户信息更新测试",
	UserDelete:       "用户删除测试",
}