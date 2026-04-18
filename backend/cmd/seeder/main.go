package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/database"
	"github.com/wangsa/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var createSuperAdmin = flag.Bool("create-admin", false, "Create default superadmin account")
	var createSampleEvents = flag.Bool("create-events", false, "Create sample events for testing")
	var createSampleNotes = flag.Bool("create-notes", false, "Create sample notes for testing")
	flag.Parse()

	if !*createSuperAdmin && !*createSampleEvents && !*createSampleNotes {
		log.Println("Usage:")
		log.Println("  go run cmd/seeder/main.go -create-admin     # Create superadmin account")
		log.Println("  go run cmd/seeder/main.go -create-events    # Create sample events")
		log.Println("  go run cmd/seeder/main.go -create-notes     # Create sample notes")
		log.Println("  go run cmd/seeder/main.go -create-admin -create-events -create-notes   # Create all")
		os.Exit(1)
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations first
	if err := database.RunMigrations(db, ""); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	noteRepo := repository.NewNoteRepository(db)

	// Create superadmin
	if *createSuperAdmin {
		if err := createSuperAdminAccount(userRepo); err != nil {
			log.Fatalf("Failed to create superadmin account: %v", err)
		}
		log.Println("✅ Superadmin account created successfully!")
		log.Println("📧 Email: admin@wangsa.internal")
		log.Println("🔑 Password: WangsaAdmin2024!")
		log.Println("⚠️  Please change the password after first login")
	}

	// Create sample events
	if *createSampleEvents {
		if err := createSampleEventsData(eventRepo, userRepo); err != nil {
			log.Fatalf("Failed to create sample events: %v", err)
		}
		log.Println("✅ Sample events created successfully!")
	}

	// Create sample notes
	if *createSampleNotes {
		if err := createSampleNotesData(noteRepo, userRepo); err != nil {
			log.Fatalf("Failed to create sample notes: %v", err)
		}
		log.Println("✅ Sample notes created successfully!")
	}
}

func createSuperAdminAccount(userRepo domain.UserRepository) error {
	// Check if admin already exists
	existingAdmin, err := userRepo.FindByEmail("admin@wangsa.internal")
	if err != nil {
		return err
	}

	if existingAdmin != nil {
		log.Println("⚠️  Superadmin account already exists, skipping creation")
		return nil
	}

	// Hash password
	password := "WangsaAdmin2024!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create superadmin user
	admin := &domain.User{
		Name:     "Super Administrator",
		Email:    "admin@wangsa.internal",
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	return userRepo.Create(admin)
}

func createSampleEventsData(eventRepo domain.EventRepository, userRepo domain.UserRepository) error {
	// Get admin user to use as creator
	admin, err := userRepo.FindByEmail("admin@wangsa.internal")
	if err != nil {
		return err
	}
	if admin == nil {
		log.Println("⚠️  Superadmin not found, creating sample events with user ID 1")
		admin = &domain.User{ID: 1} // Fallback to ID 1
	}

	// Get current date for creating events around this month
	now := time.Now()

	sampleEvents := []*domain.Event{
		{
			Title:       "Reuni Keluarga Besar",
			Description: strPtr("Acara reuni tahunan keluarga besar di rumah nenek"),
			Location:    strPtr("Rumah Nenek, Yogyakarta"),
			StartAt:     time.Date(now.Year(), now.Month(), 15, 10, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 15, 17, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#CC6649",
			Notes:       strPtr("<h3>Agenda Reuni</h3><ol><li><strong>Pembukaan</strong> - Salam dan doa pembuka</li><li><strong>Laporan Kegiatan</strong> - Update dari setiap keluarga</li><li><strong>Rencana Tahun Depan</strong> - Diskusi acara mendatang</li></ol><p><strong>Catatan:</strong> Mohon setiap keluarga membawa makanan khas daerah masing-masing.</p>"),
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Ulang Tahun Mama",
			Description: strPtr("Perayaan ulang tahun mama yang ke-65"),
			Location:    strPtr("Rumah Mama"),
			StartAt:     time.Date(now.Year(), now.Month(), 20, 18, 30, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 20, 22, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#10B981",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Arisan RT",
			Description: strPtr("Arisan bulanan RT 05/RW 02"),
			Location:    strPtr("Balai RT"),
			StartAt:     time.Date(now.Year(), now.Month(), 25, 19, 0, 0, 0, now.Location()),
			IsRecurring: false,
			Color:       "#3B82F6",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Pengajian Keluarga",
			Description: strPtr("Pengajian rutin keluarga setiap akhir bulan"),
			Location:    strPtr("Masjid Al-Ikhlas"),
			StartAt:     time.Date(now.Year(), now.Month(), 28, 19, 30, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 28, 21, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#8B5CF6",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Rapat Keluarga",
			Description: strPtr("Rapat koordinasi acara akhir tahun"),
			Location:    strPtr("Rumah Pak Lurah"),
			StartAt:     time.Date(now.Year(), now.Month(), 30, 14, 0, 0, 0, now.Location()),
			IsRecurring: false,
			Color:       "#F59E0B",
			Notes:       strPtr("<h3>Notulensi Rapat</h3><p><strong>Hadir:</strong> 15 orang</p><h4>Keputusan Rapat:</h4><ul><li>Acara akhir tahun akan diadakan tanggal 28 Desember</li><li>Budget yang disepakati: <strong>Rp 5.000.000</strong></li><li>Panitia: Pak Budi (koordinator), Bu Sari (konsumsi)</li></ul><p><em>Rapat selesai pukul 16:30 WIB</em></p>"),
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		// Events for next month
		{
			Title:       "Piknik Keluarga",
			Description: strPtr("Piknik bersama di pantai parangtritis"),
			Location:    strPtr("Pantai Parangtritis"),
			StartAt:     time.Date(now.Year(), now.Month()+1, 5, 8, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month()+1, 5, 18, 0, 0, 0, now.Location())),
			IsRecurring: false,
			Color:       "#06D6A0",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Workshop Memasak",
			Description: strPtr("Belajar masak gudeg bareng-bareng"),
			Location:    strPtr("Rumah Tante Sari"),
			StartAt:     time.Date(now.Year(), now.Month()+1, 12, 9, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month()+1, 12, 14, 0, 0, 0, now.Location())),
			IsRecurring: false,
			Color:       "#FF6B6B",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
	}

	for _, event := range sampleEvents {
		if err := eventRepo.Create(event); err != nil {
			log.Printf("⚠️  Failed to create event %s: %v", event.Title, err)
		} else {
			log.Printf("✅ Created event: %s", event.Title)
		}
	}

	return nil
}

func createSampleNotesData(noteRepo domain.NoteRepository, userRepo domain.UserRepository) error {
	// Get admin user to use as creator
	admin, err := userRepo.FindByEmail("admin@wangsa.internal")
	if err != nil {
		return err
	}
	if admin == nil {
		log.Println("⚠️  Superadmin not found, creating sample notes with user ID 1")
		admin = &domain.User{ID: 1} // Fallback to ID 1
	}

	sampleNotes := []*domain.Note{
		{
			Title:     "📞 Nomor Telepon Penting",
			Content:   `<h3>Kontak Darurat Keluarga</h3><ul><li><strong>Rumah Sakit:</strong> (0274) 123-456</li><li><strong>Puskesmas:</strong> (0274) 789-012</li><li><strong>Polsek:</strong> (0274) 345-678</li></ul><h3>Kontak Keluarga</h3><ul><li><strong>Om Budi:</strong> 0812-3456-7890</li><li><strong>Tante Sari:</strong> 0856-7890-1234</li><li><strong>Pak RT:</strong> 0878-9012-3456</li></ul><p><em>Update terakhir: Januari 2025</em></p>`,
			Category:  strPtr("Kontak"),
			IsPinned:  true,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "🏠 Alamat Anggota Keluarga",
			Content:   `<h3>Daftar Alamat Terbaru</h3><table><tr><th>Nama</th><th>Alamat</th><th>Telepon</th></tr><tr><td>Keluarga Budi</td><td>Jl. Malioboro 15, Yogyakarta</td><td>0812-1111-2222</td></tr><tr><td>Keluarga Sari</td><td>Jl. Solo 42, Klaten</td><td>0813-3333-4444</td></tr><tr><td>Keluarga Andi</td><td>Jl. Magelang 88, Magelang</td><td>0814-5555-6666</td></tr></table><p><strong>Catatan:</strong> Mohon update jika ada perubahan alamat.</p>`,
			Category:  strPtr("Kontak"),
			IsPinned:  true,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "💊 Informasi Kesehatan",
			Content:   `<h3>Catatan Medis Keluarga</h3><h4>Alergi dan Kondisi Khusus:</h4><ul><li><strong>Mama:</strong> Alergi udang, diabetes</li><li><strong>Papa:</strong> Hipertensi, pantang garam</li><li><strong>Adik Rina:</strong> Alergi debu, asma ringan</li></ul><h4>Obat Rutin:</h4><ul><li><strong>Mama:</strong> Metformin 2x sehari</li><li><strong>Papa:</strong> Amlodipin 1x pagi</li></ul><p><em>⚠️ Selalu konsultasi dengan dokter sebelum mengubah dosis</em></p>`,
			Category:  strPtr("Kesehatan"),
			IsPinned:  false,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "🎉 Daftar Ulang Tahun",
			Content:   `<h3>Kalender Ulang Tahun Keluarga</h3><h4>Januari</h4><ul><li>15 Jan - Om Budi (1975)</li><li>28 Jan - Tante Dewi (1980)</li></ul><h4>Februari</h4><ul><li>3 Feb - Adik Rina (2005)</li><li>20 Feb - Mama (1958)</li></ul><h4>Maret</h4><ul><li>12 Mar - Papa (1955)</li><li>30 Mar - Kakak Doni (1985)</li></ul><p><em>💡 Jangan lupa siapkan kado dan ucapan selamat!</em></p>`,
			Category:  strPtr("Keluarga"),
			IsPinned:  false,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "💰 Catatan Kas Keluarga",
			Content:   `<h3>Aturan Kas Keluarga</h3><ol><li><strong>Iuran bulanan:</strong> Rp 50.000 per keluarga</li><li><strong>Deadline:</strong> Setiap tanggal 10</li><li><strong>Penggunaan:</strong> Acara keluarga, santunan, dll</li></ol><h3>Rekening Kas</h3><ul><li><strong>Bank:</strong> BCA</li><li><strong>No. Rekening:</strong> 1234567890</li><li><strong>Atas Nama:</strong> Sari Wangsa</li></ul><p><strong>Saldo terakhir (Des 2024):</strong> Rp 2.500.000</p>`,
			Category:  strPtr("Keuangan"),
			IsPinned:  false,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "🍽️ Resep Masakan Keluarga",
			Content:   `<h3>Resep Gudeg Nenek</h3><h4>Bahan:</h4><ul><li>Nangka muda 1 kg</li><li>Santan kental 500ml</li><li>Gula merah 200gr</li><li>Daun salam, daun jeruk</li><li>Bumbu halus: bawang merah, bawang putih, kemiri, ketumbar</li></ul><h4>Cara masak:</h4><ol><li>Rebus nangka hingga empuk</li><li>Tumis bumbu halus</li><li>Masukkan nangka, santan, gula merah</li><li>Masak dengan api kecil 2-3 jam</li></ol><p><em>"Rahasia nenek: masak pakai kayu bakar biar wangi!" - Nenek (2020)</em></p>`,
			Category:  strPtr("Resep"),
			IsPinned:  false,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "📋 Checklist Acara Keluarga",
			Content:   `<h3>Template Persiapan Acara</h3><h4>1 Bulan Sebelum:</h4><ul><li>☐ Tentukan tanggal dan tempat</li><li>☐ Buat grup chat koordinasi</li><li>☐ Booking tempat (jika perlu)</li></ul><h4>2 Minggu Sebelum:</h4><ul><li>☐ Konfirmasi kehadiran</li><li>☐ Atur konsumsi</li><li>☐ Siapkan dokumentasi</li></ul><h4>1 Hari Sebelum:</h4><ul><li>☐ Belanja bahan makanan</li><li>☐ Persiapkan tempat</li><li>☐ Reminder ke semua anggota</li></ul><p><em>Salin checklist ini untuk setiap acara baru</em></p>`,
			Category:  strPtr("Acara"),
			IsPinned:  false,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
		{
			Title:     "🎯 Rencana 2025",
			Content:   `<h3>Target Keluarga Tahun 2025</h3><h4>Acara Rutin:</h4><ul><li>Reuni tahunan - <strong>Juli 2025</strong></li><li>Piknik keluarga - <strong>Agustus 2025</strong></li><li>Syukuran akhir tahun - <strong>Desember 2025</strong></li></ul><h4>Project Keluarga:</h4><ul><li>Renovasi rumah nenek</li><li>Buku silsilah keluarga digital</li><li>Arisan motor untuk anak muda</li></ul><h4>Target Tabungan:</h4><ul><li>Kas keluarga: <strong>Rp 10.000.000</strong></li><li>Dana darurat: <strong>Rp 5.000.000</strong></li></ul>`,
			Category:  strPtr("Rencana"),
			IsPinned:  true,
			CreatedBy: int64Ptr(admin.ID),
			UpdatedBy: int64Ptr(admin.ID),
		},
	}

	for _, note := range sampleNotes {
		if err := noteRepo.Create(note); err != nil {
			log.Printf("⚠️  Failed to create note %s: %v", note.Title, err)
		} else {
			log.Printf("✅ Created note: %s", note.Title)
		}
	}

	return nil
}

// Helper functions for pointers
func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func int64Ptr(i int64) *int64 {
	return &i
}
