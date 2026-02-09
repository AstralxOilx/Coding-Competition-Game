package model

var (
	UserRole = []string{
		"DEV",
		"GM",
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
	Rank = []string{
		// --- ระดับพื้นฐาน (The New Era) ---
		"Script Kiddie",   // Rank 1: ผู้เริ่มต้นที่ยังใช้โค้ดคนอื่น
		"Syntax Initiate", // Rank 2: เริ่มเข้าใจโครงสร้างภาษา
		"Logic Weaver",    // Rank 3: เขียนอัลกอริทึมพื้นฐานได้

		// --- ระดับเชี่ยวชาญ (The Practitioner) ---
		"Frontend Artisan",  // Rank 4: ผู้สร้างสรรค์ส่วนติดต่อผู้ใช้
		"Backend Architect", // Rank 5: ผู้ออกแบบระบบหลังบ้าน
		"Fullstack Nomad",   // Rank 6: ผู้ชำนาญการทั้งหน้าและหลังบ้าน
		"System Voyager",    // Rank 7: ผู้เริ่มเข้าใจการจัดการ Infrastructure

		// --- ระดับปรมาจารย์ (The Elite) ---
		"Kernel Shadow",    // Rank 8: เข้าใจลึกถึงระดับแกนกลางระบบ
		"Byte Reaper",      // Rank 9: จัดการหน่วยความจำและ Binary ได้แม่นยำ
		"Security Phantom", // Rank 10: ผู้เชี่ยวชาญด้านช่องโหว่และการป้องกัน

		// --- ระดับตำนาน (The Transcendence) ---
		"Protocol Warden",  // Rank 11: ผู้ออกแบบมาตรฐานโลกดิจิทัล
		"Network Overlord", // Rank 12: ผู้ควบคุมการไหลเวียนข้อมูลมหาศาล
		"Binary Oracle",    // Rank 13: ผู้หยั่งรู้ทุกการทำงานของเลขฐานสอง
		"Code Singularity", // Rank 14: จุดสูงสุดที่โค้ดและจิตวิญญาณรวมเป็นหนึ่ง
	}

	RankTier = []int{
		3, // ระดับเริ่มต้น (Initiate)
		2, // ระดับกลาง (Expert)
		1, // ระดับสูงสุด (Master/Elite) ของ Rank นั้นๆ
	}

	FriendStatus = []string{
		"PENDING",
		"ACCEPTED",
		"BLOCKED",
	}
)

const (
	// ModeClassic: เน้นความถูกต้อง แก้โจทย์ตาม Test Case ไม่จำกัดเวลา
	ModeClassic = "Classic"

	// ModeSpeedRun: เน้นความเร็ว ใครส่งโค้ดที่ถูกต้องเร็วที่สุดชนะ
	ModeSpeedRun = "SpeedRun"

	// ModeDuel: การดวล 1 ต่อ 1 แบบ Real-time ผ่าน WebSocket
	ModeDuel = "Duel"

	// ModeCyberSiege: โหมดพิเศษแบบทีม หรือการแข่งเก็บแต้มสะสมรายสัปดาห์
	ModeCyberSiege = "CyberSiege"

	FriendStatusPending  = 0
	FriendStatusAccepted = 1
	FriendStatusBlocked  = 2

	// User Role (บทบาทผู้ใช้)
	RoleDev   = 0
	RoleGM    = 1
	RoleUser  = 2
	RoleGuest = 3

	StatusOnline      = 0
	StatusOffline     = 1
	StatusUnavailable = 2
)
