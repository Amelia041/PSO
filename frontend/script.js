// Konfigurasi Mock Data
const USE_MOCK = true; 
const API = window.location.origin;

// Inisialisasi Database Lokal (Mock) di Browser
if (!localStorage.getItem('mock_notes')) {
    const initialData = [
        { 
            id: 1, 
            user_id: 1, 
            title: "[Kelas] Pemodelan Sistem Objek", 
            category: "PSO", 
            content: "Pengenalan materi semester 4.", 
            date: "2026-05-12T08:00" 
        },
        { 
            id: 2, 
            user_id: 1, 
            title: "[Tugas] Implementasi CRUD Go", 
            category: "Database", 
            content: "Mengerjakan ETL ke staging area.", 
            date: "2026-05-13T23:59" 
        }
    ];
    localStorage.setItem('mock_notes', JSON.stringify(initialData));
}

// --- FUNGSI AUTHENTICATION ---
async function login() {
    const u = document.getElementById('username').value;
    const p = document.getElementById('password').value;

    if (!u || !p) {
        alert("Username dan Password tidak boleh kosong!");
        return;
    }

    if (USE_MOCK) {
        // Simulasi Login Sukses
        localStorage.setItem('user_id', '1');
        localStorage.setItem('username', u);
        window.location.href = 'dashboard.html';
    } else {
        try {
            const res = await fetch(`${API}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: u, password: p })
            });
            const data = await res.json();
            if (res.ok) {
                localStorage.setItem('user_id', data.id);
                localStorage.setItem('username', data.username);
                window.location.href = 'dashboard.html';
            } else {
                alert(data.error || "Login Gagal!");
            }
        } catch (err) {
            alert("Backend tidak merespon. Aktifkan USE_MOCK untuk testing UI.");
        }
    }
}

function logout() {
    localStorage.clear();
    window.location.href = 'index.html';
}

// --- FUNGSI SIMPAN DATA (NOTES & SCHEDULE) ---
async function saveEntry(type) {
    const userId = localStorage.getItem('user_id');
    const titleInput = document.getElementById(type === 'note' ? 'note-title' : 'sch-title');
    const catInput = document.getElementById(type === 'note' ? 'note-category' : 'sch-category');
    
    // Khusus Schedule
    const timeInput = document.getElementById('sch-time');
    const roomInput = document.getElementById('sch-room');
    
    // Khusus Note
    const contentInput = document.getElementById('note-content');

    if (!titleInput.value) {
        alert("Judul wajib diisi!");
        return;
    }

    const payload = {
        id: Date.now(),
        user_id: parseInt(userId),
        title: type === 'note' ? titleInput.value : `[Jadwal] ${titleInput.value}`,
        category: catInput.value,
        content: type === 'note' ? contentInput.value : `Lokasi: ${roomInput.value}`,
        date: type === 'note' ? new Date().toISOString() : timeInput.value
    };

    if (USE_MOCK) {
        let notes = JSON.parse(localStorage.getItem('mock_notes'));
        notes.push(payload);
        localStorage.setItem('mock_notes', JSON.stringify(notes));
        alert("Berhasil disimpan (Mock Mode)!");
        window.location.href = 'dashboard.html';
    } else {
        try {
            const res = await fetch(`${API}/notes`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            if (res.ok) {
                alert("Berhasil disimpan ke Database!");
                window.location.href = 'dashboard.html';
            }
        } catch (err) {
            alert("Gagal terhubung ke Backend.");
        }
    }
}

// --- FUNGSI RENDER DATA ---
function renderDashboard() {
    const userName = localStorage.getItem('username');
    const uid = localStorage.getItem('user_id');
    
    if (!uid) {
        window.location.href = 'index.html';
        return;
    }

    document.getElementById('user-name').innerText = userName;
    document.getElementById('current-date').innerText = new Date().toLocaleDateString('id-ID', { 
        weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' 
    });

    const notes = JSON.parse(localStorage.getItem('mock_notes')) || [];
    const recentContainer = document.getElementById('recent-notes');

    if (recentContainer) {
        // Tampilkan 3 data terbaru
        const latest = notes.filter(n => n.user_id == uid).reverse().slice(0, 3);
        recentContainer.innerHTML = latest.map(n => `
            <div class="note-box">
                <small style="color:var(--its-blue); font-weight:bold;">${n.category}</small>
                <div style="font-weight:600;">${n.title}</div>
                <div style="font-size:0.75rem; color:#718096;">${new Date(n.date).toLocaleString('id-ID')}</div>
            </div>
        `).join('');
    }
}

function renderGroups() {
    const uid = localStorage.getItem('user_id');
    const notes = JSON.parse(localStorage.getItem('mock_notes')) || [];
    const categories = ['PSO', 'Progjar', 'Database', 'Umum'];

    categories.forEach(cat => {
        const container = document.getElementById(`group-${cat}`);
        const countBadge = document.getElementById(`count-${cat}`);
        
        const filtered = notes.filter(n => n.user_id == uid && n.category === cat);
        
        if (countBadge) countBadge.innerText = filtered.length;
        if (container) {
            container.innerHTML = filtered.length > 0 ? filtered.map(n => `
                <div class="note-box" style="background: white; border: 1px solid #edf2f7;">
                    <div style="font-weight:600; font-size:0.9rem;">${n.title}</div>
                    <small style="color:#a0aec0;">${new Date(n.date).toLocaleDateString('id-ID')}</small>
                </div>
            `).join('') : '<p style="color:#cbd5e0; font-size:0.8rem; text-align:center;">Kosong</p>';
        }
    });
}