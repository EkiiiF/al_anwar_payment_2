package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type SuperUserServiceImpl struct {
	DB           *gorm.DB
	Validate     *validator.Validate
	Repository   repository.SuperUserRepository
	SettingRepo  repository.SettingRepository
	AuditLog     LogService
	Notification NotificationService
}

func NewSuperUserService(
	db *gorm.DB,
	validate *validator.Validate,
	repo repository.SuperUserRepository,
	settingRepo repository.SettingRepository,
	logService LogService,
	notifService NotificationService,
) SuperUserService {
	return &SuperUserServiceImpl{
		DB:           db,
		Validate:     validate,
		Repository:   repo,
		SettingRepo:  settingRepo,
		AuditLog:     logService,
		Notification: notifService,
	}
}

// guardianRelation mengembalikan relasi wali, default "Orang Tua" jika kosong.
func guardianRelation(rel string) string {
	if rel == "" {
		return "Orang Tua"
	}
	return rel
}

func (s *SuperUserServiceImpl) GetDashboardStats() (response.SuperUserDashboardStats, error) {
	totalStudents, err := s.Repository.CountActiveStudents(s.DB)
	if err != nil {
		return response.SuperUserDashboardStats{}, fmt.Errorf("gagal hitung santri: %w", err)
	}

	unpaidInvoices, err := s.Repository.CountUnpaidInvoices(s.DB)
	if err != nil {
		return response.SuperUserDashboardStats{}, fmt.Errorf("gagal hitung tagihan: %w", err)
	}

	paidInvoices, err := s.Repository.CountPaidInvoices(s.DB)
	if err != nil {
		return response.SuperUserDashboardStats{}, fmt.Errorf("gagal hitung tagihan lunas: %w", err)
	}

	now := time.Now()
	totalIncome, err := s.Repository.SumIncomeThisMonth(s.DB, int(now.Month()), now.Year())
	if err != nil {
		return response.SuperUserDashboardStats{}, fmt.Errorf("gagal hitung pendapatan: %w", err)
	}

	monthlyIncomeData, err := s.Repository.GetMonthlyIncomeForYear(s.DB, now.Year())
	if err != nil {
		return response.SuperUserDashboardStats{}, fmt.Errorf("gagal hitung pendapatan bulanan: %w", err)
	}

	var monthlyStats []response.MonthlyPaymentStat
	// Initialize all 12 months
	monthMap := make(map[int]float64)
	for _, data := range monthlyIncomeData {
		monthMap[data.Month] = data.Total
	}
	for i := 1; i <= 12; i++ {
		monthlyStats = append(monthlyStats, response.MonthlyPaymentStat{
			Month: i,
			Year:  now.Year(),
			Total: monthMap[i],
		})
	}

	// ─── Hijri & Semester Info ──────────────────────────────
	hijriNow := utils.GetCurrentHijriDate()
	semesterInfo := utils.GetSemesterInfo(hijriNow.Month, hijriNow.Year)

	currentHijri := response.HijriMonthInfo{
		HijriMonth:        hijriNow.Month,
		HijriMonthName:    utils.GetHijriMonthName(hijriNow.Month),
		HijriYear:         hijriNow.Year,
		Semester:          semesterInfo.Number,
		SemesterName:      semesterInfo.Name,
		AcademicYearLabel: utils.GetAcademicYearLabel(hijriNow.Month, hijriNow.Year),
		IsExamMonth:       semesterInfo.IsExamMonth,
		IsRegistration:    semesterInfo.IsRegistration,
	}

	// Ambil statistik semester dari database
	semesterData, err := s.Repository.GetSemesterPaymentStats(s.DB, hijriNow.Year)
	var semesterStats []response.SemesterPaymentStat
	if err == nil {
		for _, sd := range semesterData {
			semName := "Semester 1"
			if sd.Semester == 2 {
				semName = "Semester 2"
			}
			semesterStats = append(semesterStats, response.SemesterPaymentStat{
				SemesterName:      semName,
				AcademicYearLabel: sd.AcademicYear,
				Total:             sd.Total,
				InvoiceCount:      sd.InvoiceCount,
				PaidCount:         sd.PaidCount,
				UnpaidCount:       sd.UnpaidCount,
			})
		}
	}
	if semesterStats == nil {
		semesterStats = []response.SemesterPaymentStat{}
	}

	return response.SuperUserDashboardStats{
		TotalStudents:   totalStudents,
		UnpaidInvoices:  unpaidInvoices,
		PaidInvoices:    paidInvoices,
		TotalIncomeMo:   totalIncome,
		MonthlyPayments: monthlyStats,
		SemesterStats:   semesterStats,
		CurrentHijri:    currentHijri,
	}, nil
}

func (s *SuperUserServiceImpl) GetStudents(search string) ([]domain.Student, error) {
	return s.Repository.FindAllStudents(s.DB, search)
}

func (s *SuperUserServiceImpl) GetStudentsPaginated(search string, page int, limit int) ([]domain.Student, int64, error) {
	return s.Repository.FindAllStudentsPaginated(s.DB, search, page, limit)
}

func (s *SuperUserServiceImpl) GetInvoices(status string, month int, year int) ([]domain.Invoice, error) {
	return s.Repository.FindAllInvoices(s.DB, status, month, year)
}

func (s *SuperUserServiceImpl) GetInvoicesPaginated(status string, month int, year int, page int, limit int) ([]domain.Invoice, int64, error) {
	return s.Repository.FindAllInvoicesPaginated(s.DB, status, month, year, page, limit)
}

func (s *SuperUserServiceImpl) GetStudentsWithInvoicesPaginated(status string, month int, year int, page int, limit int) ([]domain.Student, int64, error) {
	return s.Repository.FindStudentsWithInvoicesPaginated(s.DB, status, month, year, page, limit)
}

// parseFullName memecah full_name menjadi FirstName, MiddleName, dan LastName.
// - 1 kata: FirstName saja
// - 2 kata: FirstName + LastName
// - 3+ kata: FirstName + MiddleName(s) + LastName
func parseFullName(fullName string) (first, middle, last string) {
	parts := strings.Fields(fullName)
	switch len(parts) {
	case 0:
		return "", "", ""
	case 1:
		return parts[0], "", ""
	case 2:
		return parts[0], "", parts[1]
	default:
		return parts[0], strings.Join(parts[1:len(parts)-1], " "), parts[len(parts)-1]
	}
}

func (s *SuperUserServiceImpl) CreateStudent(ctx context.Context, req request.CreateStudentRequest, operatorID, ip, ua string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Parse nama dari full_name → first, middle, last
		first, middle, last := parseFullName(req.FullName)
		studentName := domain.StudentName{
			FirstName:  first,
			MiddleName: middle,
			LastName:   last,
		}

		// Default muhadhoroh level 1 semester 1 jika tidak dikirim
		muhadLevel := req.MuhadhorohLevel
		if muhadLevel == 0 {
			muhadLevel = 1
		}
		currSem := req.CurrentSemester
		if currSem == 0 {
			currSem = 1
		}

		// Pastikan gender valid, default "L" jika kosong
		gender := req.Gender
		if gender == "" {
			gender = "L"
		}

		student := domain.Student{
			ID:              utils.GenerateID(),
			Name:            studentName,
			Nik:             req.Nik,
			BirthDate:       req.BirthDate,
			StudentNumber:   req.NIS,
			Gender:          gender,
			MuhadhorohLevel: muhadLevel,
			CurrentSemester: currSem,
			StatusID:        req.StatusTypeID,
			IsActive:        true,
			Addresses: []domain.Address{
				{
					ID:          utils.GenerateID(),
					AddressLine: req.AddressLine,
					Village:     req.Village,
					District:    req.District,
					City:        req.City,
					Province:    req.Province,
					Country:     req.Country,
					PostalCode:  req.PostalCode,
					IsPrimary:   true,
				},
			},
		}

		if err := s.Repository.CreateStudent(tx, &student); err != nil {
			return fmt.Errorf("gagal membuat data santri: %w", err)
		}

		// ─── Opsi 1: Buat akun guardian baru ─────────────────
		if req.CreateNewGuardian {
			role, err := s.Repository.FindRoleByName(tx, "guardian")
			if err != nil {
				return fmt.Errorf("role guardian tidak ditemukan: %w", err)
			}

			// Username otomatis sama dengan NIS
			username := req.NIS

			// Password default "password123" jika kosong
			rawPassword := req.GuardianPassword
			if rawPassword == "" {
				rawPassword = "password123"
			}

			hashedPassword, err := utils.HashPassword(rawPassword)
			if err != nil {
				return fmt.Errorf("gagal hash password: %w", err)
			}

			newUser := domain.User{
				ID:       utils.GenerateID(),
				Username: username,
				Password: hashedPassword,
				RoleID:   role.ID,
				IsActive: true,
			}
			if err := s.Repository.CreateUser(tx, &newUser); err != nil {
				return fmt.Errorf("gagal membuat akun guardian: %w", err)
			}

			gFirst, gMiddle, gLast := parseFullName(req.GuardianName)
			guardianName := domain.GuardianName{
				FirstName:  gFirst,
				MiddleName: gMiddle,
				LastName:   gLast,
			}

			guardian := domain.Guardian{
				ID:        utils.GenerateID(),
				StudentID: student.ID,
				UserID:    newUser.ID,
				Name:      guardianName,
				Phone:     req.GuardianPhone,
				Email:     req.GuardianEmail,
				Relation:  guardianRelation(req.GuardianRelation),
			}
			if err := s.Repository.CreateGuardian(tx, &guardian); err != nil {
				return fmt.Errorf("gagal membuat data guardian: %w", err)
			}

		} else if req.GuardianID != "" {
			// ─── Opsi 2: Gunakan wali santri yang sudah ada ──
			existingGuardian, err := s.Repository.FindGuardianByID(tx, req.GuardianID)
			if err != nil {
				return fmt.Errorf("guardian dengan ID %s tidak ditemukan: %w", req.GuardianID, err)
			}

			guardianName := existingGuardian.Name
			if req.GuardianName != "" {
				gFirst, gMiddle, gLast := parseFullName(req.GuardianName)
				guardianName = domain.GuardianName{
					FirstName:  gFirst,
					MiddleName: gMiddle,
					LastName:   gLast,
				}
			}

			phone := existingGuardian.Phone
			if req.GuardianPhone != "" {
				phone = req.GuardianPhone
			}
			email := existingGuardian.Email
			if req.GuardianEmail != "" {
				email = req.GuardianEmail
			}
			relation := existingGuardian.Relation
			if req.GuardianRelation != "" {
				relation = guardianRelation(req.GuardianRelation)
			}

			newGuardian := domain.Guardian{
				ID:        utils.GenerateID(),
				StudentID: student.ID,
				UserID:    existingGuardian.UserID, // <-- Kunci: User yang SAMA
				Name:      guardianName,
				Phone:     phone,
				Email:     email,
				Relation:  relation,
			}
			if err := s.Repository.CreateGuardian(tx, &newGuardian); err != nil {
				return fmt.Errorf("gagal menghubungkan guardian ke santri baru: %w", err)
			}

			// ─── Auto-copy alamat dari saudara jika alamat baru kosong ──
			// Jika Bendahara tidak mengisi alamat untuk anak ke-2, salin dari
			// anak pertama (saudara) yang sudah memiliki alamat.
			if req.AddressLine == "" && len(existingGuardian.Student.Addresses) > 0 {
				for _, srcAddr := range existingGuardian.Student.Addresses {
					if srcAddr.IsPrimary {
						copiedAddr := domain.Address{
							ID:          utils.GenerateID(),
							StudentID:   student.ID,
							AddressLine: srcAddr.AddressLine,
							Village:     srcAddr.Village,
							District:    srcAddr.District,
							City:        srcAddr.City,
							Province:    srcAddr.Province,
							Country:     srcAddr.Country,
							PostalCode:  srcAddr.PostalCode,
							IsPrimary:   true,
						}
						// Hapus alamat kosong yang dibuat sebelumnya
						tx.Where("student_id = ? AND address_line = ''", student.ID).Delete(&domain.Address{})
						if err := tx.Create(&copiedAddr).Error; err != nil {
							return fmt.Errorf("gagal menyalin alamat dari saudara: %w", err)
						}
						break
					}
				}
			}
		}

		s.AuditLog.Log(ctx, operatorID, "CREATE_STUDENT", "Student", student.ID, req, ip, ua)
		return nil
	})
}

func (s *SuperUserServiceImpl) UpdateStudent(ctx context.Context, id string, req request.UpdateStudentRequest, operatorID, ip, ua string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		student, err := s.Repository.FindStudentByID(tx, id)
		if err != nil {
			return fmt.Errorf("santri tidak ditemukan: %w", err)
		}

		// ── Parse nama → first, middle, last ─────────────
		first, middle, last := parseFullName(req.FullName)

		// Pastikan gender valid
		gender := req.Gender
		if gender == "" {
			gender = student.Gender // keep existing if empty
		}

		// ── Update field student secara EKSPLISIT (tanpa associations) ──
		// PENTING: Jangan gunakan db.Save(student) dengan nested associations!
		// GORM Save() dengan associations menyebabkan data address berpindah/hilang
		// karena GORM melakukan upsert berdasarkan primary key association.
		studentUpdate := map[string]interface{}{
			"first_name":       first,
			"middle_name":      middle,
			"last_name":        last,
			"student_number":   req.NIS,
			"status_id":        req.StatusTypeID,
			"nik":              req.Nik,
			"birth_date":       req.BirthDate,
			"gender":           gender,
			"muhadhoroh_level": req.MuhadhorohLevel,
			"current_semester": req.CurrentSemester,
		}
		if err := tx.Model(&domain.Student{}).Where("id = ?", id).Updates(studentUpdate).Error; err != nil {
			return fmt.Errorf("gagal update data santri: %w", err)
		}

		// ── Update Address secara TERPISAH ─────────────────
		// Cari alamat primary yang sudah ada untuk student ini.
		var existingAddr domain.Address
		addrFound := tx.Where("student_id = ? AND is_primary = ?", id, true).First(&existingAddr).Error == nil

		if addrFound {
			// Update alamat yang sudah ada berdasarkan ID-nya sendiri
			addrUpdate := map[string]interface{}{
				"address_line": req.AddressLine,
				"village":      req.Village,
				"district":     req.District,
				"city":         req.City,
				"province":     req.Province,
				"country":      req.Country,
				"postal_code":  req.PostalCode,
			}
			if err := tx.Model(&domain.Address{}).Where("id = ?", existingAddr.ID).Updates(addrUpdate).Error; err != nil {
				return fmt.Errorf("gagal update alamat: %w", err)
			}
		} else {
			// Buat alamat baru jika belum ada
			newAddr := domain.Address{
				ID:          utils.GenerateID(),
				StudentID:   id,
				AddressLine: req.AddressLine,
				Village:     req.Village,
				District:    req.District,
				City:        req.City,
				Province:    req.Province,
				Country:     req.Country,
				PostalCode:  req.PostalCode,
				IsPrimary:   true,
			}
			if err := tx.Create(&newAddr).Error; err != nil {
				return fmt.Errorf("gagal membuat alamat: %w", err)
			}
		}

		// Update guardian data jika ada
		if len(student.Guardians) > 0 {
			guardian := &student.Guardians[0]

			if req.GuardianName != "" {
				gFirst, gMiddle, gLast := parseFullName(req.GuardianName)
				guardian.Name = domain.GuardianName{
					FirstName:  gFirst,
					MiddleName: gMiddle,
					LastName:   gLast,
				}
			}
			if req.GuardianPhone != "" {
				guardian.Phone = req.GuardianPhone
			}
			if req.GuardianRelation != "" {
				guardian.Relation = req.GuardianRelation
			}
			if req.GuardianEmail != "" {
				guardian.Email = req.GuardianEmail
			}
			if err := s.Repository.SaveGuardian(tx, guardian); err != nil {
				return fmt.Errorf("gagal update guardian: %w", err)
			}

			// Update user account jika ada perubahan
			// Selalu sinkronkan username dengan NIS terbaru
			user := guardian.User
			user.Username = req.NIS

			// Jika password diisi, update password. Jika tidak, biarkan yang lama.
			if req.GuardianPassword != "" {
				hashed, err := utils.HashPassword(req.GuardianPassword)
				if err != nil {
					return fmt.Errorf("gagal hash password: %w", err)
				}
				user.Password = hashed
			}

			if err := s.Repository.SaveUser(tx, &user); err != nil {
				return fmt.Errorf("gagal update akun guardian: %w", err)
			}
		}

		s.AuditLog.Log(ctx, operatorID, "UPDATE_STUDENT", "Student", id, req, ip, ua)
		return nil
	})
}

func (s *SuperUserServiceImpl) DeleteStudent(ctx context.Context, id string, operatorID, ip, ua string) error {
	if err := s.Repository.DeleteStudent(s.DB, id); err != nil {
		return fmt.Errorf("gagal menghapus santri: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "DELETE_STUDENT", "Student", id, nil, ip, ua)
	return nil
}

func (s *SuperUserServiceImpl) ToggleStudentStatus(ctx context.Context, id string, operatorID, ip, ua string) (string, error) {
	student, err := s.Repository.FindStudentByID(s.DB, id)
	if err != nil {
		return "", fmt.Errorf("santri tidak ditemukan: %w", err)
	}

	student.IsActive = !student.IsActive
	if err := s.Repository.SaveStudent(s.DB, &student); err != nil {
		return "", fmt.Errorf("gagal mengubah status: %w", err)
	}

	statusStr := "Tidak Aktif"
	if student.IsActive {
		statusStr = "Aktif"
	}
	s.AuditLog.Log(ctx, operatorID, "TOGGLE_STUDENT_STATUS", "Student", id, map[string]string{"status": statusStr}, ip, ua)

	return statusStr, nil
}

// ─── Status Type CRUD ──────────────────────────────────────

func (s *SuperUserServiceImpl) GetStatusTypes() ([]domain.StudentStatusType, error) {
	return s.Repository.FindAllStatusTypes(s.DB)
}

func (s *SuperUserServiceImpl) CreateStatusType(ctx context.Context, req request.CreateStatusTypeRequest, operatorID, ip, ua string) (domain.StudentStatusType, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.StudentStatusType{}, err
	}

	isActiveBilling := *req.IsActiveBilling
	discount := req.DiscountPercentage

	// Logika Abdi Dalem: otomatis tidak dikenakan biaya
	if strings.EqualFold(strings.TrimSpace(req.Name), "abdi dalem") {
		isActiveBilling = false
		discount = 100
	}

	statusType := domain.StudentStatusType{
		ID:                 utils.GenerateID(),
		Name:               req.Name,
		DiscountPercentage: discount,
		IsActiveBilling:    isActiveBilling,
		Description:        req.Description,
	}

	if err := s.Repository.CreateStatusType(s.DB, &statusType); err != nil {
		return domain.StudentStatusType{}, fmt.Errorf("gagal membuat tipe status: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "CREATE_STATUS_TYPE", "StatusType", statusType.ID, req, ip, ua)
	return statusType, nil
}

func (s *SuperUserServiceImpl) UpdateStatusType(ctx context.Context, id string, req request.CreateStatusTypeRequest, operatorID, ip, ua string) (domain.StudentStatusType, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.StudentStatusType{}, err
	}

	statusType, err := s.Repository.FindStatusTypeByID(s.DB, id)
	if err != nil {
		return domain.StudentStatusType{}, fmt.Errorf("tipe status tidak ditemukan: %w", err)
	}

	isActiveBilling := *req.IsActiveBilling
	discount := req.DiscountPercentage

	if strings.EqualFold(strings.TrimSpace(req.Name), "abdi dalem") {
		isActiveBilling = false
		discount = 100
	}

	statusType.Name = req.Name
	statusType.DiscountPercentage = discount
	statusType.IsActiveBilling = isActiveBilling
	statusType.Description = req.Description

	if err := s.Repository.SaveStatusType(s.DB, &statusType); err != nil {
		return domain.StudentStatusType{}, fmt.Errorf("gagal update tipe status: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "UPDATE_STATUS_TYPE", "StatusType", id, req, ip, ua)
	return statusType, nil
}

func (s *SuperUserServiceImpl) DeleteStatusType(ctx context.Context, id string, operatorID, ip, ua string) error {
	if err := s.Repository.DeleteStatusType(s.DB, id); err != nil {
		return fmt.Errorf("gagal menghapus tipe status: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "DELETE_STATUS_TYPE", "StatusType", id, nil, ip, ua)
	return nil
}

// ─── Category CRUD ─────────────────────────────────────────

func (s *SuperUserServiceImpl) GetCategories() ([]domain.Category, error) {
	return s.Repository.FindAllCategories(s.DB)
}

func (s *SuperUserServiceImpl) CreateCategory(ctx context.Context, req request.CreateCategoryRequest, operatorID, ip, ua string) (domain.Category, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.Category{}, err
	}

	category := domain.Category{
		ID:          utils.GenerateID(),
		Name:        req.Name,
		BaseAmount:  req.BaseAmount,
		IsFixed:     req.IsFixed,
		IsActive:    req.IsActive,
		Description: req.Description,
	}

	if err := s.Repository.CreateCategory(s.DB, &category); err != nil {
		return domain.Category{}, fmt.Errorf("gagal membuat kategori: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "CREATE_CATEGORY", "Category", category.ID, req, ip, ua)
	return category, nil
}

func (s *SuperUserServiceImpl) UpdateCategory(ctx context.Context, id string, req request.CreateCategoryRequest, operatorID, ip, ua string) (domain.Category, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.Category{}, err
	}

	category, err := s.Repository.FindCategoryByID(s.DB, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("kategori tidak ditemukan: %w", err)
	}

	category.Name = req.Name
	category.BaseAmount = req.BaseAmount
	category.IsFixed = req.IsFixed
	category.IsActive = req.IsActive
	category.Description = req.Description

	if err := s.Repository.SaveCategory(s.DB, &category); err != nil {
		return domain.Category{}, fmt.Errorf("gagal update kategori: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "UPDATE_CATEGORY", "Category", id, req, ip, ua)
	return category, nil
}

func (s *SuperUserServiceImpl) DeleteCategory(ctx context.Context, id string, operatorID, ip, ua string) error {
	if err := s.Repository.DeleteCategory(s.DB, id); err != nil {
		return fmt.Errorf("gagal menghapus kategori: %w", err)
	}
	s.AuditLog.Log(ctx, operatorID, "DELETE_CATEGORY", "Category", id, nil, ip, ua)
	return nil
}

// ─── Invoice Generation (Semester-based Hijri Calendar) ────

func (s *SuperUserServiceImpl) GenerateInvoices(ctx context.Context, req request.GenerateInvoiceRequest, operatorID, ip, ua string) (int, error) {
	if err := s.Validate.Struct(req); err != nil {
		return 0, err
	}

	// ─── Auto-detect Hijri month/year jika tidak dikirim ────
	hijriMonth := req.HijriMonth
	hijriYear := req.HijriYear
	if hijriMonth == 0 || hijriYear == 0 {
		t := time.Date(req.Year, time.Month(req.Month), 15, 0, 0, 0, 0, time.Local)
		hijri := utils.GregorianToHijri(t)
		hijriMonth = hijri.Month
		hijriYear = hijri.Year
	}

	// ─── Cek apakah bulan ini billable (bukan Ramadhan) ─────
	if !utils.IsBillableMonth(hijriMonth) {
		return 0, fmt.Errorf("bulan %s (Ramadhan) adalah libur, tidak ada tagihan", utils.GetHijriMonthName(hijriMonth))
	}

	semesterInfo := utils.GetSemesterInfo(hijriMonth, hijriYear)
	academicLabel := utils.GetAcademicYearLabel(hijriMonth, hijriYear)

	// ─── Ambil data santri yang eligible ────────────────────
	students, err := s.Repository.FindStudentsForBilling(s.DB)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data santri: %w", err)
	}

	// ─── Tentukan kategori tagihan yang berlaku bulan ini ───
	// 1. Syahriyyah: Selalu ada di setiap bulan aktif
	// 2. Daftar Ulang: Hanya di bulan Syawal (awal tahun ajaran)
	// 3. Ujian Semester: Hanya di Rabi'ul Awal (sem 1) dan Sya'ban (sem 2)
	type billingComponent struct {
		Category domain.Category
		Label    string
	}
	var components []billingComponent

	// Syahriyyah Pondok & Syahriyyah Muhadhoroh — selalu dikenakan
	pondokCat, err := s.Repository.FindCategoryByName(s.DB, "Syahriyyah Pondok")
	if err != nil {
		return 0, fmt.Errorf("kategori Syahriyyah Pondok tidak ditemukan di master data: %w", err)
	}
	components = append(components, billingComponent{Category: pondokCat, Label: "Syahriyyah Pondok"})

	muhadhorohCat, err := s.Repository.FindCategoryByName(s.DB, "Syahriyyah Muhadhoroh")
	if err != nil {
		return 0, fmt.Errorf("kategori Syahriyyah Muhadhoroh tidak ditemukan di master data: %w", err)
	}
	components = append(components, billingComponent{Category: muhadhorohCat, Label: "Syahriyyah Muhadhoroh"})

	// Daftar Ulang — hanya bulan Syawal
	if semesterInfo.IsRegistration {
		regCat, err := s.Repository.FindCategoryByName(s.DB, "Daftar Ulang")
		if err == nil {
			components = append(components, billingComponent{Category: regCat, Label: "Daftar Ulang"})
		}
	}

	// Ujian Semester — hanya bulan ujian (Rabi'ul Awal / Sya'ban)
	if semesterInfo.IsExamMonth {
		examCat, err := s.Repository.FindCategoryByName(s.DB, "Ujian Semester")
		if err == nil {
			components = append(components, billingComponent{Category: examCat, Label: "Ujian Semester"})
		}
	}

	// ─── Generate tagihan per santri ────────────────────────
	generatedCount := 0
	hijriMonthName := utils.GetHijriMonthName(hijriMonth)

	for _, student := range students {
		if !student.Status.IsActiveBilling || student.MuhadhorohLevel == 0 {
			continue // Skip non-billing dan santri lulus
		}
		discount := student.Status.DiscountPercentage

		for _, comp := range components {
			amountDue := comp.Category.BaseAmount * (1 - (discount / 100.0))

			// Cek duplikasi berdasarkan Hijri month/year (idempoten)
			count, _ := s.Repository.CountInvoiceByHijri(s.DB, student.ID, comp.Category.ID, hijriMonth, hijriYear)
			if count > 0 {
				continue
			}

			var sem int
			var acadLabel string
			if comp.Category.Name != "Syahriyyah Pondok" {
				sem = semesterInfo.Number
				acadLabel = academicLabel
			}

			dueDate := utils.HijriToGregorian(hijriYear, hijriMonth, 10)
			gregMonth := int(dueDate.Month())
			gregYear := dueDate.Year()
			invoice := domain.Invoice{
				ID:                utils.GenerateID(),
				StudentID:         student.ID,
				CategoryID:        comp.Category.ID,
				InvoiceNumber:     fmt.Sprintf("INV-%s-%dH%02d-%s", comp.Category.Name, hijriYear, hijriMonth, student.StudentNumber),
				Month:             gregMonth,
				Year:              gregYear,
				HijriMonth:        hijriMonth,
				HijriYear:         hijriYear,
				Semester:          sem,
				AcademicYearLabel: acadLabel,
				AmountDue:         amountDue,
				Status:            "unpaid",
				DueDate:           dueDate,
			}
			if err := s.Repository.CreateInvoice(s.DB, &invoice); err == nil {
				generatedCount++
				// Kirim notifikasi ke wali santri
				for _, g := range student.Guardians {
					var detailMsg string
					if comp.Category.Name == "Syahriyyah Pondok" {
						detailMsg = fmt.Sprintf("Tagihan %s bulan %s %d H untuk santri %s telah tersedia.",
							comp.Label, hijriMonthName, hijriYear, student.FullName())
					} else {
						detailMsg = fmt.Sprintf("Tagihan %s bulan %s %d H (Semester %d) untuk santri %s telah tersedia.",
							comp.Label, hijriMonthName, hijriYear, sem, student.FullName())
					}
					s.Notification.Send(ctx, g.UserID, "Tagihan Baru", detailMsg, "billing")
				}
			}
		}
	}

	if generatedCount > 0 {
		s.AuditLog.Log(ctx, operatorID, "GENERATE_INVOICES", "Invoice", "", map[string]interface{}{
			"month": req.Month, "year": req.Year,
			"hijri_month": hijriMonth, "hijri_year": hijriYear,
			"semester": semesterInfo.Number, "academic_year": academicLabel,
			"components": len(components),
		}, ip, ua)
	}
	return generatedCount, nil
}

// ─── Semester Full-Package Invoice Generation ──────────────
// Generate semua tagihan untuk 1 semester penuh di depan.
// Wali santri bisa melihat total dan bayar paket atau cicil.
func (s *SuperUserServiceImpl) GenerateSemesterInvoices(ctx context.Context, semester int, hijriYear int, operatorID, ip, ua string) (int, error) {
	if semester < 1 || semester > 2 {
		return 0, fmt.Errorf("semester harus 1 atau 2")
	}

	// Ambil bulan-bulan di semester ini
	months := utils.GetSemesterMonths(semester)
	academicLabel := utils.GetAcademicYearLabel(months[0], hijriYear)

	// Ambil santri yang eligible
	students, err := s.Repository.FindStudentsForBilling(s.DB)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data santri: %w", err)
	}

	// Ambil kategori Syahriyyah Muhadhoroh
	muhadhorohCat, err := s.Repository.FindCategoryByName(s.DB, "Syahriyyah Muhadhoroh")
	if err != nil {
		return 0, fmt.Errorf("kategori Syahriyyah Muhadhoroh tidak ditemukan: %w", err)
	}

	// Kategori tambahan
	var regCat, examCat domain.Category
	var hasReg, hasExam bool
	if cat, err := s.Repository.FindCategoryByName(s.DB, "Daftar Ulang"); err == nil {
		regCat = cat
		hasReg = true
	}
	if cat, err := s.Repository.FindCategoryByName(s.DB, "Ujian Semester"); err == nil {
		examCat = cat
		hasExam = true
	}

	generatedCount := 0

	for _, student := range students {
		if !student.Status.IsActiveBilling || student.MuhadhorohLevel == 0 {
			continue // Skip Abdi Dalem dan santri yang sudah lulus
		}
		discount := student.Status.DiscountPercentage

		for _, hijriMonth := range months {
			// Skip Ramadhan
			if !utils.IsBillableMonth(hijriMonth) {
				continue
			}

			info := utils.GetSemesterInfo(hijriMonth, hijriYear)

			// Tentukan komponen bulan ini
			type comp struct {
				cat   domain.Category
				label string
			}
			var components []comp
			components = append(components, comp{cat: muhadhorohCat, label: "Syahriyyah Muhadhoroh"})

			if info.IsRegistration && hasReg {
				components = append(components, comp{cat: regCat, label: "Daftar Ulang"})
			}
			if info.IsExamMonth && hasExam {
				components = append(components, comp{cat: examCat, label: "Ujian Semester"})
			}

			for _, c := range components {
				amountDue := c.cat.BaseAmount * (1 - (discount / 100.0))

				// Cek duplikasi
				count, _ := s.Repository.CountInvoiceByHijri(s.DB, student.ID, c.cat.ID, hijriMonth, hijriYear)
				if count > 0 {
					continue
				}

				dueDate := utils.HijriToGregorian(hijriYear, hijriMonth, 10)
				gregMonth := int(dueDate.Month())
				gregYear := dueDate.Year()
				invoice := domain.Invoice{
					ID:                utils.GenerateID(),
					StudentID:         student.ID,
					CategoryID:        c.cat.ID,
					InvoiceNumber:     fmt.Sprintf("INV-%s-%dH%02d-%s", c.cat.Name, hijriYear, hijriMonth, student.StudentNumber),
					Month:             gregMonth,
					Year:              gregYear,
					HijriMonth:        hijriMonth,
					HijriYear:         hijriYear,
					Semester:          semester,
					AcademicYearLabel: academicLabel,
					AmountDue:         amountDue,
					Status:            "unpaid",
					DueDate:           dueDate,
				}
				if err := s.Repository.CreateInvoice(s.DB, &invoice); err == nil {
					generatedCount++
				}
			}
		}

		// Kirim 1 notifikasi ringkasan semester ke wali
		if generatedCount > 0 {
			for _, g := range student.Guardians {
				s.Notification.Send(ctx, g.UserID, "Tagihan Semester",
					fmt.Sprintf("Tagihan Semester %d Tahun Ajaran %s untuk santri %s (Muhadhoroh %d) telah diterbitkan. Anda dapat membayar secara paket atau dicicil per bulan.",
						semester, academicLabel, student.Name.FirstName, student.MuhadhorohLevel),
					"billing")
			}
		}
	}

	if generatedCount > 0 {
		s.AuditLog.Log(ctx, operatorID, "GENERATE_SEMESTER_INVOICES", "Invoice", "", map[string]interface{}{
			"semester": semester, "hijri_year": hijriYear,
			"academic_year": academicLabel, "total_invoices": generatedCount,
		}, ip, ua)
	}
	return generatedCount, nil
}

func (s *SuperUserServiceImpl) GetSettings() ([]domain.Setting, error) {
	return s.SettingRepo.FindAll(s.DB)
}

func (s *SuperUserServiceImpl) UpdateSetting(ctx context.Context, key, value, operatorID, ip, ua string) error {
	setting, err := s.SettingRepo.FindByKey(s.DB, key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			setting = domain.Setting{
				Key:   key,
				Value: value,
			}
			if err := s.SettingRepo.Save(s.DB, &setting); err != nil {
				return fmt.Errorf("gagal buat setting baru: %w", err)
			}
		} else {
			return fmt.Errorf("gagal cari setting: %w", err)
		}
	} else {
		setting.Value = value
		if err := s.SettingRepo.Save(s.DB, &setting); err != nil {
			return fmt.Errorf("gagal simpan setting: %w", err)
		}
	}

	s.AuditLog.Log(ctx, operatorID, "UPDATE_SETTING", "Setting", key, map[string]string{"value": value}, ip, ua)
	return nil
}

// ─── User Management ───────────────────────────────────────

func (s *SuperUserServiceImpl) GetAllUsers() ([]domain.User, error) {
	return s.Repository.FindAllUsers(s.DB)
}

func (s *SuperUserServiceImpl) GetAllUsersPaginated(page int, limit int) ([]domain.User, int64, error) {
	return s.Repository.FindAllUsersPaginated(s.DB, page, limit)
}

func (s *SuperUserServiceImpl) GetAllRoles() ([]domain.Role, error) {
	return s.Repository.FindAllRoles(s.DB)
}

func (s *SuperUserServiceImpl) UpdateUserRole(ctx context.Context, userID, roleID string) error {
	if err := s.Repository.UpdateUserRole(s.DB, userID, roleID); err != nil {
		return fmt.Errorf("gagal mengubah role user: %w", err)
	}
	return nil
}

func (s *SuperUserServiceImpl) ToggleUserActive(ctx context.Context, userID string) (bool, error) {
	return s.Repository.ToggleUserActive(s.DB, userID)
}

func (s *SuperUserServiceImpl) GetActivityLogsPaginated(page int, limit int) ([]domain.ActivityLog, int64, error) {
	return s.Repository.FindAllActivityLogsPaginated(s.DB, page, limit)
}

func (s *SuperUserServiceImpl) SendPreBillingNotifications(ctx context.Context, nextMonth int, nextYear int, targetDay int) error {
	students, err := s.Repository.FindStudentsForBilling(s.DB)
	if err != nil {
		return fmt.Errorf("gagal mengambil data santri untuk notifikasi: %w", err)
	}

	monthName := utils.GetHijriMonthName(nextMonth)

	for _, student := range students {
		if !student.Status.IsActiveBilling || student.MuhadhorohLevel == 0 {
			continue
		}
		discount := student.Status.DiscountPercentage

		// Hitung estimasi tagihan
		pondokCat, errP := s.Repository.FindCategoryByName(s.DB, "Syahriyyah Pondok")
		muhadhorohCat, errM := s.Repository.FindCategoryByName(s.DB, "Syahriyyah Muhadhoroh")

		var estTotal float64
		if errP == nil {
			estTotal += pondokCat.BaseAmount * (1 - (discount / 100.0))
		}
		if errM == nil {
			estTotal += muhadhorohCat.BaseAmount * (1 - (discount / 100.0))
		}

		message := fmt.Sprintf(
			"Pemberitahuan: Tagihan Syahriyyah bulan mendatang (%s %d H) untuk santri %s akan otomatis diterbitkan dalam 5 hari (pada tanggal %d %s %d H). Estimasi tagihan adalah Rp %.2f (sudah termasuk potongan status '%s'). Mohon persiapkan dana pembayaran Anda.",
			monthName, nextYear, student.FullName(), targetDay, monthName, nextYear, estTotal, student.Status.Name,
		)

		for _, g := range student.Guardians {
			s.Notification.Send(ctx, g.UserID, "Pemberitahuan Tagihan Mendatang (H-5)", message, "billing")
		}
	}
	return nil
}
