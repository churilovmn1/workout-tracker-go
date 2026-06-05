const API = '/api';
let token = localStorage.getItem('token') || '';
let currentUser = null;

const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

async function api(method, path, body) {
    const opts = { method, headers: { 'Content-Type': 'application/json' } };
    if (token) opts.headers['Authorization'] = 'Bearer ' + token;
    if (body !== undefined) opts.body = JSON.stringify(body);
    const res = await fetch(API + path, opts);
    if (res.status === 204) return null;
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Request failed');
    return data;
}

function toast(msg, type = 'success') {
    const el = document.createElement('div');
    el.className = 'toast ' + type;
    el.textContent = msg;
    document.body.appendChild(el);
    setTimeout(() => el.remove(), 3000);
}

// ── Auth ──────────────────────────────────────────────────────────────────

function initAuth() {
    if (token) { showApp(); return; }
    $('#auth-screen').style.display = '';
    $('#app-screen').style.display = 'none';
}

$$('.auth-tabs button').forEach((btn) => {
    btn.addEventListener('click', () => {
        $$('.auth-tabs button').forEach((b) => b.classList.remove('active'));
        btn.classList.add('active');
        const mode = btn.dataset.mode;
        $('#register-fields').style.display = mode === 'register' ? '' : 'none';
        $('#auth-submit').textContent = mode === 'register' ? 'Зарегистрироваться' : 'Войти';
        $('#auth-form').dataset.mode = mode;
    });
});

$('#auth-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const mode = e.target.dataset.mode;
    const login = $('#auth-login').value;
    const password = $('#auth-password').value;
    try {
        if (mode === 'register') {
            await api('POST', '/auth/register', { login, email: $('#auth-email').value, password });
            toast('Регистрация успешна!');
        }
        const data = await api('POST', '/auth/login', { login, password });
        token = data.token;
        localStorage.setItem('token', token);
        showApp();
    } catch (err) {
        toast(err.message, 'error');
    }
});

function showApp() {
    $('#auth-screen').style.display = 'none';
    $('#app-screen').style.display = '';
    parseToken();   // sets currentUser and nav visibility
    const startPage = currentUser.role === 'admin' ? 'schedule' : 'workouts';
    navigate(startPage);
}

function parseToken() {
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        currentUser = { id: payload.user_id, login: payload.login, role: payload.role };
        setupNav();
    } catch {
        logout();
    }
}

function setupNav() {
    const isAdmin = currentUser.role === 'admin';
    // Workouts + Stats: only for regular users
    $('#nav-workouts').style.display = isAdmin ? 'none' : '';
    $('#nav-stats').style.display    = isAdmin ? 'none' : '';
    // Schedule + Trainer tabs: only for admins
    $('#nav-schedule').style.display = isAdmin ? '' : 'none';
    $('#nav-admin').style.display    = isAdmin ? '' : 'none';
    // Username label
    $('.user-name').textContent = isAdmin ? 'Тренер' : (currentUser.login || 'User #' + currentUser.id);
}

function logout() {
    token = '';
    currentUser = null;
    localStorage.removeItem('token');
    // Reset nav visibility
    $('#nav-workouts').style.display = '';
    $('#nav-schedule').style.display = 'none';
    $('#nav-admin').style.display    = 'none';
    $('#auth-screen').style.display = '';
    $('#app-screen').style.display  = 'none';
}

$('#logout-btn').addEventListener('click', logout);

// ── Navigation ────────────────────────────────────────────────────────────

function navigate(page) {
    $$('.page').forEach((p) => p.classList.remove('active'));
    $$('.nav-links button').forEach((b) => b.classList.remove('active'));

    const pageEl = $('#page-' + page);
    if (!pageEl) return;
    pageEl.classList.add('active');

    const btn = $(`.nav-links button[data-page="${page}"]`);
    if (btn) btn.classList.add('active');

    const loaders = {
        workouts: loadWorkouts,
        schedule: loadSchedule,
        exercises: loadExercises,
        stats:    loadStats,
        admin:    loadAdminClients,
    };
    if (loaders[page]) loaders[page]();
}

$$('.nav-links button').forEach((btn) => {
    btn.addEventListener('click', () => navigate(btn.dataset.page));
});

// ── Workouts (users) ──────────────────────────────────────────────────────

async function loadWorkouts() {
    try {
        const workouts = await api('GET', '/workouts');
        const list = $('#workouts-list');
        if (!workouts || !workouts.length) {
            list.innerHTML = '<div class="card"><p>Нет тренировок. Создайте первую!</p></div>';
            return;
        }
        list.innerHTML = workouts.map((w) => `
            <div class="card">
                <div class="card-header">
                    <div>
                        <div class="card-title">${esc(w.title)}</div>
                        <div class="card-subtitle">${formatDate(w.date)} · ${w.duration_minutes} мин</div>
                    </div>
                    <div class="card-actions">
                        <button class="btn btn-sm btn-outline" onclick="copyWorkout(${w.id})">Копировать</button>
                        <button class="btn btn-sm btn-danger" onclick="deleteWorkout(${w.id})">Удалить</button>
                    </div>
                </div>
                ${w.notes ? '<p style="color:var(--text-muted);font-size:0.85rem">' + esc(w.notes) + '</p>' : ''}
                ${w.trainer_comment ? '<div class="trainer-comment">💬 Тренер: ' + esc(w.trainer_comment) + '</div>' : ''}
            </div>
        `).join('');
    } catch (err) { toast(err.message, 'error'); }
}

$('#btn-new-workout').addEventListener('click', () => {
    // Reset exercise filter when opening modal
    $('#modal-ex-search').value = '';
    $('#modal-ex-group').value  = '';
    $('#workout-modal').classList.add('active');
    loadExerciseOptions();
});

$('#close-workout-modal').addEventListener('click', () => $('#workout-modal').classList.remove('active'));
$('#workout-modal').addEventListener('click', (e) => {
    if (e.target === $('#workout-modal')) $('#workout-modal').classList.remove('active');
});

// Returns exercises matching current modal filter
function getModalFilteredExercises() {
    const search = ($('#modal-ex-search')?.value || '').toLowerCase().trim();
    const group  = $('#modal-ex-group')?.value || '';
    return (window._exercises || []).filter((e) => {
        if (group && e.muscle_group !== group) return false;
        if (search && !e.name.toLowerCase().includes(search)) return false;
        return true;
    });
}

function buildSelectOptions(exercises) {
    if (!exercises.length) return '<option value="" disabled>Не найдено</option>';
    return exercises
        .map((e) => `<option value="${e.id}">${esc(e.name)} (${esc(e.muscle_group)})</option>`)
        .join('');
}

// Rebuilds options in all existing exercise selects after filter changes
function updateExerciseSelects() {
    const opts = buildSelectOptions(getModalFilteredExercises());
    $$('#workout-exercises .exercise-row select[name=exercise_id]').forEach((sel) => {
        sel.innerHTML = opts;
    });
}

let modalExSearchTimer = null;
$('#modal-ex-search').addEventListener('input', () => {
    clearTimeout(modalExSearchTimer);
    modalExSearchTimer = setTimeout(updateExerciseSelects, 200);
});
$('#modal-ex-group').addEventListener('change', updateExerciseSelects);

async function loadExerciseOptions() {
    try {
        const exercises = await api('GET', '/exercises');
        window._exercises = exercises || [];
        $('#workout-exercises').innerHTML = '';
        addExerciseRow();
    } catch (err) { toast(err.message, 'error'); }
}

function addExerciseRow() {
    const opts = buildSelectOptions(getModalFilteredExercises());
    const row = document.createElement('div');
    row.className = 'exercise-row';
    row.innerHTML = `
        <select name="exercise_id">${opts}</select>
        <input type="number" name="sets" placeholder="Подходы" min="1" value="3">
        <input type="number" name="reps" placeholder="Повторы" min="1" value="10">
        <input type="number" name="weight" placeholder="Вес (кг)" min="0" step="0.5" value="0">
        <button class="btn btn-sm btn-danger" onclick="this.parentElement.remove()">×</button>`;
    $('#workout-exercises').appendChild(row);
}

$('#add-exercise-row').addEventListener('click', addExerciseRow);

$('#workout-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const exercises = [];
    $$('#workout-exercises .exercise-row').forEach((row) => {
        exercises.push({
            exercise_id: parseInt(row.querySelector('[name=exercise_id]').value),
            sets: parseInt(row.querySelector('[name=sets]').value) || 0,
            reps: parseInt(row.querySelector('[name=reps]').value) || 0,
            weight_kg: parseFloat(row.querySelector('[name=weight]').value) || 0,
        });
    });
    try {
        await api('POST', '/workouts', {
            title: $('#w-title').value,
            date: $('#w-date').value || new Date().toISOString().slice(0, 10),
            duration_minutes: parseInt($('#w-duration').value) || 0,
            notes: $('#w-notes').value,
            exercises,
        });
        $('#workout-modal').classList.remove('active');
        toast('Тренировка сохранена!');
        loadWorkouts();
    } catch (err) { toast(err.message, 'error'); }
});

async function deleteWorkout(id) {
    if (!confirm('Удалить тренировку?')) return;
    try {
        await api('DELETE', '/workouts/' + id);
        toast('Тренировка удалена');
        loadWorkouts();
    } catch (err) { toast(err.message, 'error'); }
}

async function copyWorkout(id) {
    try {
        await api('POST', '/workouts/' + id + '/copy');
        toast('Тренировка скопирована!');
        loadWorkouts();
    } catch (err) { toast(err.message, 'error'); }
}

// ── Schedule (trainer) ────────────────────────────────────────────────────

let currentWeekStart = getMonday(new Date());
let scheduleEntries = [];
let detailEntry = null;

function getMonday(d) {
    const dt = new Date(d);
    const day = dt.getDay();
    const diff = dt.getDate() - day + (day === 0 ? -6 : 1);
    dt.setDate(diff);
    dt.setHours(0, 0, 0, 0);
    return dt;
}

function addDays(d, n) {
    const r = new Date(d);
    r.setDate(r.getDate() + n);
    return r;
}

// Uses LOCAL year/month/day — avoids UTC shift that pushed week to wrong Monday.
function toISODate(d) {
    return [
        d.getFullYear(),
        String(d.getMonth() + 1).padStart(2, '0'),
        String(d.getDate()).padStart(2, '0'),
    ].join('-');
}

function changeWeek(dir) {
    currentWeekStart = addDays(currentWeekStart, dir * 7);
    loadSchedule();
}

async function loadSchedule() {
    const weekParam = toISODate(currentWeekStart);
    try {
        scheduleEntries = await api('GET', '/admin/schedule?week=' + weekParam) || [];
        renderSchedule();
    } catch (err) { toast(err.message, 'error'); }
}

function renderSchedule() {
    const weekEnd = addDays(currentWeekStart, 6);
    const DAY_NAMES = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];
    const todayStr = toISODate(new Date());

    // Week label
    const fmt = (d) => d.toLocaleDateString('ru-RU', { day: '2-digit', month: 'short' });
    $('#week-label').textContent = fmt(currentWeekStart) + ' – ' + fmt(weekEnd);

    // Group entries by day index 0..6 (Mon..Sun)
    const byDay = Array.from({ length: 7 }, () => []);
    for (const e of scheduleEntries) {
        const d = new Date(e.scheduled_at);
        const wd = d.getDay(); // 0=Sun
        const idx = wd === 0 ? 6 : wd - 1; // Mon=0..Sun=6
        byDay[idx].push(e);
    }

    let html = '<div class="schedule-week">';
    for (let i = 0; i < 7; i++) {
        const day = addDays(currentWeekStart, i);
        const dayStr = toISODate(day);
        const isToday = dayStr === todayStr;
        const label = DAY_NAMES[i] + ' ' + day.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' });

        html += `<div class="schedule-day">
            <div class="schedule-day-header${isToday ? ' today' : ''}">${label}</div>`;

        for (const e of byDay[i]) {
            const at = new Date(e.scheduled_at);
            const timeStr = at.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
            const idx = scheduleEntries.indexOf(e);
            html += `<div class="schedule-slot ${e.status}" onclick="openSlotDetail(${idx})">
                <div class="schedule-slot-time">${timeStr}</div>
                <div class="schedule-slot-client">${esc(e.client_login)}</div>
                <div class="schedule-slot-title">${esc(e.title)}</div>
            </div>`;
        }

        html += '</div>';
    }

    if (scheduleEntries.length === 0) {
        html += '<div class="schedule-empty">Записей на эту неделю нет</div>';
    }

    html += '</div>';
    $('#schedule-grid').innerHTML = html;
}

// Schedule modal
async function openScheduleModal() {
    // Load clients for select
    try {
        const users = await api('GET', '/admin/users');
        const sel = $('#sch-client');
        sel.innerHTML = (users || [])
            .map((u) => `<option value="${u.id}">${esc(u.login)}</option>`)
            .join('');
    } catch (err) { toast(err.message, 'error'); return; }

    // Default date = today
    $('#sch-date').value = toISODate(new Date());
    $('#schedule-modal').classList.add('active');
}

function closeScheduleModal() { $('#schedule-modal').classList.remove('active'); }

$('#schedule-modal').addEventListener('click', (e) => {
    if (e.target === $('#schedule-modal')) closeScheduleModal();
});

$('#schedule-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const date = $('#sch-date').value;
    const time = $('#sch-time').value;
    if (!date || !time) { toast('Укажите дату и время', 'error'); return; }

    try {
        await api('POST', '/admin/schedule', {
            client_id: parseInt($('#sch-client').value),
            title: $('#sch-title').value,
            scheduled_at: date + 'T' + time,
            duration_minutes: parseInt($('#sch-duration').value) || 60,
            notes: $('#sch-notes').value,
        });
        closeScheduleModal();
        $('#schedule-form').reset();
        toast('Запись добавлена!');
        loadSchedule();
    } catch (err) { toast(err.message, 'error'); }
});

// Detail modal
function openSlotDetail(idx) {
    detailEntry = scheduleEntries[idx];
    if (!detailEntry) return;

    const at = new Date(detailEntry.scheduled_at);
    const timeStr = at.toLocaleString('ru-RU', {
        day: '2-digit', month: 'long', year: 'numeric',
        hour: '2-digit', minute: '2-digit',
    });

    const statusMap = { planned: 'Запланировано', completed: 'Завершено', cancelled: 'Отменено' };
    const statusColor = { planned: 'var(--primary)', completed: 'var(--success)', cancelled: 'var(--danger)' };

    $('#det-title').textContent = detailEntry.title;
    $('#det-info').innerHTML = `
        <div>👤 Клиент: <b>${esc(detailEntry.client_login)}</b></div>
        <div>📅 ${timeStr}</div>
        <div>⏱ ${detailEntry.duration_minutes} мин</div>
        <div>Статус: <span style="color:${statusColor[detailEntry.status]}">${statusMap[detailEntry.status] || detailEntry.status}</span></div>
        ${detailEntry.notes ? '<div>📝 ' + esc(detailEntry.notes) + '</div>' : ''}
    `;

    // Show/hide action buttons based on status
    const isPlanned = detailEntry.status === 'planned';
    $('#det-done-btn').style.display   = isPlanned ? '' : 'none';
    $('#det-cancel-btn').style.display = isPlanned ? '' : 'none';

    $('#schedule-detail-modal').classList.add('active');
}

function closeDetailModal() {
    $('#schedule-detail-modal').classList.remove('active');
    detailEntry = null;
}

$('#schedule-detail-modal').addEventListener('click', (e) => {
    if (e.target === $('#schedule-detail-modal')) closeDetailModal();
});

async function setScheduleStatus(status) {
    if (!detailEntry) return;
    try {
        await api('PUT', '/admin/schedule/' + detailEntry.id, { status });
        toast(status === 'completed' ? 'Отмечено как завершено' : 'Запись отменена');
        closeDetailModal();
        loadSchedule();
    } catch (err) { toast(err.message, 'error'); }
}

async function deleteScheduleEntry() {
    if (!detailEntry || !confirm('Удалить запись из расписания?')) return;
    try {
        await api('DELETE', '/admin/schedule/' + detailEntry.id);
        toast('Запись удалена');
        closeDetailModal();
        loadSchedule();
    } catch (err) { toast(err.message, 'error'); }
}

function goToClientFromDetail() {
    if (!detailEntry) return;
    const clientID = detailEntry.client_id;
    const login    = detailEntry.client_login;
    closeDetailModal();
    navigate('admin');
    // small delay to let the page render before calling drill-down
    setTimeout(() => loadAdminUserWorkouts(clientID, login), 50);
}

// ── Exercises ─────────────────────────────────────────────────────────────

let exerciseSearchTimer = null;

function loadExercises() {
    fetchExercises($('#exercise-filter')?.value || '', $('#exercise-search')?.value || '');
}

async function fetchExercises(group, search) {
    try {
        const params = new URLSearchParams();
        if (group)  params.set('muscle_group', group);
        if (search) params.set('search', search);
        const qs = params.toString() ? '?' + params.toString() : '';
        const exercises = await api('GET', '/exercises' + qs);
        window._exercises = exercises || [];
        renderExercises(exercises || []);
    } catch (err) { toast(err.message, 'error'); }
}

function renderExercises(exercises) {
    const list = $('#exercises-list');
    if (!exercises.length) {
        list.innerHTML = '<div class="card"><p>Упражнения не найдены.</p></div>';
    } else {
        list.innerHTML = `<div class="card"><table>
            <thead><tr><th>Название</th><th>Группа мышц</th><th>Описание</th>${currentUser?.role === 'admin' ? '<th></th>' : ''}</tr></thead>
            <tbody>${exercises.map((e) => `
                <tr>
                    <td>${esc(e.name)}</td>
                    <td><span class="badge">${esc(e.muscle_group)}</span></td>
                    <td style="color:var(--text-muted)">${esc(e.description)}</td>
                    ${currentUser?.role === 'admin' ? `<td><button class="btn btn-sm btn-danger" onclick="deleteExercise(${e.id})">×</button></td>` : ''}
                </tr>`).join('')}
            </tbody></table></div>`;
    }
    if (currentUser?.role === 'admin') $('#admin-exercise-form').style.display = '';
}

$('#exercise-search')?.addEventListener('input', () => {
    clearTimeout(exerciseSearchTimer);
    exerciseSearchTimer = setTimeout(() => {
        fetchExercises($('#exercise-filter').value, $('#exercise-search').value);
    }, 300);
});

$('#exercise-filter')?.addEventListener('change', () => {
    fetchExercises($('#exercise-filter').value, $('#exercise-search').value);
});

$('#exercise-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        await api('POST', '/exercises', {
            name: $('#ex-name').value,
            muscle_group: $('#ex-muscle').value,
            description: $('#ex-desc').value,
        });
        toast('Упражнение добавлено!');
        $('#ex-name').value = '';
        $('#ex-desc').value = '';
        loadExercises();
    } catch (err) { toast(err.message, 'error'); }
});

async function deleteExercise(id) {
    if (!confirm('Удалить упражнение?')) return;
    try {
        await api('DELETE', '/exercises/' + id);
        toast('Упражнение удалено');
        loadExercises();
    } catch (err) { toast(err.message, 'error'); }
}

// ── Stats ─────────────────────────────────────────────────────────────────

// Module-level state for charts so we can destroy/recreate them on reload.
let _chartWeight = null;
let _chartPR     = null;
// All workout exercises keyed by exercise_id for PR history lookups.
let _prHistory   = {}; // exercise_id -> [{date, weight_kg}]
// Exercise map id -> name
let _exMap       = {};

async function loadStats() {
    try {
        const [prData, volumeData, exercises, metrics] = await Promise.all([
            api('GET', '/stats/pr'),
            api('GET', '/stats/volume'),
            api('GET', '/exercises'),
            api('GET', '/metrics'),
        ]);

        // Build lookup map
        _exMap = {};
        (exercises || []).forEach((e) => { _exMap[e.id] = e.name; });

        // ── Summary cards ────────────────────────────────────────────────
        $('#stat-volume').textContent = Math.round(volumeData.weekly_volume).toLocaleString() + ' кг';

        // ── Personal records table ───────────────────────────────────────
        const prList = $('#pr-list');
        if (!prData || !prData.length) {
            prList.innerHTML = '<p style="color:var(--text-muted)">Пока нет рекордов</p>';
            $('#stat-pr-count').textContent = '0';
        } else {
            $('#stat-pr-count').textContent = prData.length;
            prList.innerHTML = `<table>
                <thead><tr><th>Упражнение</th><th>Вес</th><th>Подходы × Повторы</th></tr></thead>
                <tbody>${prData.map((r) => `<tr>
                    <td>${esc(_exMap[r.exercise_id] || '#' + r.exercise_id)}</td>
                    <td><span class="badge badge-pr">${r.weight_kg} кг</span></td>
                    <td>${r.sets}×${r.reps}</td>
                </tr>`).join('')}</tbody></table>`;
        }

        // ── PR-exercise select ───────────────────────────────────────────
        const sel = $('#pr-chart-exercise');
        const prevVal = sel.value;
        sel.innerHTML = '<option value="">— выберите упражнение —</option>';
        (prData || []).forEach((r) => {
            const opt = document.createElement('option');
            opt.value = r.exercise_id;
            opt.textContent = _exMap[r.exercise_id] || '#' + r.exercise_id;
            sel.appendChild(opt);
        });
        if (prevVal) sel.value = prevVal;

        // ── Charts ───────────────────────────────────────────────────────
        renderWeightChart(metrics || []);
        renderPRChart(sel.value ? parseInt(sel.value) : null);

        // ── Metrics table ────────────────────────────────────────────────
        renderMetricsTable(metrics || []);

    } catch (err) { toast(err.message, 'error'); }
}

// ── Body-weight chart ──────────────────────────────────────────────────────

function renderWeightChart(metrics) {
    const withWeight = metrics.filter((m) => m.weight_kg != null);
    const canvas = $('#chart-weight');
    const empty  = $('#chart-weight-empty');

    if (_chartWeight) { _chartWeight.destroy(); _chartWeight = null; }

    if (!withWeight.length) {
        canvas.style.display = 'none';
        empty.style.display  = '';
        return;
    }
    canvas.style.display = '';
    empty.style.display  = 'none';

    _chartWeight = new Chart(canvas, {
        type: 'line',
        data: {
            labels: withWeight.map((m) => formatDate(m.measured_at)),
            datasets: [{
                label: 'Вес (кг)',
                data: withWeight.map((m) => m.weight_kg),
                borderColor: '#6c63ff',
                backgroundColor: '#6c63ff22',
                tension: 0.3,
                fill: true,
                pointRadius: 4,
            }],
        },
        options: {
            responsive: true,
            plugins: { legend: { display: false } },
            scales: {
                x: { ticks: { color: '#888', maxTicksLimit: 8 }, grid: { color: '#333' } },
                y: { ticks: { color: '#888' }, grid: { color: '#333' } },
            },
        },
    });
}

// ── PR-by-exercise chart ───────────────────────────────────────────────────

async function renderPRChart(exerciseID) {
    const canvas = $('#chart-pr');
    const empty  = $('#chart-pr-empty');

    if (_chartPR) { _chartPR.destroy(); _chartPR = null; }

    if (!exerciseID) {
        canvas.style.display = 'none';
        empty.style.display  = '';
        return;
    }

    try {
        const points = await api('GET', '/stats/exercise-progress?exercise_id=' + exerciseID);

        if (!points || !points.length) {
            canvas.style.display = 'none';
            empty.style.display  = '';
            return;
        }
        canvas.style.display = '';
        empty.style.display  = 'none';

        _chartPR = new Chart(canvas, {
            type: 'line',
            data: {
                labels: points.map((p) => formatDate(p.date)),
                datasets: [{
                    label: _exMap[exerciseID] || 'Упражнение',
                    data: points.map((p) => p.max_weight),
                    borderColor: '#2ecc71',
                    backgroundColor: '#2ecc7122',
                    tension: 0.3,
                    fill: true,
                    pointRadius: 4,
                }],
            },
            options: {
                responsive: true,
                plugins: { legend: { display: false } },
                scales: {
                    x: { ticks: { color: '#888', maxTicksLimit: 8 }, grid: { color: '#333' } },
                    y: { ticks: { color: '#888', callback: (v) => v + ' кг' }, grid: { color: '#333' } },
                },
            },
        });
    } catch (err) { toast(err.message, 'error'); }
}

$('#pr-chart-exercise').addEventListener('change', (e) => {
    const id = e.target.value ? parseInt(e.target.value) : null;
    renderPRChart(id);
});

// ── Body metrics table ─────────────────────────────────────────────────────

function renderMetricsTable(metrics) {
    const wrap  = $('#metrics-table-wrap');
    const empty = $('#metrics-empty');
    if (!metrics.length) {
        wrap.innerHTML = '<p id="metrics-empty" style="color:var(--text-muted)">Замеров пока нет.</p>';
        return;
    }

    // Show newest first (metrics are sorted oldest-first by API for charts)
    const rows = [...metrics].reverse().slice(0, 10);

    // Trend arrow: compare each row to the previous chronological entry
    const trend = (curr, prev, field) => {
        if (curr[field] == null || prev == null || prev[field] == null) return '';
        const diff = curr[field] - prev[field];
        if (Math.abs(diff) < 0.05) return '';
        return diff > 0
            ? `<span class="trend-up">▲ ${diff.toFixed(1)}</span>`
            : `<span class="trend-down">▼ ${Math.abs(diff).toFixed(1)}</span>`;
    };

    wrap.innerHTML = `<table class="metrics-table">
        <thead><tr>
            <th>Дата</th>
            <th>Вес</th>
            <th>% жира</th>
            <th>Грудь</th>
            <th>Талия</th>
            <th>Бёдра</th>
            <th>Бицепс</th>
            <th></th>
        </tr></thead>
        <tbody>${rows.map((m, i) => {
            const prev = rows[i + 1] || null; // older entry (we reversed above)
            return `<tr>
                <td>${formatDate(m.measured_at)}</td>
                <td>${m.weight_kg != null ? m.weight_kg + ' кг' : '—'} ${trend(m, prev, 'weight_kg')}</td>
                <td>${m.body_fat_percent != null ? m.body_fat_percent + '%' : '—'} ${trend(m, prev, 'body_fat_percent')}</td>
                <td>${m.chest_cm != null ? m.chest_cm : '—'}</td>
                <td>${m.waist_cm != null ? m.waist_cm : '—'} ${trend(m, prev, 'waist_cm')}</td>
                <td>${m.hips_cm != null ? m.hips_cm : '—'}</td>
                <td>${m.bicep_cm != null ? m.bicep_cm : '—'} ${trend(m, prev, 'bicep_cm')}</td>
                <td><button class="btn btn-sm btn-danger" onclick="deleteMetric(${m.id})">×</button></td>
            </tr>`;
        }).join('')}</tbody>
    </table>`;
}

// ── Metric form ────────────────────────────────────────────────────────────

function toggleMetricForm() {
    const wrap = $('#metric-form-wrap');
    const show = wrap.style.display === 'none';
    wrap.style.display = show ? '' : 'none';
    if (show) {
        // Default to today
        const today = new Date();
        $('#m-date').value = toISODate ? toISODate(today) : today.toISOString().slice(0, 10);
    }
}

$('#metric-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const numOrNull = (id) => {
        const v = parseFloat($(id).value);
        return isNaN(v) ? null : v;
    };
    const body = {
        measured_at:      $('#m-date').value || new Date().toISOString().slice(0, 10),
        weight_kg:        numOrNull('#m-weight'),
        body_fat_percent: numOrNull('#m-fat'),
        chest_cm:         numOrNull('#m-chest'),
        waist_cm:         numOrNull('#m-waist'),
        hips_cm:          numOrNull('#m-hips'),
        bicep_cm:         numOrNull('#m-bicep'),
    };
    try {
        await api('POST', '/metrics', body);
        toast('Замер сохранён!');
        $('#metric-form').reset();
        toggleMetricForm();
        loadStats();
    } catch (err) { toast(err.message, 'error'); }
});

async function deleteMetric(id) {
    if (!confirm('Удалить замер?')) return;
    try {
        await api('DELETE', '/metrics/' + id);
        toast('Замер удалён');
        loadStats();
    } catch (err) { toast(err.message, 'error'); }
}

// ── Admin panel — Клиенты ─────────────────────────────────────────────────

function switchAdminTab(tab) {
    $$('.subtab').forEach((b) => b.classList.toggle('active', b.dataset.subtab === tab));
    // Currently only 'clients' exists — extend here when adding more tabs
    if (tab === 'clients') loadAdminClients();
}

async function loadAdminClients() {
    // Show users view, hide workouts drill-down
    $('#admin-users-view').style.display = '';
    $('#admin-workouts-view').style.display = 'none';

    try {
        const users = await api('GET', '/admin/users');
        const list = $('#admin-users-list');
        if (!users || !users.length) {
            list.innerHTML = '<div class="card"><p>Нет пользователей.</p></div>';
            return;
        }
        list.innerHTML = users
            .filter((u) => u.role !== 'admin') // hide trainer account from list
            .map((u) => `
            <div class="card" style="cursor:pointer" onclick="loadAdminUserWorkouts(${u.id}, '${esc(u.login)}')">
                <div class="card-header">
                    <div>
                        <div class="card-title">${esc(u.login)}</div>
                        <div class="card-subtitle">${esc(u.email)}</div>
                    </div>
                    <div class="card-actions">
                        <span class="badge">Тренировки →</span>
                    </div>
                </div>
            </div>`).join('');
    } catch (err) { toast(err.message, 'error'); }
}

// State for the currently open client drill-down (used to refresh after a comment save)
let adminWorkouts = [];
let adminCurrentUserID = null;
let adminCurrentLogin = '';

async function loadAdminUserWorkouts(userID, login) {
    $('#admin-users-view').style.display    = 'none';
    $('#admin-workouts-view').style.display = '';
    $('#admin-user-title').textContent = 'Тренировки: ' + login;
    adminCurrentUserID = userID;
    adminCurrentLogin  = login;

    // Reset to workouts tab and destroy stale charts from a previous client
    switchClientTab('workouts');
    if (_clientChartWeight) { _clientChartWeight.destroy(); _clientChartWeight = null; }
    if (_clientChartPR)     { _clientChartPR.destroy();     _clientChartPR = null; }

    try {
        const workouts = await api('GET', '/admin/users/' + userID + '/workouts');
        adminWorkouts = workouts || [];
        const list = $('#admin-workouts-list');
        if (!adminWorkouts.length) {
            list.innerHTML = '<div class="card"><p>Нет тренировок.</p></div>';
            return;
        }
        // Pass the workout's index (not the comment text) to avoid breaking the
        // onclick attribute — the comment is looked up from adminWorkouts on open.
        list.innerHTML = adminWorkouts.map((w, i) => `
            <div class="card">
                <div class="card-header">
                    <div>
                        <div class="card-title">${esc(w.title)}</div>
                        <div class="card-subtitle">${formatDate(w.date)} · ${w.duration_minutes} мин</div>
                    </div>
                    <div class="card-actions">
                        <button class="btn btn-sm btn-outline" onclick="openCommentModal(${i})">
                            ${w.trainer_comment ? 'Изменить комментарий' : '+ Комментарий'}
                        </button>
                    </div>
                </div>
                ${w.notes ? '<p style="color:var(--text-muted);font-size:0.85rem">' + esc(w.notes) + '</p>' : ''}
                ${w.trainer_comment ? '<div class="trainer-comment">💬 ' + esc(w.trainer_comment) + '</div>' : ''}
            </div>`).join('');
    } catch (err) { toast(err.message, 'error'); }
}

function showAdminUsers() {
    $('#admin-users-view').style.display    = '';
    $('#admin-workouts-view').style.display = 'none';
}

// ── Client sub-tabs (Тренировки / Прогресс) ───────────────────────────────

let _clientChartWeight = null;
let _clientChartPR     = null;

function switchClientTab(tab) {
    $$('#admin-workouts-view .subtab').forEach((b) =>
        b.classList.toggle('active', b.id === 'client-tab-' + tab));
    $('#admin-workouts-list').style.display   = tab === 'workouts' ? '' : 'none';
    $('#admin-progress-view').style.display   = tab === 'progress' ? '' : 'none';
    if (tab === 'progress') loadClientProgress();
}

async function loadClientProgress() {
    if (adminCurrentUserID === null) return;
    try {
        const [metrics, exercises] = await Promise.all([
            api('GET', '/admin/users/' + adminCurrentUserID + '/metrics'),
            api('GET', '/exercises'),
        ]);

        // Weight chart
        renderClientWeightChart(metrics || []);

        // Populate exercise select from all exercises (not just PRs — client may
        // have data for any of them)
        const sel = $('#client-pr-exercise');
        const prev = sel.value;
        sel.innerHTML = '<option value="">— выберите упражнение —</option>';
        (exercises || []).forEach((e) => {
            const opt = document.createElement('option');
            opt.value = e.id;
            opt.textContent = e.name;
            sel.appendChild(opt);
        });
        if (prev) sel.value = prev;

        if (sel.value) renderClientPRChart(parseInt(sel.value));

    } catch (err) { toast(err.message, 'error'); }
}

function renderClientWeightChart(metrics) {
    const withWeight = (metrics || []).filter((m) => m.weight_kg != null);
    const canvas = $('#chart-client-weight');
    const empty  = $('#chart-client-weight-empty');
    if (_clientChartWeight) { _clientChartWeight.destroy(); _clientChartWeight = null; }

    if (!withWeight.length) {
        canvas.style.display = 'none';
        empty.style.display  = '';
        return;
    }
    canvas.style.display = '';
    empty.style.display  = 'none';

    _clientChartWeight = new Chart(canvas, {
        type: 'line',
        data: {
            labels: withWeight.map((m) => formatDate(m.measured_at)),
            datasets: [{
                label: 'Вес (кг)',
                data: withWeight.map((m) => m.weight_kg),
                borderColor: '#6c63ff',
                backgroundColor: '#6c63ff22',
                tension: 0.3,
                fill: true,
                pointRadius: 4,
            }],
        },
        options: {
            responsive: true,
            plugins: { legend: { display: false } },
            scales: {
                x: { ticks: { color: '#888', maxTicksLimit: 8 }, grid: { color: '#333' } },
                y: { ticks: { color: '#888' }, grid: { color: '#333' } },
            },
        },
    });
}

async function renderClientPRChart(exerciseID) {
    const canvas = $('#chart-client-pr');
    const empty  = $('#chart-client-pr-empty');
    if (_clientChartPR) { _clientChartPR.destroy(); _clientChartPR = null; }

    if (!exerciseID || adminCurrentUserID === null) {
        canvas.style.display = 'none';
        empty.style.display  = '';
        return;
    }
    try {
        const points = await api('GET',
            '/admin/users/' + adminCurrentUserID + '/exercise-progress?exercise_id=' + exerciseID);

        if (!points || !points.length) {
            canvas.style.display = 'none';
            empty.style.display  = '';
            return;
        }
        canvas.style.display = '';
        empty.style.display  = 'none';

        const exName = $('#client-pr-exercise').selectedOptions[0]?.textContent || 'Упражнение';
        _clientChartPR = new Chart(canvas, {
            type: 'line',
            data: {
                labels: points.map((p) => formatDate(p.date)),
                datasets: [{
                    label: exName,
                    data: points.map((p) => p.max_weight),
                    borderColor: '#2ecc71',
                    backgroundColor: '#2ecc7122',
                    tension: 0.3,
                    fill: true,
                    pointRadius: 4,
                }],
            },
            options: {
                responsive: true,
                plugins: { legend: { display: false } },
                scales: {
                    x: { ticks: { color: '#888', maxTicksLimit: 8 }, grid: { color: '#333' } },
                    y: { ticks: { color: '#888', callback: (v) => v + ' кг' }, grid: { color: '#333' } },
                },
            },
        });
    } catch (err) { toast(err.message, 'error'); }
}

$('#client-pr-exercise').addEventListener('change', (e) => {
    const id = e.target.value ? parseInt(e.target.value) : null;
    renderClientPRChart(id);
});

// Comment modal
function openCommentModal(idx) {
    const w = adminWorkouts[idx];
    if (!w) return;
    $('#comment-workout-id').value = w.id;
    $('#comment-text').value = w.trainer_comment || '';
    $('#comment-modal').classList.add('active');
}

function closeCommentModal() { $('#comment-modal').classList.remove('active'); }

$('#comment-modal').addEventListener('click', (e) => {
    if (e.target === $('#comment-modal')) closeCommentModal();
});

async function submitComment() {
    const workoutId = $('#comment-workout-id').value;
    const comment   = $('#comment-text').value;
    try {
        await api('PUT', '/admin/workouts/' + workoutId + '/comment', { comment });
        toast('Комментарий сохранён!');
        closeCommentModal();
        // Reload the current client's workouts so the comment shows on the card.
        if (adminCurrentUserID !== null) {
            loadAdminUserWorkouts(adminCurrentUserID, adminCurrentLogin);
        }
    } catch (err) { toast(err.message, 'error'); }
}

// ── Helpers ───────────────────────────────────────────────────────────────

function esc(s) {
    if (!s) return '';
    const d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
}

function formatDate(s) {
    return new Date(s).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' });
}

// ── Init ──────────────────────────────────────────────────────────────────
initAuth();
