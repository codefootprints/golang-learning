package main

import (
	"testing"
)

func TestCreateTaskLogic(t *testing.T) {
	initDB()

	// 1. Setup data dummy
	testTask := Task{
		Title:  "Cleanup Test Task",
		Status: "Testing",
	}

	// 2. Eksekusi
	created, err := CreateTaskLogic(testTask)
	if err != nil {
		t.Fatalf("Gagal membuat task: %v", err)
	}

	// --- BAGIAN CLEANUP ---
	// Perintah ini akan dijalankan OTOMATIS setelah fungsi test selesai
	t.Cleanup(func() {
		db.Unscoped().Delete(&created)
		// Unscoped() digunakan agar data benar-benar terhapus dari tabel (bukan soft delete)
	})
	// ----------------------

	// 3. Assert (Validasi)
	if created.ID == 0 {
		t.Error("Harusnya ID tidak nol setelah disimpan")
	}

	if created.Title != testTask.Title {
		t.Errorf("Ekspektasi %s, dapat %s", testTask.Title, created.Title)
	}
}
