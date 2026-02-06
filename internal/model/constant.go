package model

var (
	UserRole = []string{
		"GM",
		"DEV",
		"USER",
		"GUEST",
	}
	UserStatus = []string{
		"ONLINE",
		"OFFLINE",
		"UNAVAILABLE",
	}
	CreatorRole = []string{
		// Color Theme:
		// Architect: สีม่วง/ทอง (ดูสูงส่ง)
		// Sentinel: สีน้ำเงิน/เงิน (ดูมั่นคง)
		// Coder: สีเขียว Neon (สีมาตรฐานของ Terminal)
		// Infiltrator: สีเทา (ดูลึกลับ)
		// Watcher: สีแดงเข้ม (ดูน่าเกรงขาม)
		"THE ARCHITECT",   // "ผู้พิทักษ์ความระเบียบและผู้คุมกฎแห่งสนามประลอง" (ผู้สร้างห้อง)
		"THE SENTINEL",    // "ผู้ดูแลความเรียบเรียงและกฎระเบียบ"  (ผู้สร้างห้อง)
		"THE CODER",       // "ผู้ขับเคลื่อนตัวเลขและอักขระสู่ชัยชนะ" (Player ทั่วไป)
		"THE INFILTRATOR", // "ผู้ลักลอบเข้าสู่ระบบเพื่อค้นหาโอกาส" (Gust)
		"THE WATCHER",     // "ผู้เฝ้ามองจากเงามืดที่ไม่ทิ้งร่องรอย" (ผู้ชม)
	}
)
